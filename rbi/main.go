package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v3/host"
)

var (
	containerID string
	hostInfo    *host.InfoStat
)

// TCPProxy tcp to websocket proxy conn object
var TCPProxy *Websockify

// Websockify holds the vcnserver tcp connection and client websocket connection
type Websockify struct {
	wsConn  *websocket.Conn
	tcpAddr *net.TCPAddr
	tcpConn *net.TCPConn
}

func initTCPVCNServer() {
	addr := "localhost:5901"
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal("failed to resolve vcn server: " + err.Error())
		return
	}

	var p = &Websockify{
		tcpAddr: tcpAddr,
	}

	tcpConn, err := net.DialTCP(p.tcpAddr.Network(), nil, p.tcpAddr)
	if err != nil {
		log.Fatal("failed to connect to vcn server: " + err.Error())
		return
	}
	p.tcpConn = tcpConn
	log.Println("connected to vcn server: " + addr)

	TCPProxy = p
}

func main() {
	initTCPVCNServer()

	addr := ":8888"
	if os.Getenv("PORT") != "" {
		addr = ":" + os.Getenv("PORT")
	}

	// read container_id from env
	containerID = os.Getenv("CONTAINER_ID")
	if containerID == "" {
		log.Fatal("container id is not provided")
	}
	hostInfo, _ = host.Info()

	router := chi.NewRouter()
	router.HandleFunc("/{container_id}/ws", websockifyHandler)

	server := http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	interruptChan := make(chan os.Signal, 1)
	go func() {
		signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		// Block until we receive our signal.
		<-interruptChan
		log.Println("os interrupt signal received, shutting down")

		// shutdown server
		if err := server.Shutdown(context.Background()); err != nil {
			log.Println("shtdown error:", err.Error())
		}
		TCPProxy.close("[os_interrupt]]")

		os.Exit(1)
	}()

	// shutdown server after 10 minute
	go func() {
		timer := time.NewTimer(10 * time.Minute)
		<-timer.C
		log.Println("shutdown server after 10 minute")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Println("shtdown error:", err.Error())
		}
		TCPProxy.close("[shutdown_after_10_minute]")

		os.Exit(1)
	}()

	log.Println("rbi listening on ", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from sqrx-angago container only
		// if r.Host != "sqrx-angago" {
		// 	return false
		// }
		return true
	},
	HandshakeTimeout: 5 * time.Second,
}

func websockifyHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Scheme != "ws" {
	// 	jsonUnauthResp(w)
	// 	return
	// }

	reqContainerID := chi.URLParam(r, "container_id")
	if reqContainerID != containerID {
		log.Println("invalid container id")
		jsonUnauthResp(w)
		return
	}

	// new websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		jsonUnauthResp(w)
		return
	}
	TCPProxy.wsConn = conn

	log.Printf("WebSocket %s connected to %+v:%d\n",
		TCPProxy.wsConn.RemoteAddr(), TCPProxy.tcpAddr.IP, TCPProxy.tcpAddr.Port)

	// start the tcp to websocket proxy
	TCPProxy.start()
}

func (p *Websockify) start() {
	go p.ReadWebSocket()
	go p.ReadTCP()
}

func (p *Websockify) close(cause string) {
	log.Println(cause, " closing tcp and websocket connection")
	p.tcpConn.Close()
	p.wsConn.Close()
}

// ReadWebSocket reads from the WebSocket and
// writes to the backend TCP connection.
func (p *Websockify) ReadWebSocket() {
	log.Println("reading websocket stream")
	for {
		_, data, err := p.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("WebSocket closed:", err)
			}
			log.Println("[error] read websocket:", err.Error())

			p.close("[error_read_websocket] " + err.Error())
			break
		}
		log.Println("[read] from websocket:", string(data))

		_, err = p.tcpConn.Write(data)
		if err != nil {
			log.Println("[error] write webSocketToTCP:", err.Error())

			// reconnect to vcn server
			tcpConn, err := net.DialTCP(p.tcpAddr.Network(), nil, p.tcpAddr)
			if err != nil {
				p.close("[error_reconn_vcn] " + err.Error())
				log.Println("failed to connect to vcn server: " + err.Error())
				return
			}
			p.tcpConn = tcpConn

			p.tcpConn.Write(data)
		}
	}
}

// ReadTCP reads from the backend TCP connection and
// writes to the WebSocket.
func (p *Websockify) ReadTCP() {
	log.Println("reading tcp stream")
	buffer := make([]byte, 6000)

	for {
		bytesRead, err := p.tcpConn.Read(buffer)
		if err != nil {
			log.Println("[error] read tcp:", err.Error())
			p.close("[error_read_tcp] " + err.Error())
			break
		}
		log.Println("[read] from VCN server: ", bytesRead)

		if err := p.wsConn.WriteMessage(websocket.BinaryMessage, buffer[:bytesRead]); err != nil {
			p.close("[error_write_to_ws] " + err.Error())
			log.Println("[error] write tcpToWebSocket:", err.Error())
			break
		}
	}
}

// non vcn websocket handler
func rbiHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Scheme != "ws" {
	// 	jsonUnauthResp(w)
	// 	return
	// }

	reqContainerID := chi.URLParam(r, "container_id")
	if reqContainerID != containerID {
		log.Println("invalid container id")
		jsonUnauthResp(w)
		return
	}

	// new websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		jsonUnauthResp(w)
		return
	}
	defer conn.Close()
	TCPProxy.wsConn = conn

	log.Printf("WebSocket %s connected to %+v:%d\n",
		TCPProxy.wsConn.RemoteAddr(), TCPProxy.tcpAddr.IP, TCPProxy.tcpAddr.Port)

	for {
		// Read message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("WebSocket closed:", err)
				os.Exit(1)
			}

			log.Println("WebSocket read error:", err)
			break
		}

		// Process message
		log.Printf("Received message: %s\n", message)

		res := map[string]interface{}{
			"data": map[string]interface{}{
				"reqContainerID": reqContainerID,
				"host":           r.Host,
				"host_info":      hostInfo,
				"message":        "you're in safe zone with rbi",
				"client_message": string(message),
			},
			"meta": map[string]interface{}{
				"status":      "ok",
				"status_code": http.StatusOK,
			},
		}

		resp, err := json.Marshal(&res)
		if err != nil {
			log.Println("error:", err.Error())
			return
		}

		// Write message back to client
		// render the remote browser frame
		err = conn.WriteMessage(websocket.TextMessage, resp)
		if err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}
}

func jsonUnauthResp(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	resp := map[string]interface{}{
		"message": "unauthorized access",
	}
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Println("error:", err.Error())
	}
	return
}
