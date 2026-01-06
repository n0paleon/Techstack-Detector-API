package domain

import "context"

type DNSFetcher interface {
	Fetch(ctx context.Context, host string) *DNSResult

	// Query custom dns query
	Query(
		ctx context.Context,
		host string,
		recordType string, // "A", "TXT", "CNAME", etc...
	) ([]DNSRecord, error)
}
