package fetcher

import (
	"context"
	"net"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DNSFetcherTestSuite struct {
	suite.Suite
	server  *dns.Server
	addr    string
	fetcher *DNSFetcher
}

// SetupSuite menjalankan DNS server lokal untuk testing
func (s *DNSFetcherTestSuite) SetupSuite() {
	// Pilih port acak yang tersedia
	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	s.NoError(err)
	s.addr = pc.LocalAddr().String()
	_ = pc.Close()

	// Handler untuk menjawab query DNS dummy
	dns.HandleFunc("example.com.", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)

		qtype := r.Question[0].Qtype
		switch qtype {
		case dns.TypeA:
			rr, _ := dns.NewRR("example.com. 3600 IN A 1.2.3.4")
			m.Answer = append(m.Answer, rr)
		case dns.TypeTXT:
			rr, _ := dns.NewRR(`example.com. 3600 IN TXT "v=spf1 include:_spf.google.com ~all"`)
			m.Answer = append(m.Answer, rr)
		case dns.TypeMX:
			rr, _ := dns.NewRR("example.com. 3600 IN MX 10 mail.example.com.")
			m.Answer = append(m.Answer, rr)
		}
		_ = w.WriteMsg(m)
	})

	s.server = &dns.Server{Addr: s.addr, Net: "udp"}
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	s.fetcher = NewDNSFetcher(s.addr, 10)
}

func (s *DNSFetcherTestSuite) TearDownSuite() {
	_ = s.server.Shutdown()
}

func (s *DNSFetcherTestSuite) TestFetch_Success() {
	ctx := context.Background()
	host := "example.com"

	result := s.fetcher.Fetch(ctx, host)

	assert.NotNil(s.T(), result)

	// Cek Record A
	assert.Contains(s.T(), result.Records, "A")
	assert.Equal(s.T(), "1.2.3.4", result.Records["A"][0].Value)

	// Cek Record TXT
	assert.Contains(s.T(), result.Records, "TXT")
	assert.Contains(s.T(), result.Records["TXT"][0].Value, "v=spf1")

	// Cek Record MX
	assert.Contains(s.T(), result.Records, "MX")
	assert.Equal(s.T(), "mail.example.com.", result.Records["MX"][0].Value)
}

func (s *DNSFetcherTestSuite) TestFetch_NoRecords() {
	ctx := context.Background()
	// Domain yang tidak didaftarkan di handler dns.HandleFunc
	host := "nonexistent.com"

	result := s.fetcher.Fetch(ctx, host)

	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.Records)
}

func TestDNSFetcherSuite(t *testing.T) {
	suite.Run(t, new(DNSFetcherTestSuite))
}
