package whats

import (
	"go.mau.fi/whatsmeow/binary/proto"
)

func NewProduct(header interface{}, body, footer string, ctx *proto.ContextInfo) (*proto.ProductMessage, error) {
	var message = &proto.ProductMessage{
		Product:          &proto.ProductMessage_ProductSnapshot{},
		BusinessOwnerJid: StringP("owner@jid"),
		Catalog:          &proto.ProductMessage_CatalogSnapshot{},
		Body:             &body,
		Footer:           &footer,
		ContextInfo:      ctx,
	}

	return message, nil
}

func NewProductMessage(header interface{}, body, footer string, ctx *proto.ContextInfo) (*proto.Message, error) {

	var message, err = NewProduct(header, body, footer, ctx)

	return &proto.Message{ProductMessage: message}, err
}
