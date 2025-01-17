// Copyright 2019 PingCAP, Inc.
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

package verification

import (
	"fmt"
	"hash/crc64"

	"github.com/tikv/client-go/v2/tikv"
	"github.com/twotigers93/tidb/br/pkg/lightning/common"
	"go.uber.org/zap/zapcore"
)

var ecmaTable = crc64.MakeTable(crc64.ECMA)

// KVChecksum is the checksum of a collection of key-value pairs.
type KVChecksum struct {
	base      uint64
	prefixLen int
	bytes     uint64
	kvs       uint64
	checksum  uint64
}

// NewKVChecksum creates a new KVChecksum with the given checksum.
func NewKVChecksum(checksum uint64) *KVChecksum {
	return &KVChecksum{
		checksum: checksum,
	}
}

// NewKVChecksumWithKeyspace creates a new KVChecksum with the given checksum and keyspace.
func NewKVChecksumWithKeyspace(k tikv.Codec) *KVChecksum {
	ks := k.GetKeyspace()
	return &KVChecksum{
		base:      crc64.Update(0, ecmaTable, ks),
		prefixLen: len(ks),
	}
}

// MakeKVChecksum creates a new KVChecksum with the given checksum.
func MakeKVChecksum(bytes uint64, kvs uint64, checksum uint64) KVChecksum {
	return KVChecksum{
		bytes:    bytes,
		kvs:      kvs,
		checksum: checksum,
	}
}

// UpdateOne updates the checksum with a single key-value pair.
func (c *KVChecksum) UpdateOne(kv common.KvPair) {
	sum := crc64.Update(c.base, ecmaTable, kv.Key)
	sum = crc64.Update(sum, ecmaTable, kv.Val)

	c.bytes += uint64(c.prefixLen + len(kv.Key) + len(kv.Val))
	c.kvs++
	c.checksum ^= sum
}

// Update updates the checksum with a batch of key-value pairs.
func (c *KVChecksum) Update(kvs []common.KvPair) {
	var (
		checksum uint64
		sum      uint64
		kvNum    int
		bytes    int
	)

	for _, pair := range kvs {
		sum = crc64.Update(c.base, ecmaTable, pair.Key)
		sum = crc64.Update(sum, ecmaTable, pair.Val)
		checksum ^= sum
		kvNum++
		bytes += c.prefixLen
		bytes += len(pair.Key) + len(pair.Val)
	}

	c.bytes += uint64(bytes)
	c.kvs += uint64(kvNum)
	c.checksum ^= checksum
}

// Add adds the checksum of another KVChecksum.
func (c *KVChecksum) Add(other *KVChecksum) {
	c.bytes += other.bytes
	c.kvs += other.kvs
	c.checksum ^= other.checksum
}

// Sum returns the checksum.
func (c *KVChecksum) Sum() uint64 {
	return c.checksum
}

// SumSize returns the total size of the key-value pairs.
func (c *KVChecksum) SumSize() uint64 {
	return c.bytes
}

// SumKVS returns the total number of key-value pairs.
func (c *KVChecksum) SumKVS() uint64 {
	return c.kvs
}

// MarshalLogObject implements the zapcore.ObjectMarshaler interface.
func (c *KVChecksum) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddUint64("cksum", c.checksum)
	encoder.AddUint64("size", c.bytes)
	encoder.AddUint64("kvs", c.kvs)
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (c *KVChecksum) MarshalJSON() ([]byte, error) {
	result := fmt.Sprintf(`{"checksum":%d,"size":%d,"kvs":%d}`, c.checksum, c.bytes, c.kvs)
	return []byte(result), nil
}
