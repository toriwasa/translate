package webview

import "github.com/jchv/go-webview2"

// webview2の管理用の構造体
type WebViewManager struct {
	// WebView2オブジェクト
	w       webview2.WebView
	isAlive bool
}

// NewWebViewManager はWebViewManagerを生成する
func NewWebViewManager(w webview2.WebView) WebViewManager {
	return WebViewManager{
		w:       w,
		isAlive: true,
	}
}

// webview2を終了するメソッド
func (m *WebViewManager) Kill() {
	m.w.Destroy()
	m.isAlive = false
}

// isAliveを参照するメソッド
func (m *WebViewManager) IsAlive() bool {
	return m.isAlive
}
