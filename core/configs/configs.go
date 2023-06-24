package configs

import (
	"encoding/json"
	"io/fs"
	"os"
)

type Config struct {
	Name      string `json:"name"`
	Driver    string `json:"driver"`
	Address   string `json:"address"`
	ShowQR    bool   `json:"showqr"`
	OS        string `json:"os"`
	Platform  int32  `json:"platform"`
	Log       string `json:"log"`
	ClientLog string `json:"client_log"`
}

func Load(p string) (Config, error) {
	var c Config
	if config_bytes, err := os.ReadFile(p); err != nil {
		return c, err
	} else {
		if err := json.Unmarshal(config_bytes, &c); err != nil {
			return c, err
		} else {
			return c, err
		}
	}
}

func Save(c Config, p string) error {
	if b, err := json.MarshalIndent(c, "", "  "); err != nil {
		return err
	} else {
		return os.WriteFile(p, b, fs.ModeAppend)
	}
}
