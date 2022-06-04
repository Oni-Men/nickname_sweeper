# ニックネーム除去Bot

サーバーに

## Botの導入

[ここ](https://discord.com/developers/applications)で新しいBotを作成する。

Botタブで **SERVER MEMBERS INTENT** を有効にする。

OAuth2 URL Generator で以下にチェック。

### SCOPE
- bot
- applications.commands

### Bot permission
- Manage Nicknames
- Send Messages
- Use Slash Commands

下の方で生成されたURLをコピーしてブラウザでアクセスし、あとは指示に従ってサーバーに導入する。

## Botの起動

Botタブに表示されているTokenをコピーし、コマンドライン引数で以下のようにして与える。

`./bot.exe -token="<your token>"`

コマンドを利用したいサーバーのIDを以下のどちらかの方法で取得し、
  - 開発者設定を有効にしてるとき、サーバー名を右クリックして「IDをコピー」
  - そうでないとき、サーバー設定/ウィジェットからサーバーIDをコピーする

コマンドライン引数で以下のようにして与える。

`./bot.exe -guild="<your server id>`

上述の例でもあるようにサーバーの起動は以下のようにして行う

`./bot.exe -token="<token>" -guild="<guild>"`

## 使用方法

Botがオンライン状態になれば起動が成功している。

その状態でチャット欄に「/」を打ち込むとコマンド一覧が表示される。

そのまま「/sweep-nickname」と打ち込んでエンターで送信するとコマンドが実行される。

### オプションについて

高度なオプションはめんどくさいので実装してない。

現在、「ロール付きユーザーを除外対象にする」オプションのみが利用できる。

このオプションを有効にして実行すると、ロールのついているユーザーからはニックネームを削除しない。

