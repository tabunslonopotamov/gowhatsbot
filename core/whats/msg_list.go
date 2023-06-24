package whats

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
)

func NewListMessage(title, description, toggle, footer string, section []*proto.ListMessage_Section, ctx *proto.ContextInfo, client *whatsmeow.Client) (*proto.Message, error) {
	var message, err = NewList(title, description, toggle, footer, section, ctx, client)
	return &proto.Message{ListMessage: message}, err

}

func NewList(title, description, toggle, footer string, section []*proto.ListMessage_Section, ctx *proto.ContextInfo, client *whatsmeow.Client) (*proto.ListMessage, error) {

	var message = &proto.ListMessage{
		Title:       &title,
		Description: &description,
		ButtonText:  &toggle,
		FooterText:  &footer,
		ListType:    proto.ListMessage_SINGLE_SELECT.Enum(),
		Sections:    section,
		ContextInfo: ctx,
	}

	return message, nil
}
