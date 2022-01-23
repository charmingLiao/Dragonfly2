/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"time"

	"github.com/go-http-utils/headers"
	"go.opentelemetry.io/otel/propagation"

	"d7y.io/dragonfly/v2/client/config"
	"d7y.io/dragonfly/v2/client/daemon/metrics"
	"d7y.io/dragonfly/v2/client/daemon/peer"
	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/idgen"
	"d7y.io/dragonfly/v2/pkg/rpc/base"
	"d7y.io/dragonfly/v2/pkg/rpc/scheduler"
	"d7y.io/dragonfly/v2/pkg/util/net/httputils"
)

var _ *logger.SugaredLoggerOnWith // pin this package for no log code generation

var (
	// layerReg the regex to determine if it is an image download
	layerReg     = regexp.MustCompile("^.+/blobs/sha256.*$")
	traceContext = propagation.TraceContext{}
)

// transport implements RoundTripper for dragonfly.
// It uses http.fileTransport to serve requests that need to use dragonfly,
// and uses http.Transport to serve the other requests.
type transport struct {
	// baseRoundTripper is an implementation of RoundTripper that supports HTTP
	baseRoundTripper http.RoundTripper

	// shouldUseDragonfly is used to determine to download resources with or without dragonfly
	shouldUseDragonfly func(req *http.Request) bool

	// peerTaskManager is the peer task manager
	peerTaskManager peer.TaskManager

	// peerHost is the peer host info
	peerHost *scheduler.PeerHost

	// defaultFilter is used when http request without X-Dragonfly-Filter Header
	defaultFilter string

	// defaultBiz is used when http request without X-Dragonfly-Biz Header
	defaultBiz string

	// dumpHTTPContent indicates to dump http request header and response header
	dumpHTTPContent bool
}

// Option is functional config for transport.
type Option func(rt *transport) *transport

// WithPeerHost sets the peerHost for transport
func WithPeerHost(peerHost *scheduler.PeerHost) Option {
	return func(rt *transport) *transport {
		rt.peerHost = peerHost
		return rt
	}
}

// WithPeerTaskManager sets the peerTaskManager for transport
func WithPeerTaskManager(peerTaskManager peer.TaskManager) Option {
	return func(rt *transport) *transport {
		rt.peerTaskManager = peerTaskManager
		return rt
	}
}

// WithTLS configures TLS config used for http transport.
func WithTLS(cfg *tls.Config) Option {
	return func(rt *transport) *transport {
		rt.baseRoundTripper = defaultHTTPTransport(cfg)
		return rt
	}
}

// WithCondition configures how to decide whether to use dragonfly or not.
func WithCondition(c func(r *http.Request) bool) Option {
	return func(rt *transport) *transport {
		rt.shouldUseDragonfly = c
		return rt
	}
}

// WithDefaultFilter sets default filter for http requests with X-Dragonfly-Filter Header
func WithDefaultFilter(f string) Option {
	return func(rt *transport) *transport {
		rt.defaultFilter = f
		return rt
	}
}

// WithDefaultBiz sets default biz for http requests with X-Dragonfly-Biz Header
func WithDefaultBiz(b string) Option {
	return func(rt *transport) *transport {
		rt.defaultBiz = b
		return rt
	}
}

func WithDumpHTTPContent(b bool) Option {
	return func(rt *transport) *transport {
		rt.dumpHTTPContent = b
		return rt
	}
}

// New constructs a new instance of a RoundTripper with additional options.
func New(options ...Option) (http.RoundTripper, error) {
	rt := &transport{
		baseRoundTripper:   defaultHTTPTransport(nil),
		shouldUseDragonfly: NeedUseDragonfly,
	}

	for _, opt := range options {
		opt(rt)
	}

	return rt, nil
}

// RoundTrip only process first redirect at present
func (rt *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if rt.shouldUseDragonfly(req) {
		// delete the Accept-Encoding header to avoid returning the same cached
		// result for different requests
		req.Header.Del("Accept-Encoding")

		ctx := req.Context()
		if req.URL.Scheme == "https" {
			// for https, the trace info is in request header
			ctx = traceContext.Extract(req.Context(), propagation.HeaderCarrier(req.Header))
		}

		logger.Debugf("round trip with dragonfly: %s", req.URL.String())
		metrics.ProxyRequestViaDragonflyCount.Add(1)
		resp, err = rt.download(ctx, req)
	} else {
		logger.Debugf("round trip directly, method: %s, url: %s", req.Method, req.URL.String())
		req.Host = req.URL.Host
		req.Header.Set("Host", req.Host)
		metrics.ProxyRequestNotViaDragonflyCount.Add(1)
		resp, err = rt.baseRoundTripper.RoundTrip(req)
	}

	if err != nil {
		return resp, err
	}

	if resp.ContentLength > 0 {
		metrics.ProxyRequestBytesCount.WithLabelValues(req.Method).Add(float64(resp.ContentLength))
	}
	if err != nil {
		logger.With("method", req.Method, "url", req.URL.String()).
			Errorf("round trip error: %s", err)
	}
	rt.processDumpHTTPContent(req, resp)
	return resp, err
}

