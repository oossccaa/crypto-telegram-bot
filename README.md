# 虛擬貨幣價格監控 Telegram Bot

這是一個用 Go 語言開發的 Telegram Bot，用於監控 ADA（Cardano）和 ETH（Ethereum）的價格變動，當價格變動超過 0.2% 時會自動發送通知。

## 功能特色

- 🔍 監控 ADA 和 ETH 價格
- 📊 從 CoinGecko API 獲取即時價格資料
- 📱 透過 Telegram Bot 發送價格警報
- ⏰ 每 5 分鐘自動檢查一次價格變動
- 🚨 當價格變動超過 0.2% 時自動推播通知

## 安裝和設定

### 1. 系統需求

- Go 1.19 或更新版本
- 網路連線（用於 API 請求）

### 2. 下載和安裝依賴套件

```bash
# 克隆專案（如果從 git 取得）
git clone <repository-url>
cd crypto-telegram-bot

# 安裝依賴套件
go mod tidy
```

### 3. 建立 Telegram Bot

1. 在 Telegram 中搜尋 `@BotFather`
2. 發送 `/newbot` 建立新的機器人
3. 按照指示設定機器人名稱和用戶名
4. 取得 Bot Token（格式類似：`123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ`）

### 4. 取得 Chat ID

1. 將您的機器人加入到目標群組或開始私聊
2. 發送一則測試訊息給機器人
3. 在瀏覽器中訪問：`https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
4. 在回應的 JSON 中找到 `"chat"` -> `"id"` 的數值

### 5. 設定配置檔案

```bash
# 複製範例配置檔案
cp config.yaml.example config.yaml

# 編輯配置檔案，填入您的 Bot Token 和 Chat ID
nano config.yaml
```

配置檔案範例：
```yaml
bot_token: "123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ"
chat_id: -1234567890
```

## 使用方法

### 方法一：Docker 部署（推薦）

使用 Docker 是最簡單的部署方式：

```bash
# 確保已設定好 config.yaml 檔案
cp config.yaml.example config.yaml
# 編輯 config.yaml 填入您的 Bot Token 和 Chat ID

# 使用 docker-compose 啟動
docker-compose up -d

# 查看運行日誌
docker-compose logs -f

# 停止服務
docker-compose down

# 重新啟動服務
docker-compose restart
```

Docker 部署的優點：
- 無需安裝 Go 環境
- 容器化運行，環境一致
- 自動重啟功能
- 簡單的日誌管理

### 方法二：直接執行 Bot

```bash
# 編譯並執行
go run main.go
```

或者編譯成執行檔：

```bash
# 編譯
go build -o crypto-bot

# 執行
./crypto-bot
```

### 預期行為

- Bot 啟動後會每 5 分鐘檢查一次 ADA 和 ETH 的價格
- 如果價格變動超過 0.2%（上漲或下跌），會發送包含以下資訊的通知：
  - 當前價格
  - 前次價格
  - 變動金額和百分比
  - 時間戳記

### 通知訊息範例

```
📈 ADA Price Alert

💰 Current Price: $0.123456
📊 Previous Price: $0.120000
📈 Change: 2.88% (0.003456)
⏰ Time: now
```

## 專案結構

```
crypto-telegram-bot/
├── main.go                    # 主程式入口點
├── config/
│   └── config.go             # 配置管理
├── models/
│   └── price.go              # 價格資料模型
├── services/
│   ├── coingecko.go          # CoinGecko API 服務
│   ├── telegram.go           # Telegram Bot 服務
│   └── monitor.go            # 價格監控服務
├── config.yaml.example       # 配置檔案範例
├── config.yaml               # 實際配置檔案（需要手動建立）
├── Dockerfile                # Docker 建置檔案
├── docker-compose.yml        # Docker Compose 配置
├── go.mod                    # Go 模組檔案
└── README.md                 # 說明文件
```

## 技術細節

- **語言**: Go 1.19+
- **API**: CoinGecko API（免費版）
- **Telegram Bot**: go-telegram-bot-api/v5
- **配置格式**: YAML
- **監控間隔**: 5 分鐘
- **警報閾值**: 0.2% 價格變動

## 故障排除

### 常見問題

1. **Bot Token 錯誤**
   - 確認 Bot Token 格式正確
   - 確認機器人已啟用

2. **Chat ID 錯誤**
   - 確認 Chat ID 為數字格式
   - 確認機器人已加入目標群組或私聊

3. **API 請求失敗**
   - 檢查網路連線
   - CoinGecko API 可能有速率限制

4. **配置檔案讀取失敗**
   - 確認 `config.yaml` 檔案存在
   - 確認 YAML 格式正確

### 日誌訊息

Bot 會在控制台輸出以下訊息：
- 啟動確認
- 價格檢查錯誤
- API 請求失敗

## 授權

本專案僅供學習和個人使用。

## 注意事項

- 請勿頻繁請求 CoinGecko API，以免觸發速率限制
- 保護好您的 Bot Token，不要分享給他人
- 本 Bot 僅供參考，投資決策請謹慎評估