// line_message.go
// LINE Messaging APIを利用して任意のメッセージを送信するサンプル
// セキュリティベストプラクティスに従い、アクセストークン等は環境変数から取得
// 実行例: LINE_CHANNEL_ACCESS_TOKEN=xxxx LINE_USER_ID=xxxx go run line_message.go

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// LINE Messaging APIのエンドポイント
const lineAPIEndpoint = "https://api.line.me/v2/bot/message/push"

// 環境変数からアクセストークンとユーザーIDを取得
func getEnvVars() (string, string, error) {
	// ログの設定
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// .envファイルを読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".envファイルの読み込みに失敗しました: %v", err)
	}

	token := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	userID := os.Getenv("LINE_USER_ID")
	if token == "" {
		return "", "", errors.New("環境変数 LINE_CHANNEL_ACCESS_TOKEN が設定されていません")
	}
	if userID == "" {
		return "", "", errors.New("環境変数 LINE_USER_ID が設定されていません")
	}
	return token, userID, nil
}

// LINEにメッセージを送信
func sendLineMessage(token, userID, message string) error {
	payload := map[string]interface{}{
		"to": userID,
		"messages": []map[string]string{
			{
				"type": "text",
				"text": message,
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("リクエストボディの生成に失敗: %w", err)
	}

	req, err := http.NewRequest("POST", lineAPIEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("HTTPリクエスト生成に失敗: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("LINE APIへのリクエストに失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LINE APIエラー: %s", resp.Status)
	}
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	token, userID, err := getEnvVars()
	if err != nil {
		log.Fatalf("環境変数取得エラー: %v", err)
	}

	message := "こんにちは！これはLINE Messaging APIから送信されたテストメッセージです。"
	if err := sendLineMessage(token, userID, message); err != nil {
		log.Fatalf("メッセージ送信エラー: %v", err)
	}
	fmt.Println("メッセージ送信に成功しました。")
}
