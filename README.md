# Weather Notification

本プロジェクトは、OpenWeatherMap API を利用して指定都市の今日の天気予報を取得する Go プログラムです。  
API キーや都市名は環境変数で管理し、セキュリティベストプラクティスに従っています。

## 機能概要

- OpenWeatherMap API から本日の天気予報（天気、気温、湿度）を取得
- 日本語での出力
- API キーや都市名は環境変数から取得
- エラーハンドリング・詳細なログ出力

## ディレクトリ構成

```
.
├── src/
│   └── weather_today.go
├── README.md
├── .clinerules
├── .gitignore
├── LICENSE
```

## 必要要件

- Go 1.18 以上
- OpenWeatherMap API キー

## 環境変数

| 変数名              | 説明                       | 例       |
| ------------------- | -------------------------- | -------- |
| OPENWEATHER_API_KEY | OpenWeatherMap の API キー | xxxxxxxx |
| WEATHER_CITY        | 天気を取得する都市名       | Tokyo    |

## 使い方

1. OpenWeatherMap で API キーを取得  
   https://openweathermap.org/api

2. 環境変数を設定し、プログラムを実行

```sh
export OPENWEATHER_API_KEY=あなたのAPIキー
export WEATHER_CITY=Tokyo
go run src/weather_today.go
```

3. 出力例

```
【Tokyo の天気予報 (2025-04-27 15:00)】
天気: 曇り (曇りがち)
気温: 18.5℃
湿度: 60%
```

## セキュリティ

- API キーや都市名などの機密情報は、必ず環境変数で管理してください
- .env ファイルや API キーを含むファイルは Git 管理対象外としてください

## ライセンス

本リポジトリは MIT ライセンスです。
