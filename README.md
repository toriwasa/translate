# README

## What's This?
- ChatGPT API にアクセスして、指定した文字列を日本語に翻訳した結果を返すCLIツール
- 有効なChatGPTのAPIキーが必要

## Usage
- とりあえずLinuxでの実行方法

```bash
# 環境変数 OPENAI_API_KEY にAPIキーを指定する
export OPENAI_API_KEY=sk-xxxxxxxxxx

# 翻訳して欲しい文字列を指定する
./translate -t "Hello, World!"
translated: こんにちは、世界！
```
