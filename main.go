package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/toriwasa/translate/app/htmlgenerator"
	"github.com/toriwasa/translate/app/translate"
	"github.com/toriwasa/translate/infrastructure/webview2viewer"
)

func main() {
	// DEBUGログのフォーマットを設定
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// DEBUGログのプレフィックスを設定
	log.SetPrefix("DEBUG: ")

	// コマンドライン引数を解析する。 -t, -v オプションを定義する
	var t string
	var isVerbose bool

	// t は 翻訳対象の文字列を表す
	// デフォルト値は "Hello, World!" である
	flag.StringVar(&t, "t", "Hello, World!", "translate target string")

	// v はログを冗長に出力するモードを表す
	// デフォルト値は false である
	flag.BoolVar(&isVerbose, "v", false, "output verbose log")

	// --help オプションをカスタマイズする
	flag.Usage = func() {
		println("Usage: translate -t <translate target text>")
		println("Example: translate -t \"Target Text\"")
		println("Description: translate Text to Japanese")
		println("Options:")
		flag.PrintDefaults()
	}

	// コマンドライン引数を解析する
	flag.Parse()

	// verbose モードでない場合はログを出力しない
	if !isVerbose {
		log.SetOutput(io.Discard)
	}

	// コマンドライン引数を出力する
	log.Printf("t: %s, isVerbose: %t", t, isVerbose)

	// 翻訳する
	translated, err := translate.Translate(t)
	if err != nil {
		panic(err)
	}

	// 一時ディレクトリパスを取得する
	tempDir := os.TempDir()

	// 一時ファイルを生成する
	// webview2でHTMLとして表示するために拡張子を .html にする
	tempFile, err := os.CreateTemp(tempDir, "translated*.html")
	if err != nil {
		panic(err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	// HTMLGeneratorを生成する
	g := htmlgenerator.NewHTMLGenerator(t, translated, tempFile)

	// HTMLGeneratorが保持する情報を元にHTMLを生成する
	err = g.Generate()
	if err != nil {
		panic(err)
	}

	// 一時ファイルのパスを取得する
	tempFilePath := tempFile.Name()
	log.Printf("tempFilePath: %s", tempFilePath)

	// 一時ファイルを閲覧する
	err = webview2viewer.OpenFileWithWebview2(tempFilePath)
	if err != nil {
		panic(err)
	}
}