// NeedUseDragonfly is the default value for shouldUseDragonfly, which downloads all
// images layers with dragonfly.
func NeedUseDragonfly(req *http.Request) bool {
	return req.Method == http.MethodGet && layerReg.MatchString(req.URL.Path)
}

// download uses dragonfly to download.
// the ctx has span info from transport, did not use the ctx from request
func (rt *transport) download(ctx context.Context, req *http.Request) (*http.Response, error) {
	url := req.URL.String()
	peerID := idgen.PeerID(rt.peerHost.Ip)
	log := logger.With("peer", peerID, "component", "transport")
	log.Infof("start download with url: %s", url)

	// Init meta value
	meta := &base.UrlMeta{Header: map[string]string{}}

	// Set meta range's value
	if rg := req.Header.Get("Range"); len(rg) > 0 {
		meta.Digest = ""
		meta.Range = rg
	}

	// Pick header's parameters
	filter := httputils.PickHeader(req.Header, config.HeaderDragonflyFilter, rt.defaultFilter)
	tag := httputils.PickHeader(req.Header, config.HeaderDragonflyBiz, rt.defaultBiz)

	// Delete hop-by-hop headers
	delHopHeaders(req.Header)

	meta.Header = httputils.HeaderToMap(req.Header)
	meta.Tag = tag
	meta.Filter = filter

	body, attr, err := rt.peerTaskManager.StartStreamTask(
		ctx,
		&peer.StreamTaskRequest{
			URL:     url,
			URLMeta: meta,
			PeerID:  peerID,
		},
	)
	if err != nil {
		log.Errorf("download fail: %v", err)
		// add more info for debugging
		if attr != nil {
			err = fmt.Errorf("task: %s\npeer: %s\nerror: %s",
				attr[config.HeaderDragonflyTask], attr[config.HeaderDragonflyPeer], err)
		}
		return nil, err
	}

	hdr := httputils.MapToHeader(attr)
	log.Infof("download stream attribute: %v", hdr)

	var contentLength int64 = -1
	if l, ok := attr[headers.ContentLength]; ok {
		if i, e := strconv.ParseInt(l, 10, 64); e == nil {
			contentLength = i
		}
	}

	resp := &http.Response{
		StatusCode:    http.StatusOK,
		Body:          body,
		Header:        hdr,
		ContentLength: contentLength,

		Proto:      req.Proto,
		ProtoMajor: req.ProtoMajor,
		ProtoMinor: req.ProtoMinor,
	}
	return resp, nil
}

func (rt *transport) processDumpHTTPContent(req *http.Request, resp *http.Response) {
	if !rt.dumpHTTPContent {
		return
	}
	if out, e := httputil.DumpRequest(req, false); e == nil {
		logger.Debugf("dump request in transport: %s", string(out))
	} else {
		logger.Errorf("dump request in transport error: %s", e)
	}
	if resp == nil {
		return
	}
	if out, e := httputil.DumpResponse(resp, false); e == nil {
		logger.Debugf("dump response in transport: %s", string(out))
	} else {
		logger.Errorf("dump response in transport error: %s", e)
	}
}

func defaultHTTPTransport(cfg *tls.Config) *http.Transport {
	if cfg == nil {
		cfg = &tls.Config{InsecureSkipVerify: true}
	}

	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       cfg,
	}
}

// Hop-by-hop headers. These are removed when sent to the backend.
// As of RFC 7230, hop-by-hop headers are required to appear in the
// Connection header field. These are the headers defined by the
// obsoleted RFC 2616 (section 13.5.1) and are used for backward
// compatibility.
// copy from net/http/httputil/reverseproxy.go

// dragonfly need generate task id with header, need to remove some other headers
var hopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",

	// remove by dragonfly
	"Accept",
	"User-Agent",
	"X-Forwarded-For",
}

// delHopHeaders delete hop-by-hop headers.
func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
	// remove correlation with trace header
	for _, h := range traceContext.Fields() {
		header.Del(h)
	}
}
