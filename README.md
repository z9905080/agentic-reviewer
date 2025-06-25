# CV Analyzer API

## 專案目標

本專案旨在開發一套基於 Go 語言的 API 服務，用於對上傳的履歷 (CV) 和職位描述 (JD) 文件進行初步的 AI 分析和評分。系統遵循 Clean Architecture 和 Domain-Driven Design (DDD) 原則，並使用 `google/wire` 進行依賴注入。其設計目標是實現一個可擴展、可維護且 LLM 後端可抽換的框架。

## 主要功能

*   **API 端點**: 提供 HTTP API 端點接收 CV 和 JD 文件。
*   **文件解析**: 支援解析純文本 (.txt) 和 PDF (.pdf) 格式的文件以提取內容。
*   **AI 分析 (骨架)**: 整合 LLM (大型語言模型) 進行分析的核心框架已搭建，目前支援 OpenAI 和 Claude (骨架，需填充實際 API 調用邏ics)。LLM Repository 可抽換。
*   **評分與報告 (骨架)**: 生成分析結果，包含評分、摘要等 (目前為 mock 數據)。
*   **Clean Architecture**: 清晰分層的代碼結構，易於理解和擴展。
*   **DDD 原則**: 專注於核心業務領域的建模。
*   **依賴注入**: 使用 `google/wire` 管理組件依賴。

## 架構概覽

本專案採用 Clean Architecture 設計，主要分為以下幾層：

*   **Domain Layer (`internal/domain`)**: 包含核心業務實體 (如 `CV`, `JobDescription`, `AnalysisResult`) 和業務邏輯。此層不依賴任何其他層。
*   **Application Layer (`internal/application`)**: 包含應用程式的特定業務規則 (Use Cases，如 `AnalyzeCVUseCase`) 和介面定義 (Ports，如 `FileParser`, `AnalysisService`, `StorageService`)。此層依賴 Domain Layer。
*   **Infrastructure Layer (`internal/infrastructure`)**: 提供外部服務的具體實作 (Adapters)，例如文件解析器 (`PlainTextParser`, `PdfParser`)、LLM 客戶端 (`OpenAIAdapter`, `ClaudeAdapter`) 和持久化儲存 (`InMemoryStorage`)。此層實現 Application Layer 定義的介面。
*   **Interfaces Layer (`internal/interfaces`)**: 處理與外部世界的交互，主要是 API Handler (`AnalysisHandler`)、路由 (`routes.go`) 和數據傳輸對象 (DTOs)。此層依賴 Application Layer。

依賴關係遵循單向原則，從外層指向內層。`google/wire` 用於在 `cmd/cv-analyzer/main.go` 中組裝這些依賴。

## 設定與安裝

### 前置需求

*   **Go**: 版本 1.19 或更高。
*   **(可選) LLM API Keys**: 如果要啟用真實的 LLM 分析 (目前 LLM Adapter 為骨架)，需要設定以下環境變數：
    *   `OPENAI_API_KEY`: 您的 OpenAI API 金鑰。
    *   `CLAUDE_API_KEY`: 您的 Anthropic Claude API 金鑰。
    (注意：目前 LLM Adapters 的 `Analyze` 方法返回的是 mock 數據，尚未完全實現對這些金鑰的真實 API 調用。)

### 安裝步驟

1.  **克隆倉儲**:
    ```bash
    git clone <repository-url>
    cd cv-analyzer
    ```

2.  **安裝依賴**:
    ```bash
    go mod tidy
    ```

3.  **產生 Wire 依賴注入代碼**:
    ```bash
    go generate ./...
    # 或者直接執行 wire (需先安裝: go install github.com/google/wire/cmd/wire@latest)
    # wire ./cmd/cv-analyzer
    ```
    這將會在 `cmd/cv-analyzer/` 目錄下生成 `wire_gen.go` 文件。

## 如何執行應用程式

完成安裝步驟後，可以執行以下命令來啟動 API 伺服器：

```bash
go run cmd/cv-analyzer/main.go
```

預設情況下，伺服器會在 `http://localhost:8080` (或在 `main.go` 中配置的端口) 啟動。

## API 端點

### 分析履歷與職位描述

*   **請求**: `POST /api/v1/analysis/cv`
*   **Content-Type**: `multipart/form-data`
*   **表單參數**:
    *   `cv`: 履歷文件 (必需, .txt 或 .pdf)
    *   `jd`: 職位描述文件 (必需, .txt 或 .pdf)

*   **成功回應 (200 OK)**:
    ```json
    {
        "id": "analysis-uuid",
        "cv_id": "cv-uuid",
        "cv_path": "/tmp/unique-cv-file.pdf",
        "jd_id": "jd-uuid",
        "jd_path": "/tmp/unique-jd-file.txt",
        "score": 85.5,
        "summary": "This is a mock summary from the LLM.",
        "matching_keywords": ["Go", "Clean Architecture"],
        "missing_keywords": ["Kubernetes"],
        "analyzed_at": "YYYY-MM-DDTHH:MM:SSZ"
    }
    ```
    (注意: `cv_path` 和 `jd_path` 是伺服器上的臨時文件路徑，主要用於調試。`id`, `cv_id`, `jd_id` 等 UUID 由系統生成。評分和摘要內容目前由 mock LLM adapter 提供。)

*   **錯誤回應 (400 Bad Request / 500 Internal Server Error)**:
    ```json
    {
        "error": "Error message describing the issue"
    }
    ```

## 如何執行測試

1.  **執行所有單元測試**:
    ```bash
    go test ./...
    ```

2.  **執行特定包的測試**:
    例如，執行 `application/usecases` 包下的測試：
    ```bash
    go test ./internal/application/usecases/...
    ```

## 目前狀態與未來改進

### 目前已完成:
*   核心框架和目錄結構。
*   Domain entities 和 Application use cases/ports 定義。
*   Txt 和 Pdf 文件解析器。
*   ParserRouter 動態選擇解析器。
*   OpenAI 和 Claude LLM Adapter 的骨架。
*   In-memory storage service。
*   API handler for CV/JD upload and analysis invocation.
*   使用 `google/wire` 進行依賴注入。
*   `AnalyzeCVUseCase` 的單元測試和 mocks。

### 未來可能的改進方向:
*   **完成 LLM Adapters**: 實現 `OpenAIAdapter` 和 `ClaudeAdapter` 中真實的 API 調用邏輯，包括 prompt engineering 和回應解析。
*   **增強文件解析**:
    *   支援 DOCX 文件格式。
    *   考慮使用更強大的 PDF 解析庫，或整合 OCR 功能處理掃描型 PDF。
*   **錯誤處理與日誌**: 引入結構化日誌 (如 `zerolog`, `zap`)，並建立更完善的錯誤處理機制和 API 錯誤碼。
*   **配置管理**: 從設定檔 (如 YAML/JSON) 或環境變數中讀取 LLM API Keys 和其他應用程式配置。
*   **資料庫整合**: 替換 `InMemoryStorage` 為真實的資料庫儲存 (如 PostgreSQL, MySQL)。
*   **整合測試**: 為 API Endpoints 和 Infrastructure components 編寫更全面的整合測試。
*   **Swagger/OpenAPI**: 解決 Swagger UI 的整合問題，提供動態的 API 文件。
*   **安全性**: 增加 API 安全性措施 (如認證、授權)。
*   **異步處理**: 對於耗時的 LLM 分析，可以考慮改為異步處理模式。

## 貢獻

歡迎對本專案進行貢獻！請 Fork 本倉儲，建立您的分支，提交您的變更，並發起 Pull Request。
