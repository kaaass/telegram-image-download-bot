package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var token string
var allowedChatID int64
var downloadPath string
var httpProxy string

func main() {
	// 加载环境变量
	readEnvVars()

	// 创建 Bot
	httpClient := createHTTPClient(httpProxy)
	bot, err := tgbotapi.NewBotAPIWithClient(token, httpClient)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// 等待消息
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// 只允许指定的 chat ID
		if update.Message.Chat.ID != allowedChatID {
			log.Printf("Unauthorized access attempt from chat ID: %d\n", update.Message.Chat.ID)
			continue
		}

		// 判断消息是否包含图片
		var fileID string
		var ext string

		if update.Message.Photo != nil {
			photo := (*update.Message.Photo)[len(*update.Message.Photo)-1]
			fileID = photo.FileID
			ext = "jpg"
		} else if update.Message.Document != nil && isImageMimeType(update.Message.Document.MimeType) {
			fileID = update.Message.Document.FileID
			ext = getExtensionByMimeType(update.Message.Document.MimeType)
		}
		if fileID == "" {
			sendMessage(bot, update.Message.Chat.ID, "消息中不包含图片")
			continue
		}

		// 获取图片 URL
		fileURL, err := bot.GetFileDirectURL(fileID)
		if err != nil {
			log.Printf("Error getting file url: %v\n", err)
			sendMessage(bot, update.Message.Chat.ID, "获取图片 URL 时发生错误，请稍后重试。")
			continue
		}

		// 下载并保存图片
		currentTime := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("image_%s.%s", currentTime, ext)
		sendMessage(bot, update.Message.Chat.ID, "开始下载图片，请稍候...")
		err = downloadAndSaveFile(httpClient, fileURL, filepath.Join(downloadPath, filename))
		if err != nil {
			log.Printf("Error downloading and saving image: %v\n", err)
			sendMessage(bot, update.Message.Chat.ID, "下载并保存图片时发生错误，请稍后重试。")
			continue
		}

		sendMessage(bot, update.Message.Chat.ID, fmt.Sprintf("图片已成功下载并保存为：%s", filename))
	}
}

func readEnvVars() {
	token = os.Getenv("TELEGRAM_API_TOKEN")
	if token == "" {
		log.Fatal("Please set TELEGRAM_API_TOKEN environment variable")
	}

	allowedChatIDStr := os.Getenv("ALLOWED_CHAT_ID")
	if allowedChatIDStr == "" {
		log.Fatal("Please set ALLOWED_CHAT_ID environment variable")
	}

	downloadPath = os.Getenv("DOWNLOAD_PATH")
	if downloadPath == "" {
		log.Fatal("Please set DOWNLOAD_PATH environment variable")
	}

	httpProxy = os.Getenv("HTTP_PROXY")

	parsedID, err := strconv.ParseInt(allowedChatIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid chat ID: %v", err)
	}
	allowedChatID = parsedID
}

func createHTTPClient(httpProxy string) *http.Client {
	httpClient := &http.Client{}

	if httpProxy != "" {
		proxyURL, err := url.Parse(httpProxy)
		if err != nil {
			log.Fatalf("Invalid proxy URL: %v", err)
		}

		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	return httpClient
}

func downloadAndSaveFile(httpClient *http.Client, url string, filepath string) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %v\n", err)
		}
	}(file)

	_, err = io.Copy(file, resp.Body)
	return err
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v\n", err)
	}
}

func isImageMimeType(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}

func getExtensionByMimeType(mimeType string) string {
	ext := strings.Split(mimeType, "/")[1]
	if ext == "jpeg" {
		ext = "jpg"
	}
	return ext
}
