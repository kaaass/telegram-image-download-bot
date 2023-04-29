# Telegram Image Download Bot

这是一个使用 Go 语言编写的 Telegram bot，可以自动下载发送给 bot 的图片（包括作为文件发送的未经压缩的图片）并将它们保存到本地。

## Disclaimer

本项目的代码（包括本 README.md）均使用 ChatGPT 与 Copilot 生成。请注意，本项目的代码质量可能不高，且可能存在潜在的安全风险。请在使用本项目之前自行评估风险。

## 依赖

该项目依赖于 `github.com/go-telegram-bot-api/telegram-bot-api` 包。在开始使用之前，请确保已安装此依赖包：

```sh
go get -u github.com/go-telegram-bot-api/telegram-bot-api
```

## 使用方法

1. 使用 [@BotFather](https://t.me/BotFather) 在 Telegram 上创建一个新的 bot，并获取 API 令牌。
2. 设置以下环境变量：

    - `TELEGRAM_API_TOKEN`：你的 Telegram bot API 令牌。
    - `ALLOWED_CHAT_ID`：允许访问 bot 的 Telegram chat ID。
    - `DOWNLOAD_PATH`：图片下载到本地的文件夹路径。
    - `HTTP_PROXY`：（可选）用于访问 Telegram 和下载图片的 HTTP 代理。

3. 编译并运行项目：

```sh
go build
./telegram-image-download-bot
```

或者，你可以使用 Docker 构建并运行项目。请参考 [Dockerfile](#Dockerfile) 部分了解详细信息。

4. 在 Telegram 中向 bot 发送图片。下载的图片将保存到本地的 `DOWNLOAD_PATH` 文件夹中。

## Dockerfile

要使用 Docker 运行此项目，请按照以下步骤操作：

1. 构建 Docker 镜像：

```sh
docker build -t your-image-name .
```

2. 运行 Docker 镜像：

```sh
docker run -d --name your-container-name \
  -v /path/to/download/folder:/downloads \
  -e TELEGRAM_API_TOKEN="your-api-token" \
  -e ALLOWED_CHAT_ID="your-chat-id" \
  -e HTTP_PROXY="http://proxy.example.com:8080" \
  your-image-name
```

请使用你自己的镜像名、容器名、下载文件夹路径、API 令牌和代理地址替换示例中的占位符。

## 注意事项

- 请确保 `DOWNLOAD_PATH` 文件夹具有正确的读写权限。
- 如有需要，可以根据实际情况调整 `HTTP_PROXY` 环境变量。
- 请注意，仅支持 `ALLOWED_CHAT_ID` 中指定的 chat ID 访问 bot。请确保已正确设置 chat ID。
