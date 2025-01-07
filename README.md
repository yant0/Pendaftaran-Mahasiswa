# Sistem Manajemen Mahasiswa dan Jurusan

Aplikasi ini adalah program berbasis command-line untuk mengelola data mahasiswa dan jurusan. Mendukung operasi CRUD untuk mahasiswa dan jurusan, serta fitur tambahan seperti pengurutan, penyaringan, impor/ekspor data, dan autentikasi pengguna.

---

## Fitur

1. **Autentikasi**:

    - Role admin dan mahasiswa dengan fungsi berbeda.
    - Kredensial admin default: `admin` / `admin123`.

2. **Manajemen Jurusan**:

    - Menampilkan, menambah, mengedit, dan menghapus jurusan.
    - Impor/ekspor data jurusan dari/ke file CSV.

3. **Manajemen Mahasiswa**:

    - Menampilkan, menambah, mengedit, dan menghapus mahasiswa.
    - Menampilkan mahasiswa berdasarkan jurusan atau diurutkan berdasarkan nilai tes.
    - Impor/ekspor data mahasiswa dari/ke file CSV.

4. **Ekspor Data**:

    - Ekspor semua data (jurusan dan mahasiswa) ke file teks.

5. **Utilitas**:
    - Pembuatan data dummy untuk pengujian.
    - Antarmuka dinamis menggunakan termbox untuk navigasi.

---

## Cara Penggunaan

1. **Menjalankan Program**:

    - Pastikan Go telah terinstal di sistem Anda.
    - Letakkan file data yang diperlukan (`jurusan.csv` dan `mahasiswa.csv`) di direktori yang sama.
    - Jalankan perintah:
        ```bash
        go run tubes.go
        ```

2. **Autentikasi**:

    - Masukkan username dan password untuk login.
    - Admin memiliki akses penuh ke semua fitur; mahasiswa hanya memiliki akses terbatas.

3. **Pilihan Menu Utama**:

    - `Jurusan : Tampilkan` - Menampilkan semua jurusan.
    - `Jurusan : Tambah` - Menambahkan jurusan baru.
    - `Mahasiswa : Tampilkan` - Menampilkan semua mahasiswa.
    - `Mahasiswa : Tambah` - Menambahkan data mahasiswa baru.
    - `Export : Jurusan ke CSV` - Mengekspor data jurusan ke file CSV.
    - `Export : Mahasiswa ke CSV` - Mengekspor data mahasiswa ke file CSV.
    - `Extra : ganti data dengan dummy` - Mengganti data dengan nilai uji.

4. **Pengurutan Data**:
    - Mengurutkan mahasiswa berdasarkan nilai tes secara ascending atau descending.

---

## Deskripsi File

-   **`jurusan.csv`**: Berisi data jurusan.
-   **`mahasiswa.csv`**: Berisi data mahasiswa.
-   **`tubes.go`**: File program utama yang ditulis dalam bahasa Go.

---

## Instalasi

1. Clone repositori:
    ```bash
    git clone <repository_url>
    ```
2. Masuk ke direktori proyek:
    ```bash
    cd nama_direktori_proyek
    ```
3. Instal dependensi:
    - Program menggunakan `termbox-go` dan `github.com/aquasecurity/table`. Instal dengan perintah:
      `bash
go get github.com/nsf/termbox-go
go get github.com/aquasecurity/table
`
      Jika `go.mod` dan `go.sum` ada di dalam direktori, dependensi automatis terdownload dan bisa langsung dijalankan.

---
