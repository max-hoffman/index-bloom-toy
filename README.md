# index-bloom-toy

sandboxing different data structures for chunk lookup

Search benchmarks:
```
goos: darwin
goarch: amd64
pkg: github.com/max-hoffman/index-bloom-toy
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkTableIndexSearch
BenchmarkTableIndexSearch/binaryIndex_search
BenchmarkTableIndexSearch/binaryIndex_search-12         	   58119	     20228 ns/op
BenchmarkTableIndexSearch/hashIndex_search
BenchmarkTableIndexSearch/hashIndex_search-12           	  367981	      3431 ns/op
BenchmarkTableIndexSearch/bloomIndex_search
BenchmarkTableIndexSearch/bloomIndex_search-12          	  100902	     12379 ns/op
```

Setup benchmarks:
```
goos: darwin
goarch: amd64
pkg: github.com/max-hoffman/index-bloom-toy
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkTableIndexSetup
BenchmarkTableIndexSetup/binaryIndex_setup
BenchmarkTableIndexSetup/binaryIndex_setup-12         	       1	3625448809 ns/op	160606008 B/op	    1611 allocs/op
BenchmarkTableIndexSetup/hashIndex_setup
BenchmarkTableIndexSetup/hashIndex_setup-12           	       1	4275124567 ns/op	481545736 B/op	  340631 allocs/op
BenchmarkTableIndexSetup/bloomIndex_setup
BenchmarkTableIndexSetup/bloomIndex_setup-12          	       1	5721139682 ns/op	209779368 B/op	    2804 allocs/op
```
