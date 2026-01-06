package domain

import (
	"net/http"
)

type FetchContext struct {
	Target string

	DNS  *DNSResult
	TLS  *TLSResult
	HTTP map[string]*HTTPResult // key = request ID

	DNSFetcher DNSFetcher // custom DNS Fetcher interface
}

type FetchPlan struct {
	BaseURL string

	Requests []FetchRequest
	TLS      bool
	JS       bool
}

type FetchRequest struct {
	ID          string // "root", "wp-json", "admin"
	Path        string
	Method      string
	Description string
}

func (r FetchRequest) HTTPMethod() string {
	if r.Method == "" {
		return http.MethodGet
	}
	return r.Method
}

type HTTPResult struct {
	RequestID string
	URL       string
	FinalURL  string

	StatusCode int
	Headers    http.Header
	Body       []byte

	Error error // per-request error
}

type DNSResult struct {
	Records map[string][]DNSRecord
}

type DNSRecord struct {
	Name  string
	Value string
	TTL   uint32
}

type TLSResult struct {
	Version     string
	CipherSuite string
	Issuer      string
	Subject     string
}
