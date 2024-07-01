package utils

import (
	"bytes"
	"github.com/gogf/gf/v2/errors/gerror"
	"io"
	"net/http"
)

func Req(client *http.Client, host string, content []byte) (result []byte, err error) {
	req, err := http.NewRequest(
		http.MethodPost,
		"https://"+host,
		bytes.NewReader(content),
	)
	if err != nil {
		err = gerror.Wrapf(err, "new http request failed")
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		err = gerror.Wrapf(err, "do http request failed")
		return
	}

	result, err = io.ReadAll(resp.Body)
	if err != nil {
		err = gerror.Wrapf(err, "read http response failed")
		return
	}

	return
}
