package captcha

import (
	"errors"
	"github.com/mojocn/base64Captcha"
	"image/color"
	"time"
)

// CaptchaType 验证码类型
type CaptchaType string

const (
	TypeDigit  CaptchaType = "digit"  // 数字验证码
	TypeString CaptchaType = "string" // 字符串验证码
	TypeMath   CaptchaType = "math"   // 数学验证码
	TypeAudio  CaptchaType = "audio"  // 音频验证码
)

// CreateDriver 根据类型创建验证码驱动
func CreateDriver(captchaType CaptchaType, config *CaptchaConfig) base64Captcha.Driver {
	// 设置默认背景颜色（白色）
	bgColor := &color.RGBA{R: 255, G: 255, B: 255, A: 255}

	switch captchaType {
	case TypeString:
		// 字符串验证码驱动
		return base64Captcha.NewDriverString(
			config.Height,         // height
			config.Width,          // width
			config.NoiseCount,     // noiseCount
			config.ShowLineOption, // showLineOptions
			config.Length,         // length
			config.Source,         // source
			bgColor,               // bgColor
			nil,                   // fontsStorage (使用默认)
			config.Fonts,          // fonts
		)
	case TypeMath:
		// 数学验证码驱动 - 修正参数
		return base64Captcha.NewDriverMath(
			config.Height,         // height
			config.Width,          // width
			config.NoiseCount,     // noiseCount
			config.ShowLineOption, // showLineOptions
			bgColor,               // bgColor
			nil,                   // fontsStorage
			config.Fonts,          // fonts
		)
	case TypeAudio:
		// 音频验证码驱动
		return base64Captcha.NewDriverAudio(
			config.Length,   // length
			config.Language, // language
		)
	default: // TypeDigit
		// 数字验证码驱动
		return base64Captcha.NewDriverDigit(
			config.Width,    // width
			config.Height,   // height
			config.Length,   // length
			config.MaxSkew,  // maxSkew
			config.DotCount, // dotCount
		)
	}
}

// GenerateByType 根据类型生成验证码
func GenerateByType(captchaType CaptchaType) (*CaptchaResult, error) {
	if captchaManager == nil {
		return nil, errors.New("验证码管理器未初始化")
	}

	// 使用管理器的配置
	driver := CreateDriver(captchaType, captchaManager.config)

	// 临时创建验证码实例
	tempCaptcha := base64Captcha.NewCaptcha(driver, captchaManager.store)
	id, b64s, answer, err := tempCaptcha.Generate()
	if err != nil {
		return nil, err
	}

	return &CaptchaResult{
		CaptchaID:    id,
		CaptchaImage: b64s,
		Answer:       answer,
		ExpireTime:   time.Now().Add(time.Duration(captchaManager.config.ExpireSeconds) * time.Second).Unix(),
	}, nil
}

// GenerateStringCaptcha 生成字符串验证码（便捷方法）
func GenerateStringCaptcha(source string, length int) (*CaptchaResult, error) {
	if captchaManager == nil {
		return nil, errors.New("验证码管理器未初始化")
	}

	// 临时修改配置
	tempConfig := *captchaManager.config
	if source != "" {
		tempConfig.Source = source
	}
	if length > 0 {
		tempConfig.Length = length
	}

	driver := CreateDriver(TypeString, &tempConfig)
	tempCaptcha := base64Captcha.NewCaptcha(driver, captchaManager.store)
	id, b64s, answer, err := tempCaptcha.Generate()
	if err != nil {
		return nil, err
	}

	return &CaptchaResult{
		CaptchaID:    id,
		CaptchaImage: b64s,
		Answer:       answer,
		ExpireTime:   time.Now().Add(time.Duration(tempConfig.ExpireSeconds) * time.Second).Unix(),
	}, nil
}

// GenerateMathCaptcha 生成数学验证码
func GenerateMathCaptcha() (*CaptchaResult, error) {
	return GenerateByType(TypeMath)
}

// GenerateAudioCaptcha 生成音频验证码
func GenerateAudioCaptcha() (*CaptchaResult, error) {
	return GenerateByType(TypeAudio)
}

// GenerateDigitCaptcha 生成数字验证码
func GenerateDigitCaptcha() (*CaptchaResult, error) {
	return GenerateByType(TypeDigit)
}
