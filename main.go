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
var binaryTypes = []string{
	"image/",
	"audio/",
	"video/",
}

func isBinary(contentType string) bool {
	if contentType == "" {
		return false
	}
	for _, pre := range binaryTypes {
		if strings.HasPrefix(contentType, pre) {
			return true
		}
	}
	return false
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
	// record request, should be done before passing to the transport
	bodyReqRaw, _ := httputil.DumpRequest(r, true)
	bodyReq := normaliseTrailing(bodyReqRaw)
	// make the request
	w, err := l.trip.RoundTrip(r)
	// record response
	defer func() {
		isBinary := isBinary(w.Header.Get("content-type"))
		bodyRespRaw, err := httputil.DumpResponse(w, !isBinary)
		if err != nil {
			fmt.Println("ok", err)
		}
		bodyResp := normaliseTrailing(bodyRespRaw)
		if isBinary {
			bodyResp = append(bodyResp, "\n<< BINARY DATA >>\n"...)
		}
		seq := atomic.AddUint64(l.seq, 1)
		fmt.Fprint(l.out, withHeader(seq, "request", bodyReq))
		fmt.Fprint(l.out, withHeader(seq, "response", bodyResp))
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
	flagToURL, err := url.Parse(*flagTo)
	if err != nil {
		log.Fatalf("error parsing `-to`: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(flagToURL)
	proxy.Transport = logTransport{
		out:  os.Stdout,
		seq:  new(uint64),
		trip: http.DefaultTransport,
	}
	http.Handle("/", proxy)
	log.Printf("listening on %q", *flagListenAddr)
	if err := http.ListenAndServe(*flagListenAddr, nil); err != nil {
		log.Fatalf("error starting: %v", err)
	}
}
