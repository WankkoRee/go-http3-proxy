package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"io"
	"net/http"
)

func (client *Client) handler() (mux *http.ServeMux) {
	mux = http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintln(w, "hello, there is client!")
		case http.MethodPost:
			result, err := io.ReadAll(r.Body)
			if err != nil {
				g.Log().Errorf(context.Background(), "%+v", err)
				return
			}
			fmt.Fprintf(w, "%s> client -", string(result))
		default:
			g.Log().Errorf(context.Background(), "%+v", "method is not supported")
		}
	})

	return
}
