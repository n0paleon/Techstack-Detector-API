# Global Detector FetchPlan

Folder ini berisi **detector global** yang bertugas **mendefinisikan FetchPlan umum dan generik**  
yang dapat digunakan oleh banyak detector lain (webserver, CDN, WAF, framework, dsb).

## Tujuan

`global` **BUKAN** detector teknologi tertentu.  
Folder ini hanya mengatur **apa saja request HTTP dasar** yang sebaiknya selalu dilakukan
untuk mendapatkan sinyal teknologi secara umum.

Contoh use-case:
- Mengambil homepage (`/`)
- Memicu halaman 404 untuk melihat error page default
- Request path acak untuk melihat fallback behavior server
- Request ringan yang aman dan low-cost

---

## Aturan Penting (WAJIB)

### ✅ Yang BOLEH diedit
Hanya bagian **FetchPlan**.

FetchPlan di folder ini menentukan:
- Path yang akan di-fetch
- Method HTTP
- Header default
- Tujuan fetch (homepage, error page, dll)

### ❌ Yang TIDAK BOLEH diedit
- Logic deteksi (`Detect`)
- Parsing header / body
- Penentuan teknologi
- Skoring / confidence
- Side-effect apa pun

Folder ini **tidak boleh**:
- Mengembalikan `Technology`
- Mengandung logika fingerprint
- Bergantung pada detector lain

---

## Prinsip Desain

- **Generic**  
  Tidak spesifik ke teknologi tertentu (nginx, apache, litespeed, dll)

- **Reusable**  
  FetchPlan di sini dipakai lintas detector

- **Low Risk & Low Noise**  
  Tidak agresif, tidak menyerupai scanning, aman untuk target publik

- **Deterministic**  
  Request yang konsisten agar hasil bisa dibandingkan antar target

---

## Contoh FetchPlan yang Valid

- `GET /` → homepage
- `GET /<random-string>` → memicu 404
- `HEAD /` → header-only inspection

---

## Contoh FetchPlan yang TIDAK Valid

❌ Request ke:
- `/wp-admin`
- `/admin`
- `/phpmyadmin`
- Path eksploit / probing spesifik

❌ Conditional fetch berbasis teknologi tertentu

---

## Catatan Arsitektur

- Detector spesifik **HARUS** mendefinisikan logic deteksinya sendiri
- `global` hanya bertindak sebagai **shared fetch orchestration**
- Jika butuh fetch tambahan yang spesifik teknologi → buat FetchPlan di detector masing-masing

---

## Ringkasannya

> **Global detector = fetch WHAT, bukan detect WHAT**

Jika ragu apakah sebuah perubahan layak masuk folder ini, kemungkinan besar **tidak**.

