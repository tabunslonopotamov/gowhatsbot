package utils

import (
	"fmt"
	"math"
	"os/exec"
	"path"
	"runtime"
)

func GetMyPath() string {
	if _, filename, _, ok := runtime.Caller(1); ok {
		return path.Dir(filename)
	}

	return ""
}

func RunCmd(cmdline string, combine bool) ([]byte, error) {
	var command = exec.Command("/bin/sh", "-c", cmdline)

	if combine {
		return command.CombinedOutput()
	} else {
		return command.Output()
	}

}

const KILO float64 = 1024

func FormatBytes(bytes float64, decimals int) string {
	var sizes = []string{"Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	if decimals < 0 {
		decimals = 0
	}

	var i = math.Floor(math.Log(bytes) / math.Log(KILO))
	var tmpl = fmt.Sprintf("%s0.%df %s", "%", decimals, "%s")
	return fmt.Sprintf(tmpl, bytes/math.Pow(KILO, i), sizes[int(i)])
}
