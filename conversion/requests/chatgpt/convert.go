package chatgpt

import (
	chatgpt_types "freechatgpt/internal/chatgpt"
	"freechatgpt/internal/tokens"
	official_types "freechatgpt/typings/official"
	"regexp"
)

var gptsRegexp = regexp.MustCompile(`-gizmo-g-(\w+)`)

func ConvertAPIRequest(api_request official_types.APIRequest, account string, secret *tokens.Secret, deviceId string, proxy string) chatgpt_types.ChatGPTRequest {
	chatgpt_request := chatgpt_types.NewChatGPTRequest()
	chatgpt_request.Model = api_request.Model

	matches := gptsRegexp.FindStringSubmatch(api_request.Model)
	if len(matches) == 2 {
		chatgpt_request.ConversationMode.Kind = "gizmo_interaction"
		chatgpt_request.ConversationMode.GizmoId = "g-" + matches[1]
	}
	ifMultimodel := secret.Token != ""
	for _, api_message := range api_request.Messages {
		if api_message.Role == "system" {
			api_message.Role = "critic"
		}
		chatgpt_request.AddMessage(api_message.Role, api_message.Content, ifMultimodel, account, secret, deviceId, proxy)
	}
	return chatgpt_request
}

func ConvertTTSAPIRequest(input string) chatgpt_types.ChatGPTRequest {
	chatgpt_request := chatgpt_types.NewChatGPTRequest()
	chatgpt_request.HistoryAndTrainingDisabled = false
	chatgpt_request.AddAssistantMessage(input)
	return chatgpt_request
}
