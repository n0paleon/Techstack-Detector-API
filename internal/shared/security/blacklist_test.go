package security

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlacklist(t *testing.T) {
	rules := []string{
		"127.0.0.1",
		"192.168.1.0/24",
		"evil.com",
		"internal", // Suffix check
	}
	bl := NewBlacklist(rules)

	t.Run("IP Blocking", func(t *testing.T) {
		assert.True(t, bl.IsIPBlocked(net.ParseIP("127.0.0.1")))
		assert.False(t, bl.IsIPBlocked(net.ParseIP("8.8.8.8")))
	})

	t.Run("CIDR Range Blocking", func(t *testing.T) {
		assert.True(t, bl.IsIPBlocked(net.ParseIP("192.168.1.50")))
		assert.True(t, bl.IsIPBlocked(net.ParseIP("192.168.1.255")))
		assert.False(t, bl.IsIPBlocked(net.ParseIP("192.168.2.1")))
	})

	t.Run("Hostname and Suffix Blocking", func(t *testing.T) {
		assert.True(t, bl.IsHostBlocked("evil.com"))
		assert.False(t, bl.IsHostBlocked("service.internal"))
		assert.True(t, bl.IsHostBlocked("INTERNAL")) // Case-insensitive
		assert.False(t, bl.IsHostBlocked("google.com"))
	})
}
