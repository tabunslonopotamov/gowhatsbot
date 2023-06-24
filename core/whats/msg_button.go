package whats

import (
	"go.mau.fi/whatsmeow/binary/proto"
)

func NewButtonsMessage(header interface{}, content, footer string, buttons []*proto.ButtonsMessage_Button, ctx *proto.ContextInfo) (*proto.Message, error) {
	var message, err = NewButtons(header, content, footer, buttons, ctx)
	return &proto.Message{ButtonsMessage: message}, err
}

func NewButtons(header interface{}, content, footer string, buttons []*proto.ButtonsMessage_Button, ctx *proto.ContextInfo) (*proto.ButtonsMessage, error) {

	var message = &proto.ButtonsMessage{
		ContentText: &content,
		FooterText:  &footer,
		Buttons:     buttons,
		ContextInfo: ctx,
	}

	switch hd := header.(type) {
	case *proto.ButtonsMessage_DocumentMessage:
		message.HeaderType = proto.ButtonsMessage_DOCUMENT.Enum()
		message.Header = hd

	case *proto.ButtonsMessage_ImageMessage:
		message.HeaderType = proto.ButtonsMessage_IMAGE.Enum()
		message.Header = hd

	case *proto.ButtonsMessage_VideoMessage:
		message.HeaderType = proto.ButtonsMessage_VIDEO.Enum()
		message.Header = hd

	case *proto.ButtonsMessage_Text:
		message.HeaderType = proto.ButtonsMessage_TEXT.Enum()
		message.Header = hd

	case *proto.ButtonsMessage_LocationMessage:
		message.HeaderType = proto.ButtonsMessage_LOCATION.Enum()
		message.Header = hd
	default:
		message.HeaderType = proto.ButtonsMessage_EMPTY.Enum()
		message.Header = nil

	}

	return message, nil

}
