package server

import (
	"github.com/uber-go/zap"
	"time"
)

type ServerOption struct {
	HandshakeTimeout time.Duration
	TcpKeepalive     bool
	TcpRcvBuf        int
	TcpSenBuf        int
	Loger            zap.Logger
}

type Server struct {
	Buckets []*Bucket
	Options ServerOption
}

func New(opt ServerOption) *Server {
	s := new(Server)
	s.Options = opt
	return s
}

func (s *Server) Run() {
}
