// weather_today.go
// OpenWeatherMap APIを利用して本日の天気予報を取得するサンプル
// セキュリティベストプラクティスに従い、APIキー等は環境変数から取得
// 実行例: OPENWEATHER_API_KEY=xxxx WEATHER_CITY=Tokyo go run weather_today.go

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// OpenWeatherMap APIのレスポンス構造体
type WeatherResponse struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Name string `json:"name"`
	Dt   int64  `json:"dt"`
}

// 環境変数からAPIキーと都市名を取得
func getEnvVars() (string, string, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	city := os.Getenv("WEATHER_CITY")
	if apiKey == "" {
		return "", "", errors.New("環境変数 OPENWEATHER_API_KEY が設定されていません")
	}
	if city == "" {
		return "", "", errors.New("環境変数 WEATHER_CITY が設定されていません")
	}
	return apiKey, city, nil
}

// OpenWeatherMap APIから天気情報を取得
func fetchWeather(apiKey, city string) (*WeatherResponse, error) {
	// APIエンドポイント
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&lang=ja&units=metric",
		city, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("APIリクエストに失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("APIエラー: %s, レスポンス: %s", resp.Status, string(body))
	}

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)
	}
	return &weather, nil
}

func main() {
	// ログの設定
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	apiKey, city, err := getEnvVars()
	if err != nil {
		log.Fatalf("環境変数取得エラー: %v", err)
	}

	weather, err := fetchWeather(apiKey, city)
	if err != nil {
		log.Fatalf("天気情報取得エラー: %v", err)
	}

	// 日付を日本時間で表示
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	date := time.Unix(weather.Dt, 0).In(jst).Format("2006-01-02 15:04")

	fmt.Printf("【%s の天気予報 (%s)】\n", weather.Name, date)
	if len(weather.Weather) > 0 {
		fmt.Printf("天気: %s (%s)\n", weather.Weather[0].Main, weather.Weather[0].Description)
	}
	fmt.Printf("気温: %.1f℃\n", weather.Main.Temp)
	fmt.Printf("湿度: %d%%\n", weather.Main.Humidity)
}
