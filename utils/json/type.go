package json

// to json struct
type Thinking struct {
	Type string `json:"type"`
}

type ChatReq struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	Thinking  Thinking  `json:"thinking"`
	Stream    bool      `json:"stream"`
	MaxTokens int       `json:"max_tokens"`
	Temp      float32   `json:"temperature"`
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
	// ReasoningContent string `json:"reasoning_content"`
}

// parse json struct
type ChatResponse struct {
	Choices []Choice `json:"choices"`
	Created int64    `json:"created"`
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	Message      Message `json:"message"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// type Message struct {
// 	Role    string `json:"role"`
// 	Content string `json:"content"`
// }
