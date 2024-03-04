# 纯真IP数据库QQWry.dat读取工具

## 使用
`go get github.com/ryanexo/QQWry-Reader`

## 示例
```Go
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

***
## 并发安全
测试设备i5-12490F，内存3600Mhz 16G*2，已进行并发测试未发现问题

## 查询效率（基于上述设备）

文件模式查询10万次耗时: **7.6896057s**

文件模式并发查询10万次耗时: **13.1974191s**

内存模式查询1百万次耗时: **3.5238938s**

内存模式并发查询1百万次耗时: **2.2537558s**
