# Implementasi WebRTC Data Channel dengan Pion (Golang)

Proyek ini adalah implementasi sederhana namun fungsional dari teknologi **WebRTC** menggunakan library **Pion** di bahasa pemrograman Go (Golang).

Aplikasi ini mendemonstrasikan bagaimana server Golang dapat bertindak sebagai *Peer* WebRTC untuk melakukan streaming data secara *real-time* (bi-directional) ke browser klien tanpa menggunakan HTTP Polling atau WebSocket standar, melainkan menggunakan protokol SCTP/UDP melalui WebRTC.

## Prasyarat

- **Go (Golang)** versi 1.20 atau yang lebih baru.
- Web Browser modern (Chrome, Firefox, Edge, atau Safari).

## Cara Menjalankan Aplikasi

1. **Clone atau Download** repository ini.
2. Buka terminal di dalam folder proyek.
3. Unduh dependency yang diperlukan:
   ```bash
   go mod tidy
Jalankan server:

Bash

go run main.go
Buka browser dan akses:

http://localhost:8080
Klik tombol "Start Connection" pada halaman web.

Cara Menggunakan
Koneksi: Setelah tombol "Start Connection" diklik, perhatikan status berubah menjadi Connected.

Streaming Server: Di kotak Log, Anda akan melihat server secara otomatis mengirimkan Waktu Server setiap detik. Ini membuktikan adanya koneksi persisten.

Kirim Pesan: Ketik pesan di kolom input dan tekan Send. Server akan membalas (echo) pesan Anda melalui jalur data WebRTC yang sama.

Penjelasan Teknis (Cara Kerja)
Aplikasi ini bekerja dalam dua fase utama: Signaling (Sinyalisasi) dan Data Transfer.

1. Signaling (HTTP Handshake)
Sebelum WebRTC dapat terhubung, kedua belah pihak (Browser dan Server) harus saling bertukar informasi jaringan dan format data. Proses ini disebut Signaling.

Offer: Browser membuat SDP (Session Description Protocol) Offer dan mengirimkannya ke server melalui HTTP POST ke endpoint /sdp.

Answer: Server Golang menerima Offer tersebut, menyiapkan PeerConnection menggunakan library Pion, dan mengembalikan SDP Answer.

2. Peer-to-Peer Connection (WebRTC)
Setelah pertukaran SDP selesai:

Koneksi HTTP diputus.

Browser dan Server membangun koneksi langsung (Peer-to-Peer) menggunakan protokol transportasi WebRTC (biasanya di atas UDP).

Data Channel dibuka. Server tidak lagi menunggu permintaan (request) dari klien. Sebaliknya, server menjalankan Goroutine yang secara proaktif "mem-push" data waktu ke klien setiap detik.

Berbeda dengan aplikasi web biasa (Request-Response), aplikasi ini mendemonstrasikan:

Server-Initiated Events: Server dapat mengirim data kapan saja tanpa diminta klien.

UDP Transport: Penggunaan protokol yang lebih cepat untuk transmisi data real-time.

Pion Library: Implementasi native WebRTC di sisi backend (Go), bukan sekadar perantara (signaling server) saja.
