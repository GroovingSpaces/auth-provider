# Auth Provider

Modul Go untuk berkomunikasi dengan Auth Service (Grooving). Menyediakan fungsi verifikasi token, mengambil data user, dan mengelola roles melalui HTTP client dengan timeout 10 detik.

**Versi:** [v1.0.2](https://github.com/GroovingSpaces/auth-provider/releases/tag/v1.0.2)

## Persyaratan

- Go 1.25 atau lebih baru
- Auth Service yang dapat diakses (base URL)

## Instalasi

```bash
# Versi terbaru
go get github.com/GroovingSpaces/auth-provider

# Versi tertentu (contoh: v1.0.2)
go get github.com/GroovingSpaces/auth-provider@v1.0.2
```

## Konfigurasi

Panggil `authprovider.Init(host)` sekali saat aplikasi mulai (misalnya di `main` atau saat setup). Parameter `host` adalah base URL Auth Service (tanpa trailing slash).

```go
package main

import (
	"os"

	"github.com/GroovingSpaces/auth-provider"
)

func main() {
	host := os.Getenv("AUTH_SERVICE_HOST")
	if host == "" {
		host = "http://localhost:8080"
	}
	authprovider.Init(host)

	// Gunakan authprovider.VerifyToken, authprovider.VerifyTokenWithMiddleware, authprovider.GetCurrentUser, authprovider.GetRoles di sini...
}
```

### Variabel lingkungan (opsional)

| Variabel              | Deskripsi                    | Contoh              |
|-----------------------|-----------------------------|---------------------|
| `AUTH_SERVICE_HOST`   | Base URL Auth Service       | `https://auth.example.com` |

## Penggunaan

Semua fungsi membutuhkan **Bearer token** (JWT) di header `Authorization`. Pastikan token valid dari Auth Service.

### 1. Verifikasi Token

Memverifikasi token JWT dan mengembalikan data user serta claims jika valid.

```go
resp, err := authprovider.VerifyToken(token)
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

### 2. Verifikasi Token dengan Middleware (Permission)

Memverifikasi token dan memastikan user punya permission tertentu (role aktif, permission aktif). Cocok untuk middleware HTTP yang membatasi akses per menu/aksi.

```go
data, err := authprovider.VerifyTokenWithMiddleware(token, "users.create")
if err != nil {
	// err bisa dto.ErrInvalidToken, dto.ErrRoleInactive, dto.ErrPermissionInactive, dto.ErrRoleForbidden
	return
}
// data berisi VerifyTokenData (User, Claims, Token) — akses boleh
```

**Parameter:** `token` (JWT), `permissionName` (slug permission, misalnya `"users.create"`).

**Response:** `dto.VerifyTokenData` — berisi `User`, `Claims`, `Token` jika akses boleh.

**Error:** `dto.ErrInvalidToken`, `dto.ErrRoleInactive`, `dto.ErrPermissionInactive`, `dto.ErrRoleForbidden`.

### 3. Get Current User

Mengambil data user yang sedang login berdasarkan token.

```go
resp, err := authprovider.GetCurrentUser(token)
if err != nil {
	return
}
if resp.Status == "OK" {
	user := resp.Data
	// user.ID, user.Email, user.Username
}
```

**Response:** `dto.GetCurrentUserResponse` — berisi `Data` (ID, Email, Username).

### 4. Get Roles

Mengambil daftar roles dari Auth Service (untuk mapping role_name ke role_id, dll.).

```go
resp, err := authprovider.GetRoles(token)
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
// dto.ErrRoleInactive
// dto.ErrPermissionInactive
// dto.ErrRoleForbidden
// dto.ErrTimeoutError
```

Contoh penanganan timeout:

```go
resp, err := authprovider.VerifyToken(token)
if err != nil {
	if errors.Is(err, dto.ErrTimeoutError) {
		// Timeout ke Auth Service
	}
	return
}
```

## Endpoint Auth Service yang dipanggil

| Fungsi                       | Method | Path                          |
|-----------------------------|--------|-------------------------------|
| `VerifyToken`               | POST   | `/api/v1/auth/verify-token`   |
| `VerifyTokenWithMiddleware` | (memakai `VerifyToken`) | — |
| `GetCurrentUser`            | GET    | `/api/v1/auth/me`             |
| `GetRoles`                  | GET    | `/api/v1/roles`               |

Host/base URL di-set melalui `authprovider.Init(host)`.

## Lisensi

Digunakan dalam proyek GroovingSpaces.
