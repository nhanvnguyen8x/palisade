package transport

import "github.com/nhanvnguyen8x/palisade/internal/features/chat/application"

func ToChatResponse(answer *application.ChatAnswer) ChatResponse {
	return ChatResponse{
		Answer:   answer.Answer,
		Contexts: answer.Contexts,
	}
}

