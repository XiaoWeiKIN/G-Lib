# G-Lib

Go 基础工具库，移植自 [VictoriaMetrics](https://github.com/VictoriaMetrics/VictoriaMetrics) 的 `lib/` 系列工具包。

模块路径：`github.com/XiaoWeiKIN/G-Lib`，仅依赖标准库。

## 包结构

| 包 | 上游 | 说明 |
| --- | --- | --- |
| `bytesutil` | `lib/bytesutil` | 完整移植：零拷贝字符串/字节转换、`ByteBuffer` 与对象池、字符串驻留、带缓存的字符串匹配/转换器 |
| `slicesutil` | `lib/slicesutil` | 完整移植：泛型切片扩容工具与 `Buffer[T]`/`BufferPool[T]`（`bytesutil` 依赖） |
| `fasttime` | `lib/fasttime` | 完整移植：秒级缓存的当前时间戳（`bytesutil` 依赖） |
| `atomicutil` | `lib/atomicutil` | 部分移植：`Uint64`（防伪共享）、`Slice[T]`、`CacheLineSize`（`fasttime` 依赖） |
| `timeutil` | `lib/timeutil` | 部分移植：仅 `AddJitterToDuration`（`bytesutil` 依赖） |

## bytesutil API

```go
// 零拷贝转换（结果的有效期取决于原始数据是否可达且未被修改）
s := bytesutil.ToUnsafeString(b)
b := bytesutil.ToUnsafeBytes(s)

// 四种 resize 策略：是否拷贝原内容 × 是否按 2 的幂超额分配
b = bytesutil.ResizeWithCopyMayOverallocate(b, n)
b = bytesutil.ResizeWithCopyNoOverallocate(b, n)
b = bytesutil.ResizeNoCopyMayOverallocate(b, n)
b = bytesutil.ResizeNoCopyNoOverallocate(b, n)

// 字节缓冲与池
var bbp bytesutil.ByteBufferPool
bb := bbp.Get()
bb.MustWrite([]byte("foo"))       // 也实现了 io.Writer / io.WriterTo / io.ReaderFrom
r := bb.NewReader()               // 返回 bytesutil.ReadCloser
bbp.Put(bb)

// 字符串驻留（string interning），降低重复字符串的内存占用
s = bytesutil.InternString(s)
s = bytesutil.InternBytes(b)
s = bytesutil.Itoa(n)             // 重复调用同一个 n 不再分配内存

// 带缓存的匹配器 / 转换器，避免对相同输入重复执行昂贵函数
m := bytesutil.NewFastStringMatcher(func(s string) bool { return strings.HasPrefix(s, "foo") })
ok := m.Match("foobar")
tr := bytesutil.NewFastStringTransformer(strings.ToLower)
out := tr.Transform("FooBar")
```

### 命令行参数

`bytesutil` 沿用上游行为，在 `init` 阶段向标准库 `flag` 注册三个参数，用于控制驻留缓存（同时作用于 `FastStringMatcher`/`FastStringTransformer` 的缓存过期时间）：

- `-internStringMaxLen`（默认 `500`）：允许驻留的最大字符串长度
- `-internStringDisableCache`（默认 `false`）：禁用驻留缓存
- `-internStringCacheExpireDuration`（默认 `6m`）：缓存过期时间

不调用 `flag.Parse()` 时全部使用默认值。

## 移植改动

相对上游源码的改动仅限于让代码脱离 VictoriaMetrics 主仓库：

1. **导入路径**：`github.com/VictoriaMetrics/VictoriaMetrics/lib/*` → `github.com/XiaoWeiKIN/G-Lib/*`。
2. **去掉 `lib/filestream` 依赖**：`ByteBuffer.NewReader()` 原返回 `filestream.ReadCloser`，现返回本包定义的同签名接口 `bytesutil.ReadCloser`（`Path() string` / `Read([]byte) (int, error)` / `MustClose()`）。Go 接口是结构化匹配，返回值仍可直接赋给 `filestream.ReadCloser`。
3. **去掉 `golang.org/x/sys/cpu` 依赖**：`atomicutil.CacheLineSize` 改为按 `GOARCH` 的构建标签定义（与 `x/sys/cpu` 取值一致：arm/mips 系 32、arm64/ppc64 系 128、s390x 256、其余 64；上游 wasm 为 0，此处取 64 以避免填充计算除零）。
4. **去掉 `github.com/valyala/fastrand` 依赖**：`timeutil.AddJitterToDuration` 改用标准库 `math/rand/v2`。
5. `atomicutil` 与 `timeutil` 只移植了被依赖的部分，其余文件（`timeutil` 的 `duration.go`/`time.go`/`timezone.go` 等）未引入，因为它们会牵出 `lib/logger` 等更多依赖。

其余逻辑、注释、测试（含 `_timing_test.go` 基准测试）与上游保持一致。

## 验证

```
go vet ./... && go test ./...              # 全部通过
go test -race ./...                        # 全部通过
go build -tags synctest ./...              # fasttime 的 synctest 变体可编译
```

交叉编译验证：`linux/{amd64,386,arm,arm64,s390x,ppc64le}`、`darwin/amd64`、`windows/amd64`。

## 许可证

Apache License 2.0，与上游一致。详见 [LICENSE](LICENSE) 与 [NOTICE](NOTICE)。
