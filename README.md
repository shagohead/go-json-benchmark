Benchmark low-level Go JSON decoders
====================================

```
goos: darwin
goarch: arm64
pkg: github.com/shagohead/go-json-benchmark
cpu: Apple M4 Pro
BenchmarkJXDecoder
BenchmarkJXDecoder/concrete
BenchmarkJXDecoder/concrete-14            255651              4663 ns/op            3914 B/op        123 allocs/op
BenchmarkJXDecoder/generics
BenchmarkJXDecoder/generics-14            259596              4596 ns/op            3994 B/op        124 allocs/op
BenchmarkJSONv2Decoder
BenchmarkJSONv2Decoder-14                 170206              7017 ns/op            8173 B/op        155 allocs/op
```

Concrete & generics distincts only in usage of generics for collection types.
JSONv2 decoder simply uses concrete types like in concrete version of jx.
