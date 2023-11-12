package main

import (
	"flag"
	"io"
	"log"

	"github.com/toriwasa/translate/util"
)

func main() {
	// DEBUGログのフォーマットを設定
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// DEBUGログのプレフィックスを設定
	log.SetPrefix("DEBUG: ")

	// コマンドライン引数を解析する。 -t, -c, -d, -v オプションを定義する
	var t string
	var useClipBoard, isDryRun, isVerbose bool

	// t は 翻訳対象の文字列を表す
	// デフォルト値は "Hello, World!" である
	flag.StringVar(&t, "t", "Hello, World!", "translate target string")

	// c は クリップボードの内容を翻訳対象の文字列として利用するモードを表す
	// デフォルト値は false である
	flag.BoolVar(&useClipBoard, "c", false, "use clipboard text as translate target string")

	// d は dry run モードを表す
	// デフォルト値は false である
	flag.BoolVar(&isDryRun, "d", false, "execute dry run mode")

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
	log.Printf("t: %s, useClipBoard: %t isDryRun: %t isVerbose: %t", t, useClipBoard, isDryRun, isVerbose)

	// ここからメイン処理
	targetText := t
	// useClipBoard モードの場合はクリップボードの内容を翻訳対象の文字列として利用する
	if useClipBoard {
		// クリップボードの内容を取得する
		clipBoardText, err := util.GetClipBoardText()
		if err != nil {
			panic(err)
		}
		targetText = clipBoardText
	}

	// DryRunモードを考慮して翻訳結果を生成する
	translated, err := util.GenerateTranslated(targetText, isDryRun)
	if err != nil {
		panic(err)
	}

	// 翻訳結果をwebview2で表示する
	err = util.ShowTranslatedWithWebview2(targetText, translated)
	if err != nil {
		panic(err)
	}

}
