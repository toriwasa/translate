package translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/toriwasa/translate/domain/models/chat/completion"
)

// 環境変数 OPENAI_API_KEY から ChatGPT API の APIキーを取得する
// APIキーが取得できない場合はエラーを返す
func getAPIKey() (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("env: OPENAI_API_KEY is empty")
	}
	return apiKey, nil
}

// ChatpGPT APIにPOSTリクエストする
func requestChatGPT(text string) (string, error) {
	// ChatGPT APIに翻訳を依頼する文字列を生成する
	q := fmt.Sprintf("Translate this into Japanese. \n\n %s", text)
	// 依頼文字列をログに出力する
	log.Printf("q: %s", q)

	// 以下のJSONをバイト文字列としてリクエストボディに設定する
	// {"model": "gpt-3.5-turbo-1106", "messages": [{"role": "user", "content": q}] }
	jsonBytes, err := json.Marshal(completion.Request{
		Model: "gpt-3.5-turbo-1106",
		Messages: []completion.Messages{
			{
				Role:    "user",
				Content: q,
			},
		},
	})
	if err != nil {
		return "", err
	}

	// https://api.openai.com/v1/chat/completions エンドポイントにPOSTリクエストする
	// HTTPリクエストを作成する
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}

	// ヘッダーに Content-Type: application/json を設定する
	req.Header.Set("Content-Type", "application/json")

	// ヘッダーに Authorization: Bearer apiKey を設定する
	apiKey, err := getAPIKey()
	if err != nil {
		return "", err
	}
	authToken := fmt.Sprintf("Bearer %s", apiKey)
	req.Header.Set("Authorization", authToken)

	// HTTPリクエストを実行する
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// ステータスコードをログに出力する
	log.Printf("status: %s", resp.Status)

	// HTTPレスポンスを読み込む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// レスポンスボディをログに出力する
	log.Printf("body: %s", body)

	// ステータスコードが 200 以外の場合はエラーを返す
	if resp.StatusCode != http.StatusOK {
		// レスポンスJSONを構造体に変換する
		var result completion.Error
		err = json.Unmarshal(body, &result)
		if err != nil {
			return "", err
		}

		// 構造体の中身をログに出力する
		log.Printf("errorResult: %+v", result)

		// 構造体に含まれるエラーメッセージを返す
		return "", fmt.Errorf("errorResult: %+v", result)
	}

	// レスポンスJSONを意図した結果の構造体に変換する
	var result completion.Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	// 構造体の中身をログに出力する
	log.Printf("chatCompletion: %+v", result)

	// 構造体に含まれる翻訳結果を返す
	return result.Choices[0].Message.Content, nil

}

// Translate は引数で与えられた文字列を日本語に翻訳して返す
// 翻訳に失敗した場合はエラーを返す
func Translate(text string) (string, error) {
	return requestChatGPT(text)
}
