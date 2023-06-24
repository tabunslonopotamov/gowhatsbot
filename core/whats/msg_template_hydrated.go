package whats

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
)

func NewHydratedTemplateMessage(title interface{}, content, footer, template_id string, buttons []*proto.HydratedTemplateButton, ctx *proto.ContextInfo, client *whatsmeow.Client) (*proto.Message, error) {

	if message, err := NewHydratedTemplate(title, content, footer, template_id, buttons, ctx, client); err != nil {
		return nil, err
	} else {
		return &proto.Message{TemplateMessage: message}, err
	}

}

func NewHydratedTemplate(title interface{}, content, footer, template_id string, buttons []*proto.HydratedTemplateButton, ctx *proto.ContextInfo, client *whatsmeow.Client) (*proto.TemplateMessage, error) {

	var message = &proto.TemplateMessage{
		HydratedTemplate: &proto.TemplateMessage_HydratedFourRowTemplate{
			HydratedContentText: &content,
			HydratedFooterText:  &footer,
			HydratedButtons:     buttons,
			TemplateId:          &template_id,
			Title:               nil,
		},
		// ContextInfo: ctx,
	}

	switch the_title := title.(type) {
	case *proto.TemplateMessage_HydratedFourRowTemplate_HydratedTitleText:
		message.HydratedTemplate.Title = the_title

	case *proto.TemplateMessage_HydratedFourRowTemplate_ImageMessage:
		message.HydratedTemplate.Title = the_title

	case *proto.TemplateMessage_HydratedFourRowTemplate_DocumentMessage:
		message.HydratedTemplate.Title = the_title

	case *proto.TemplateMessage_HydratedFourRowTemplate_VideoMessage:
		message.HydratedTemplate.Title = the_title

	case *proto.TemplateMessage_HydratedFourRowTemplate_LocationMessage:
		message.HydratedTemplate.Title = the_title

	}

	return message, nil
}
