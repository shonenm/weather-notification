// weather.go
// OpenWeatherMap APIを利用した天気情報取得ロジック（再利用用）
// セキュリティベストプラクティスに従い、APIキー等は環境変数から取得

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// 天気API用の環境変数を取得
func GetWeatherEnvVars() (string, string, string, error) {
	_ = godotenv.Load()
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

// OpenWeatherMap APIから天気情報を取得
func FetchWeather(apiKey, lat, lon string) (*ForecastResponse, error) {
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

// 24時間以内に雨が降るか判定し、雨の時間帯リストを返す
func NeedUmbrella(forecast *ForecastResponse) (bool, string) {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)
	horizon := now.Add(24 * time.Hour)
	var rainTimes []string

	for _, entry := range forecast.List {
		entryTime := time.Unix(entry.Dt, 0).In(jst)
		if entryTime.After(now) && entryTime.Before(horizon) {
			for _, w := range entry.Weather {
				if w.Main == "Rain" || w.Main == "雨" ||
					w.Description == "雨" || w.Description == "小雨" || w.Description == "強い雨" {
					rainTimes = append(rainTimes, entryTime.Format("01/02 15:04"))
					break
				}
			}
		}
	}
	if len(rainTimes) > 0 {
		return true, fmt.Sprintf("本日(%s)は以下の時間帯で雨の予報があります: %v", now.Format("2006-01-02"), rainTimes)
	}
	return false, ""
}
