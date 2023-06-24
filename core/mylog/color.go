package mylog

import (
	"github.com/fatih/color"
)

func RedString(s string) string {
	return ColoringString(color.RedString, s)
}

func BlueString(s string) string {
	return ColoringString(color.BlueString, s)
}

func CyanString(s string) string {
	return ColoringString(color.CyanString, s)
}

func YellowString(s string) string {
	return ColoringString(color.YellowString, s)
}

func WihteString(s string) string {
	return ColoringString(color.WhiteString, s)
}

func BlackString(s string) string {
	return ColoringString(color.BlackString, s)
}

func MagentaString(s string) string {
	return ColoringString(color.MagentaString, s)
}

func GreenString(s string) string {
	return ColoringString(color.GreenString, s)
}

type colorF func(string, ...interface{}) string

func ColoringString(f colorF, s ...interface{}) string {
	return f("%s", s...)
}
