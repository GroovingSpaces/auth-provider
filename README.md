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

	// Gunakan authprovider.VerifyToken, VerifyTokenWithMiddleware, GetCurrentUser, GetRoles, CreateUser, UpdateUser, DeleteUser di sini...
}
```

### Variabel lingkungan (opsional)

| Variabel              | Deskripsi                    | Contoh              |
|-----------------------|-----------------------------|---------------------|
| `AUTH_SERVICE_HOST`   | Base URL Auth Service       | `https://auth.example.com` |

## Penggunaan

Semua fungsi (kecuali `Init`) membutuhkan **Bearer token** (JWT) di header `Authorization`. Pastikan token valid dari Auth Service.

---

### Init

Menginisialisasi client dan base URL Auth Service. **Wajib dipanggil sekali** sebelum memakai fungsi lain.

| Parameter | Tipe   | Keterangan                                      |
|-----------|--------|--------------------------------------------------|
| `host`    | string | Base URL Auth Service (tanpa trailing slash)     |

```go
authprovider.Init("https://auth.example.com")
```

---

### VerifyToken

Memverifikasi token JWT dan mengembalikan data user serta claims jika valid.

**Signature:** `VerifyToken(token string) (dto.VerifyTokenResponse, error)`

| Parameter | Tipe   | Keterangan   |
|-----------|--------|---------------|
| `token`   | string | JWT Bearer    |

```go
resp, err := authprovider.VerifyToken(token)
if err != nil {
	return
}
if resp.Status == "OK" && resp.Data.Valid {
	user := resp.Data.User
	claims := resp.Data.Claims
}
```

**Response:** `dto.VerifyTokenResponse` — `Data.Valid`, `Data.User`, `Data.Claims`, `RequestAPICallResult`.

---

### VerifyTokenWithMiddleware

Memverifikasi token dan memastikan user punya permission tertentu (role aktif, permission aktif). Cocok untuk middleware HTTP yang membatasi akses per menu/aksi.

**Signature:** `VerifyTokenWithMiddleware(token string, permissionName string) (dto.VerifyTokenData, error)`

| Parameter        | Tipe   | Keterangan              |
|------------------|--------|--------------------------|
| `token`          | string | JWT Bearer               |
| `permissionName`  | string | Slug permission (contoh: `"users.create"`) |

```go
data, err := authprovider.VerifyTokenWithMiddleware(token, "users.create")
if err != nil {
	// err: dto.ErrInvalidToken, dto.ErrRoleInactive, dto.ErrPermissionInactive, dto.ErrRoleForbidden
	return
}
// data = VerifyTokenData (User, Claims, Token)
```

**Response:** `dto.VerifyTokenData` — `User`, `Claims`, `Token` jika akses boleh.

---

### GetCurrentUser

Mengambil data user yang sedang login berdasarkan token.

**Signature:** `GetCurrentUser(token string) (dto.GetCurrentUserResponse, error)`

| Parameter | Tipe   | Keterangan |
|-----------|--------|------------|
| `token`   | string | JWT Bearer |

```go
resp, err := authprovider.GetCurrentUser(token)
if err != nil {
	return
}
if resp.Status == "OK" {
	user := resp.Data // ID, Email, Username
}
```

**Response:** `dto.GetCurrentUserResponse` — `Data` berisi `ID`, `Email`, `Username`.

---

### GetRoles

Mengambil daftar roles dari Auth Service (mapping role_name ke role_id, dll.).

**Signature:** `GetRoles(token string) (dto.GetRolesResponse, error)`

| Parameter | Tipe   | Keterangan |
|-----------|--------|------------|
| `token`   | string | JWT Bearer |

```go
resp, err := authprovider.GetRoles(token)
if err != nil {
	return
}
if resp.Status == "OK" {
	roles := resp.Data // []RoleData
}
```

**Response:** `dto.GetRolesResponse` — `Data` berisi slice `RoleData` (ID, Name, Description, IsActive, Permissions).

---

### CreateUser

Membuat user baru di Auth Service.

**Signature:** `CreateUser(token string, payload dto.CreateUserRequest) (dto.CreateUserResponse, error)`

