package chat

import (
	"Backend/config"
	"Backend/validator"
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/openai"
)

type IService interface {
	getTitles(string) ([]Chat, error)
	getChatHistory(int32) ([]llms.MessageContent, error)
	processOutput(string, int32, string, string, string, string) (string, error)
	generateTitle(string, string, string, string, string) (int32, string, error)
	deleteChat(string, int32) error
	checkInput(string, string, string) (bool, map[string]string)
}

type service struct {
	chatRepo repo
	multiLLM *config.MultiLLM
}

func NewService(chatRepo repo, multiLLM *config.MultiLLM) IService {
	return &service{
		chatRepo: chatRepo,
		multiLLM: multiLLM,
	}
}

func (s *service) getTitles(userID string) ([]Chat, error) {

	return s.chatRepo.getTitles(userID)
}

func (s *service) getChatHistory(chatID int32) ([]llms.MessageContent, error) {
	return s.chatRepo.getMessageHistory(chatID)
}

func withAPIKey(modelType string, apiKey string) (llms.Model, error) {
	if modelType == "OpenAI" {
		return openai.New(openai.WithToken(apiKey))
	} else if modelType == "Google" {
		return googleai.New(context.Background(), googleai.WithAPIKey(apiKey))
	} else if modelType == "Anthropic" {
		return anthropic.New(anthropic.WithToken(apiKey))
	}
	return nil, fmt.Errorf("invalid model type: %s", modelType)
}

func withoutAPIKey(modelType string, multiLLM *config.MultiLLM) llms.Model {
	if modelType == "OpenAI" {
		return multiLLM.OpenAI
	} else if modelType == "Google" {
		return multiLLM.GoogleAI
	} else if modelType == "Anthropic" {
		return multiLLM.AntAI
	}
	return nil
}

func getModel(modelType string, apiKey string, multiLLM *config.MultiLLM) (llms.Model, error) {
	if apiKey != "" {
		opt, err := withAPIKey(modelType, apiKey)
		if err != nil {
			return nil, err
		}
		return opt, nil
	}
	return withoutAPIKey(modelType, multiLLM), nil
}

func (s *service) generateTitle(userID string, modelType string, modelName string, apiKey string, prompt string) (int32, string, error) {
	option, err := getModel(modelType, apiKey, s.multiLLM)
	if err != nil {
		return 0, "", err
	}

	var opts []llms.CallOption
	if modelName != "" {
		opts = append(opts, llms.WithModel(modelName))
	}

	titlePrompt := fmt.Sprintf(
		"Based on the following initial prompt, generate a concise and descriptive title for the conversation:\n\n%s",
		prompt,
	)
	titles, err := option.GenerateContent(context.Background(), []llms.MessageContent{llms.TextParts(llms.ChatMessageTypeHuman, titlePrompt)}, opts...)
	if err != nil {
		return 0, "", err
	}

	if userID == "" {
		return 0, titles.Choices[0].Content, nil
	}

	return s.chatRepo.insertTitle(userID, titles.Choices[0].Content)
}

func (s *service) processOutput(userID string, chatID int32, modelType string, modelName string, apiKey string, prompt string) (string, error) {
	option, err := getModel(modelType, apiKey, s.multiLLM)
	if err != nil {
		return "", err
	}

	conversation, err := s.chatRepo.getMessageHistory(chatID)
	if err != nil {
		return "", err
	}

	conversation = append(conversation, llms.TextParts(llms.ChatMessageTypeHuman, prompt))
	var opts []llms.CallOption
	if modelName != "" {
		opts = append(opts, llms.WithModel(modelName))
	}

	content, err := option.GenerateContent(context.Background(), conversation, opts...)
	if err != nil {
		return "", err
	}

	if userID != "" {
		if err := s.chatRepo.insertLatestMessage(chatID, prompt, content.Choices[0].Content); err != nil {
			return "", err
		}
	}

	return content.Choices[0].Content, nil
}

func (s *service) deleteChat(userID string, chatID int32) error {
	return s.chatRepo.deleteChat(userID, chatID)
}

func (s *service) checkInput(modelType string, model string, prompt string) (bool, map[string]string) {
	v := validator.New()

	v.Check(prompt != "", "prompt", "Empty prompt")
	v.Check(validator.In(modelType, "OpenAI", "Google", "Anthropic"), "modelType", "Invalid model type")
	v.Check(validator.In(model, append(config.LLMLists[modelType], "")...), "model", "Invalid model")

	return v.Valid(), v.Errors
}
