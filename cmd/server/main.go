package main

import (
	"context"
	"crypto/tls"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/number571/go-http3-proxy/utils"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"net"
	"net/http"
)

type Server struct {
	server    http3.Server
	transport *quic.Transport
	tlsConfig *tls.Config
}

func NewServer() (server *Server) {
	tlsConfig := utils.GenerateTLSConfig()

	packetConn, err := net.ListenPacket("udp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}

	transport := &quic.Transport{Conn: packetConn}

	server = &Server{
		server: http3.Server{
			TLSConfig: tlsConfig,
		},
		transport: transport,
		tlsConfig: tlsConfig,
	}

	return
}

func (server *Server) Start() (err error) {
	defer server.transport.Conn.Close()

	server.server.Handler = server.handler()

	ln, err := server.transport.ListenEarly(server.tlsConfig, nil)
	if err != nil {
		panic(err)
	}

	err = server.server.ServeListener(ln)
	if err != nil {
		err = gerror.Wrapf(err, "http3 server.ServeListener failed")
	}

	return
}

func (server *Server) client() (client *http.Client) {
	dial := func(ctx context.Context, addr string, tlsCfg *tls.Config, cfg *quic.Config) (quic.EarlyConnection, error) {
		remoteAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			return nil, err
		}

		return server.transport.DialEarly(ctx, remoteAddr, tlsCfg, cfg)
	}

	client = &http.Client{
		Transport: &http3.RoundTripper{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Dial: dial,
		},
	}

	return
}

func main() {
	server := NewServer()

	g.Log().Info(context.Background(), "Server is listening...")
	g.Log().Errorf(context.Background(), "%+v", server.Start())
}
