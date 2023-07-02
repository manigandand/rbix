package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var (
	defaultUpgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow connections from sqrx-client container only
			// if r.Host != "sqrx-client" {
			// 	return false
			// }
			return true
		},
		HandshakeTimeout: 5 * time.Second,
	}

	defaultDialer = websocket.DefaultDialer
)

func getContainerInfo(cintainerID string) (time.Time, error) {
	now := time.Now()
	url := fmt.Sprintf("%s/v1/status/%s", SqrxAPIServer, cintainerID)
	resp, err := http.Get(url)
	if err != nil {
		return now, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return now, fmt.Errorf("invalid container id")
	}

	res := map[string]time.Time{}
	if err := json.NewDecoder(resp.Body).Decode(&resp); err != nil {
		return now, err
	}

	return res["valid_till"], nil
}

func angagoServer(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "container_id")
	// check if the container id valid or not from sqrx-api
	// validTill, err := getContainerInfo(containerID)
	// if err != nil {
	// 	log.Println("invalid container id:", err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// log.Println("session valid till:", validTill)

	targetURL := fmt.Sprintf("ws://%s:%s/%s/ws", containerID, SqrxRBIBoxPort, containerID)
	remote, _ := url.Parse(targetURL)

	log.Printf("Proxying the Host %s to %s\n", r.Host, targetURL)
	angagoProxy(w, r, remote)
}

// angagoProxy Serve a reverse proxy for a given url
func angagoProxy(w http.ResponseWriter, r *http.Request, remote *url.URL) {
	// prepare remote(target) websocket connection
	// this is act as client connextion(sqrx-proxy) to the remote server(sqrx-rbi-box)
	remoteWSConn, resp, err := defaultDialer.Dial(remote.String(), remoteReqHeaders(r))
	if err != nil {
		log.Printf("couldn't dial to remote backend url %s", err)
		if resp != nil {
			// handshake failed so we have to return the error state to client
			if err := copyResponse(w, resp); err != nil {
				log.Printf("failed to write error state to client: %s", err)
				return
			}
		}
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	defer remoteWSConn.Close()

	// upgrade client websocket connection
	// this is act as server connection(sqrx-proxy) to the client(sqrx-client)
	clientWSConn, err := defaultUpgrader.Upgrade(w, r, clientReqHeaders(resp))
	if err != nil {
		log.Printf("websocketproxy: couldn't upgrade %s", err)
		return
	}
	defer clientWSConn.Close()

	// proxy the messages between client and remote
	errClient := make(chan error, 1)
	errRemote := make(chan error, 1)

	go echo(clientWSConn, remoteWSConn, errClient)
	go echo(remoteWSConn, clientWSConn, errRemote)

	select {
	case err = <-errClient:
		log.Printf("echo: Error when copying from client to backend: %v\n", err)
	case err = <-errRemote:
		log.Printf("echo: Error when copying from backend to client: %v\n", err)
	}

	log.Println("closing client and remote connections")
	return
}

// echo copies messages from src to dst and back until an error occurs.
func echo(srcConn, dstConn *websocket.Conn, errCh chan error) {
	for {
		// read from srcConn
		msgType, msg, err := srcConn.ReadMessage()
		if err != nil {
			m := websocket.FormatCloseMessage(websocket.CloseNormalClosure, fmt.Sprintf("%v", err))
			if e, ok := err.(*websocket.CloseError); ok {
				m = websocket.FormatCloseMessage(e.Code, e.Text)
			}

			dstConn.WriteMessage(websocket.CloseMessage, m)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("src conn closed:", err, " closing dst conn")
				dstConn.Close()
			}
			errCh <- err
			break
		}

		// and write to dstConn
		err = dstConn.WriteMessage(msgType, msg)
		if err != nil {
			errCh <- err
			break
		}
	}
}

func clientReqHeaders(resp *http.Response) http.Header {
	upgradeHeader := http.Header{}
	if hdr := resp.Header.Get("Sec-Websocket-Protocol"); hdr != "" {
		upgradeHeader.Set("Sec-Websocket-Protocol", hdr)
	}
	if hdr := resp.Header.Get("Set-Cookie"); hdr != "" {
		upgradeHeader.Set("Set-Cookie", hdr)
	}
	return upgradeHeader
}

func remoteReqHeaders(req *http.Request) http.Header {
	requestHeader := http.Header{}
	if origin := req.Header.Get("Origin"); origin != "" {
		requestHeader.Add("Origin", origin)
	}
	for _, prot := range req.Header[http.CanonicalHeaderKey("Sec-WebSocket-Protocol")] {
		requestHeader.Add("Sec-WebSocket-Protocol", prot)
	}
	for _, cookie := range req.Header[http.CanonicalHeaderKey("Cookie")] {
		requestHeader.Add("Cookie", cookie)
	}
	if req.Host != "" {
		requestHeader.Set("Host", req.Host)
	}

	requestHeader.Set("X-Forwarded-Host", req.Header.Get("Host"))
	requestHeader.Set("X-Forwarded-Proto", "http")
	if req.TLS != nil {
		requestHeader.Set("X-Forwarded-Proto", "https")
	}

	return requestHeader
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		fmt.Println(k, vv)
		dst.Set(k, vv[0])
	}
}

func copyResponse(rw http.ResponseWriter, resp *http.Response) error {
	copyHeader(rw.Header(), resp.Header)
	rw.WriteHeader(resp.StatusCode)
	defer resp.Body.Close()

	_, err := io.Copy(rw, resp.Body)
	return err
}

// --- test ---
func angagoServerHttp(w http.ResponseWriter, r *http.Request) {
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
