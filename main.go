package index_bloom_toy

import (
	"encoding/binary"
	"github.com/bits-and-blooms/bloom/v3"
)

type prefix [8]byte

type indexType uint8

const (
	unknownIndexType indexType = iota
	hashIndexType
	binaryIndexType
	bloomIndexType
)

func (t indexType) String() string {
	switch t {
	case hashIndexType:
		return "hashIndex"
	case binaryIndexType:
		return "binaryIndex"
	case bloomIndexType:
		return "bloomIndex"
	default:
		return "unknown"
	}
}

type tableIndex interface {
	has(prefix) bool
	prefixes() []prefix
}

type binaryTableIndex struct {
	p []prefix
}

var _ tableIndex = (*binaryTableIndex)(nil)

func newBinaryTableIndex(p []prefix) *binaryTableIndex {
	return &binaryTableIndex{
		p: p,
	}
}

func (b *binaryTableIndex) has(p prefix) bool {
	idx := binSearch(binary.BigEndian.Uint64(p[:]), len(b.p), b.p)
	return idx < len(b.p) && b.p[idx] == p
}

func (b *binaryTableIndex) prefixes() []prefix {
	return b.p
}

type hashTableIndex struct {
	p []prefix
	m map[prefix]struct{}
}

var _ tableIndex = (*hashTableIndex)(nil)

func newHashTableIndex(p []prefix) *hashTableIndex {
	m := make(map[prefix]struct{}, len(p))
	for _, p := range p {
		m[p] = struct{}{}
	}
	return &hashTableIndex{
		p: p,
		m: m,
	}
}

func (b *hashTableIndex) has(p prefix) bool {
	_, ok := b.m[p]
	return ok
}

func (b *hashTableIndex) prefixes() []prefix {
	return b.p
}

type bloomTableIndex struct {
	p     []prefix
	bloom *bloom.BloomFilter
}

var _ tableIndex = (*bloomTableIndex)(nil)

const bloomCapacity = 100000
const bloomFalsePositiveRate = 0.01

func newBloomTableIndex(p []prefix) *bloomTableIndex {
	filter := bloom.NewWithEstimates(bloomCapacity, bloomFalsePositiveRate)
	for _, p := range p {
		filter.Add(p[:])
	}
	return &bloomTableIndex{
		p:     p,
		bloom: filter,
	}
}

func (b *bloomTableIndex) has(p prefix) bool {
	if b.bloom.Test(p[:]) {
		// binary search
		idx := binSearch(binary.BigEndian.Uint64(p[:]), len(b.p), b.p)
		return idx < len(b.p) && b.p[idx] == p
	}
	return false
}

func (b *bloomTableIndex) prefixes() []prefix {
	return b.p
}

func binSearch(query uint64, count int, prefixes []prefix) int {
	idx, j := 0, count
	for idx < j {
		h := idx + (j-idx)/2 // avoid overflow when computing h
		// i ≤ h < j
		tmp := binary.BigEndian.Uint64(prefixes[h][:])
		if tmp < query {
			idx = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	return idx
}

type tableIndexes struct {
	indexes []tableIndex
}

func (ti *tableIndexes) has(p prefix) bool {
	for _, i := range ti.indexes {
		if i.has(p) {
			return true
		}
	}
	return false
}
