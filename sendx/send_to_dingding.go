package sendx

import (
	"base/httpx"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type dingDingConfig struct {
	webhookURL  string
	accessToken string
	secret      string
}

type DingDingMsgConfigOption func(*dingDingConfig)

func WithDingDingConfig(webhookURL, accessToken, secret string) MsgConfigOption {
	return func(c *MsgConfig) {
		c.DingDing = &dingDingConfig{
			webhookURL:  webhookURL,
			accessToken: accessToken,
			secret:      secret,
		}
	}
}

func WithDefaultDingDingConfig() MsgConfigOption {
	return func(c *MsgConfig) {
		c.DingDing = &dingDingConfig{
			webhookURL:  "https://oapi.dingtalk.com/robot/send",
			accessToken: "5eadc983521dd210f4fbabb3ecd72a88c339dc06efbdf7d46443c8e534ef8575",
			secret:      "SECf4a385fa5466f4983ff465af072321f2f4a929ba9a980c264844009e3b0d29ba",
		}
	}
}

func (d dingDingConfig) SendMsg(ctx context.Context, msg string) error {
	timestamp := time.Now().UnixNano() / 1e6
	timestampStr := strconv.FormatInt(timestamp, 10)

	stringToSign := timestampStr + "\n" + d.secret

	h := hmac.New(sha256.New, []byte(d.secret))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	signedURL := fmt.Sprintf("%s?access_token=%s&timestamp=%s&sign=%s",
		d.webhookURL,
		d.accessToken,
		timestampStr,
		url.QueryEscape(signature),
	)

	message := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": msg,
		},
	}

	//jsonData, err := json.Marshal(message)
	//if err != nil {
	//	return err
	//}

	httpc := httpx.NewHttpClient(10*time.Second, 0, 0)
	var result = make(map[string]interface{})
	err := httpc.Post(ctx, signedURL, nil, message, &result)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
