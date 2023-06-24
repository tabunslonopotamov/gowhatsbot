package whats

import (
	"context"
	"os"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
)

func NewStickerMessage(data []byte, animated bool, ctx *proto.ContextInfo, client *whatsmeow.Client) (*proto.Message, error) {
	var message, err = NewSticker(data, animated, ctx, client)
	return &proto.Message{StickerMessage: message}, err
}

func NewSticker(data []byte, animated bool, ctx *proto.ContextInfo, client *whatsmeow.Client) (*proto.StickerMessage, error) {

	var mimetype = "image/webp"
	if up, err := client.Upload(context.Background(), data, whatsmeow.MediaImage); err != nil {
		return nil, err
	} else {

		var message = &proto.StickerMessage{
			Url:           &up.URL,
			FileSha256:    up.FileSHA256,
			FileEncSha256: up.FileEncSHA256,
			MediaKey:      up.MediaKey,
			Mimetype:      StringP(mimetype),
			DirectPath:    &up.DirectPath,
			FileLength:    Uint64P(uint64(len(data))),
			IsAnimated:    &animated,
			ContextInfo:   ctx,
		}

		return message, err
	}
}

func NewStickerMessageFile(filename string, animated bool, ctx *proto.ContextInfo, client *whatsmeow.Client) (*proto.Message, error) {

	if _, err := os.Stat(filename); err != nil {
		return nil, err
	} else {
		if data, err := os.ReadFile(filename); err != nil {
			return nil, err
		} else {
			return NewStickerMessage(data, animated, ctx, client)
		}
	}
}
