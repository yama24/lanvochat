### INSTRUKSIN: EPIC 4 - VIDEO CALL (Kelancaran & RTP)

**Fokus Task:** Implementasi protokol *real-time* Audio/Video.

**PRINSIP UTAMA:** Lancar (UDP/RTP) dan Aman (Signaling via TCP).

| Stack Kunci | Versi/Spesifikasi |
| :--- | :--- |
| **Go Backend** | Go 1.25.3 |
| **RTP Library** | Pion WebRTC (atau setara Go) |
| **Jaringan** | UDP/RTP dan TCP |

**Task Rinci:**
1.  **Signaling Protocol over TCP:** Definisikan dan implementasikan pesan *signaling* (e.g., CALL_INVITE) menggunakan koneksi **TCP** yang andal.
2.  **RTP/UDP Core Setup:** Integrasikan library Go **RTP**. Buat logika Go untuk mengelola *socket* **UDP** terpisah untuk *streaming* media.
3.  **Media Capture & Stream:** Buat fungsi Go yang menangkap data dari kamera/mikrofon dan mengirimkan paket **RTP** ke *peer* tujuan.
4.  **Video Call UI Component:** Buat antarmuka React untuk panggilan video, termasuk tampilan *video feed* dan kontrol (*mute*, *end call*).

Berikan kerangka kode Go untuk setup RTP dan fungsi *signaling* TCP.