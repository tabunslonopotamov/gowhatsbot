<div align="center">
    <img src="gowhatsbot.png" width="40%" /><br>
    <h1>GoWhatsBot</h1>
    <br>
</div>

# GoWhatsBot : Apa itu ?
GoWhatsBot adalah Bot WhatsApp yang dibangun dengan Go-lang bebasiskan library [` whatsmeow `](https://github.com/tulir/whatsmeow).


# Quick Setup

## Konfigurasi
Untuk menjalankan Bot, kita perlu untuk mengatur konfigurasi database pada berkas ` gowhatsbot.json `. Jika tidak terdapat pada direktori repo, maka kita bisa membuatnya dengan contoh isi konfigurasi sebagai berikut :
``` json
{
    "driver": "sqlite3",
    "address": "file:soursop.db?_foreign_keys=on",
    "name": "JustSoursop",
    "os": "MacOS",
    "platform": 7,
    "log": "debug",
    "client_log": "error",
    "showqr": true
}
```
Pada contoh diatas, driver yang akan di gunakan adalah ` sqlite3 ` dengan alamat ` file:godev.db?__foreign_keys=on `.

Secara default ada 2 library driver database yang tesedia yaitu ` pgx ` untuk database PostgreSQL dan ` go-sqlite3 ` untuk SQLite.

_(perlu untuk menambahkan baris kode jika ingin menambahkan dukungan layanan database lainnya)_

## Autoload

Buatlah berkas ` ./addons/loader.go ` jika belum ada didalam direktori `./addons` dan tambahkan `import` paket yang ingin di `build` bersama aplikasi sebagai `addons`:

```go
package addons

import (
	"main/core/addon"
	"main/core/listener"
	"main/core/mylog"

	// Ini adalah contoh standar mengimport paket addons
	_ "main/addons/help"
)

var Call = addon.Calls

func init() {

	// Menampilkan daftar addons yang di load beserta fiturnya
	for _, a := range addon.AddOnList {
		mylog.Info("Add-on", a.Name, "loaded with", len(a.Features), "feature(s).")

	}

	// Menampilkan daftar listener yang diregistrasi
	for k, v := range listener.EventListener {
		mylog.Info(len(v), k, "has registered.")

	}
}
```

## Menjalankan & Kompilasi

### Menjalankan
Untuk menjalankan bot tanpa kompilasi cukup untuk menjalankan perintah :
```sh
CGO_ENABLED=1 go run .
```

### Kompilasi
Untuk kompilasi agar dapat mendukung driver database maupun library yang menggunakan sumber program dari bahasa C dengan menambahkan parameter `CGO_ENABLED=1` :
```sh
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -v -o ./main .

```

# Kontribusi ?
Silahkan jika ingin melakukan kontribusi dengan membuka issue, pull request maupun diskusi.

# Library ?
- [whatsmeow](https://go.mau.fi/whatsmeow)
- [qrterminal](https://github.com/mdp/qrterminal)
- [pgx](https://github.com/jackc/pgx)
- [go-sqlite3](https://github.com/mattn/go-sqlite3)
- [go-qrcode](https://github.com/skip2/go-qrcode)


# Terimakasih (Thanks)
Terimakasih kepada semua orang yang berkontribusi dalam projek ini.