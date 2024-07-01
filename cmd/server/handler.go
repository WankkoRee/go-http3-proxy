package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/number571/go-http3-proxy/utils"
	"io"
	"net/http"
)

func (server *Server) handler() (mux *http.ServeMux) {
	mux = http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintln(w, "hello, there is server!")
		case http.MethodPost:
			result, err := io.ReadAll(r.Body)
			if err != nil {
				g.Log().Errorf(context.Background(), "%+v", err)
				return
			}
			fmt.Fprintf(w, "%s> server -", string(result))

			g.Log().Debugf(context.Background(), "remote: %s", r.RemoteAddr)
			result, err = utils.Req(server.client(), r.RemoteAddr, []byte(`server -`))
			if err != nil {
				g.Log().Errorf(context.Background(), "%+v", err)
				return
			}
			g.Log().Infof(context.Background(), "%s> server", string(result))
		default:
			g.Log().Errorf(context.Background(), "%+v", "method is not supported")
		}
	})

	return
}
