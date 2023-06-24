package listener

import (
	"main/core/mylog"
)

func LogLoaded() {
	// Menampilkan daftar listener yang diregistrasi
	for k, v := range EventListener {
		mylog.Info(len(v), k, "has registered.")
	}
}
