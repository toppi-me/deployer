package bot

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/toppi-me/deployer/internal/github"
	"github.com/toppi-me/deployer/internal/log"
)

var bot *tgbotapi.BotAPI
var chatID int64

// InitTelegramBot initialize bot instance, need to run only once
func InitTelegramBot() (err error) {
	chatID, err = strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		return err
	}

	bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		return err
	}

	return
}

// SendDeployInfo send to initialized chatID message with success deploy info
func SendDeployInfo(gitEvent *github.PushEvent, outs [][]byte) {
	file := tgbotapi.FileBytes{
		Name: strconv.FormatInt(time.Now().UnixNano(), 10) + ".txt",
		Bytes: func() []byte {
			bytesBuff := bytes.Buffer{}
			for _, out := range outs {
				bytesBuff.Write(out)
				bytesBuff.WriteString("\n\n")
			}
			return bytesBuff.Bytes()
		}(),
	}

	msg := tgbotapi.NewDocument(chatID, file)
	msg.Caption = fmt.Sprintf(
		"↗️ New Deploy From *%s* To *%s* / _%s_",
		gitEvent.AuthorName,
		gitEvent.Repository,
		gitEvent.Branch,
	)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	return
}

// SendErrorMsg send to initialized chatID message with error deploy info
func SendErrorMsg(gitEvent *github.PushEvent, outs [][]byte) {
	file := tgbotapi.FileBytes{
		Name: strconv.FormatInt(time.Now().UnixNano(), 10) + ".txt",
		Bytes: func() []byte {
			bytesBuff := bytes.Buffer{}
			for _, out := range outs {
				bytesBuff.Write(out)
				bytesBuff.WriteString("\n\n")
			}
			return bytesBuff.Bytes()
		}(),
	}

	msg := tgbotapi.NewDocument(chatID, file)
	msg.Caption = fmt.Sprintf(
		"❌ Fail Deploy From *%s* To *%s* / _%s_",
		gitEvent.AuthorName,
		gitEvent.Repository,
		gitEvent.Branch,
	)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	return
}
