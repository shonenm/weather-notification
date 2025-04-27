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

	"github.com/joho/godotenv"
)

// OpenWeatherMap APIのレスポンス構造体
type ForecastResponse struct {
	List []struct {
		Dt   int64 `json:"dt"`
		Main struct {
			Temp     float64 `json:"temp"`
			Humidity int     `json:"humidity"`
		} `json:"main"`
		Weather []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
		} `json:"weather"`
	} `json:"list"`
	City struct {
		Name string `json:"name"`
	} `json:"city"`
}

// 環境変数からAPIキーと都市名を取得
func getEnvVars() (string, string, error) {
	// ログの設定
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// .envファイルを読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".envファイルの読み込みに失敗しました: %v", err)
	}

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
func fetchWeather(apiKey, city string) (*ForecastResponse, error) {
	// APIエンドポイント
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&lang=ja&units=metric",
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

	var weather ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)
	}


	if len(weather.List) == 0 {
		return nil, errors.New("天気予報のリストが空です")
	}
	
	return &weather, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	apiKey, city, err := getEnvVars()
	if err != nil {
		log.Fatalf("環境変数取得エラー: %v", err)
	}

	forecast, err := fetchWeather(apiKey, city)
	if err != nil {
		log.Fatalf("天気情報取得エラー: %v", err)
	}

	// JSTで今日の日付を取得
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)
	today := now.Format("2006-01-02")

	fmt.Printf("【%s の本日(%s)3時間ごとの天気予報】\n", forecast.City.Name, today)

	horizon := now.Add(24 * time.Hour)

	fmt.Printf("【%s のこれから24時間の天気予報】\n", forecast.City.Name)

	found := false
	for _, entry := range forecast.List {
		entryTime := time.Unix(entry.Dt, 0).In(jst)
		if entryTime.After(now) && entryTime.Before(horizon) {
			found = true
			fmt.Printf("\n[%s]\n", entryTime.Format("01/02 15:04"))
			if len(entry.Weather) > 0 {
				fmt.Printf("天気: %s (%s)\n", entry.Weather[0].Main, entry.Weather[0].Description)
			}
			fmt.Printf("気温: %.1f℃\n", entry.Main.Temp)
			fmt.Printf("湿度: %d%%\n", entry.Main.Humidity)
		}
	}
	if !found {
		fmt.Println("これから24時間の天気予報データが見つかりませんでした。")
	}
}
