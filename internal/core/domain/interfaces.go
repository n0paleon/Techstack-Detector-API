package domain

import (
	"context"
	"net"
)

type DNSFetcher interface {
	Fetch(ctx context.Context, host string) *DNSResult

	// Query custom dns query
	Query(
		ctx context.Context,
		host string,
		recordType string, // "A", "TXT", "CNAME", etc...
	) ([]DNSRecord, error)
}

type TargetValidator interface {
	Validate(ctx context.Context, target string) (*ResolvedTarget, error)
}

type Blacklist interface {
	IsIPBlocked(ip net.IP) bool
	IsHostBlocked(host string) bool
}
