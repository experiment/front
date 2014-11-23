package main

import (
	"io"
	"log"
	"net/http"
	"regexp"
)

type Proxy struct {
	config           Config
	contentTypeRegex *regexp.Regexp
}

func newProxy(config Config) *Proxy {
	// Todo, add config error checking

	contentTypeRegex, err := regexp.Compile(config.allowedContentTypes)
	if err != nil {
		log.Panic(err)
	}

	return &Proxy{
		config:           config,
		contentTypeRegex: contentTypeRegex,
	}
}

func (p Proxy) handler(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query().Get("url")
	if request != "" {
		p.proxyRequest(w, request)
	} else {
		log.Println("No request url to proxy")
		http.Error(w, "No request url to proxy", http.StatusBadRequest)
	}
}

func (p Proxy) proxyRequest(w http.ResponseWriter, request string) {
	resp, err := http.Get(request) // http.Get follows up to 10 redirects
	if err != nil {
		log.Print(err)
		// Todo, handle specific errors
		http.Error(w, "Could not proxy", http.StatusInternalServerError)
	}

	if p.contentTypeRegex.MatchString(resp.Header.Get("Content-Type")) {
		p.writeResponse(w, resp)
	} else {
		log.Println("Upstream content doesn't match configured allowedContentTypes")
		http.Error(w, "Upstream content doesn't match configured allowedContentTypes", http.StatusBadRequest)
	}
}

func (p Proxy) writeResponse(w http.ResponseWriter, resp *http.Response) {
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	io.Copy(w, resp.Body)
	resp.Body.Close()
}
