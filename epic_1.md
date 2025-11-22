### INSTRUKSIN: EPIC 1 - SETUP PROYEK & BASIS DATA SQLITE

**Fokus Task:** Menyiapkan lingkungan dasar (Wails/React) dan lapisan data andal (SQLite) untuk aplikasi chat desktop LAN P2P.

**STACK WAJIB:** Go 1.25.3, Wails, React/TypeScript (Yarn), SQLite.

**Task Rinci:**
1.  **Project Initialization & Config:** Buat kerangka proyek Wails baru dengan template React/TypeScript. Konfigurasi `wails.json` untuk **Frameless Window** dan mengaktifkan **System Tray**.
2.  **SQLite DB Implementation:** Gunakan driver Go **SQLite** (`database/sql` & `go-sqlite3`). Buat modul Go yang menangani inisialisasi koneksi DB dan membuat skema tabel `messages` dan `peers`.
3.  **Data API Implementation:** Definisikan fungsi Go yang menggunakan SQL untuk menyimpan pesan (`SaveMessage`) dan mengambil riwayat pesan (`GetMessageHistory`).
4.  **Background Operation:** Konfigurasi *System Tray* agar aplikasi tetap *listening* di *background*.

Berikan kerangka kode Go untuk inisialisasi Wails, setup SQLite, dan fungsi `SaveMessage`.