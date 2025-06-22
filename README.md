# Todoアプリケーション
## 概要
Goを用いて開発した、シンプルなTodo管理アプリケーションです。
Todoの一覧表示や詳細編集、削除機能などがあります。

## 動作イメージ
![動作イメージ動画](./docks/todo_go_app.gif)

# 技術スタック
* 言語：Go
  * フロント側もTemplateで作成しています（一部JavaScriptあり）。
* Webフレームワーク：gin
* DB：MySQL
* ORM：GORM
* その他
  * gomock
  * txdb
  * testify
  * Air
  * Delve

## 解説記事
**[Todoアプリを作ってGoへ入門！：シリーズまとめ](https://qiita.com/MinadukiSekina/items/3f6e4420cd493703bb73)**

## 環境構築
DevContainerを使用しています。`db/db.env`を作成して下記の環境変数を記入した後、「Reopen in Container」を実行してください。

```text
MYSQL_ROOT_PASSWORD={お好みのパスワード}
MYSQL_DATABASE={お好みの名称}
TEST_DATABASE={お好みの名称}
MYSQL_USER={お好みのユーザー名}
MYSQL_PASSWORD={お好みのパスワード}
```

## 動作確認
DevContainerを立ち上げた後、http://localhost:3000 にアクセスしてください。

## テストについて
`make gotest`を実行してください。