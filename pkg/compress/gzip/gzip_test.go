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

package gzip

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestUnCompression(t *testing.T) {
	// aaaaaabbbbbbccccccc
	fileName := "./compression"
	compressionFile(fileName, t)
	defer os.Remove(fileName)
	gzipFile, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	gzipCom := NewGzipCompress()
	gzipReader, err := gzipCom.UnCompression(gzipFile)
	if err != nil {
		t.Fatal(err)
	}
	defer gzipReader.Close()
	var bf bytes.Buffer
	_, err = io.Copy(&bf, gzipReader)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("read content %s", bf.String())
	var dst bytes.Buffer
	reader := bytes.NewReader(bf.Bytes())
	_, err = reader.Seek(5, io.SeekStart)
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.CopyN(&dst, reader, 3)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("read content %s", dst.String())
}

func TestCompression(t *testing.T) {
	fileName := "./compression"
	compressionFile(fileName, t)
	defer os.Remove(fileName)
}

func compressionFile(fileName string, t *testing.T) {
	f, err := os.Create(fileName)
	if err != nil {
		t.Fatal(err)
	}
	gzipCom := NewGzipCompress()
	compression, err := gzipCom.Compression(f)
	if err != nil {
		t.Fatal(err)
	}
	defer compression.Close()
	//reader := bytes.NewReader([]byte("aaaaaabbbbbbccccccc"))
	_, err = compression.Write([]byte("aaaaaabbbbbbccccccc"))
	//_, err = io.Copy(compression, reader)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompressRatio(t *testing.T) {
	gzipCompress := NewGzipCompress()
	file, err := os.Open("../testdata/issue6550.gz.base64")
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	ratio, err := gzipCompress.CompressRatio(data)
	if err != nil {
		t.Fatal()
	}
	t.Logf("ratio %+v", ratio)
	if ratio < 2.3 {
		t.Fatal("Compress Ratio Too small")
	}
}
