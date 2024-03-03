# 纯真IP数据库QQWry.dat读取工具

## Install
`go get github.com/newbieleee/QQWry-Reader`

## Example
```Go
func main() {
    q, err := New(`./qqwry.dat`)
    if err != nil {
        panic(err)
    }
    r, err := q.Query("1.0.0.1")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("国家: %s\n地区: %s\n合并: %s", r.Country, r.Region, r.Location)
}
```
