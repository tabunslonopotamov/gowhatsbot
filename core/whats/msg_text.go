package whats

import (
	_ "image/jpeg"
	_ "image/png"

	"go.mau.fi/whatsmeow/binary/proto"
)

func NewConversation(text string) (*proto.Message, error) {

	return &proto.Message{Conversation: &text}, nil
}

func NewExtendedMessage(text string, ctx *proto.ContextInfo) (*proto.Message, error) {

	var extendedtext_message = &proto.ExtendedTextMessage{
		Text:        &text,
		ContextInfo: ctx,
	}

	return &proto.Message{ExtendedTextMessage: extendedtext_message}, nil
}
