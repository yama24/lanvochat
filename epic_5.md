### EPIC 5: TESTING & FINALIZATION

**Fokus Task:** Pengujian, validasi keandalan, dan persiapan *deployment*.

**PRINSIP UTAMA:** Keandalan (Testing) dan Keringanan (Final Build).

**Task Rinci:**
1.  **Unit Testing:** Tulis *unit test* Go yang berfokus pada modul **SQLite** dan fungsi **TCP Message Handling** untuk menjamin integritas data.
2.  **Integration Testing:** Buat rencana pengujian penuh di lingkungan **LAN nyata** untuk memverifikasi Peer Discovery, Chat, dan Video Call.
3.  **Final Build & Deployment:** Jalankan perintah `wails build` untuk menghasilkan **single binary executable** yang **ringan** untuk platform target (Windows, macOS, Linux).
4.  **Dokumentasi:** Selesaikan *README* yang mencakup prasyarat dan instruksi *deployment*.

Berikan contoh struktur rencana *Integration Testing* yang fokus pada *LAN environment*.