package main

import (
	"os"
)

// docker use
func init() {
	switch {
	case len(os.Args) != 2:
		return
	case os.Args[1] != "docker":
		return
	}

	proxyHost = "host.docker.internal:1080"
	remoteHost = "host.docker.internal:8080"
}
