package server

import (
	"time"

	"github.com/Terry-Mao/goconf"
)

var (
	Conf  *Config
	gconf *goconf.Config
)

type Config struct {
	// base
	PidFile    string   `goconf:"base:pidfile"`
	Log        string   `goconf:"base:log"`
	ServerId   int      `goconf:"base:serverid"`
	Debug      bool     `goconf:"base:debug"`
	PprofAddrs []string `goconf:"base:pprof.addrs"`

	// tcp
	HandshakeTimeout time.Duration `goconf:"tcp:handshaketimeout"`
	TcpBind          []string      `goconf:"tcp:bind"`
	TcpSndbuf        int           `goconf:"tcp:sndbufsize"`
	TcpRcvbuf        int           `goconf:"tcp:rcvbufsize"`
	TcpKeepalive     bool          `goconf:"tcp:keepalive"`

	// log
	LogLever int `goconf:"log:level"`
}

func NewConfig() *Config {
	Conf = &Config{
		PidFile:    "/tmp/server-chat.pid",
		Log:        "./",
		ServerId:   1,
		Debug:      true,
		PprofAddrs: []string{"0.0.0.0:20000"},

		HandshakeTimeout: 1,
		TcpBind:          []string{"0.0.0.0:60020"},
		TcpSndbuf:        128,
		TcpRcvbuf:        128,
		LogLever:         1,
	}
	return Conf
}

// init the global config.
func InitConfig(cfgpath string) (err error) {
	Conf = NewConfig()
	gconf = goconf.New()
	if err = gconf.Parse(cfgpath); err != nil {
		return err
	}
	if err := gconf.Unmarshal(Conf); err != nil {
		return err
	}
	return nil
}

func ReloadConfig() (*Config, error) {
	conf := NewConfig()
	ngconf, err := gconf.Reload()
	if err != nil {
		return nil, err
	}
	if err := ngconf.Unmarshal(conf); err != nil {
		return nil, err
	}
	gconf = ngconf
	return conf, nil
}
