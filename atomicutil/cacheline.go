//go:build !arm && !mips && !mipsle && !mips64 && !mips64le && !arm64 && !ppc64 && !ppc64le && !s390x

package atomicutil

// CacheLineSize is the size of a CPU cache line.
const CacheLineSize = 64
