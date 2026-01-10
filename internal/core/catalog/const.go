package catalog

// this file defines constants for each adapter, as writing string keys directly to the registry would be expensive if changes were made
// use the constants provided here to avoid data mismatches between adapters.

//go:generate stringer -type=DetectorID -linecomment
type DetectorID int

const (
	GLOBAL      DetectorID = iota // GLOBAL
	Cloudflare                    // cloudflare
	PHP                           // php-programming-language
	ApacheHTTPD                   // apache-httpd
	LiteSpeed                     // litespeed
	Nginx                         // nginx
	WordPress                     // wordpress
)
