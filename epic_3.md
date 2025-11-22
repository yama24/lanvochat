### EPIC 3: FRONTEND & UI/UX (Cantik & Modern)

**Fokus Task:** Mengembangkan antarmuka pengguna React yang modern, menghubungkannya ke backend Go, dan mengimplementasikan *popup* notifikasi.

**PRINSIP UTAMA:** Cantik & Modern (React/TS) dan *Frameless Window*.

| Stack Kunci | Versi/Spesifikasi |
| :--- | :--- |
| **Go Backend** | Go 1.25.3 |
| **Frontend** | React/TypeScript |
| **JS Package Manager** | Yarn (Node.js 22.21.0) |

**Task Rinci:**
1.  **Frontend Setup & Styling:** Setup *styling* (misalnya Tailwind CSS) dan buat komponen React dasar yang modern.
2.  **Peer List Component:** Buat komponen React yang mendengarkan Wails Events dari UDP untuk menampilkan dan memperbarui daftar *peer online*.
3.  **Chat Window Component & History:** Buat komponen yang memanggil Go API untuk memuat riwayat dari **SQLite** dan menampilkan *message bubbles*.
4.  **Chat Popup Window:** Di Go backend, buat fungsi yang menggunakan `wails.Runtime.Window.NewBrowserWindow` untuk membuat **jendela kecil baru** (**tanpa bingkai**) yang diposisikan di **sudut kanan bawah layar** saat pesan masuk.