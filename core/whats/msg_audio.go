package whats

import (
	"context"
	"fmt"
	"log"
	"main/core/media"
	"os"
	"path"
	"strconv"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
)

func NewAudioMessage(data []byte, mimetype string, ptt bool, seconds int, ctx *proto.ContextInfo, client *whatsmeow.Client) (*proto.Message, error) {
	var message, err = NewAudio(data, mimetype, ptt, seconds, ctx, client)
	return &proto.Message{AudioMessage: message}, err

}

func NewAudio(data []byte, mimetype string, ptt bool, seconds int, ctx *proto.ContextInfo, client *whatsmeow.Client) (*proto.AudioMessage, error) {

	if up, err := client.Upload(context.Background(), data, whatsmeow.MediaAudio); err != nil {
		return nil, err
	} else {
		var message = &proto.AudioMessage{
			Url:           &up.URL,
			Mimetype:      &mimetype,
			FileLength:    &up.FileLength,
			FileSha256:    up.FileSHA256,
			FileEncSha256: up.FileEncSHA256,
			MediaKey:      up.MediaKey,
			DirectPath:    &up.DirectPath,
			Ptt:           &ptt,
			ContextInfo:   ctx,
			Seconds:       Uint32P(uint32(seconds)),
		}

		return message, err
	}

}

func NewAudioMessageFile(filename string, ctx *proto.ContextInfo, client *whatsmeow.Client) (*proto.Message, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	} else {

		var seconds = 0
		var mimetype = "audio/mpeg"
		var ptt bool = true

		switch strings.ToLower(path.Ext(filename)) {
		case ".mp3":
			{
				mimetype = "audio/mpeg"

				var args_probe = []string{
					"-v quiet",
					"-print_format json",
					"-show_streams",
				}

				if map_stream, output, err := media.FfProbe(fmt.Sprintf(`"%s"`, filename), true, args_probe); err != nil {
					log.Println(err, string(output), filename)
				} else {

					if tmpvar, ok := map_stream["duration"]; ok {
						if tmpval, err := strconv.ParseFloat(tmpvar.(string), 32); err == nil {
							seconds = int(tmpval)
						} else {
							log.Printf("duration %v", EventToString(tmpvar))
						}
					}
				}

			}
		case ".ogg":
			{
				mimetype = "audio/ogg; codecs=opus"
			}

		}

		if data, err := os.ReadFile(filename); err != nil {
			return nil, err
		} else {

			return NewAudioMessage(data, mimetype, ptt, seconds, ctx, client)
		}
	}
}
