# ChatGPT-to-API

從 ChatGPT 網站模擬 API

## Usage

### 1. 設置帳號密碼

#### OpenAI 帳號
建立 `accounts.txt` 存放 Email 和 Password，格式參考:
```
Email-A:password
Email-B:password:2
Email-C:password:2/5
```

密碼後的數字表示輪詢次數，默認為 1 次。上例表示第一次對話使用帳戶 A，而後兩次對話使用帳戶 B，帳戶 C 為 Teams 帳戶，接著五次對話使用帳戶 C 的 Teams，然後兩次使用帳戶 C 的個人，如此循環。

登入後的 Access tokens 和 PUID 會存放在 `access_tokens.json`，每天自動更新。

注意:
1. 如果有修改帳號密碼，請將 `access_tokens.json` 刪除，重新運行即可自動生成
2. 如果使用第三方登入，請往下看

#### 第三方登入
如果使用第三方登入，請在 `accounts.txt` 添加第三方帳戶和任意密碼，並建立 `cookies.json` 以存放登入 cookies，格式參考:
```json
{
    "第三方帳號名稱": [
        {
            "Name": "__Secure-next-auth.session-token",
            "Value": "網頁登入第三方帳戶後，cookies 中的 __Secure-next-auth.session-token 值",
            "Path": "/",
            "Domain": "",
            "Expires": "0001-01-01T00:00:00Z",
            "MaxAge": 0,
            "Secure": true,
            "HttpOnly": true,
            "SameSite": 2,
            "Unparsed": null
        }
    ]
}
```
#### API Key (Optional)
如 OpenAI 官方 API 一樣，可給模擬的 API 添加 API 密鑰認證。

建立 `api_keys.txt` 以存放 API 密鑰，格式參考:
```
sk-123456
88888888
```

### 2. 運行
#### 本地
**需安裝 Go 環境**
```
go build
./freechatgpt
```

#### Docker
**需安裝 Docker**
```bash
docker compose up -d
```

## 環境變數
- `SERVER_HOST` - 預設 127.0.0.1
- `SERVER_PORT` - 預設 8080
- `ENABLE_HISTORY` - 預設 false，不允許網頁儲存紀錄 (即臨時交談)
