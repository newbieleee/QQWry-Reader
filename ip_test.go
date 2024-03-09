package ip

import (
	`math/rand`
	`net`
	`sync`
	`testing`
)

var (
	ipFile        = "qqwry.dat"
	ipAddr        = "1.0.0.1"
	expectCountry = "美国"
	expectRegion  = "APNIC&CloudFlare公共DNS服务器"
)

func BenchmarkData_Query(b *testing.B) {
	q, _ := New(WithFileMode(ipFile))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ip := net.IPv4(
			byte(rand.Intn(254)+1),
			byte(rand.Intn(255)),
			byte(rand.Intn(255)),
			byte(rand.Intn(255)),
		)
		_, _ = q.Query(ip.String())
	}
	b.StopTimer()
	_ = q.Close()
}

func BenchmarkData_Query_Mem(b *testing.B) {
	q, _ := New(WithMemoryMode(ipFile))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ip := net.IPv4(
			byte(rand.Intn(254)+1),
			byte(rand.Intn(255)),
			byte(rand.Intn(255)),
			byte(rand.Intn(255)),
		)
		_, _ = q.Query(ip.String())
	}
	b.StopTimer()
	_ = q.Close()
}

func TestFileMode_Query(t *testing.T) {
	q, _ := New(WithFileMode(ipFile))
	r, _ := q.Query(ipAddr)
	if string(r.Country) != expectCountry {
		t.Error("unexpected country")
	}
	if string(r.Region) != expectRegion {
		t.Error("unexpected region")
	}
	_ = q.Close()
}

func TestFileMode_Query_Concurrency_Times10000(t *testing.T) {
	q, _ := New(WithFileMode(ipFile))
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
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
	_ = q.Close()
}

func TestMemMode_Query(t *testing.T) {
	q, _ := New(WithMemoryMode(ipFile))
	r, _ := q.Query(ipAddr)
	if string(r.Country) != expectCountry {
		t.Error("unexpected country")
	}
	if string(r.Region) != expectRegion {
		t.Error("unexpected region")
	}
	_ = q.Close()
}

func TestMemMode_Query_Concurrency_Times100000(t *testing.T) {
	q, _ := New(WithMemoryMode(ipFile))
	var wg sync.WaitGroup
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
	_ = q.Close()
}
