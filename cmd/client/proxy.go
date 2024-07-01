package main

import (
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go"
	"github.com/wzshiming/socks5"
	"net"
)

// for QUIC protocol (quic-go)
type sConnWrapper struct {
	net.PacketConn
}

func (c *sConnWrapper) SetReadBuffer(bytes int) error {
	socks5udpConn := c.PacketConn.(*socks5.UDPConn)
	udpConn := socks5udpConn.PacketConn.(*net.UDPConn)

	return udpConn.SetReadBuffer(bytes)
}

func (c *sConnWrapper) SetWriteBuffer(bytes int) error {
	socks5udpConn := c.PacketConn.(*socks5.UDPConn)
	udpConn := socks5udpConn.PacketConn.(*net.UDPConn)

	return udpConn.SetWriteBuffer(bytes)
}

func proxyDialer(proxyURL string) func(ctx context.Context, addr string, tlsCfg *tls.Config, cfg *quic.Config) (quic.EarlyConnection, error) {
	dialer, err := socks5.NewDialer(proxyURL)
	if err != nil {
		panic(err)
	}

	return func(ctx context.Context, addr string, tlsCfg *tls.Config, cfg *quic.Config) (quic.EarlyConnection, error) {
		proxyConn, err := dialer.DialContext(ctx, "udp", addr)
		if err != nil {
			return nil, err
		}

		remoteAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			return nil, err
		}

		connWrapper := &sConnWrapper{proxyConn.(net.PacketConn)}
		earlyConn, err := quic.DialEarly(ctx, connWrapper, remoteAddr, tlsCfg, cfg)
		if err != nil {
			return nil, err
		}

		return earlyConn, nil
	}
}
