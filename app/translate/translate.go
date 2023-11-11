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

// リクエスト用の構造体を生成する
func newRequest(text string) completion.Request {
	// ChatGPT APIに翻訳を依頼する文字列を生成する
	q := fmt.Sprintf("Translate this into Japanese. \n\n%s", text)
	log.Printf("q: %s", q)

	// リクエスト用の構造体を生成する
	req := completion.Request{
		Model: "gpt-3.5-turbo-1106",
		Messages: []completion.Messages{
			{
				Role:    "user",
				Content: q,
			},
		},
	}

	// リクエスト用の構造体を返す
	return req
}

// APIリクエスト用の http.Request を生成する
func newHTTPRequest(req completion.Request, apiKey string) (*http.Request, error) {
	// リクエスト用の構造体をJSONに変換する
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// JSONをバイト文字列としてリクエストボディに設定する
	body := bytes.NewBuffer(jsonBytes)

	// Chat Completion APIのエンドポイントへのPOSTリクエストオブジェクトを生成する
	httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", body)
	if err != nil {
		return nil, err
	}

	// ヘッダーに Content-Type: application/json を設定する
	httpReq.Header.Set("Content-Type", "application/json")

	// ヘッダーに Authorization: Bearer apiKey を設定する
	authToken := fmt.Sprintf("Bearer %s", apiKey)
	httpReq.Header.Set("Authorization", authToken)

	// HTTPリクエストを返す
	return httpReq, nil
}

// ChatpGPT APIにPOSTリクエストする
func requestChatGPT(req *http.Request) (string, error) {

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
	// APIキーを取得する
	apiKey, err := getAPIKey()
	if err != nil {
		return "", err
	}

	// リクエスト用の構造体を生成する
	req := newRequest(text)
	httpReq, error := newHTTPRequest(req, apiKey)
	if err != nil {
		return "", error
	}

	// ChatGPT APIにPOSTリクエストする
	translated, err := requestChatGPT(httpReq)
	if err != nil {
		return "", err
	}

	// 翻訳結果を返す
	return translated, nil

}
