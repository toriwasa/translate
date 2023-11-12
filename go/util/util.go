package util

import (
	"os"

	"github.com/atotto/clipboard"
	"github.com/toriwasa/translate/go/app/translate"
)

// Dry Run モードを考慮して翻訳結果を生成する
func GenerateTranslated(target string, isDryRun bool) (string, error) {
	if isDryRun {
		return "dry run モード 翻訳結果", nil
	}

	// 翻訳結果を生成する
	translated, err := translate.Translate(target)
	if err != nil {
		return "", err
	}
	return translated, nil
}

// 一時ファイルを生成してファイルオブジェクトを返す
func CreateTempFile() (*os.File, error) {
	// 一時ディレクトリパスを取得する
	tempDir := os.TempDir()

	// 一時ファイルを生成する
	tempFile, err := os.CreateTemp(tempDir, "translated*.html")
	if err != nil {
		return nil, err
	}
	return tempFile, nil
}

// クリップボードの内容を取得する
func GetClipBoardText() (string, error) {
	// クリップボードの内容を取得する
	clipBoardText, err := clipboard.ReadAll()
	if err != nil {
		return "", err
	}
	return clipBoardText, nil
}
