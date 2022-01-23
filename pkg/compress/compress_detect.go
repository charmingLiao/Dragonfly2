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

package compress

import (
	"fmt"
	"io"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/synclock"
)

const (
	// DefaultMinPiecesNum min pieces num
	// because the first three are used for detection
	DefaultMinPiecesNum = 6

	// DefaultDetectPiecesNum detect pieces num
	DefaultDetectPiecesNum = 3

	// KeyPrefix cache key
	KeyPrefix = "CompressDetect"
)

var DefaultCompress *CompressDetect

type CompressConfig struct {
	// Ratio Compress ratio
	Ratio float32 `yaml:"ratio" mapstructure:"ratio"`

	// Algorithm compress algorithm
	Algorithm Algorithm `yaml:"algorithm" mapstructure:"algorithm"`

	// DetectChanSize
	DetectChanSize int `yaml:"detectChanSize" mapstructure:"detectChanSize"`

	// ConcurrentSize compress detect goroutine size
	ConcurrentSize int `yaml:"concurrentSize" mapstructure:"concurrentSize"`
}

type CompressDetect struct {
	dataCh                 chan CompressData
	config                 CompressConfig
	detectResult           cache.Cache
	compressMap            map[Algorithm]Compressor
	compress               Compressor
	concurrentSize         int
	recordDetectNumPerTask cache.Cache
}

type CompressData struct {
	taskID string
	data   *[]byte
}

func NewDefaultCompressDetect(compressConfig CompressConfig) *CompressDetect {
	compress := NewCompress(compressConfig.Algorithm)
	compressMap := make(map[CompressAlgorithm]Compress)
	compressMap[compressConfig.Algorithm] = compress
	DefaultCompress = &CompressDetect{
		dataCh:                 make(chan CompressData, compressConfig.DetectChanSize),
		config:                 compressConfig,
		detectResult:           cache.New(0, 0),
		compressMap:            compressMap,
		compress:               compress,
		concurrentSize:         compressConfig.ConcurrentSize,
		recordDetectNumPerTask: cache.New(0, 0),
	}
	return DefaultCompress
}

func NewCompressData(taskId string, data *[]byte) CompressData {
	return CompressData{
		data:   data,
		taskID: taskId,
	}
}

func (cd *CompressDetect) Send(data CompressData) {
	cd.dataCh <- data
}

func (cd *CompressDetect) CompressAlgorithm() Algorithm {
	return cd.config.Algorithm
}

func (cd *CompressDetect) IsCompress(Key string) bool {
	_, b := cd.detectResult.Get(Key)
	return b
}

func (cd *CompressDetect) Decompress(reader io.Reader, algorithm Algorithm) (io.ReadCloser, error) {
	var compressAlgorithmInstance Compressor
	if algorithm != "" {
		synclock.Lock(string(algorithm), true)
		if compressAlgorithm, ok := cd.compressMap[algorithm]; !ok {
			compressAlgorithmInstance = NewCompress(algorithm)
			synclock.UnLock(string(algorithm), true)

			synclock.Lock(string(algorithm), false)
			cd.compressMap[algorithm] = compressAlgorithmInstance
			synclock.UnLock(string(algorithm), false)
		} else {
			compressAlgorithmInstance = compressAlgorithm
			synclock.Lock(string(algorithm), true)
		}
	} else {
		compressAlgorithmInstance = cd.compress
	}
	readCloser, err := compressAlgorithmInstance.Decompress(reader)
	if err != nil {
		return readCloser, errors.Wrap(err, fmt.Sprintf("Dynamic uncompress for %s algorithm error", cd.config.Algorithm))
	}
	return readCloser, nil
}

func (cd *CompressDetect) Compress(writer io.Writer) (io.WriteCloser, error) {
	writeCloser, err := cd.compress.Compress(writer)
	if err != nil {
		return writeCloser, errors.Wrap(err, fmt.Sprintf("Dynamic compress for %s algorithm error", cd.config.Algorithm))
	}
	return writeCloser, nil
}

func (cd *CompressDetect) Run() {
	lockerPool := synclock.NewLockerPool()
	for i := 0; i < cd.concurrentSize; i++ {
		go func() {
			confRatio := cd.config.Ratio
			for data := range cd.dataCh {
				ratio, err := cd.compress.CompressRatio(*data.data)
				if err != nil {
					continue
				}
				if _, b := cd.detectResult.Get(data.taskID); b {
					continue
				}
				var detectResult bool
				if ratio > confRatio {
					detectResult = true
				} else {
					detectResult = false
				}
				lockerPool.Lock(KeyJoin(data.taskID), true)
				if detectNum, b := cd.recordDetectNumPerTask.Get(KeyJoin(data.taskID)); !b {
					cd.recordDetectNumPerTask.SetDefault(KeyJoin(data.taskID), []bool{detectResult})
					lockerPool.UnLock(KeyJoin(data.taskID), true)
					continue
				} else {
					detectResults := detectNum.([]bool)
					if len(detectResults) == DefaultDetectPiecesNum {
						// check all compress ratio detect result is qualified
						var allQualifiedRatio bool = true
						for _, b := range detectResults {
							if !b {
								allQualifiedRatio = false
								break
							}
						}
						lockerPool.Lock(data.taskID, false)
						if allQualifiedRatio {
							cd.detectResult.SetDefault(data.taskID, true)
							logger.Infof("task %s compress detect is true", data.taskID)
						} else {
							cd.detectResult.SetDefault(data.taskID, false)
							logger.Infof("task %s compress detect is false", data.taskID)
						}
						lockerPool.UnLock(data.taskID, false)
					}
					detectResults = append(detectResults, detectResult)
					cd.recordDetectNumPerTask.SetDefault(KeyJoin(data.taskID), detectResults)
					lockerPool.UnLock(KeyJoin(data.taskID), true)
				}
			}
		}()
	}
	logger.Infof("Dynamic compress detect run success,compress ratio %d, compress algorithm %s",
		cd.config.Ratio, cd.config.Algorithm)
}

func KeyJoin(taskId string) string {
	return fmt.Sprintf("%s-%s", KeyPrefix, taskId)
}
