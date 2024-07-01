package main

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/wzshiming/socks5"
)

func main() {
	// TCP is used only to create a `client-proxy` connection,
	// but all other traffic (client-proxy-server) already uses the UDP protocol
	g.Log().Info(context.Background(), "Proxy is listening...")
	g.Log().Errorf(context.Background(), "%+v", socks5.NewServer().ListenAndServe("tcp", "0.0.0.0:1080"))
}
