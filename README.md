# 纯真IP数据库QQWry.dat读取工具

## Install
`go get github.com/newbieleee/QQWry-Reader`

## Example
```Go
func main() {
    data, err := ip.New(`./qqwry.dat`)
    if err != nil {
        panic(err)
    }
    r, err := data.Query("1.0.0.1")
    if err != nil {
        panic(err)
    }
    // 按字段输出
    fmt.Printf("国家: %s\n地区: %s", r.Country, r.Region)
    
    // 合并输出
    location := make([]byte, len(r.Country)+len(r.Region))
    n := copy(location, r.Country)
    copy(location[n:], r.Region)
    fmt.Println(string(location))
}
```