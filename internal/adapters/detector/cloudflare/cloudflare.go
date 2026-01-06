package cloudflare

import (
	"TechstackDetectorAPI/internal/core/catalog"
	"TechstackDetectorAPI/internal/core/domain"
	"context"
	"net"
	"strings"

	"github.com/go-resty/resty/v2"
)

type CloudFlare struct {
	ipv4Cidr []*net.IPNet
	ipv6Cidr []*net.IPNet
}

func NewCloudFlare() *CloudFlare {
	cf := &CloudFlare{}
	client := resty.New()

	var apiRes struct {
		Result struct {
			IPv4Cidrs []string `json:"ipv4_cidrs"`
			IPv6Cidrs []string `json:"ipv6_cidrs"`
		} `json:"result"`
		Success bool `json:"success"`
	}

	_, err := client.R().
		SetResult(&apiRes).
		Get("https://api.cloudflare.com/client/v4/ips")

	if err == nil && apiRes.Success {
		for _, cidr := range apiRes.Result.IPv4Cidrs {
			if _, ipNet, err := net.ParseCIDR(cidr); err == nil {
				cf.ipv4Cidr = append(cf.ipv4Cidr, ipNet)
			}
		}
		for _, cidr := range apiRes.Result.IPv6Cidrs {
			if _, ipNet, err := net.ParseCIDR(cidr); err == nil {
				cf.ipv6Cidr = append(cf.ipv6Cidr, ipNet)
			}
		}
	}

	return cf
}

func (d *CloudFlare) Detect(ctx *domain.FetchContext) ([]domain.Technology, error) {
	if ctx == nil {
		return nil, nil
	}

	// 1. NS Record (Passive DNS)
	if ctx.DNS != nil {
		if nsRecords, ok := ctx.DNS.Records["NS"]; ok {
			for _, record := range nsRecords {

				// a. Direct NS hostname match
				if strings.Contains(strings.ToLower(record.Value), "cloudflare.com") {
					return d.buildResult(), nil
				}

				// b. Resolve NS target via DNSFetcher (A / AAAA)
				if ctx.DNSFetcher != nil {
					if d.matchCloudflareByHost(ctx, record.Value) {
						return d.buildResult(), nil
					}
				}
			}
		}
	}

	// 2. SOA Record (Passive DNS)
	if ctx.DNS != nil {
		if soaRecords, ok := ctx.DNS.Records["SOA"]; ok {
			for _, record := range soaRecords {
				if strings.Contains(strings.ToLower(record.Value), "cloudflare.com") {
					return d.buildResult(), nil
				}
			}
		}
	}

	// 3. Fallback: HTTP Headers
	for _, res := range ctx.HTTP {
		if res == nil {
			continue
		}

		if res.Headers.Get("Cf-Ray") != "" ||
			strings.EqualFold(res.Headers.Get("Server"), "cloudflare") {
			return d.buildResult(), nil
		}
	}

	return nil, nil
}

func (d *CloudFlare) matchCloudflareByHost(
	ctx *domain.FetchContext,
	host string,
) bool {

	// Query A
	aRecords, err := ctx.DNSFetcher.Query(
		context.Background(),
		host,
		"A",
	)
	if err == nil {
		for _, r := range aRecords {
			if d.isCloudflareIP(r.Value) {
				return true
			}
		}
	}

	// Query AAAA
	aaaaRecords, err := ctx.DNSFetcher.Query(
		context.Background(),
		host,
		"AAAA",
	)
	if err == nil {
		for _, r := range aaaaRecords {
			if d.isCloudflareIP(r.Value) {
				return true
			}
		}
	}

	return false
}

func (d *CloudFlare) buildResult() []domain.Technology {
	return []domain.Technology{
		*catalog.CloudFlare(),
	}
}

func (d *CloudFlare) isCloudflareIP(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}

	ranges := d.ipv4Cidr
	if parsed.To4() == nil {
		ranges = d.ipv6Cidr
	}

	for _, ipNet := range ranges {
		if ipNet.Contains(parsed) {
			return true
		}
	}
	return false
}

func (d *CloudFlare) Name() string                              { return "cloudflare" }
func (d *CloudFlare) FetchPlan(target string) *domain.FetchPlan { return nil }
