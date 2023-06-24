package whats

import (
	"context"
	"errors"
	"fmt"
	"main/core/configs"
	"main/core/emoji"
	"main/core/mylog"
	"os"
	"strings"
	"time"

	qrterminal "github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"

	// Imported for postgres dialect
	_ "github.com/jackc/pgx/v4/stdlib"

	// Imported for sqlite dialect
	_ "github.com/mattn/go-sqlite3"
)

type SentHandler func(types.JID, *proto.Message, whatsmeow.SendResponse, error)

// GoWhatsBot is a whatsmeow client container
type GoWhatsBot struct {
	*whatsmeow.Client

	Config *configs.Config

	Container   *sqlstore.Container
	DeviceStore *store.Device

	// Logging output for Database and Client
	DatabaseLog waLog.Logger
	ClientLog   waLog.Logger

	QRCode    string
	Connected bool
	Kill      chan bool

	// Statistics

	// Handler that executed when message sended
	SentHandlers []SentHandler
}

// Create new GoWhatsBot instance
func NewGoWhatsBot(cfg *configs.Config) (*GoWhatsBot, error) {

	var that_empty = []string{}
	if cfg.Name == "" {
		that_empty = append(that_empty, "Name")
	}

	if cfg.Driver == "" {
		that_empty = append(that_empty, "Driver")
	}

	if cfg.Address == "" {
		that_empty = append(that_empty, "Address")
	}

	if len(that_empty) > 0 {
		var error_text = strings.Join(that_empty, ", ")
		return nil, errors.New("some configuration is empty " + error_text)
	}

	var newGWB = GoWhatsBot{}

	var clientLevel = "DEBUG"
	if len(cfg.ClientLog) > 0 {
		clientLevel = strings.ToUpper(cfg.ClientLog)
	}

	if len(cfg.OS) > 0 {
		store.DeviceProps.Os = StringP(cfg.OS)
	}

	if cfg.Platform > 0 {
		store.DeviceProps.PlatformType = (*proto.DeviceProps_PlatformType)(Int32P(cfg.Platform))
	}

	newGWB.DatabaseLog = waLog.Stdout("Database", "ERROR", true)
	newGWB.ClientLog = waLog.Stdout("Client", clientLevel, true)
	newGWB.Config = cfg

	if ctr, err := sqlstore.New(cfg.Driver, cfg.Address, newGWB.DatabaseLog); err != nil {
		return nil, err
	} else {
		newGWB.Container = ctr
	}

	if dvc, err := newGWB.Container.GetFirstDevice(); err != nil {
		return nil, err
	} else {
		newGWB.DeviceStore = dvc
	}

	newGWB.Client = whatsmeow.NewClient(newGWB.DeviceStore, newGWB.ClientLog)

	return &newGWB, nil
}

func (g *GoWhatsBot) GetDeviceProps() *proto.DeviceProps {
	return store.DeviceProps
}

func (g *GoWhatsBot) GetConfig() *configs.Config {

	return g.Config
}

func (g *GoWhatsBot) Run() {
	if g.Client.Store.ID == nil {
		if qrchan, err := g.Client.GetQRChannel(context.Background()); err == nil {
			if err := g.Client.Connect(); err != nil {
				panic(err)
			}

			for qritem := range qrchan {
				if qritem.Event == "code" {
					g.QRCode = qritem.Code
					if g.Config.ShowQR {
						// qrterminal.GenerateHalfBlock(qritem.Code, qrterminal.M, os.Stdout)

						var qrterminalConfig = qrterminal.Config{
							HalfBlocks:     true,
							Level:          qrterminal.M,
							Writer:         os.Stdout,
							BlackChar:      qrterminal.BLACK_BLACK,
							WhiteChar:      qrterminal.WHITE_WHITE,
							BlackWhiteChar: qrterminal.BLACK_WHITE,
							WhiteBlackChar: qrterminal.WHITE_BLACK,
						}

						qrterminal.GenerateWithConfig(qritem.Code, qrterminalConfig)
					}
				} else {
					g.Connected = qritem.Event == "success"
				}
			}
		}
	} else {
		if err := g.Client.Connect(); err != nil {
			panic(err)
		}
	}

	for {
		select {
		case <-g.Kill:
			mylog.Info(fmt.Sprintf("Client %s disconnected", g.Config.Name))
			g.Stop()
			return
		default:
			time.Sleep(1000 * time.Millisecond)

		}

	}
}

func (g *GoWhatsBot) Stop() {
	g.Client.Disconnect()
	close(g.Kill)
}

func (g *GoWhatsBot) AddSentHandler(h SentHandler) {
	g.SentHandlers = append(g.SentHandlers, h)
}

func (g *GoWhatsBot) SendMessage(chat types.JID, message *proto.Message) (whatsmeow.SendResponse, error) {
	var resp, err = g.Client.SendMessage(context.Background(), chat, message)

	// Execute every Sent handlers
	for _, handler := range g.SentHandlers {
		handler(chat, message, resp, err)
	}

	return resp, err
}

func (g *GoWhatsBot) SendChatPresence(chat types.JID, media types.ChatPresenceMedia, state types.ChatPresence) error {
	return g.Client.SendChatPresence(chat, state, media)
}

func (g *GoWhatsBot) SendTyping(chat types.JID) error {
	return g.SendChatPresence(chat, types.ChatPresenceMediaText, types.ChatPresenceComposing)
}

func (g *GoWhatsBot) SendStopTyping(chat types.JID) error {
	return g.SendChatPresence(chat, types.ChatPresenceMediaText, types.ChatPresencePaused)
}

func (g *GoWhatsBot) SendRecording(chat types.JID) error {
	return g.SendChatPresence(chat, types.ChatPresenceMediaAudio, types.ChatPresenceComposing)
}

func (g *GoWhatsBot) SendStopRecording(chat types.JID) error {
	return g.SendChatPresence(chat, types.ChatPresenceMediaAudio, types.ChatPresencePaused)
}

func (g *GoWhatsBot) SendReact(info types.MessageInfo, react emoji.Emoji) (whatsmeow.SendResponse, error) {
	this_message := &proto.Message{
		ReactionMessage: &proto.ReactionMessage{
			Key: &proto.MessageKey{
				RemoteJid:   StringP(info.Chat.String()),
				Participant: StringP(info.Sender.String()),
				FromMe:      BoolP(info.IsFromMe),
				Id:          &info.ID,
			},
			Text: StringP(string(react)),
		},
	}

	return g.SendMessage(info.Chat, this_message)
}

func (g *GoWhatsBot) RemoveReact(info types.MessageInfo) (whatsmeow.SendResponse, error) {
	return g.SendReact(info, "")
}

func (g *GoWhatsBot) SendText(chat types.JID, text string, newctx *proto.ContextInfo) (whatsmeow.SendResponse, error) {
	if newmsg, err := NewExtendedMessage(text, newctx); err != nil {
		return whatsmeow.SendResponse{}, err
	} else {
		return g.SendMessage(chat, newmsg)
	}
}

func (g *GoWhatsBot) MarkReadOne(info types.MessageInfo) error {
	return g.Client.MarkRead([]string{info.ID}, time.Now(), info.Chat, info.Sender)
}
