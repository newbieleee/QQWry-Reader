### 简介
高性能的离线IP归属地查询库，数据源来自[纯真IP数据库](https://www.cz88.net/)。

### 使用
```
go get github.com/ryanexo/QQWry-Reader
```

### 示例
```go
func main() {
    q, err := ip.New(ip.WithMemoryMode("./qqwry.dat"))
    if err != nil {
        panic(err)
    }
    r, err := q.Query("1.0.0.1")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("国家: %s\n地区: %s\n合并: %s\n", r.Country, r.Region, r.Location)
    fmt.Printf("IP库版本号: %s", q.Version())
}
```

### 并发安全
使用goroutine模拟并发未发生竟态

### 测试

内存模式下单次查询为10微秒级别，文件模式下单次查询为100微秒级别

```text
=== RUN   TestFileMode_Query
--- PASS: TestFileMode_Query (0.00s)
=== RUN   TestFileMode_Query_Concurrency_Times10000
--- PASS: TestFileMode_Query_Concurrency_Times10000 (1.46s)
=== RUN   TestMemMode_Query
--- PASS: TestMemMode_Query (0.00s)
=== RUN   TestMemMode_Query_Concurrency_Times10000
--- PASS: TestMemMode_Query_Concurrency_Times10000 (0.06s)
```

```text
goos: windows
goarch: amd64
pkg: github.com/ryanexo/QQWry-Reader
cpu: 12th Gen Intel(R) Core(TM) i5-12400F
BenchmarkData_Query
BenchmarkData_Query-12             16491             72876 ns/op            8675 B/op         44 allocs/op
BenchmarkData_Query_Mem
BenchmarkData_Query_Mem-12        527658              2158 ns/op            8675 B/op         44 allocs/op
```