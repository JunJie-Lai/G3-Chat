package config

import (
	"context"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/openai"
	"os"
)

var LLMLists = map[string][]string{
	"OpenAI": {
		"gpt-4.1-nano",
		"gpt-4.1-mini",
		"gpt-4.1",
		"gpt-4o",
		"gpt-4o-mini",
		"o4-mini",
		"o3",
		"o3-mini",
		"o3-pro",
		"gpt-4.5-preview",
	},
	"Google": {
		"gemini-2.5-flash-preview-05-20",
		"gemini-2.5-pro-preview-06-05",
		"gemini-2.0-flash",
		"gemini-2.0-flash-lite",
	},
	"Anthropic": {
		"claude-sonnet-4-0",
		"claude-opus-4-0",
		"claude-3-7-sonnet-latest",
		"claude-3-5-sonnet-latest",
	},
}

type MultiLLM struct {
	OpenAI   *openai.LLM
	GoogleAI *googleai.GoogleAI
	AntAI    *anthropic.LLM
}

func NewAI() (*MultiLLM, error) {
	openAI, err := openai.New(openai.WithToken(os.Getenv("OPENAI_API_KEY")), openai.WithModel(LLMLists["OpenAI"][0]))
	if err != nil {
		return nil, err
	}

	googleAI, err := googleai.New(context.Background(), googleai.WithAPIKey(os.Getenv("GEMINI_API_KEY")), googleai.WithDefaultModel(LLMLists["Google"][0]))
	if err != nil {
		return nil, err
	}

	antAI, err := anthropic.New(anthropic.WithToken(os.Getenv("ANTHROPIC_API_KEY")), anthropic.WithModel(LLMLists["Anthropic"][0]))
	if err != nil {
		return nil, err
	}

	return &MultiLLM{
		OpenAI:   openAI,
		GoogleAI: googleAI,
		AntAI:    antAI,
	}, nil
}
