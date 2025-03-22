package otp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// Redis 客戶端
	redisClient *redis.Client

	// 上下文
	ctx = context.Background()

	// HTTP 客戶端
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	// OTP 請求等待映射
	otpRequests     = make(map[string]chan string)
	otpRequestsLock sync.Mutex

	// Line Bot API URL
	lineBotURL string

	// 超時設置（秒）
	otpTimeout = 300 // 5 分鐘
)

// OTPRequest 表示一個 OTP 請求
type OTPRequest struct {
	Email     string `json:"email"`
	RequestID string `json:"request_id"`
}

// OTPResponse 表示一個 OTP 響應
type OTPResponse struct {
	Email     string `json:"email"`
	OTP       string `json:"otp"`
	RequestID string `json:"request_id"`
}

// InitOTP 初始化 OTP 服務
func InitOTP() error {
	// 從環境變量獲取 Line Bot API URL
	lineBotURL = os.Getenv("LINEBOT_API_URL")
	if lineBotURL == "" {
		lineBotURL = "http://localhost:4000/ask_otp"
	}

	// 從環境變量獲取 Redis 配置
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0

	// 初始化 Redis 客戶端
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// 測試連接
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	return nil
}

func CloseOTP() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

func NotifyOTPRequired(email string) (string, error) {
	if redisClient == nil {
		return "", errors.New("redis client not initialized")
	}

	requestID := fmt.Sprintf("%s-%d", email, time.Now().UnixNano())

	request := OTPRequest{
		Email:     email,
		RequestID: requestID,
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := httpClient.Post(
		lineBotURL,
		"application/json",
		bytes.NewBuffer(requestJSON),
	)
	if err != nil {
		return "", fmt.Errorf("failed to send OTP request to Line Bot: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Line Bot returned non-200 status: %d", resp.StatusCode)
	}

	return requestID, nil
}

func WaitForOTP(email, requestID string, timeout time.Duration) (string, error) {
	if redisClient == nil {
		return "", errors.New("redis client not initialized")
	}

	redisKey := fmt.Sprintf("otp:response:%s:%s", email, requestID)

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		respJSON, err := redisClient.Get(ctx, redisKey).Result()
		if err == redis.Nil {
			time.Sleep(2 * time.Second)
			continue
		} else if err != nil {
			return "", fmt.Errorf("redis error: %w", err)
		}

		var response OTPResponse
		if err := json.Unmarshal([]byte(respJSON), &response); err != nil {
			return "", fmt.Errorf("failed to unmarshal OTP response: %w", err)
		}

		return response.OTP, nil
	}

	return "", errors.New("timeout waiting for OTP")
}
