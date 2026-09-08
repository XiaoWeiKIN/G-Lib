//go:build arm || mips || mipsle || mips64 || mips64le

package atomicutil

// CacheLineSize is the size of a CPU cache line.
const CacheLineSize = 32
