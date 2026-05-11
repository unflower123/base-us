package consts

import "fmt"

const (
	AGENT_USER_TOKEN = "agent_user_token"
	AGENT_USER_INFO  = "agent_user_info"
)

func GenerateAgentUserTokenKey(serverName string, agentUserID uint64) string {
	return fmt.Sprintf("%s:%s:%v", serverName, AGENT_USER_TOKEN, agentUserID)
}

func GenerateAgentUserInfoKey(serverName string, agentUserID uint64) string {
	return fmt.Sprintf("%s:%s:%v", serverName, AGENT_USER_INFO, agentUserID)
}
