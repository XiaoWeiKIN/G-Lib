//go:build arm64 || ppc64 || ppc64le

package atomicutil

// CacheLineSize is the size of a CPU cache line.
const CacheLineSize = 128
