package completion

// ChatGPの Chat Completion API レスポンスを表す構造体
// 以下のJSONを元に構造体を定義する
// {
//   "id": "chatcmpl-xxxxxxxxxxxxxxxx",
//   "object": "chat.completion",
//   "created": 1699695973,
//   "model": "gpt-3.5-turbo-1106",
//   "choices": [
//     {
//       "index": 0,
//       "message": {
//         "role": "assistant",
//         "content": "こんにちは、世界！今朝はどうですか？"
//       },
//       "finish_reason": "stop"
//     }
//   ],
//   "usage": {
//     "prompt_tokens": 23,
//     "completion_tokens": 15,
//     "total_tokens": 38
//   },
//   "system_fingerprint": "fp_eeff13170a"
// }

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Result struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}
