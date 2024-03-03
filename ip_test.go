package ip

import (
    `fmt`
    `math/rand`
    `net`
    `sync`
    `testing`
    `time`
)

func TestFileMode_Query(t *testing.T) {
    q, err := New(WithFileMode("./qqwry.dat"))
    if err != nil {
        panic(err)
    }
    r, err := q.Query("1.0.0.1")
    if err != nil {
        panic(err)
    }
    if string(r.Country) != "美国" {
        t.Error("unexpected country")
    }
    if string(r.Region) != "APNIC&CloudFlare公共DNS服务器" {
        t.Error("unexpected region")
    }

    // 10万次查询
    s := time.Now()
    for i := 0; i < 100000; i++ {
        ip := net.IPv4(
            byte(rand.Intn(254)+1),
            byte(rand.Intn(255)),
            byte(rand.Intn(255)),
            byte(rand.Intn(255)),
        )
        _, err := q.Query(ip.String())
        if err != nil {
            t.Error(err)
        }
    }
    e := time.Since(s)

    fmt.Println("文件模式查询10万次耗时:", e)
    _ = q.Close()
}

func TestFileMode_Query_Concurrency(t *testing.T) {
    q, err := New(WithFileMode("./qqwry.dat"))
    if err != nil {
        panic(err)
    }

    var wg sync.WaitGroup

    // 10万次查询
    s := time.Now()
    for i := 0; i < 100000; i++ {
        wg.Add(1)
        go func() {
            ip := net.IPv4(
                byte(rand.Intn(254)+1),
                byte(rand.Intn(255)),
                byte(rand.Intn(255)),
                byte(rand.Intn(255)),
            )
            _, err := q.Query(ip.String())
            if err != nil {
                t.Error(err)
            }
            wg.Done()
        }()
    }
    wg.Wait()
    e := time.Since(s)

    fmt.Println("文件模式并发查询10万次耗时:", e)
    _ = q.Close()
}

func TestMemMode_Query(t *testing.T) {
    q, err := New(WithMemoryMode(`./qqwry.dat`))
    if err != nil {
        panic(err)
    }
    r, err := q.Query("1.0.0.1")
    if err != nil {
        panic(err)
    }
    if string(r.Country) != "美国" {
        t.Error("unexpected country")
    }
    if string(r.Region) != "APNIC&CloudFlare公共DNS服务器" {
        t.Error("unexpected region")
    }

    // 1百万次查询
    s := time.Now()
    for i := 0; i < 1000000; i++ {
        ip := net.IPv4(
            byte(rand.Intn(254)+1),
            byte(rand.Intn(255)),
            byte(rand.Intn(255)),
            byte(rand.Intn(255)),
        )
        _, err := q.Query(ip.String())
        if err != nil {
            t.Error(err)
        }
    }
    e := time.Since(s)

    fmt.Println("内存模式查询1百万次耗时:", e)
    _ = q.Close()
}

func TestMemMode_Query_Concurrency(t *testing.T) {
    q, err := New(WithMemoryMode(`./qqwry.dat`))
    if err != nil {
        panic(err)
    }

    var wg sync.WaitGroup

    // 1百万次查询
    s := time.Now()
    for i := 0; i < 1000000; i++ {
        wg.Add(1)
        go func() {
            ip := net.IPv4(
                byte(rand.Intn(254)+1),
                byte(rand.Intn(255)),
                byte(rand.Intn(255)),
                byte(rand.Intn(255)),
            )
            _, err := q.Query(ip.String())
            if err != nil {
                t.Error(err)
            }
            wg.Done()
        }()
    }
    wg.Wait()
    e := time.Since(s)

    fmt.Println("内存模式并发查询1百万次耗时:", e)
    _ = q.Close()
}
