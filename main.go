package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/toriwasa/translate/app/translate"
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
		log.Fatal(err)
	}

	// 翻訳結果を標準出力する
	fmt.Printf("translated: %s\n", translated)

}
