package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
)

// https://www.iana.org/assignments/media-types/media-types.xhtml
func isBinary(contentType string) bool {
	group := strings.SplitN(contentType, "/", 2)[0]
	switch group {
	case "image", "audio", "video":
		return true
	default:
		return false
	}
}

func normaliseTrailing(in []byte) []byte {
	return append(bytes.TrimRight(in, "\n\r"), '\n')
}

func withHeader(seq uint64, head string, body []byte) string {
	const margin = "#########"
	return fmt.Sprintf("\n%s (%d) %s %s\n%s", margin, seq, head, margin, body)
}

type logTransport struct {
	out  io.Writer
	seq  *uint64
	trip http.RoundTripper
}

func (l logTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	bodyReqRaw, _ := httputil.DumpRequest(r, true)
	bodyReq := normaliseTrailing(bodyReqRaw)

	w, err := l.trip.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	isBinary := isBinary(w.Header.Get("content-type"))
	bodyRespRaw, _ := httputil.DumpResponse(w, !isBinary)
	bodyResp := normaliseTrailing(bodyRespRaw)
	if isBinary {
		bodyResp = append(bodyResp, "\nBINARY DATA\n"...)
	}

	seq := atomic.AddUint64(l.seq, 1)
	fmt.Fprint(l.out, withHeader(seq, "request", bodyReq))
	fmt.Fprint(l.out, withHeader(seq, "response", bodyResp))
	return w, nil
}

var _ http.RoundTripper = (*logTransport)(nil)

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
	toURL, err := url.Parse(*flagTo)
	if err != nil {
		log.Fatalf("error parsing `-to`: %v", err)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(toURL)
			r.Out.Host = toURL.Host
		},
		Transport: logTransport{
			out:  os.Stdout,
			seq:  new(uint64),
			trip: http.DefaultTransport,
		},
	}

	log.Printf("listening on %q", *flagListenAddr)
	if err := http.ListenAndServe(*flagListenAddr, proxy); err != nil {
		log.Fatalf("error starting: %v", err)
	}
}
