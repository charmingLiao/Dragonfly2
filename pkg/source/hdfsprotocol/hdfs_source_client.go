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

package hdfsprotocol

import (
	"io"
	"net/url"
	"os/user"
	"strings"
	"sync"
	"time"

	"github.com/colinmarc/hdfs/v2"
	"github.com/pkg/errors"

	"d7y.io/dragonfly/v2/pkg/source"
	"d7y.io/dragonfly/v2/pkg/util/rangeutils"
	"d7y.io/dragonfly/v2/pkg/util/timeutils"
)

const (
	HDFSClient = "hdfs"
)
const (
	// hdfsUseDataNodeHostName set hdfs client whether user hostname connect to datanode
	hdfsUseDataNodeHostName = "dfs.client.use.datanode.hostname"
	// hdfsUseDataNodeHostNameValue set value is true
	hdfsUseDataNodeHostNameValue = "true"
)

func init() {
	if err := source.Register(HDFSClient, NewHDFSSourceClient(), adapter); err != nil {
		panic(err)
	}
}

func adapter(request *source.Request) *source.Request {
	clonedRequest := request.Clone(request.Context())
	return clonedRequest
}

// hdfsSourceClient is an implementation of the interface of SourceClient.
type hdfsSourceClient struct {
	sync.RWMutex
	clientMap map[string]*hdfs.Client
}

// hdfsFileReaderClose is a combination object of the  io.LimitedReader and io.Closer
type hdfsFileReaderClose struct {
	limitedReader io.Reader
	closer        io.Closer
}

func newHdfsFileReaderClose(r io.ReadCloser, n int64) io.ReadCloser {
	return &hdfsFileReaderClose{
		limitedReader: io.LimitReader(r, n),
		closer:        r,
	}
}

type HDFSSourceClientOption func(p *hdfsSourceClient)

func (h *hdfsSourceClient) GetContentLength(request *source.Request) (int64, error) {
	hdfsClient, path, err := h.getHDFSClientAndPath(request.URL)
	if err != nil {
		return source.UnknownSourceFileLen, err
	}
	info, err := hdfsClient.Stat(path)
	if err != nil {
		return source.UnknownSourceFileLen, err
	}
	return info.Size(), nil
}

func (h *hdfsSourceClient) IsSupportRange(request *source.Request) (bool, error) {
	hdfsClient, path, err := h.getHDFSClientAndPath(request.URL)
	if err != nil {
		return false, err
	}
	_, err = hdfsClient.Stat(path)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *hdfsSourceClient) IsExpired(request *source.Request, info *source.ExpireInfo) (bool, error) {
	hdfsClient, path, err := h.getHDFSClientAndPath(request.URL)
	if err != nil {
		return false, err
	}

	fileInfo, err := hdfsClient.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.ModTime().Format(source.LastModifiedLayout) != info.LastModified, nil
}

func (h *hdfsSourceClient) Download(request *source.Request) (*source.Response, error) {
	hdfsClient, path, err := h.getHDFSClientAndPath(request.URL)
	if err != nil {
		return nil, err
	}

	hdfsFile, err := hdfsClient.Open(path)
	if err != nil {
		return nil, err
	}

	fileInfo := hdfsFile.Stat()

	// default read all data when rang is nil
	var limitReadN = fileInfo.Size()
	if limitReadN < 0 {
		return nil, errors.Errorf("file length is illegal, length: %d", limitReadN)
	}

	if request.Header.Get(source.Range) != "" {
		requestRange, err := rangeutils.ParseRange(request.Header.Get(source.Range), uint64(limitReadN))
		if err != nil {
			return nil, err
		}
		_, err = hdfsFile.Seek(int64(requestRange.StartIndex), 0)
		if err != nil {
			hdfsFile.Close()
			return nil, err
		}
		limitReadN = int64(requestRange.Length())
	}

	response := source.NewResponse(
		newHdfsFileReaderClose(hdfsFile, limitReadN),
		source.WithExpireInfo(source.ExpireInfo{
			LastModified: timeutils.Format(fileInfo.ModTime()),
		}))
	return response, nil
}

func (h *hdfsSourceClient) GetLastModified(request *source.Request) (int64, error) {

	hdfsClient, path, err := h.getHDFSClientAndPath(request.URL)
	if err != nil {
		return -1, err
	}

	info, err := hdfsClient.Stat(path)
	if err != nil {
		return -1, err
	}

	return info.ModTime().UnixNano() / time.Millisecond.Nanoseconds(), nil
}

// getHDFSClient return hdfs client
func (h *hdfsSourceClient) getHDFSClient(url *url.URL) (*hdfs.Client, error) {
	// get client for map
	h.RWMutex.RLock()
	if client, ok := h.clientMap[url.Host]; ok {
		h.RWMutex.RUnlock()
		return client, nil
	}
	h.RWMutex.RUnlock()

	// create client option
	options := hdfs.ClientOptionsFromConf(map[string]string{
		hdfsUseDataNodeHostName: hdfsUseDataNodeHostNameValue,
	})
	options.Addresses = strings.Split(url.Host, ",")
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	options.User = u.Username

	// create hdfs client and put map
	h.RWMutex.Lock()
	client, err := hdfs.NewClient(options)
	if err != nil {
		h.RWMutex.Unlock()
		return nil, err
	}
	h.clientMap[url.Host] = client
	h.RWMutex.Unlock()
	return client, err
}

// getHDFSClientAndPath return client and path
func (h *hdfsSourceClient) getHDFSClientAndPath(url *url.URL) (*hdfs.Client, string, error) {
	client, err := h.getHDFSClient(url)
	if err != nil {
		return nil, "", errors.Errorf("hdfs create client failed, url is %s", url)
	}
	return client, url.Path, nil
}

func NewHDFSSourceClient(opts ...HDFSSourceClientOption) source.ResourceClient {
	return newHDFSSourceClient(opts...)
}

func newHDFSSourceClient(opts ...HDFSSourceClientOption) *hdfsSourceClient {
	sourceClient := &hdfsSourceClient{
		clientMap: make(map[string]*hdfs.Client),
	}
	for i := range opts {
		opts[i](sourceClient)
	}
	return sourceClient
}

var _ source.ResourceClient = (*hdfsSourceClient)(nil)

func (rc *hdfsFileReaderClose) Read(p []byte) (n int, err error) {
	return rc.limitedReader.Read(p)
}

func (rc *hdfsFileReaderClose) Close() error {
	return rc.closer.Close()
}
