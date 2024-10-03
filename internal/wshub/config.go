package wshub

import "time"

// Config is the configuration for the websocket hub.
type Config struct {
	// Time allowed to write a message to the peer.
	WriteWaitSec int
	// Time allowed to read the next pong message from the peer.
	PongWaitSec int
	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriodSec int
	// App Domain (for Websocket origin check)
	AppDomain string
	// Websocket Subdomain (for Websocket origin check)
	WebsocketSubdomain string
}

// WriteWait returns the write wait duration.
func (c *Config) WriteWait() time.Duration {
	waitSec := c.WriteWaitSec
	if waitSec <= 0 {
		waitSec = 10 // prevents panic: non-positive interval for NewTicker
	}
	return time.Duration(waitSec) * time.Second
}

// PingPeriod returns the ping period duration.
func (c *Config) PingPeriod() time.Duration {
	periodSec := c.PingPeriodSec
	if periodSec <= 0 {
		periodSec = 54 // prevents panic: non-positive interval for NewTicker
	}
	return time.Duration(periodSec) * time.Second
}

// PongWait returns the pong wait duration.
func (c *Config) PongWait() time.Duration {
	waitSec := c.PongWaitSec
	if waitSec <= 0 {
		waitSec = 60 // prevents panic: non-positive interval for NewTicker
	}
	return time.Duration(waitSec) * time.Second
}
