package addon

import (
	"errors"
	"main/core/whats"

	"go.mau.fi/whatsmeow/types/events"
)

func FeaturePatternValid(feature *Feature, i interface{}, gwb *whats.GoWhatsBot) error {
	_, _, pattern, _, err := whats.PrepareExec(i)

	if err != nil {
		return err
	}

	if feature.IsMy(pattern) {
		return nil
	}

	return errors.New("feature not invalid")
}

func AddonValidFromMe(ao *AddOn, i interface{}, gwb *whats.GoWhatsBot) error {

	if eMsg, ok := i.(*events.Message); ok {
		if eMsg.Info.IsFromMe {
			return nil
		}
	}

	return errors.New("user not permitted")
}
