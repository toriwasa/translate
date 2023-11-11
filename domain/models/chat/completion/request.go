package completion

// 以下のJSONを元に構造体を定義する
// {"model": "gpt-3.5-turbo-1106", "messages": [{"role": "user", "content": q}] }
type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string     `json:"model"`
	Messages []Messages `json:"messages"`
}
