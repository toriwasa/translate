package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/jchv/go-webview2"
	"github.com/toriwasa/translate/app/htmlgenerator"
	"github.com/toriwasa/translate/infrastructure/webview"
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

	// 結果待機用の一時ファイルを生成する
	loadingHTMLFile, err := util.CreateTempFile()
	if err != nil {
		panic(err)
	}
	defer loadingHTMLFile.Close()
	defer os.Remove(loadingHTMLFile.Name())
	log.Printf("loadingHTMLFile: %s", loadingHTMLFile.Name())

	// 結果待機用の htmlgenerator を生成する
	loadingHTMLGenerator := htmlgenerator.NewHTMLGenerator(targetText, "翻訳結果を取得中...", loadingHTMLFile)

	// HTMLGeneratorが保持する情報を元にHTMLを生成する
	err = loadingHTMLGenerator.Generate()
	if err != nil {
		panic(err)
	}

	// 結果待機用のWebview2を起動する
	// WebView2オブジェクトを生成する
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     true,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "Translate Result",
			Width:  800,
			Height: 600,
			IconId: 2, // icon resource id
			Center: true,
		},
	})
	if w == nil {
		panic("failed to load webview")
	}
	manager := webview.NewWebViewManager(w)

	// WebView2オブジェクトのウィンドウサイズを変更不可能にする設定を追加する
	// w.SetSize(800, 600, webview2.HintFixed)

	// WebView2オブジェクトで指定されたHTMLを開く
	w.Navigate(loadingHTMLFile.Name())

	// WebView2のウィンドウを閉じる関数をJavaScriptに公開する
	w.Bind("closeWebView", func() {
		manager.Kill()
	})

	// ESC, Enter, Spaceキーでウィンドウを閉じる機能を追加する
	w.Init(`
	document.addEventListener("keydown", function(event) {
		if (event.keyCode === 27 // ESC
			|| event.keyCode === 13  // Enter
			|| event.keyCode === 32 // Space
		) {
			closeWebView(); // Goの関数を呼び出す
		}
	});
	`)

	// 別スレッドで翻訳処理を実行する
	go func() {
		// 翻訳結果表示用の一時ファイルを生成する
		translatedHTMLFile, err := util.CreateTempFile()
		if err != nil {
			panic(err)
		}
		defer translatedHTMLFile.Close()
		defer os.Remove(translatedHTMLFile.Name())
		log.Printf("translatedHTMLFile: %s", translatedHTMLFile.Name())

		// DryRunモードを考慮して翻訳結果を生成する
		translated, err := util.GenerateTranslated(targetText, isDryRun)
		if err != nil {
			panic(err)
		}

		// 翻訳結果用の htmlgenerator を生成する
		resultHTMLGenerator := htmlgenerator.NewHTMLGenerator(targetText, translated, translatedHTMLFile)
		// HTMLGeneratorが保持する情報を元にHTMLを生成する
		err = resultHTMLGenerator.Generate()
		if err != nil {
			panic(err)
		}

		// WebView2で翻訳結果HTMLに遷移する
		// WebView2オブジェクトが終了している場合は何もしない
		if !manager.IsAlive() {
			return
		}
		w.Dispatch(func() {
			w.Navigate(translatedHTMLFile.Name())
		})
	}()

	// WebView2オブジェクトを起動する
	w.Run()

}
