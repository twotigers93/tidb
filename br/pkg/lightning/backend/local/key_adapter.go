// Copyright 2021 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package local

import (
	"math"

	"github.com/pingcap/errors"
	"github.com/twotigers93/tidb/br/pkg/lightning/common"
	"github.com/twotigers93/tidb/util/codec"
)

// KeyAdapter is used to encode and decode keys.
type KeyAdapter interface {
	// Encode encodes the key with its corresponding rowID. It appends the encoded key to dst and returns the
	// resulting slice. The encoded key is guaranteed to be in ascending order for comparison.
	Encode(dst []byte, key []byte, rowID []byte) []byte

	// Decode decodes the original key to dst. It appends the encoded key to dst and returns the resulting slice.
	Decode(dst []byte, data []byte) ([]byte, error)

	// EncodedLen returns the encoded key length.
	EncodedLen(key []byte, rowID []byte) int
}

func reallocBytes(b []byte, n int) []byte {
	newSize := len(b) + n
	if cap(b) < newSize {
		bs := make([]byte, len(b), newSize)
		copy(bs, b)
		return bs
	}
	return b
}

type noopKeyAdapter struct{}

func (noopKeyAdapter) Encode(dst []byte, key []byte, _ []byte) []byte {
	return append(dst, key...)
}

func (noopKeyAdapter) Decode(dst []byte, data []byte) ([]byte, error) {
	return append(dst, data...), nil
}

func (noopKeyAdapter) EncodedLen(key []byte, _ []byte) int {
	return len(key)
}

var _ KeyAdapter = noopKeyAdapter{}

type dupDetectKeyAdapter struct{}

func (dupDetectKeyAdapter) Encode(dst []byte, key []byte, rowID []byte) []byte {
	dst = codec.EncodeBytes(dst, key)
	dst = reallocBytes(dst, len(rowID)+2)
	dst = append(dst, rowID...)
	rowIDLen := uint16(len(rowID))
	dst = append(dst, byte(rowIDLen>>8), byte(rowIDLen))
	return dst
}

func (dupDetectKeyAdapter) Decode(dst []byte, data []byte) ([]byte, error) {
	if len(data) < 2 {
		return nil, errors.New("insufficient bytes to decode value")
	}
	rowIDLen := uint16(data[len(data)-2])<<8 | uint16(data[len(data)-1])
	tailLen := int(rowIDLen + 2)
	if len(data) < tailLen {
		return nil, errors.New("insufficient bytes to decode value")
	}
	_, key, err := codec.DecodeBytes(data[:len(data)-tailLen], dst[len(dst):cap(dst)])
	if err != nil {
		return nil, err
	}
	if len(dst) == 0 {
		return key, nil
	}
	if len(dst)+len(key) <= cap(dst) {
		dst = dst[:len(dst)+len(key)]
		return dst, nil
	}
	// New slice is allocated, append key to dst manually.
	return append(dst, key...), nil
}

func (dupDetectKeyAdapter) EncodedLen(key []byte, rowID []byte) int {
	return codec.EncodedBytesLength(len(key)) + len(rowID) + 2
}

var _ KeyAdapter = dupDetectKeyAdapter{}

// static vars for rowID
var (
	MinRowID  = common.EncodeIntRowID(math.MinInt64)
	ZeroRowID = common.EncodeIntRowID(0)
)
