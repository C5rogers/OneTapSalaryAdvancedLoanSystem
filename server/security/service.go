package security

import (
	"net"
	"time"
)

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// filter out old attempts
	history := rl.attempts[ip]
	valid := []time.Time{}
	for _, t := range history {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}

	// check limit
	if len(valid) >= rl.limit {
		rl.attempts[ip] = valid
		return false
	}

	valid = append(valid, now)
	rl.attempts[ip] = valid
	return true
}

// Helper to extract IP from request
func GetIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}
