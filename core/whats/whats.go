package whats

import (
	"reflect"
	"strings"

	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

var expirationChat map[string]uint32 = map[string]uint32{}

func GetTextContext(msg *proto.Message) (string, *proto.ContextInfo) {
	var text string = ""
	var ctx *proto.ContextInfo

	if msg != nil {

		if msg_type := msg.GetAudioMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()

		} else if msg_type := msg.GetButtonsMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetContentText()

		} else if msg_type := msg.GetButtonsResponseMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetSelectedButtonId()

		} else if msg_type := msg.GetContactMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetVcard()

		} else if msg_type := msg.GetConversation(); msg_type != "" {
			text = msg_type

		} else if msg_type := msg.GetDocumentMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetCaption()

		} else if msg_type := msg.GetExtendedTextMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetText()

		} else if msg_type := msg.GetImageMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetCaption()

		} else if msg_type := msg.GetListMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetDescription()

		} else if msg_type := msg.GetListResponseMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.SingleSelectReply.GetSelectedRowId()

		} else if msg_type := msg.GetProductMessage(); msg_type != nil {
			ctx = msg_type.ContextInfo
			text = msg_type.GetBody()

		} else if msg_type := msg.GetReactionMessage(); msg_type != nil {
			text = msg_type.GetText()

		} else if msg_type := msg.GetStickerMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()

		} else if msg_type := msg.GetTemplateButtonReplyMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetSelectedId()

		} else if msg_type := msg.GetTemplateMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()

			if msg_subtype := msg_type.GetFourRowTemplate(); msg_subtype != nil {
				text = msg_subtype.Content.GetNamespace()

			} else if msg_subtype := msg_type.GetHydratedTemplate(); msg_subtype != nil {
				text = msg_subtype.GetTemplateId()

			} else if msg_subtype := msg_type.GetHydratedFourRowTemplate(); msg_subtype != nil {
				text = msg_subtype.GetTemplateId()

			}

		} else if msg_type := msg.GetVideoMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetCaption()

		}

	}

	return text, ctx
}

type ContextMode int

const (
	// Don't apply anything
	CTX_NONE ContextMode = iota

	// Always quote everywhere
	CTX_QUOTE

	// Only mention everywhere (without quote)
	CTX_ONLYMENTION

	// Only quote if it Chat is a group
	CTX_GROUP_QUOTE

	// Only mention if it Chat in a group (without quote)
	CTX_GROUP_ONLYMENTION
)

func NewContext(e *events.Message, mode ContextMode) *proto.ContextInfo {
	_, ctx := GetTextContext(e.Message)

	var newctx proto.ContextInfo
	if ctx != nil {
		newctx = proto.ContextInfo{
			Expiration: ctx.Expiration,
		}
	}

	if exp, ok := expirationChat[e.Info.Chat.User]; ok {
		newctx.Expiration = Uint32P(exp)
	} else {
		if newctx.Expiration != nil {
			expirationChat[e.Info.Chat.User] = *newctx.Expiration
		}
	}

	var quote bool
	var mention bool

	switch mode {

	case CTX_QUOTE:
		quote = true
		mention = false

	case CTX_ONLYMENTION:
		quote = false
		mention = true

	case CTX_GROUP_QUOTE:
		quote = e.Info.IsGroup
		mention = false

	case CTX_GROUP_ONLYMENTION:
		quote = false
		mention = e.Info.IsGroup

	}

	if quote {
		newctx.StanzaId = StringP(e.Info.ID)
		newctx.Participant = StringP(e.Info.Sender.ToNonAD().String())
		newctx.QuotedMessage = e.Message
		newctx.RemoteJid = StringP(e.Info.Chat.String())
	}

	if mention {
		newctx.Participant = StringP(e.Info.Sender.String())
	}

	return &newctx
}

func PrepareExec(e interface{}) (*events.Message, *proto.ContextInfo, string, []string, error) {
	var event = e.(*events.Message)

	var text, ctx = GetTextContext(event.Message)
	text = strings.TrimSpace(text)

	if len(text) > 0 {
		var args = strings.Split(text, " ")
		return event, ctx, args[0], args[1:], nil
	} else {
		return event, ctx, "", []string{}, nil
	}

}

func ParticipantUser(s *string) string {
	var tag = "@"
	if strings.Contains(*s, tag) {
		var users = strings.Split(*s, tag)
		if len(users) > 0 {
			return users[0]
		}
	}

	return *s
}

func EventToString(i interface{}) string {
	return reflect.TypeOf(i).String()
}
