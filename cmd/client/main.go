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
	"time"
)

var (
	// can be overwritten if used docker-mode (init.go)
	proxyHost  = "127.0.0.1:1080"
	remoteHost = "127.0.0.1:8080"
)

type Client struct {
	client    http.Client
	server    http3.Server
	transport *quic.Transport
	tlsConfig *tls.Config
}

func NewClient() (client *Client) {
	tlsConfig := utils.GenerateTLSConfig()

	packetConn, err := net.ListenPacket("udp", "0.0.0.0:8081")
	if err != nil {
		panic(err)
	}

	transport := &quic.Transport{Conn: packetConn}

	client = &Client{
		client: http.Client{
			Transport: &http3.RoundTripper{
				//Dial: proxyDialer("socks5://" + proxyHost),
				Dial: func(ctx context.Context, addr string, tlsCfg *tls.Config, cfg *quic.Config) (quic.EarlyConnection, error) {
					remoteAddr, err := net.ResolveUDPAddr("udp", addr)
					if err != nil {
						return nil, err
					}

					return transport.DialEarly(ctx, remoteAddr, tlsCfg, cfg)
				},
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		server: http3.Server{
			TLSConfig: tlsConfig,
		},
		transport: transport,
		tlsConfig: tlsConfig,
	}

	return
}

func (client *Client) Start() (err error) {
	defer client.transport.Conn.Close()

	client.server.Handler = client.handler()

	ln, err := client.transport.ListenEarly(client.tlsConfig, nil)
	if err != nil {
		panic(err)
	}

	err = client.server.ServeListener(ln)
	if err != nil {
		err = gerror.Wrapf(err, "http3 server.ServeListener failed")
	}

	return
}

func main() {
	client := NewClient()

	go func() {
		g.Log().Info(context.Background(), "Client is listening...")
		g.Log().Errorf(context.Background(), "%+v", client.Start())
	}()

	for ; ; time.Sleep(time.Second) {
		result, err := utils.Req(&client.client, remoteHost, []byte(`client -`))
		if err != nil {
			g.Log().Errorf(context.Background(), "%+v", err)
			continue
		}
		g.Log().Infof(context.Background(), "%s> client", string(result))
	}
}
