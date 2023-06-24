package main

import (
	_ "main/addons"
	"main/core/addon"
	"main/core/addon/listener"
	"main/core/configs"
	"main/core/mylog"
	"main/core/whats"
)

var MainConfig configs.Config
var ConfigName = "gowhatsbot.json"

func init() {

	if mc, err := configs.Load(ConfigName); err != nil {
		panic(err)
	} else {
		MainConfig = mc
		mylog.Level = mylog.GetLevelBy(MainConfig.Log)

		mylog.Info("Loading config file", ConfigName)
	}
}

func main() {

	mylog.Info("Loading ", MainConfig.Name)
	if gwb, err := whats.NewGoWhatsBot(&MainConfig); err != nil {
		mylog.Error(err)
		panic(err)

	} else {

		mylog.Info("Name :", MainConfig.Name)
		mylog.Info("Driver :", MainConfig.Driver)
		mylog.Info("Address :", MainConfig.Address)
		mylog.Info("Client OS :", *gwb.GetDeviceProps().Os)
		mylog.Info("Client PlatformType :", gwb.GetDeviceProps().PlatformType.String())
		mylog.Info("Client Version :", gwb.GetDeviceProps().Version)

		gwb.AddEventHandler(func(evt interface{}) {
			go addon.Calls(evt, gwb)
			go listener.Call(evt, gwb)
		})

		addon.LogLoaded()
		listener.LogLoaded()

		gwb.Run()
	}

}
