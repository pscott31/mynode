package config

const (
	MAGIC_MAIN           uint32 = 0xD9B4BEF9
	DEFAULT_MAGIC               = MAGIC_MAIN
	DEFAULT_REMOTE_ADDR         = "127.0.0.1:8333"
	DEFAULT_VERSION      int32  = 31900 // TODO
	DEFAULT_SERVICES     uint64 = 1     // NODE_NETWORK
	DEFAULT_START_HEIGHT int32  = 0
)

type Config struct {
	RemoteAddr  string
	Magic       uint32
	Version     int32
	Services    uint64
	StartHeight int32
}

func Default() *Config {
	return &Config{
		RemoteAddr:  DEFAULT_REMOTE_ADDR,
		Magic:       DEFAULT_MAGIC,
		Version:     DEFAULT_VERSION,
		Services:    DEFAULT_SERVICES,
		StartHeight: DEFAULT_START_HEIGHT,
	}
}
