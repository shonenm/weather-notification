// main.go
package main

import (
	"fmt"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 環境変数を取得
	apiKey, lat, lon, err := GetWeatherEnvVars()
	if err != nil {
		log.Fatalf("天気API環境変数取得エラー: %v", err)
	}
	token, userID, err := GetLineEnvVars()
	if err != nil {
		log.Fatalf("LINE環境変数取得エラー: %v", err)
	}

	// 天気予報を取得
	forecast, err := FetchWeather(apiKey, lat, lon)
	if err != nil {
		log.Fatalf("天気情報取得エラー: %v", err)
	}

	// 雨が降るか判定
	needUmbrella, rainMessage := NeedUmbrella(forecast)

	if needUmbrella {
		// 雨情報を送信
		if err := SendTextMessage(token, userID, rainMessage); err != nil {
			log.Fatalf("LINEメッセージ送信エラー: %v", err)
		}
		fmt.Println("雨予報メッセージをLINEに送信しました。")
	} else {
		fmt.Println("今日は傘の必要はなさそうです。LINE送信はしません。")
	}
}
