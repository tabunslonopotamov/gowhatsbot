## GoWhatsBot `AddOns`, `Listener` dan cara menambahkannya kedalam `./addons/loader.go`

### AddOns

#### 1. Bagaimana cara membuat `AddOn` ?

Buat `package` direktori `ping` dalam `addons`, lalu buat sebuah berkas `ping.go`:

```
addons/
└── ping/
    └── ping.go
```

<details>
    <summary>Berikut contoh isi  berkas ping.go :</summary>


```go
package ping

import (
    "fmt"
    "main/core/addon"
    "main/core/mylog"
    "main/core/whats"
    "time"

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

// Menambahkan AddOn serta mendaftarkannya ke daftar eksekusi
var ItAddOn.NewAddOnRegistered("Pinger", "Melakukan ping")

func init() {
    // Menggunakan validator bawaan agar hanya nomor user yang valid
    ItAddOn.Validator = addon.AddonValidFromMe

    // Menambahkan fitur beserta perintah eksekusinya
    ItAddOn.AddFeatures(&addon.Feature{
        Patterns:    []string{".ping", ".p"},
        Description: "Get test ping time.",
        Usage:       "{pattern}",
        Execute:     executePing,

        // Menggunakan validator bawaan untuk mencocokkan patern
        Validator:   addon.FeaturePatternValid, 
        
        // Disini si fitur hanya mengharapkan event 
        Expecting:   &events.Message{}, 
        
    })

}

func executePing(a *addon.AddOn, f *addon.Feature, e interface{}, c *whatsmeow.Client) error {
    
    // Mempersiapkan event, pattern, args (teks setelah kata perintah)
    var eventMsg, _, pattern, args, err = whats.PrepareExec(e)
    var newctx = whats.NewContext(eventMsg, true, true, c)

    if err != nil {
        mylog.Error(a.Name, f.Name, pattern, args, err)
        return err
    }

    var text_result = fmt.Sprint(time.Since(eventMsg.Info.Timestamp).Milliseconds(), "ms")
    if resp, err := whats.SendTextMessage(eventMsg.Info.Chat, text_result, newctx, c); err != nil {
        return err
    } else {
        mylog.Info(pattern, resp.ID)
    }

    return nil
}
```

</details>
<br>

Setelah selesai membuat berkas `ping.go`, maka paket tersebut harus didaftarkan pada import di `./addons/loader.go` :
 _(jika belum ada, pembuatannya ada pada README.md utama)_
```go
package addons

import (
 ...

    _ "main/addons/ping"
)

 ...
```


#### 2. Bagaimana cara membuat `Listener` ?
Sama halnya dengan cara pembuatan `AddOn` diatas, mulai dari membuat direktori hingga berkas `ping.go`. Hanya saja ada sedikit perbedaan cara mendaftarkan listenernya, berikut contohnya :

<details>
    <summary>Contoh Listener</summary>

```go
package ping

import (
    "main/core/listener"
    "main/core/mylog"
    "main/core/whats"

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

func init() {

    // Mendaftarkan untuk memonitor event pesan masuk
    listener.AddEvent("*events.Message", executePing)
}

func executePing(e interface{}, c *whatsmeow.Client) error {
    

    // Mempersiapkan event, pattern, args (teks setelah kata perintah)
    var eventMsg, _, pattern, args, err = whats.PrepareExec(e)
    var newctx = whats.NewContext(eventMsg, true, true, c)

    if err != nil {
        mylog.Error(a.Name, f.Name, pattern, args, err)
        return err
    }

    // Melakukan cek manual pada pattern
    if pattern == ".p" || pattern == ".ping" {
        var text_result = fmt.Sprint(time.Since(eventMsg.Info.Timestamp).Milliseconds(), "ms")
        if resp, err := whats.SendTextMessage(eventMsg.Info.Chat, text_result, newctx, c); err != nil {
            return err
        } else {
            mylog.Info(pattern, resp.ID)
        }        
    }

    return nil
}
```

</details>
