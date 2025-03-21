package chatgpt

import (
	"freechatgpt/typings"
	chatgpt_types "freechatgpt/typings/chatgpt"
	official_types "freechatgpt/typings/official"
	"strings"
)

func ConvertToString(chatgpt_response *chatgpt_types.ChatGPTResponse, previous_text *typings.StringStruct, role bool) string {
    currentContent := chatgpt_response.Message.Content.Parts[0].(string)
    
    if role {
        translated_response := official_types.NewChatCompletionChunk("")
        translated_response.Choices[0].Delta.Role = chatgpt_response.Message.Author.Role
        previous_text.Text = currentContent
        return "data: " + translated_response.String() + "\n\n"
    } else {
        newContent := strings.Replace(currentContent, previous_text.Text, "", 1)
        if newContent == "" {
            return ""
        }
        
        translated_response := official_types.NewChatCompletionChunk(newContent)
        previous_text.Text = currentContent
        return "data: " + translated_response.String() + "\n\n"
    }
}
