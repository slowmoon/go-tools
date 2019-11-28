package config

import (
	"github.com/slowmoon/go-tools/xtime"
)

type Config struct {
	// Report network e.g. unixgram, tcp, udp
	Network string `dsn:"network"`
	// For TCP and UDP networks, the addr has the form "host:port".
	// For Unix networks, the address must be a file system path.
	Addr string `dsn:"address"`
	// DEPRECATED
	Proto string `dsn:"network"`
	// DEPRECATED
	Chan int `dsn:"query.chan,"`
	// Report timeout
	Timeout xtime.Duration `dsn:"query.timeout,200ms"`
	// DisableSample
	DisableSample bool `dsn:"query.disable_sample"`
	// probabilitySampling
	Probability float32 `dsn:"-"`
	// ProtocolVersion
	ProtocolVersion int32 `dsn:"query.protocol_version,2"`
}
