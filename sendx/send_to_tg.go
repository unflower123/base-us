package sendx

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type tgConfig struct {
	botToken string
	chatID   int64
	tBot     *tgbotapi.BotAPI
}

type tgMsgConfigOption func(*tgConfig)

func WithTgConfig(token string, chatID int64) MsgConfigOption {
	return func(c *MsgConfig) {
		c.TgConfig = &tgConfig{
			botToken: token,
			chatID:   chatID,
		}
	}
}

func (t *tgConfig) SendMsg(ctx context.Context, text string) (err error) {
	if t.tBot == nil {
		t.tBot, err = tgbotapi.NewBotAPI(t.botToken)
		if err != nil {
			logx.Errorf("new bot api failed: %v", err)
			return err
		}
	}
	msg := tgbotapi.NewMessage(t.chatID, text)

	_, err = t.tBot.Send(msg)
	if err != nil {
		logx.Errorf("send message failed: %v", err)
		return err
	}
	return nil
}

func (t *tgConfig) SetChatId(chatId int64) *tgConfig {
	t.chatID = chatId
	return t
}

func (t *tgConfig) SendReplyMsg(ctx context.Context, text string, msgId int) (err error) {
	if t.tBot == nil {
		t.tBot, err = tgbotapi.NewBotAPI(t.botToken)
		if err != nil {
			logx.Errorf("new bot api failed: %v", err)
			return err
		}
	}
	msg := tgbotapi.NewMessage(t.chatID, text)
	msg.ReplyToMessageID = msgId

	_, err = t.tBot.Send(msg)
	if err != nil {
		logx.Errorf("send message failed: %v", err)
		return err
	}
	return nil
}

func (t *tgConfig) SendPhotoURLWithCaption(ctx context.Context, imageURL string, caption string) (err error) {
	if t.tBot == nil {
		// 确保 BotAPI 客户端已初始化 (与您的 SendMsg 保持一致)
		t.tBot, err = tgbotapi.NewBotAPI(t.botToken)
		if err != nil {
			logx.Errorf("new bot api failed: %v", err)
			return err
		}
	}

	file := tgbotapi.FileURL(imageURL)

	photoMsg := tgbotapi.NewPhoto(t.chatID, file)
	photoMsg.Caption = caption

	_, err = t.tBot.Send(photoMsg)
	if err != nil {
		logx.WithContext(ctx).Errorf("send photo message via URL failed: %v", err)
		return fmt.Errorf("failed to send photo message: %w", err)
	}

	return nil
}
