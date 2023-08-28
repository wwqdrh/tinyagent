package swarm

import (
	"context"
	"net"
	"time"
)

func GetServiceIP(domain string) string {
	// Create a custom resolver with a dialer that connects to 127.0.0.11:53
	resolver := &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.DialTimeout("udp", "127.0.0.11:53", 1*time.Second)
		},
	}

	// Use the resolver to lookup the host
	ips, err := resolver.LookupHost(context.TODO(), domain)
	if err != nil || len(ips) == 0 {
		return "127.0.0.1"
	}
	return ips[0]
}
