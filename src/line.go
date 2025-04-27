// line.go
// LINE Messaging APIを利用して任意のメッセージを送信するサンプル
// セキュリティベストプラクティスに従い、アクセストークン等は環境変数から取得

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const lineAPIEndpoint = "https://api.line.me/v2/bot/message/push"

// アクセストークンとユーザーIDを取得
func GetLineEnvVars() (string, string, error) {
	_ = godotenv.Load()

	token := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	userID := os.Getenv("MY_USER_ID")
	if token == "" {
		return "", "", errors.New("環境変数 LINE_CHANNEL_ACCESS_TOKEN が設定されていません")
	}
	if userID == "" {
		return "", "", errors.New("環境変数 MY_USER_ID が設定されていません")
	}
	return token, userID, nil
}

// LINEにテキストメッセージを送信
func SendTextMessage(token, userID, message string) error {
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
