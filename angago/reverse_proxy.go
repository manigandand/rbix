package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func angagoServer(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "container_id")

	target := fmt.Sprintf("http://%s:%s/%s/ws", containerID, SqrxRBIBoxPort, containerID)
	resp, err := http.Get(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}
	w.WriteHeader(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
func angagoServer(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "container_id")
	// check if the container id valid or not

	targetURL = "http://sqrx-api:8080"
	log.Printf("Proxying the Host %s to %s\n", r.Host, targetURL)
	angagoProxy(w, r, angagoConfig.AbsTargetURL(targetURL))
}

// angagoProxy Serve a reverse proxy for a given url
func angagoProxy(w http.ResponseWriter, r *http.Request, target string) {
	remote, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Transport = &customTransport{http.DefaultTransport}
	proxy.ModifyResponse = modifyProxyResponse

	r.URL.Scheme = remote.Scheme
	r.URL.Host = remote.Host
	// r.Host = remote.Host
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

	proxy.ServeHTTP(w, r)
}
*/
