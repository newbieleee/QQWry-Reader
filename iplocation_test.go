package ip

import (
    `testing`
)

func TestData_Query(t *testing.T) {
    q, err := New(`./qqwry.dat`)
    if err != nil {
        panic(err)
    }
    r, err := q.Query("1.0.0.1")
    if err != nil {
        panic(err)
    }
    if string(r.country) != "美国" {
        t.Error("unexpected country")
    }
    if string(r.region) != "APNIC&CloudFlare公共DNS服务器" {
        t.Error("unexpected region")
    }
}
