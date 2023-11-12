package translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/toriwasa/translate/go/domain/models/chat/completion"
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

// モデル応答を生成するための構造体を生成する
func newCreate(text string) completion.Create {
	// ChatGPT APIに翻訳を依頼する文字列を生成する
	q := fmt.Sprintf("Translate this into Japanese. \n\n%s", text)
	log.Printf("q: %s", q)

	// モデル応答を生成するための構造体を生成する
	create := completion.Create{
		Model: "gpt-3.5-turbo-1106",
		Messages: []completion.Messages{
			{
				Role:    "user",
				Content: q,
			},
		},
	}

	// リクエスト用の構造体を返す
	return create
}

// Chat Completion APIリクエスト用の http.Request を生成する
func newRequest(create completion.Create, apiKey string) (*http.Request, error) {
	// モデル応答生成用の構造体をバイト型JSONに変換する
	jsonBytes, err := json.Marshal(create)
	if err != nil {
		return nil, err
	}

	// JSONをバイト文字列としてリクエストボディに設定する
	body := bytes.NewBuffer(jsonBytes)

	// Chat Completion APIのエンドポイントへのPOSTリクエストオブジェクトを生成する
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", body)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}

	// ヘッダーに Content-Type: application/json を設定する
	req.Header.Set("Content-Type", "application/json")

	// ヘッダーに Authorization: Bearer apiKey を設定する
	authToken := fmt.Sprintf("Bearer %s", apiKey)
	req.Header.Set("Authorization", authToken)

	// HTTPリクエストを返す
	return req, nil
}

// Chat Completion APIへのリクエストを実行する
func createChatCompletionHttpRequest(req *http.Request) (*http.Response, error) {
	// HTTPリクエストを実行する
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("resp is nil")
	}
	return resp, nil
}

// APIへのリクエスト結果JSONを構造体に変換する
func convertResponseToStruct(resp *http.Response) (completion.Result, error) {
	// HTTPレスポンスを読み込む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return completion.Result{}, err
	}

	// ステータスコードが 200 以外の場合はエラーを返す
	if resp.StatusCode != http.StatusOK {
		// エラーレスポンスJSONを構造体に変換する
		var completionError completion.CompletionError
		err = json.Unmarshal(body, &completionError)
		if err != nil {
			return completion.Result{}, err
		}
		// エラー構造体を返す
		return completion.Result{}, completionError
	}

	// ステータスコードが 200 の場合は正常レスポンスJSONを構造体に変換する
	var result completion.Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		return completion.Result{}, err
	}

	// 構造体を返す
	return result, nil
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
	create := newCreate(text)
	httpReq, error := newRequest(create, apiKey)
	if err != nil {
		return "", error
	}

	// ChatGPT APIにPOSTリクエストする
	res, err := createChatCompletionHttpRequest(httpReq)
	if err != nil {
		return "", err
	}

	// リクエスト結果を構造体に変換する
	result, err := convertResponseToStruct(res)
	if err != nil {
		return "", err
	}

	// 翻訳結果を返す
	return result.Choices[0].Message.Content, nil

}
