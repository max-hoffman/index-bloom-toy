package index_bloom_toy

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

// TODO benchmark lookup w and w/o bloom

// TODO test CPU and memory usage

func newTestTableIndexes(indexCnt, addrCnt int, typ indexType) *tableIndexes {
	indexes := make([]tableIndex, indexCnt)
	for i := range indexes {
		addrs := make([]prefix, addrCnt)
		for i := range addrs {
			addrs[i] = newAddr()
		}
		sort.Slice(addrs, func(i int, j int) bool {
			return bytes.Compare(addrs[i][:], addrs[j][:]) <= 0
		})

		var index tableIndex
		switch typ {
		case binaryIndexType:
			index = newBinaryTableIndex(addrs)
		case hashIndexType:
			index = newHashTableIndex(addrs)
		case bloomIndexType:
			index = newBloomTableIndex(addrs)
		}
		indexes[i] = index
	}
	return &tableIndexes{
		indexes: indexes,
	}
}

func newAddr() prefix {
	p := new(prefix)
	rand.Read(p[:])
	return *p
}

var result bool

func BenchmarkTableIndexSearch(b *testing.B) {
	var r bool
	indexTypes := []indexType{binaryIndexType, hashIndexType, bloomIndexType}
	numTableFiles := 400
	numAddrsPerIndex := 50_000

	for _, typ := range indexTypes {
		b.Run(fmt.Sprintf("%s search", typ), func(b *testing.B) {
			ti := newTestTableIndexes(numTableFiles, numAddrsPerIndex, typ)
			toLookup := ti.indexes[numTableFiles/2].prefixes()[numAddrsPerIndex/2]
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				// always record the result to prevent
				// the compiler eliminating the function call.
				r = ti.has(toLookup)
				if !r {
					panic("address not found")
				}
			}
			// always store the result to a package level variable
			// so the compiler cannot eliminate the Benchmark itself.
			result = r
		})
	}
}

var result2 *tableIndexes

func BenchmarkTableIndexSetup(b *testing.B) {
	var r *tableIndexes
	indexTypes := []indexType{binaryIndexType, hashIndexType, bloomIndexType}
	numTableFiles := 400
	numAddrsPerIndex := 50_000

	for _, typ := range indexTypes {
		b.Run(fmt.Sprintf("%s setup", typ), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				// always record the result to prevent
				// the compiler eliminating the function call.
				r = newTestTableIndexes(numTableFiles, numAddrsPerIndex, typ)
			}
			// always store the result to a package level variable
			// so the compiler cannot eliminate the Benchmark itself.
			result2 = r
		})
	}
}
