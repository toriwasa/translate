package util

import (
	"log"
	"os"

	"github.com/atotto/clipboard"
	"github.com/toriwasa/translate/app/htmlgenerator"
	"github.com/toriwasa/translate/app/translate"
	"github.com/toriwasa/translate/infrastructure/webview2viewer"
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

// 翻訳結果をwebview2で表示する
func ShowTranslatedWithWebview2(target, translated string) error {
	// 一時ディレクトリパスを取得する
	tempDir := os.TempDir()

	// 一時ファイルを生成する
	// webview2でHTMLとして表示するために拡張子を .html にする
	tempFile, err := os.CreateTemp(tempDir, "translated*.html")
	if err != nil {
		return err
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	// HTMLGeneratorを生成する
	g := htmlgenerator.NewHTMLGenerator(target, translated, tempFile)

	// HTMLGeneratorが保持する情報を元にHTMLを生成する
	err = g.Generate()
	if err != nil {
		return err
	}

	// 一時ファイルのパスを取得する
	tempFilePath := tempFile.Name()
	log.Printf("tempFilePath: %s", tempFilePath)

	// 一時ファイルを閲覧する
	err = webview2viewer.OpenFileWithWebview2(tempFilePath)
	if err != nil {
		return err
	}
	return nil
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
