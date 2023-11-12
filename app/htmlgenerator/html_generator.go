package htmlgenerator

import (
	"embed"
	"html/template"
	"io"
)

// HTMLテンプレートファイルをビルド時に埋め込む
//
//go:embed translate_result.html
var embedFileSystem embed.FS

// 生成内容を表す構造体
type HTMLGenerator struct {
	// 翻訳対象の文字列
	Target string
	// 翻訳結果の文字列
	Translated string
	// HTMLファイルの書き込みを担当するWriter
	Writer *io.Writer
}

// NewHTMLGenerator は HTMLGenerator を生成する
func NewHTMLGenerator(target, translated string, writer io.Writer) HTMLGenerator {
	return HTMLGenerator{
		Target:     target,
		Translated: translated,
		Writer:     &writer,
	}
}

// Generate は HTMLGenerator が保持する情報を元にHTMLを生成する
func (g HTMLGenerator) Generate() error {
	// テンプレートオブジェクトを生成する
	t := template.New("translateResult")

	// 埋め込んだHTMLテンプレート定義を読み込む
	templateFile, _ := embedFileSystem.ReadFile("translate_result.html")

	// テンプレートオブジェクトにHTMLテンプレートのパース結果を上書きする
	// tへの副作用が発生する
	t, err := t.Parse(string(templateFile))
	if err != nil {
		return err
	}

	// Writerの向き先にテンプレートを元にしたHTMLを生成する
	// Writerへの副作用が発生する
	err = t.Execute(*g.Writer, g)
	if err != nil {
		return err
	}

	// 正常終了
	return nil
}
