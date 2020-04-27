package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync/atomic"
)

func withHeader(seq uint64, head string, body []byte) string {
	const margin = "#########"
	return fmt.Sprintf("%s (%d) %s %s\n%s", margin, seq, head, margin, body)
}

type logTransport struct {
	out  io.Writer
	seq  *uint64
	trip http.RoundTripper
}

func (l logTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	w, err := l.trip.RoundTrip(r)
	defer func() {
		bodyReq, _ := httputil.DumpRequest(r, true)
		bodyResp, _ := httputil.DumpResponse(w, true)
		seq := atomic.AddUint64(l.seq, 1)
		fmt.Fprintf(l.out, withHeader(seq, "request", bodyReq))
		fmt.Fprintf(l.out, withHeader(seq, "response", bodyResp))
	}()
	return w, err
}

func main() {
	flagListenAddr := flag.String("listen-addr", "", "address to listen on, eg. :5050")
	flagTo := flag.String("to", "", "address to proxy to, eg. http://localhost:4040")
	flag.Parse()
	if *flagListenAddr == "" {
		log.Fatalf("please provide `-listen-addr`")
	}
	if *flagTo == "" {
		log.Fatalf("please provide `-to`")
	}
	transport := logTransport{
		out:  os.Stdout,
		seq:  new(uint64),
		trip: http.DefaultTransport,
	}
	toURL, err := url.Parse(*flagTo)
	if err != nil {
		log.Fatalf("error parsing `-to`: %v", err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println()
		proxy := httputil.NewSingleHostReverseProxy(toURL)
		proxy.Transport = transport
		r.URL.Host = toURL.Host
		r.URL.Scheme = toURL.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = toURL.Host
		proxy.ServeHTTP(w, r)
	})
	log.Printf("listening on %q", *flagListenAddr)
	if err := http.ListenAndServe(*flagListenAddr, nil); err != nil {
		log.Fatalf("error starting: %v", err)
	}
}
