package domain

import "net"

type ResolvedTarget struct {
	RawURL string
	Scheme string
	Host   string
	IPs    []net.IP
}
