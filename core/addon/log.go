package addon

import "main/core/mylog"

func LogLoaded() {
	// Menampilkan daftar addons yang di load beserta fiturnya
	for _, a := range AddOnList {
		mylog.Info("Add-on", a.Name, "loaded with", len(a.Features), "feature(s).")

	}
}
