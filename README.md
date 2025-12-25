Benchmark low-level Go JSON decoders
====================================

```
goos: darwin
goarch: arm64
pkg: github.com/shagohead/go-json-benchmark
cpu: Apple M4 Pro
BenchmarkJXDecoder
BenchmarkJXDecoder-14             243699              4860 ns/op            3914 B/op        123 allocs/op
BenchmarkJSONV2Decoder
BenchmarkJSONV2Decoder-14         163411              7218 ns/op            8173 B/op        155 allocs/op
```
