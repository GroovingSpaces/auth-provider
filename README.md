# Auth Provider

Modul Go untuk berkomunikasi dengan Auth Service (Grooving). Menyediakan fungsi verifikasi token, mengambil data user, dan mengelola roles melalui HTTP client dengan timeout 10 detik.

## Persyaratan

- Go 1.25 atau lebih baru
- Auth Service yang dapat diakses (base URL)

## Instalasi

```bash
go get github.com/GroovingSpaces/auth-provider
```

## Konfigurasi

Panggil `Init(host)` sekali saat aplikasi mulai (misalnya di `main` atau saat setup). Parameter `host` adalah base URL Auth Service (tanpa trailing slash).

```go
package main

import (
	"os"
)

func main() {
	host := os.Getenv("AUTH_SERVICE_HOST")
	if host == "" {
		host = "http://localhost:8080"
	}
	Init(host)

	// Gunakan VerifyToken, GetCurrentUser, GetRoles di sini...
}
```

> **Catatan:** Kode provider berada di `package main`. Jika modul ini dipakai sebagai dependency dari proyek lain, gunakan `replace` di `go.mod` atau pindahkan kode ke subpackage (misalnya `provider`) agar bisa di-import.

### Variabel lingkungan (opsional)

| Variabel              | Deskripsi                    | Contoh              |
|-----------------------|-----------------------------|---------------------|
| `AUTH_SERVICE_HOST`   | Base URL Auth Service       | `https://auth.example.com` |

## Penggunaan

Semua fungsi membutuhkan **Bearer token** (JWT) di header `Authorization`. Pastikan token valid dari Auth Service.

### 1. Verifikasi Token

Memverifikasi token JWT dan mengembalikan data user serta claims jika valid.

```go
resp, err := VerifyToken(token)
if err != nil {
	// Handle error (invalid token, timeout, dll.)
	return
}
if resp.Status == "OK" && resp.Data.Valid {
	user := resp.Data.User
	claims := resp.Data.Claims
	// Gunakan user / claims
}
```

**Response:** `dto.VerifyTokenResponse` — berisi `Data.Valid`, `Data.User`, `Data.Claims`, dan `RequestAPICallResult` untuk debugging.

### 2. Get Current User

Mengambil data user yang sedang login berdasarkan token.

```go
resp, err := GetCurrentUser(token)
if err != nil {
	return
}
if resp.Status == "OK" {
	user := resp.Data
	// user.ID, user.Email, user.Username
}
```

**Response:** `dto.GetCurrentUserResponse` — berisi `Data` (ID, Email, Username).

### 3. Get Roles

Mengambil daftar roles dari Auth Service (untuk mapping role_name ke role_id, dll.).

```go
resp, err := GetRoles(token)
if err != nil {
	return
}
if resp.Status == "OK" {
	roles := resp.Data
	// Proses daftar roles
}
```

**Response:** `dto.GetRolesResponse` — berisi data roles dan metadata response.

## DTO (Data Transfer Object)

Struktur request/response ada di package `dto`:

- **Auth:** `VerifyTokenResponse`, `VerifyTokenData`, `GetCurrentUserResponse`, `UserData`, `ClaimsResponse`
- **Request:** `CreateUserRequest`, `UpdateUserRequest` (jika nanti dipakai oleh handler lain)
- **Response:** berbagai tipe `*Response` yang mengikuti format API Auth Service

Semua response yang memanggil API menyertakan `RequestAPICallResult` (URL, method, headers, body, status code, latency) untuk keperluan logging atau debugging.

## Penanganan error

Modul mengembalikan error standar Go. Konstanta error umum ada di `dto`:

```go
import "github.com/GroovingSpaces/auth-provider/dto"

// dto.ErrInvalidToken
// dto.ErrUserInactive
// dto.ErrRoleForbidden
// dto.ErrTimeoutError
```

Contoh penanganan timeout:

```go
resp, err := VerifyToken(token)
if err != nil {
	if errors.Is(err, dto.ErrTimeoutError) {
		// Timeout ke Auth Service
	}
	return
}
```

## Endpoint Auth Service yang dipanggil

| Fungsi           | Method | Path                          |
|------------------|--------|-------------------------------|
| `VerifyToken`    | POST   | `/api/v1/auth/verify-token`   |
| `GetCurrentUser` | GET    | `/api/v1/auth/me`            |
| `GetRoles`       | GET    | `/api/v1/roles`              |

Host/base URL di-set melalui `Init(host)`.

## Lisensi

Digunakan dalam proyek GroovingSpaces.
