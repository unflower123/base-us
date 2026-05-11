/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/5/26 15:37
 */
package etcd

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/zeromicro/go-zero/core/discov"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	dialTimeout      = 5 * time.Second
	autoSyncInterval = time.Minute
)

var (
	DialTimeout = dialTimeout
)

// Config 结构体用于存储从 etcd 获取的配置和客户端
type Config struct {
	client *clientv3.Client
}

// NewConfig 创建一个新的 Config 实例
func NewConfig(config discov.EtcdConf) (*Config, error) {
	// 创建 etcd 客户端配置
	cfg := clientv3.Config{
		Endpoints:           config.Hosts,
		AutoSyncInterval:    autoSyncInterval,
		DialTimeout:         DialTimeout,
		RejectOldCluster:    true,
		PermitWithoutStream: true,
	}

	if len(config.User) > 0 && len(config.Pass) > 0 {
		cfg.Username = config.User
		cfg.Password = config.Pass
	}

	// 连接到 etcd
	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %v", err)
	}

	return &Config{client: client}, nil
}

// GetConfig 从 etcd 获取配置
func (c *Config) GetConfig(ctx context.Context, key string) (string, error) {
	// 获取配置
	resp, err := c.client.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get config: %v", err)
	}

	// 检查是否找到配置
	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("no config found for key: %s", key)
	}

	// 返回配置值
	return string(resp.Kvs[0].Value), nil
}

// WatchConfig 监听指定 key 的变更并触发服务重启
func (c *Config) WatchConfig(ctx context.Context, key, serverName string) {
	// 创建 watcher
	watchChan := c.client.Watch(ctx, key)

	// 监听变更
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				log.Printf("Key %s updated to value: %s", key, string(event.Kv.Value))
				restartService(serverName)
			case mvccpb.DELETE:
				log.Printf("Key %s deleted", key)
				restartService(serverName)
			default:
				log.Printf("Unknown event type for key %s: %v", key, event.Type)
			}
		}
	}
}

// restartService 服务重启逻辑
func restartService(serverName string) {
	log.Println("Restarting service due to configuration change...")

	// 使用 systemctl 重启服务（替换 "my-service" 为实际服务名）
	cmd := exec.Command("sudo", "systemctl", "restart", serverName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to restart service: %v, output: %s", err, string(output))
		return
	}

	log.Println("Service restarted successfully")
}

// Close 关闭 etcd 客户端连接
func (c *Config) Close() error {
	return c.client.Close()
}
