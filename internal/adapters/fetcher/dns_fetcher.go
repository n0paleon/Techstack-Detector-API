package fetcher

import (
	"TechstackDetectorAPI/internal/core/domain"
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/miekg/dns"
	"golang.org/x/sync/errgroup"
)

type DNSFetcher struct {
	resolver    string
	client      *dns.Client
	maxParallel int
}

func (f *DNSFetcher) Fetch(ctx context.Context, host string) *domain.DNSResult {
	result := &domain.DNSResult{
		Records: make(map[string][]domain.DNSRecord),
	}

	var mu sync.Mutex
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(f.maxParallel)

	for recordType := range recordTypeMap {
		recordType := recordType // capture

		g.Go(func() error {
			records, err := f.Query(ctx, host, recordType)
			if err != nil || len(records) == 0 {
				return nil
			}

			mu.Lock()
			result.Records[recordType] = records
			mu.Unlock()

			return nil
		})
	}

	_ = g.Wait()
	return result
}

func (f *DNSFetcher) Query(
	ctx context.Context,
	host string,
	recordType string,
) ([]domain.DNSRecord, error) {

	qtype, ok := recordTypeMap[strings.ToUpper(recordType)]
	if !ok {
		return nil, errors.New("unsupported DNS record type")
	}

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(host), qtype)

	resp, _, err := f.client.ExchangeContext(ctx, m, f.resolver)
	if err != nil || resp == nil {
		return nil, err
	}

	return parseAnswers(resp.Answer), nil
}

func parseAnswers(answers []dns.RR) []domain.DNSRecord {
	var records []domain.DNSRecord

	for _, ans := range answers {
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

func NewDNSFetcher(resolver string, maxParallel int) *DNSFetcher {
	return &DNSFetcher{
		resolver:    resolver,
		client:      &dns.Client{},
		maxParallel: maxParallel,
	}
}
