package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type HttpClient struct {
	Client *http.Client
	Retry  int
	Delay  time.Duration
}

func NewHttpClient(timeout time.Duration, retry int, delay time.Duration) *HttpClient {
	if retry == 0 {
		retry = 1
	}
	return &HttpClient{
		Client: &http.Client{
			Timeout: timeout,
		},
		Retry: retry,
		Delay: delay,
	}
}

// Request Sending HTTP requests, supporting retry mechanism, exponential backoff, and context
func (c *HttpClient) Request(ctx context.Context, method, url string, headers map[string]string, body io.Reader, result interface{}) error {
	var lastErr error

	for attempt := 0; attempt < c.Retry; attempt++ {
		if ctx.Err() != nil {
			return fmt.Errorf("request canceled: %v", ctx.Err())
		}

		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %v", err)
			log.Printf("Attempt %d failed: %v\n", attempt+1, lastErr)
			c.waitForRetry(attempt) // 指数退避
			continue
		}

		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := c.Client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to send request: %v", err)
			log.Printf("Attempt %d failed: %v\n", attempt+1, lastErr)
			c.waitForRetry(attempt)
			continue
		}

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			log.Printf("Attempt %d failed: %v\n", attempt+1, lastErr)
			resp.Body.Close()
			c.waitForRetry(attempt)
			continue
		}

		if result != nil {
			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				body, _ := io.ReadAll(resp.Body)
				lastErr = fmt.Errorf("failed to decode response url : %s , body : %s , error : %v", url, body, err)
				log.Printf("Attempt %d failed: %v\n", attempt+1, lastErr)
				resp.Body.Close()
				c.waitForRetry(attempt)
				continue
			}
		}

		resp.Body.Close()
		return nil
	}

	// All retry attempts have failed
	return fmt.Errorf("after %d attempts, last error: %v", c.Retry+1, lastErr)
}

// waitForRetry Implement index backoff
func (c *HttpClient) waitForRetry(attempt int) {
	if attempt < c.Retry {
		delay := c.Delay * time.Duration(1<<attempt)
		log.Printf("Waiting %v before retry...\n", delay)
		time.Sleep(delay)
	}
}

func (c *HttpClient) Get(ctx context.Context, url string, headers map[string]string, result interface{}) error {
	return c.Request(ctx, http.MethodGet, url, headers, nil, result)
}

func (c *HttpClient) Post(ctx context.Context, url string, headers map[string]string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	if headers == nil {
		headers = make(map[string]string)
	}
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json"
	}

	return c.Request(ctx, http.MethodPost, url, headers, bodyReader, result)
}

func (c *HttpClient) RequestNotRetry(method, url string, headers map[string]string, body interface{}) ([]byte, int, error) {
	var reqBody []byte
	var err error

	// 如果有请求体，将其序列化为 JSON
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal request body : %v", err)
		}
	}

	// 创建请求
	var req *http.Request

	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(method, url, nil)
	case http.MethodPost:
		req, err = http.NewRequest(method, url, bytes.NewBuffer(reqBody))

	}

	if err != nil {
		return nil, 0, fmt.Errorf("creage request error : %v", err)
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 如果是 POST 请求，设置 Content-Type 为 application/json
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	// 发送请求
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request send error : %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("reqd response error : %v", err)
	}

	return respBody, resp.StatusCode, nil
}
