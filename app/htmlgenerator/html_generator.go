package htmlgenerator

import (
	"embed"
	"html/template"
	"io"
)

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
	// 基礎となるテンプレートオブジェクトを生成する
	baseTemplate := template.New("translateResult")

	// 埋め込みファイルシステム経由でHTMLテンプレート定義を読み込む
	templateFile, err := embedFileSystem.ReadFile("translate_result.html")
	if err != nil {
		return err
	}

	// HTMLテンプレート定義を元にした新しいテンプレートオブジェクトを生成する
	t, err := baseTemplate.Parse(string(templateFile))
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
