### EPIC 2: NETWORK DISCOVERY & CORE CHAT (TCP/UDP)

**Fokus Task:** Implementasi Komunikasi Dasar (Epik 2).

**PRINSIP UTAMA:** Efisiensi Jaringan (UDP Multicast) dan Keandalan (TCP).

**Task Rinci:**
1.  **Peer Discovery (UDP Multicast):** Buat **goroutine** Go untuk mengirim dan menerima paket "I'm Alive" melalui **UDP Multicast** untuk mendeteksi peers di LAN.
2.  **TCP Listening Handler:** Buat *goroutine* Go untuk memulai *listening socket* **TCP** dan *handler* untuk membaca pesan JSON masuk. Pesan yang diterima harus diteruskan ke fungsi `SaveMessage` dari Epik 1.
3.  **Message Sending Function (TCP Client):** Buat fungsi Go (`SendMessageToPeer(ip, message)`) yang membuka koneksi TCP keluar, mengirimkan *payload* pesan JSON, dan menutup koneksi.
4.  **Wails Event Bridge:** Gunakan `wails.Runtime.EventsEmit()` di Go *backend* untuk mengirim *event* (`messageReceived`) ke *frontend* setelah pesan baru diterima.

Berikan kerangka kode Go untuk modul UDP Multicast dan fungsi TCP *listening* utama.