package main

import (
	"context"
	"encoding/json"
	"log"
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

func main() {
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
	router.HandleFunc("/{container_id}/ws", rbiHandler)

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
