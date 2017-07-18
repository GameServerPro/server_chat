package main

import (
	"flag"
	"os"
	"runtime/pprof"
	"syscall"

	"github.com/GameServerPro/glog"
	"github.com/GameServerPro/util/netutil"
	"github.com/GameServerPro/util/perf"

	"server_chat/comet/server"
)

var (
	DefaultServer *server.Server
	confFile      string
)

func init() {
	flag.StringVar(&confFile, "c", "./comet.conf", " set comet config file path")
}

func main() {
	flag.Parse()
	defer glog.Flush()
	err := server.InitConfig(confFile)
	if err != nil {
		panic(err)
	}

	DefaultServer = server.New(server.ServerOption{
		HandshakeTimeout: server.Conf.HandshakeTimeout,
		TcpKeepalive:     server.Conf.TcpKeepalive,
		TcpSenBuf:        server.Conf.TcpSndbuf,
		TcpRcvBuf:        server.Conf.TcpRcvbuf,
	})

	if server.Conf.Debug {
		glog.Infof("pprof listen addr:%v", server.Conf.PprofAddrs)
		perf.Init(server.Conf.PprofAddrs)
	}
	glog.Warningf("pprof listen addr:%v", server.Conf.PprofAddrs)
	glog.Errorf("pprof listen addr:%v", server.Conf.PprofAddrs)

	DefaultServer.Run()

	// wait for signal
	quitChan := make(chan struct{})
	netutil.ListenSignal(func(sig os.Signal) bool {
		glog.Info("recv signal:", sig)
		switch sig {
		case syscall.SIGINT, syscall.SIGQUIT:
			quitChan <- struct{}{}
			return true
		case syscall.SIGUSR2:
			// dump goroutine stack traces
			dumpOut, err := os.OpenFile("server_chat.dump", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.FileMode(0777))
			if err != nil {
				glog.Warning("dump error:", err)
			} else {
				pprof.Lookup("goroutine").WriteTo(dumpOut, 2)
				dumpOut.Close()
			}
		}
		return false
	}, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGUSR2)
	<-quitChan
}
