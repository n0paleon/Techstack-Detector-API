package fetcher

import (
	"context"
	"testing"
)

func BenchmarkDNSFetcher_Remote(b *testing.B) {
	resolvers := []struct {
		name string
		addr string
	}{
		{"Cloudflare", "1.1.1.1:53"},
		{"Google", "8.8.8.8:53"},
	}

	for _, res := range resolvers {
		b.Run(res.name, func(b *testing.B) {
			fetcher := NewDNSFetcher(res.addr, 10)
			ctx := context.Background()
			host := "google.com"

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = fetcher.Fetch(ctx, host)
			}
		})
	}
}
