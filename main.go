package main

import (
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

func main() {
	usage := `ishi.

Usage:
  ishi [-l=<port>] <upstream>
  ishi -h | --help
  ishi --version

Arguments:
  upstream  Upstream host.

Options:
  -h --help             Show help.
  --version             Show version.
  -l --listen=<port>    Specify port to listen.

Examples:
  ishi 192.168.10.2
  ishi http://192.168.10.2
  ishi https://secure.example.com`

	arguments, _ := docopt.Parse(usage, nil, true, "1.0", false)
	upstream := arguments["<upstream>"]
	listen := arguments["--listen"]

	var err error

	// define port to listen
	var port int
	if listen == nil {
		port, err = findAvailablePort()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	} else {
		port, err = strconv.Atoi(listen.(string))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid listening port\n")
			os.Exit(1)
		}
	}

	// define upstream host and scheme
	u, err := url.Parse(upstream.(string))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	var upstreamHost string
	if u.Host == "" {
		upstreamHost = u.Path
	} else {
		upstreamHost = u.Host
	}
	scheme := "http"
	if u.Scheme != "" {
		scheme = u.Scheme
	}

	// start reverse proxy server
	fmt.Printf("Listening on 127.0.0.1:%d\nFowarding to %s\n", port, upstream)
	err = httpfwd(fmt.Sprintf(":%d", port), scheme, upstreamHost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func findAvailablePort() (int, error) {
	for port := 8000; port < 9000; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			defer ln.Close()
			return port, nil
		}
	}
	return 0, errors.New("There is no available port to listen")
}

func httpfwd(listenAddr, scheme, remoteHost string) error {
	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			r.Host = remoteHost
			p := httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: scheme,
				Host:   remoteHost,
			})
			p.ServeHTTP(w, r)
		})
	return http.ListenAndServe(listenAddr, nil)
}
