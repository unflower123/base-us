package captcha

import (
	"errors"
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"time"
)

// CaptchaResult 验证码生成结果
type CaptchaResult struct {
	CaptchaID    string `json:"captcha_id"`
	CaptchaImage string `json:"captcha_image"`
	Answer       string `json:"-"`           // 不序列化到JSON，内部使用
	ExpireTime   int64  `json:"expire_time"` // 过期时间戳
}

// VerifyResult 验证结果
type VerifyResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CaptchaManager 验证码管理器
type CaptchaManager struct {
	captcha *base64Captcha.Captcha
	store   base64Captcha.Store
	config  *CaptchaConfig
}

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	Width         int     // 宽度
	Height        int     // 高度
	Length        int     // 长度
	MaxSkew       float64 // 最大倾斜
	DotCount      int     // 点数
	ExpireSeconds int     // 过期时间（秒）

	// 字符串和数学验证码共用配置
	NoiseCount     int      // 干扰线数量
	ShowLineOption int      // 显示线选项
	Source         string   // 字符源
	Fonts          []string // 字体列表

	// 音频验证码配置
	Language string // 语言
}

// DefaultConfig 默认配置
func DefaultConfig() *CaptchaConfig {
	return &CaptchaConfig{
		Width:         80,
		Height:        240,
		Length:        6,
		MaxSkew:       0.7,
		DotCount:      80,
		ExpireSeconds: 300, // 5分钟

		NoiseCount:     5,
		ShowLineOption: 2,
		Source:         "1234567890abcdefghijklmnopqrstuvwxyz",
		Fonts:          []string{"wqy-microhei.ttc"},

		Language: "zh", // 中文
	}
}

// DigitConfig 数字验证码专用配置
func DigitConfig() *CaptchaConfig {
	return &CaptchaConfig{
		Width:         80,
		Height:        240,
		Length:        6,
		MaxSkew:       0.7,
		DotCount:      80,
		ExpireSeconds: 300,
	}
}

// StringConfig 字符串验证码专用配置
func StringConfig() *CaptchaConfig {
	return &CaptchaConfig{
		Width:          100,
		Height:         40,
		Length:         4,
		ExpireSeconds:  300,
		NoiseCount:     3,
		ShowLineOption: 2,
		Source:         "ABCDEFGHJKLMNPQRSTUVWXYZ23456789",
		Fonts:          []string{"wqy-microhei.ttc"},
	}
}

// MathConfig 数学验证码专用配置
func MathConfig() *CaptchaConfig {
	return &CaptchaConfig{
		Width:          100,
		Height:         40,
		ExpireSeconds:  300,
		NoiseCount:     2,
		ShowLineOption: 1,
		Fonts:          []string{"wqy-microhei.ttc"},
	}
}

// AudioConfig 音频验证码专用配置
func AudioConfig() *CaptchaConfig {
	return &CaptchaConfig{
		Length:        4,
		ExpireSeconds: 300,
		Language:      "zh",
	}
}

// 全局验证码管理器实例
var captchaManager *CaptchaManager

// InitWithRedis 使用Redis存储初始化验证码
// InitDigitCaptchaWithRedis 初始化数字验证码（Redis存储）
func InitDigitCaptchaWithRedis(redisClient *redis.Client, config *CaptchaConfig) error {
	if redisClient == nil {
		return errors.New("Redis is null")
	}

	if config == nil {
		config = DigitConfig()
	}

	store := NewRedisStore(redisClient, config.ExpireSeconds)
	driver := CreateDriver(TypeDigit, config)

	captchaManager = &CaptchaManager{
		captcha: base64Captcha.NewCaptcha(driver, store),
		store:   store,
		config:  config,
	}

	return nil
}

// InitStringCaptchaWithRedis 初始化字符串验证码（Redis存储）
func InitStringCaptchaWithRedis(redisClient *redis.Client, config *CaptchaConfig) error {
	if redisClient == nil {
		return errors.New("Redis is null")
	}

	if config == nil {
		config = StringConfig()
	}

	store := NewRedisStore(redisClient, config.ExpireSeconds)
	driver := CreateDriver(TypeString, config)

	captchaManager = &CaptchaManager{
		captcha: base64Captcha.NewCaptcha(driver, store),
		store:   store,
		config:  config,
	}

	return nil
}

// InitMathCaptchaWithRedis 初始化数学验证码（Redis存储）
func InitMathCaptchaWithRedis(redisClient *redis.Client, config *CaptchaConfig) error {
	if redisClient == nil {
		return errors.New("Redis is null")
	}

	if config == nil {
		config = MathConfig()
	}

	store := NewRedisStore(redisClient, config.ExpireSeconds)
	driver := CreateDriver(TypeMath, config)

	captchaManager = &CaptchaManager{
		captcha: base64Captcha.NewCaptcha(driver, store),
		store:   store,
		config:  config,
	}

	return nil
}

