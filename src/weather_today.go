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

// 環境変数からAPIキーと緯度・経度を取得
func getEnvVars() (string, string, string, error) {
	// ログの設定
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// .envファイルを読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".envファイルの読み込みに失敗しました: %v", err)
	}

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	lat := os.Getenv("WEATHER_LAT")
	lon := os.Getenv("WEATHER_LON")
	if apiKey == "" {
		return "", "", "", errors.New("環境変数 OPENWEATHER_API_KEY が設定されていません")
	}
	if lat == "" {
		return "", "", "", errors.New("環境変数 WEATHER_LAT が設定されていません")
	}
	if lon == "" {
		return "", "", "", errors.New("環境変数 WEATHER_LON が設定されていません")
	}
	return apiKey, lat, lon, nil
}

// OpenWeatherMap APIから天気情報を取得（緯度・経度指定）
func fetchWeather(apiKey, lat, lon string) (*ForecastResponse, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/forecast?lat=%s&lon=%s&appid=%s&lang=ja&units=metric",
		lat, lon, apiKey,
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

	apiKey, lat, lon, err := getEnvVars()
	if err != nil {
		log.Fatalf("環境変数取得エラー: %v", err)
	}

	forecast, err := fetchWeather(apiKey, lat, lon)
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
