package fetcher

import (
	"TechstackDetectorAPI/internal/core/domain"
	"context"
	"strings"

	"github.com/miekg/dns"
)

type DNSFetcher struct {
	resolver string // dns server resolver, eg: 1.1.1.1:53 or 8.8.8.8:53
}

func (f *DNSFetcher) Fetch(ctx context.Context, host string) *domain.DNSResult {
	result := &domain.DNSResult{
		Records: make(map[string][]domain.DNSRecord),
	}

	recordTypes := map[string]uint16{
		"A":     dns.TypeA,
		"AAAA":  dns.TypeAAAA,
		"CNAME": dns.TypeCNAME,
		"MX":    dns.TypeMX,
		"TXT":   dns.TypeTXT,
		"NS":    dns.TypeNS,
		"SOA":   dns.TypeSOA,
		"SRV":   dns.TypeSRV,
	}

	for name, qtype := range recordTypes {
		records := f.query(ctx, host, qtype)
		if len(records) > 0 {
			result.Records[name] = records
		}
	}

	return result
}

func (f *DNSFetcher) query(
	ctx context.Context,
	host string,
	qtype uint16,
) []domain.DNSRecord {

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(host), qtype)

	c := new(dns.Client)
	resp, _, err := c.ExchangeContext(ctx, m, f.resolver)
	if err != nil || resp == nil {
		return nil
	}

	var records []domain.DNSRecord

	for _, ans := range resp.Answer {
		h := ans.Header()

		switch v := ans.(type) {

		case *dns.A:
			records = append(records, domain.DNSRecord{
				Name:  h.Name,
				Value: v.A.String(),
				TTL:   h.Ttl,
			})

		case *dns.AAAA:
			records = append(records, domain.DNSRecord{
				Name:  h.Name,
				Value: v.AAAA.String(),
				TTL:   h.Ttl,
			})

		case *dns.CNAME:
			records = append(records, domain.DNSRecord{
				Name:  h.Name,
				Value: v.Target,
				TTL:   h.Ttl,
			})

		case *dns.MX:
			records = append(records, domain.DNSRecord{
				Name:  h.Name,
				Value: v.Mx,
				TTL:   h.Ttl,
			})

		case *dns.NS:
			records = append(records, domain.DNSRecord{
				Name:  h.Name,
				Value: v.Ns,
				TTL:   h.Ttl,
			})

		case *dns.TXT:
			records = append(records, domain.DNSRecord{
				Name:  h.Name,
				Value: strings.Join(v.Txt, " "),
				TTL:   h.Ttl,
			})

		case *dns.SOA:
			records = append(records, domain.DNSRecord{
				Name:  h.Name,
				Value: v.Ns + " " + v.Mbox,
				TTL:   h.Ttl,
			})

		case *dns.SRV:
			records = append(records, domain.DNSRecord{
				Name:  h.Name,
				Value: v.Target,
				TTL:   h.Ttl,
			})
		}
	}

	return records
}

func NewDNSFetcher(resolver string) *DNSFetcher {
	return &DNSFetcher{
		resolver: resolver,
	}
}
