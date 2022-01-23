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

package peer

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"d7y.io/dragonfly/v2/client/daemon/storage"
	"d7y.io/dragonfly/v2/client/daemon/upload"
	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/rpc/base"
	"d7y.io/dragonfly/v2/pkg/util/digestutils"
)

type DownloadPieceRequest struct {
	piece      *base.PieceInfo
	log        *logger.SugaredLoggerOnWith
	storage    storage.TaskStorageDriver
	TaskID     string
	PeerID     string
	DstPid     string
	DstAddr    string
	CalcDigest bool
}

type DownloadPieceResult struct {
	// Size of piece
	Size int64
	// BeginTime nanosecond
	BeginTime int64
	// FinishTime nanosecond
	FinishTime int64
}

//go:generate mockgen -source piece_downloader.go -destination ../test/mock/peer/piece_downloader.go
type PieceDownloader interface {
	DownloadPiece(context.Context, *DownloadPieceRequest) (io.Reader, io.Closer, error)
}

type pieceDownloader struct {
	transport  http.RoundTripper
	httpClient *http.Client
}

type pieceDownloadError struct {
	connectionError bool
	status          string
	statusCode      int
	target          string
	err             error
}

func isConnectionError(err error) bool {
	if e, ok := err.(*pieceDownloadError); ok {
		return e.connectionError
	}
	return false
}

func isPieceNotFound(err error) bool {
	if e, ok := err.(*pieceDownloadError); ok {
		return e.statusCode == http.StatusNotFound
	}
	return false
}

func (e *pieceDownloadError) Error() string {
	if e.connectionError {
		return fmt.Sprintf("connect with %s with error: %s", e.target, e.err)
	}
	return fmt.Sprintf("download %s with error status: %s", e.target, e.status)
}

var _ PieceDownloader = (*pieceDownloader)(nil)

var defaultTransport http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   2 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	ResponseHeaderTimeout: 2 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 2 * time.Second,
}

func NewPieceDownloader(timeout time.Duration, opts ...func(*pieceDownloader) error) (PieceDownloader, error) {
	pd := &pieceDownloader{}

	for _, opt := range opts {
		if err := opt(pd); err != nil {
			return nil, err
		}
	}

	if pd.transport == nil {
		pd.transport = defaultTransport
	}

	pd.httpClient = &http.Client{
		Transport: pd.transport,
		Timeout:   timeout,
	}
	return pd, nil
}

func WithTransport(rt http.RoundTripper) func(*pieceDownloader) error {
	return func(d *pieceDownloader) error {
		d.transport = rt
		return nil
	}
}

func (p *pieceDownloader) DownloadPiece(ctx context.Context, req *DownloadPieceRequest) (io.Reader, io.Closer, error) {
	resp, err := p.httpClient.Do(buildDownloadPieceHTTPRequest(ctx, req))
	if err != nil {
		logger.Errorf("task id: %s, piece num: %d, dst: %s, download piece failed: %s",
			req.TaskID, req.piece.PieceNum, req.DstAddr, err)
		return nil, nil, &pieceDownloadError{err: err, connectionError: true}
	}
	if resp.StatusCode > 299 {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		return nil, nil, &pieceDownloadError{err: err, connectionError: false, status: resp.Status, statusCode: resp.StatusCode}
	}
	reader, closer := resp.Body.(io.Reader), resp.Body.(io.Closer)
	if req.CalcDigest {
		req.log.Debugf("calculate digest for piece %d, digest: %s", req.piece.PieceNum, req.piece.PieceMd5)
		reader = digestutils.NewDigestReader(req.log, io.LimitReader(resp.Body, int64(req.piece.RangeSize)), req.piece.PieceMd5)
	}
	return reader, closer, nil
}

func buildDownloadPieceHTTPRequest(ctx context.Context, d *DownloadPieceRequest) *http.Request {
	b := strings.Builder{}
	// FIXME switch to https when tls enabled
	b.WriteString("http://")
	b.WriteString(d.DstAddr)
	b.WriteString(upload.PeerDownloadHTTPPathPrefix)
	b.Write([]byte(d.TaskID)[:3])
	b.Write([]byte("/"))
	b.WriteString(d.TaskID)
	b.Write([]byte("?peerId="))
	b.WriteString(d.DstPid)

	u := b.String()
	logger.Debugf("built request url: %s", u)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	// TODO use string.Builder
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d",
		d.piece.RangeStart, d.piece.RangeStart+uint64(d.piece.RangeSize)-1))
	return req
}
