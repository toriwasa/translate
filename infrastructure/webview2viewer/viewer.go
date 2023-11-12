package webview2viewer

import (
	"log"

	"github.com/jchv/go-webview2"
)

// 指定されたファイルをwebview2で開く
func OpenFileWithWebview2(filePath string) error {
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
		log.Fatalln("Failed to load webview.")
	}
	// 関数を抜けたらWebView2オブジェクトを破棄する
	defer w.Destroy()

	// WebView2オブジェクトのウィンドウサイズを設定する
	w.SetSize(800, 600, webview2.HintFixed)

	// 指定されたファイルパスをURLに変換する
	url := "file:///" + filePath

	// WebView2オブジェクトで指定されたURLを開く
	w.Navigate(url)

	// WebView2のウィンドウを閉じる関数をJavaScriptに公開する
	w.Bind("closeWebView", func() {
		w.Destroy()
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

	// WebView2オブジェクトを起動する
	w.Run()

	return nil
}
