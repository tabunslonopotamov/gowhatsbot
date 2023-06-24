package whats

import (
	"context"
	"main/core/media"
	"net/http"
	"os"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
)

func NewVideoMessage(data []byte, caption string, ctx *proto.ContextInfo, client *whatsmeow.Client) (*proto.Message, error) {
	var message, err = NewVideo(data, caption, ctx, client)
	return &proto.Message{VideoMessage: message}, err
}

func NewVideo(data []byte, caption string, ctx *proto.ContextInfo, client *whatsmeow.Client) (*proto.VideoMessage, error) {
	var mimetype = http.DetectContentType(data)
	var info media.MediaInfo
	if strings.HasPrefix(mimetype, "image/gif") {
		if md, err := media.GifToMP4(data); err != nil {
			return nil, err
		} else {
			data = info.Bytes
			info = md
		}
	}

	if up, err := client.Upload(context.Background(), data, whatsmeow.MediaVideo); err != nil {
		return nil, err
	} else {

		mimetype = http.DetectContentType(data)

		message := &proto.VideoMessage{
			Caption:       &caption,
			Url:           &up.URL,
			Mimetype:      &mimetype,
			FileSha256:    up.FileSHA256,
			FileLength:    &up.FileLength,
			MediaKey:      up.MediaKey,
			Height:        Uint32P(uint32(info.Height)),
			Width:         Uint32P(uint32(info.Width)),
			FileEncSha256: up.FileEncSHA256,
			DirectPath:    &up.DirectPath,
			JpegThumbnail: info.JpegThumbail,
			Seconds:       Uint32P(uint32(info.Seconds)),
			ContextInfo:   ctx,
		}

		if strings.HasPrefix(mimetype, "image/gif") {
			message.GifAttribution = proto.VideoMessage_NONE.Enum()
			message.GifPlayback = BoolP(true)
		}

		if thumb_up, err := client.Upload(context.Background(), info.JpegThumbail, whatsmeow.MediaImage); err != nil {
			return message, err
		} else {
			message.ThumbnailDirectPath = &thumb_up.DirectPath
			message.ThumbnailSha256 = thumb_up.FileSHA256
			message.ThumbnailEncSha256 = thumb_up.FileEncSHA256
		}

		return message, err
	}
}

func NewVideoMessageFile(filename, caption string, ctx *proto.ContextInfo, client *whatsmeow.Client) (*proto.Message, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	} else {
		if data, err := os.ReadFile(filename); err != nil {
			return nil, err
		} else {
			return NewVideoMessage(data, caption, ctx, client)
		}
	}
}
