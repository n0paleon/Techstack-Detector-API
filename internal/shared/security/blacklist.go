package security

import (
	"net"
	"strings"
)

type Blacklist struct {
	hostRules []string
	ipRules   map[string]struct{}
	cidrRules []*net.IPNet
}

func NewBlacklist(rules []string) *Blacklist {
	bl := &Blacklist{
		ipRules: make(map[string]struct{}),
	}

	for _, r := range rules {

		// CIDR
		if _, cidr, err := net.ParseCIDR(r); err == nil {
			bl.cidrRules = append(bl.cidrRules, cidr)
			continue
		}

		// IP
		if ip := net.ParseIP(r); ip != nil {
			bl.ipRules[ip.String()] = struct{}{}
			continue
		}

		// hostname / suffix
		bl.hostRules = append(bl.hostRules, strings.ToLower(r))
	}

	return bl
}

// IsHostBlocked is case-insensitive (google.com != www.google.com)
func (b *Blacklist) IsHostBlocked(host string) bool {
	h := strings.ToLower(strings.TrimSuffix(host, "."))

	for _, rule := range b.hostRules {
		if h == rule {
			return true
		}
	}

	return false
}

func (b *Blacklist) IsIPBlocked(ip net.IP) bool {
	if ip == nil {
		return false
	}

	if _, ok := b.ipRules[ip.String()]; ok {
		return true
	}

	for _, cidr := range b.cidrRules {
		if cidr.Contains(ip) {
			return true
		}
	}

	return false
}
