# Weather Notification

本プロジェクトは、OpenWeatherMap API を利用して天気予報を取得し、  
雨の予報がある場合に LINE メッセージで通知する Go プログラムです。  
API キーや位置情報、LINE アクセストークンは環境変数で安全に管理しています。

## 機能概要

- OpenWeatherMap API から本日の 24 時間以内の天気予報を取得
- 雨が予想される場合のみ LINE に通知
- 取得データは日本語で出力
- 環境変数によるセキュリティ管理
- エラーハンドリング・詳細なログ出力

## ディレクトリ構成

```md
.
├── src/
│ ├── main.go
│ ├── weather/
│ │ └── weather.go
│ └── line/
│ └── line.go
├── README.md
├── .gitignor
├── go.mod
├── go.sum
└── .env (Git 管理対象外推奨)
```

## 必要要件

- Go 1.18 以上
- OpenWeatherMap API キー
- LINE Messaging API アクセストークン

## 環境変数

| 変数名                    | 説明                                        | 例         |
| ------------------------- | ------------------------------------------- | ---------- |
| OPENWEATHER_API_KEY       | OpenWeatherMap の API キー                  | xxxxxxxx   |
| WEATHER_LAT               | 天気を取得する緯度                          | 35.682839  |
| WEATHER_LON               | 天気を取得する経度                          | 139.759455 |
| LINE_CHANNEL_ACCESS_TOKEN | LINE Messaging API チャネルアクセストークン | xxxxxxxx   |
| MY_USER_ID                | 通知を送る LINE ユーザー ID                 | xxxxxxxx   |

## 使い方

1. OpenWeatherMap に登録して API キーを取得  
   [https://openweathermap.org/api](https://openweathermap.org/api)

2. LINE Developers コンソールで Messaging API を設定し、チャネルアクセストークンを取得  
   [https://developers.line.biz/console/](https://developers.line.biz/console/)

3. `.env`ファイルを作成して環境変数をセットするか、もしくは直接エクスポート

例：

```sh
export OPENWEATHER_API_KEY=あなたのAPIキー
export WEATHER_LAT=35.682839
export WEATHER_LON=139.759455
export LINE_CHANNEL_ACCESS_TOKEN=あなたのチャネルアクセストークン
export MY_USER_ID=あなたのLINEユーザーID
```
