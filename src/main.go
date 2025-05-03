package main

import (
	"fmt"
	"log"

	"weather-notification/src/line"
	"weather-notification/src/weather"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 環境変数を取得
	apiKey, lat, lon, err := weather.GetWeatherEnvVars()
	if err != nil {
		log.Fatalf("天気API環境変数取得エラー: %v", err)
	}
	token, userID, err := line.GetLineEnvVars()
	if err != nil {
		log.Fatalf("LINE環境変数取得エラー: %v", err)
	}

	// 天気予報を取得
	forecast, err := weather.FetchWeather(apiKey, lat, lon)
	if err != nil {
		log.Fatalf("天気情報取得エラー: %v", err)
	}

	// 傘が必要か判定
	needUmbrella, rainMessage := weather.NeedUmbrella(forecast)

	var message string
	if needUmbrella {
		message = rainMessage
	} else {
		message = "今日は傘の必要はなさそうです☀️"
	}

	// メッセージを送信
	if err := line.SendTextMessage(token, userID, message); err != nil {
		log.Fatalf("LINEメッセージ送信エラー: %v", err)
	}
	fmt.Println("天気通知をLINEに送信しました。")
}
