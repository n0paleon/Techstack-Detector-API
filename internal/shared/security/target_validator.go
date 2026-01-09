package security

import (
	"TechstackDetectorAPI/internal/core/domain"
	"context"
	"net"
	"net/url"
)

type TargetValidator struct {
	blacklist  domain.Blacklist
	dnsFetcher domain.DNSFetcher
}

func (v *TargetValidator) Validate(ctx context.Context, target string) (*domain.ResolvedTarget, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, domain.ErrInvalidTarget
	}

	// scheme check
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, domain.ErrInvalidTarget
	}

	// host check
	host := u.Hostname()
	if host == "" {
		return nil, domain.ErrInvalidTarget
	}
	if v.blacklist.IsHostBlocked(host) {
		return nil, domain.ErrBlockedTarget
	}

	// DNS resolve
	ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
	if err != nil || len(ips) == 0 {
		return nil, domain.ErrInvalidTarget
	}

	// ip check
	for _, ip := range ips {
		if v.blacklist.IsIPBlocked(ip) {
			return nil, domain.ErrBlockedTarget
		}
	}

	return &domain.ResolvedTarget{
		RawURL: target,
		Scheme: u.Scheme,
		Host:   host,
		IPs:    ips,
	}, nil
}

func NewTargetValidator(blacklist domain.Blacklist, dnsFetcher domain.DNSFetcher) *TargetValidator {
	return &TargetValidator{
		blacklist:  blacklist,
		dnsFetcher: dnsFetcher,
	}
}
