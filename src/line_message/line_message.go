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
	userID := os.Getenv("MY_USER_ID")
	if token == "" {
		return "", "", errors.New("環境変数 LINE_CHANNEL_ACCESS_TOKEN が設定されていません")
	}
	if userID == "" {
		return "", "", errors.New("環境変数 MY_USER_ID が設定されていません")
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

	// ここでrain_infoを定義する（本来は天気予報APIの結果を受け取る想定）
	rainInfo := "今日は傘を持ちましょう。\n"

	// 送りたいテキスト
	finalText := "今日は雨が降りそうです！傘を忘れずに！"

	// 雨の情報がある場合だけ送信
	if rainInfo == "今日は傘を持ちましょう。\n" {
		if err := sendLineMessage(token, userID, finalText); err != nil {
			log.Fatalf("メッセージ送信エラー: %v", err)
		}
		fmt.Println("メッセージ送信に成功しました。")
	} else {
		fmt.Println("今日は傘の必要はなさそうです。メッセージは送信しません。")
	}
}
