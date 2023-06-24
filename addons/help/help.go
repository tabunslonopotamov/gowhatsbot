package help

import (
	"fmt"
	"main/core/addon"
	"main/core/mylog"
	"main/core/whats"
	"strings"

	"go.mau.fi/whatsmeow/types/events"
)

// ItAddOn is a plugin that contain a help command
var ItAddOn = addon.NewAddOnRegistered("Helper", "Everything that help.")

func init() {

	ItAddOn.Validator = addon.AddonValidFromMe

	ItAddOn.AddFeatures(&addon.Feature{
		Name:      "Help",
		Patterns:  []string{".h", ".help"},
		Execute:   executeHelp,
		Expecting: &events.Message{},
		Validator: addon.FeaturePatternValid,
	})
}

func executeHelp(a *addon.AddOn, f *addon.Feature, e interface{}, gwb *whats.GoWhatsBot) error {

	var eventMsg, _, pattern, args, err = whats.PrepareExec(e)
	var newctx = whats.NewContext(eventMsg, whats.CTX_GROUP_QUOTE)

	if err != nil {
		mylog.Error(f.Name, pattern, err)
		return err
	}

	var list_text = []string{}

	for _, a := range addon.AddOnList {
		list_text = append(list_text, fmt.Sprintf("*# %s* _(%s)_", a.Name, a.Description))

		for _, f := range a.Features {
			list_text = append(list_text, fmt.Sprintf(" %s > %s", f.Name, strings.Join(f.Patterns, ", ")))

		}
		list_text = append(list_text, "")
	}

	if len(list_text) > 0 {

		var text = strings.Join(list_text, "\n")
		if resp, err := gwb.SendText(eventMsg.Info.Chat, text, newctx); err != nil {
			mylog.Error(f.Name, pattern, err)
			return err
		} else {
			mylog.Info(a.Name, f.Name, resp.ID, pattern, args, err)
		}
	}

	return nil
}
