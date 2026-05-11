package sendx

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
)

type MsgConfig struct {
	Mail     *gmailConfig
	DingDing *dingDingConfig
	TgConfig *tgConfig
}

type MsgConfigOption func(*MsgConfig)

func NewSendMsgEngine(opts ...MsgConfigOption) MsgConfig {
	conf := &MsgConfig{}

	if opts != nil {
		for _, opt := range opts {
			opt(conf)
		}
	} else {
		log.Fatal("opts is nil")
	}

	if conf.TgConfig != nil {
		bot, err := tgbotapi.NewBotAPI(conf.TgConfig.botToken)
		if err != nil {
			logx.Errorf("new bot api failed: %v", err)
		}
		conf.TgConfig.tBot = bot
	}
	return MsgConfig{
		Mail:     conf.Mail,
		DingDing: conf.DingDing,
		TgConfig: conf.TgConfig,
	}
}
