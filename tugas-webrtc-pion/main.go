package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pion/webrtc/v3"
)

// SdpContainer used for HTTP signaling
type SdpContainer struct {
	Sdp string `json:"sdp"`
}

func main() {
	// 1. Setup WebRTC Configuration (Public STUN server usually needed for real internet,
	// but localhost works without it. We add google's just in case).
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// 2. Serve the HTML file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// 3. Handle the Signaling (The "Handshake")
	http.HandleFunc("/sdp", func(w http.ResponseWriter, r *http.Request) {
		// Create a new PeerConnection
		peerConnection, err := webrtc.NewPeerConnection(config)
		if err != nil {
			http.Error(w, "Failed to create PeerConnection", http.StatusInternalServerError)
			return
		}

		// --- ADVANCED PART: DATA CHANNEL HANDLING ---
		
		// Event: When the Browser creates a data channel, the Server accepts it here
		peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
			fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

			// Register channel opening handling
			d.OnOpen(func() {
				fmt.Printf("Data channel '%s'-'%d' open. \n", d.Label(), d.ID())

				// Send an initial message
				d.SendText("Connected to Golang Pion Server!")

				// ADVANCED: Start a goroutine to stream server time every second
				// This proves the connection is persistent and bi-directional
				go func() {
					ticker := time.NewTicker(1 * time.Second)
					for range ticker.C {
						// Check if channel is still open before sending
						if d.ReadyState() == webrtc.DataChannelStateOpen {
							msg := fmt.Sprintf("Server Time: %s", time.Now().Format(time.RFC1123))
							d.SendText(msg)
						} else {
							ticker.Stop()
							return
						}
					}
				}()
			})

			// Register message handling (Echo logic)
			d.OnMessage(func(msg webrtc.DataChannelMessage) {
				log.Printf("Message from Client: '%s'\n", string(msg.Data))
				
				// Reply to the client
				response := fmt.Sprintf("Server received: '%s'", string(msg.Data))
				d.SendText(response)
			})
		})

		// --- SIGNALING LOGIC (SDP EXCHANGE) ---

		// Parse the Offer from the Browser
		var offer SdpContainer
		if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Set the remote SessionDescription (The Browser's SDP)
		if err := peerConnection.SetRemoteDescription(webrtc.SessionDescription{
			Type: webrtc.SDPTypeOffer,
			SDP:  offer.Sdp,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create an Answer
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Sets the LocalDescription, and starts our UDP listeners
		if err := peerConnection.SetLocalDescription(answer); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Wait for ICE Gathering to complete (Simple way for assignment)
		// In production, we would exchange ICE candidates via trickling.
		<-webrtc.GatheringCompletePromise(peerConnection)

		// Send the Answer back to the Browser
		response := SdpContainer{
			Sdp: peerConnection.LocalDescription().SDP,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Server starting on http://localhost:8080")
	panic(http.ListenAndServe(":8080", nil))
}