package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	addr         = flag.String("addr", ":4443", "http service address")
	wssAddr      = flag.String("wssAddr", ":4445", "websocket service address")
	peerConns    = make(map[string]*websocket.Conn)
	upgrader     = websocket.Upgrader{} // use default options
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "public/index.html")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	peerId := r.URL.Query().Get("peerId")
	if peerId == "" {
		c.Close()
		return
	}

	peerConns[peerId] = c
	defer delete(peerConns, peerId)

	// Send current peers to the connected peer
	sendCurrentPeers(c, peerId)

	// Broadcast new peer join to all other peers
	broadcastPeerJoin(peerId)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		handleMessage(peerId, message)
	}

	// Broadcast peer leave to all other peers
	broadcastPeerLeave(peerId)
}

func sendCurrentPeers(conn *websocket.Conn, excludedPeerId string) {
	var peerList []string
	for id := range peerConns {
		if id != excludedPeerId {
			peerList = append(peerList, id)
		}
	}
	message := map[string]interface{}{
		"messageId": "CURRENT_PEERS",
		"messageData": map[string]interface{}{
			"peerList": peerList,
		},
	}
	msg, _ := json.Marshal(message)
	conn.WriteMessage(websocket.TextMessage, msg)
}

func broadcastPeerJoin(peerId string) {
	message := map[string]interface{}{
		"messageId": "PEER_JOIN",
		"messageData": map[string]interface{}{
			"peerId": peerId,
		},
	}
	msg, _ := json.Marshal(message)
	for _, conn := range peerConns {
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func broadcastPeerLeave(peerId string) {
	message := map[string]interface{}{
		"messageId": "PEER_LEAVE",
		"messageData": map[string]interface{}{
			"peerId": peerId,
		},
	}
	msg, _ := json.Marshal(message)
	for _, conn := range peerConns {
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func handleMessage(peerId string, message []byte) {
	var msg map[string]interface{}
	json.Unmarshal(message, &msg)
	switch msg["messageId"] {
	case "PROXY":
		toPeerId := msg["toPeerId"].(string)
		if conn, ok := peerConns[toPeerId]; ok {
			conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	// Setup TLS
	cert, _ := ioutil.ReadFile("certificate/server.crt")
	key, _ := ioutil.ReadFile("certificate/server.key")
	certPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{certPair}}
	server := &http.Server{
		Addr:      *addr,
		Handler:   http.HandlerFunc(serveHome),
		TLSConfig: tlsConfig,
	}

	go server.ListenAndServeTLS("", "")

	// WebSocket server setup
	wssServer := &http.Server{
		Addr:      *wssAddr,
		Handler:   http.HandlerFunc(handleConnections),
		TLSConfig: tlsConfig,
	}
	log.Fatal(wssServer.ListenAndServeTLS("", ""))
}