| Parameter | Tipe                   | Keterangan        |
|-----------|------------------------|-------------------|
| `token`   | string                 | JWT Bearer        |
| `payload` | dto.CreateUserRequest | Email, Password, Name, RoleIds |

```go
payload := dto.CreateUserRequest{
	Email:    "user@example.com",
	Password: "secret123",
	Name:     "John Doe",
	RoleIds:  []string{"role-uuid-1"},
}
resp, err := authprovider.CreateUser(token, payload)
if err != nil {
	return
}
if resp.Status == "OK" {
	user := resp.Data
}
```

**Request:** `dto.CreateUserRequest` — `Email` (required, email), `Password` (required, min 8), `Name` (required, min 3), `RoleIds` (required, min 1).

**Response:** `dto.CreateUserResponse` — `Data` berisi `UserData`.

---

### UpdateUser

Memperbarui data user berdasarkan ID.

**Signature:** `UpdateUser(token string, id string, payload dto.UpdateUserRequest) (dto.UpdateUserResponse, error)`

| Parameter | Tipe                    | Keterangan                 |
|-----------|-------------------------|----------------------------|
| `token`   | string                  | JWT Bearer                 |
| `id`      | string                  | ID user yang akan diupdate |
| `payload` | dto.UpdateUserRequest   | Email, Password, Name, RoleIds, IsActive |

```go
payload := dto.UpdateUserRequest{
	Email:    "new@example.com",
	Password: "newpass123",
	Name:     "Jane Doe",
	RoleIds:  []string{"role-uuid-1"},
	IsActive: true,
}
resp, err := authprovider.UpdateUser(token, userID, payload)
if err != nil {
	return
}
if resp.Status == "OK" {
	user := resp.Data
}
```

**Request:** `dto.UpdateUserRequest` — sama seperti CreateUser + `IsActive` (required).

**Response:** `dto.UpdateUserResponse` — `Data` berisi `UserData`.

---

### DeleteUser

Menghapus user berdasarkan ID.

**Signature:** `DeleteUser(token string, id string) (dto.DeleteUserResponse, error)`

| Parameter | Tipe   | Keterangan                |
|-----------|--------|---------------------------|
| `token`   | string | JWT Bearer                |
| `id`      | string | ID user yang akan dihapus |

```go
resp, err := authprovider.DeleteUser(token, userID)
if err != nil {
	return
}
if resp.Status == "OK" {
	// User berhasil dihapus
}
```

**Response:** `dto.DeleteUserResponse` — `Status`, `Message`, `RequestAPICallResult`.

## DTO (Data Transfer Object)

Struktur request/response ada di package `dto`:

| Tipe | Keterangan |
|------|-------------|
| `VerifyTokenResponse`, `VerifyTokenData` | Verifikasi token |
| `GetCurrentUserResponse`, `GetCurrentUserData` | Data user saat ini |
| `GetRolesResponse`, `RoleData` | Daftar roles |
| `CreateUserRequest`, `CreateUserResponse` | Buat user |
| `UpdateUserRequest`, `UpdateUserResponse` | Update user |
| `DeleteUserResponse` | Hapus user |
| `UserData`, `AuthUserData`, `AuthClaims`, `ClaimsResponse` | Data user/claims |

Semua response yang memanggil API menyertakan `RequestAPICallResult` (URL, method, headers, body, status code, latency) untuk logging/debugging.

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

| Fungsi                       | Method   | Path                    |
|-----------------------------|----------|--------------------------|
| `VerifyToken`               | POST     | `/api/v1/auth/verify-token` |
| `VerifyTokenWithMiddleware` | (internal: memakai `VerifyToken`) | — |
| `GetCurrentUser`            | GET      | `/api/v1/auth/me`        |
| `GetRoles`                  | GET      | `/api/v1/roles`          |
| `CreateUser`                | POST     | `/api/v1/users`          |
| `UpdateUser`                | PUT      | `/api/v1/users/{id}`     |
| `DeleteUser`                | DELETE   | `/api/v1/users/{id}`     |

Host/base URL di-set melalui `authprovider.Init(host)`.

## Lisensi

Digunakan dalam proyek GroovingSpaces.
