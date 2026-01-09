package security

import (
	"TechstackDetectorAPI/internal/core/domain" // Pastikan path ini sesuai
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTargetValidator_Validate(t *testing.T) {
	// Setup Blacklist
	rules := []string{
		"127.0.0.1",
		"10.0.0.0/8",
		"localhost",
		"malicious.com",
	}
	bl := NewBlacklist(rules)
	validator := NewTargetValidator(bl, nil) // DNSFetcher sementara nil jika pakai net default

	tests := []struct {
		name          string
		targetURL     string
		expectedError error
	}{
		{
			name:          "Valid Public URL",
			targetURL:     "https://google.com",
			expectedError: nil,
		},
		{
			name:          "Blocked Hostname",
			targetURL:     "http://malicious.com",
			expectedError: domain.ErrBlockedTarget,
		},
		{
			name:          "Blocked by IP (Localhost)",
			targetURL:     "http://127.0.0.1",
			expectedError: domain.ErrBlockedTarget,
		},
		{
			name:          "Blocked by CIDR (Private Network)",
			targetURL:     "http://10.1.1.5",
			expectedError: domain.ErrBlockedTarget,
		},
		{
			name:          "Invalid Scheme (FTP)",
			targetURL:     "ftp://google.com",
			expectedError: domain.ErrInvalidTarget,
		},
		{
			name:          "Malformed URL",
			targetURL:     "://invalid-url",
			expectedError: domain.ErrInvalidTarget,
		},
		{
			name:          "Empty Host",
			targetURL:     "http://",
			expectedError: domain.ErrInvalidTarget,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.Validate(context.Background(), tt.targetURL)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.targetURL, result.RawURL)
			}
		})
	}
}

func TestBlacklist_Detailed(t *testing.T) {
	rules := []string{
		"127.0.0.1",
		"192.168.1.0/24",
		"evil.com",
		"internal",
	}
	bl := NewBlacklist(rules)

	t.Run("IP and CIDR Rules", func(t *testing.T) {
		tests := []struct {
			ip       string
			expected bool
		}{
			{"127.0.0.1", true},
			{"127.0.0.2", false},
			{"192.168.1.50", true},
			{"192.168.2.1", false},
		}

		for _, tt := range tests {
			assert.Equal(t, tt.expected, bl.IsIPBlocked(net.ParseIP(tt.ip)), "IP: %s", tt.ip)
		}
	})

	t.Run("Hostname and Suffix Rules", func(t *testing.T) {
		tests := []struct {
			host     string
			expected bool
		}{
			{"evil.com", true},          // Exact match
			{"SUB.EVIL.COM", false},     // Case-insensitive & subdomain
			{"another-evil.com", false}, // Should not block (different domain)
			{"internal", true},          // Exact match .internal logic
			{"myinternal", false},       // Should not block partial word
			{"google.com", false},       // Clean domain
		}

		for _, tt := range tests {
			res := bl.IsHostBlocked(tt.host)
			assert.Equal(t, tt.expected, res, "Host: %s", tt.host)
		}
	})
}