// InitAudioCaptchaWithRedis 初始化音频验证码（Redis存储）
func InitAudioCaptchaWithRedis(redisClient *redis.Client, config *CaptchaConfig) error {
	if redisClient == nil {
		return errors.New("Redis is null")
	}

	if config == nil {
		config = AudioConfig()
	}

	store := NewRedisStore(redisClient, config.ExpireSeconds)
	driver := CreateDriver(TypeAudio, config)

	captchaManager = &CaptchaManager{
		captcha: base64Captcha.NewCaptcha(driver, store),
		store:   store,
		config:  config,
	}

	return nil
}

// 对应的内存存储初始化方法
func InitDigitCaptchaWithMemory(config *CaptchaConfig) error {
	if config == nil {
		config = DigitConfig()
	}

	store := base64Captcha.DefaultMemStore
	driver := CreateDriver(TypeDigit, config)

	captchaManager = &CaptchaManager{
		captcha: base64Captcha.NewCaptcha(driver, store),
		store:   store,
		config:  config,
	}

	return nil
}

func InitStringCaptchaWithMemory(config *CaptchaConfig) error {
	if config == nil {
		config = StringConfig()
	}

	store := base64Captcha.DefaultMemStore
	driver := CreateDriver(TypeString, config)

	captchaManager = &CaptchaManager{
		captcha: base64Captcha.NewCaptcha(driver, store),
		store:   store,
		config:  config,
	}

	return nil
}

func InitMathCaptchaWithMemory(config *CaptchaConfig) error {
	if config == nil {
		config = MathConfig()
	}

	store := base64Captcha.DefaultMemStore
	driver := CreateDriver(TypeMath, config)

	captchaManager = &CaptchaManager{
		captcha: base64Captcha.NewCaptcha(driver, store),
		store:   store,
		config:  config,
	}

	return nil
}

func InitAudioCaptchaWithMemory(config *CaptchaConfig) error {
	if config == nil {
		config = AudioConfig()
	}

	store := base64Captcha.DefaultMemStore
	driver := CreateDriver(TypeAudio, config)

	captchaManager = &CaptchaManager{
		captcha: base64Captcha.NewCaptcha(driver, store),
		store:   store,
		config:  config,
	}

	return nil
}

// Generate 生成验证码
func Generate() (*CaptchaResult, error) {
	if captchaManager == nil {
		//return nil, errors.New("验证码管理器未初始化，请先调用InitWithRedis或InitWithMemory")
		return nil, errors.New("captcha init fail")
	}

	id, b64s, answer, err := captchaManager.captcha.Generate()
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

// GenerateWithAnswer 生成验证码并返回答案（用于测试或特殊场景）
func GenerateWithAnswer() (*CaptchaResult, error) {
	result, err := Generate()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GenerateWithoutAnswer 生成验证码但不返回答案（安全场景）
func GenerateWithoutAnswer() (*CaptchaResult, error) {
	result, err := Generate()
	if err != nil {
		return nil, err
	}
	// 清空答案
	result.Answer = ""
	return result, nil
}

// Verify 验证验证码
func Verify(captchaID, code string) error {
	if captchaManager == nil {
		return errors.New("The CAPTCA manager is not initialized")
	}

	if captchaID == "" || code == "" {
		return errors.New("The verification code ID or verification code cannot be empty")
	}

	success := captchaManager.store.Verify(captchaID, code, true)
	if success {
		return nil
	}
	return errors.New("The verification code is incorrect or has expired")
}

// VerifyWithoutClear 验证验证码但不清除（用于多次验证）
func VerifyWithoutClear(captchaID, code string) *VerifyResult {
	if captchaManager == nil {
		return &VerifyResult{Success: false, Message: "The CAPTCA manager is not initialized"}
	}

	success := captchaManager.store.Verify(captchaID, code, false)
	if success {
		return &VerifyResult{Success: true, Message: "verify success"}
	}
	return &VerifyResult{Success: false, Message: "verify fail"}
}

// GetAnswer 获取验证码答案（谨慎使用，主要用于测试）
func GetAnswer(captchaID string) (string, error) {
	if captchaManager == nil {
		return "", errors.New("The CAPTCA manager is not initialized")
	}

	answer := captchaManager.store.Get(captchaID, false)
	if answer == "" {
		return "", errors.New("The verification code does not exist or has expired")
	}

	return answer, nil
}

// IsInitialized 检查验证码管理器是否已初始化
func IsInitialized() bool {
	return captchaManager != nil
}
