package htmlgenerator

import (
	"html/template"
	"io"
)

// HTMLGenerator は html/template パッケージを使ってHTMLを生成する

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

	// テンプレート定義ファイルパスを定義する
	templateFilePath := "app/htmlgenerator/translate_result.html"

	// テンプレートオブジェクトにHTMLテンプレートのパース結果を上書きする
	// tへの副作用が発生する
	t, err := t.ParseFiles(templateFilePath)
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
