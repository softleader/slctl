# 規格：遷移 Wiki 內容至 docs 目錄並更新連結

## 1. 概述
本軌道的目標是將原本存放於 GitHub Wiki (`https://github.com/softleader/slctl/wiki`) 的內容，從本地的 `.wiki` 目錄完整遷移至專案內的 `docs` 目錄下，並更新整個專案中所有指向舊 Wiki 的連結，使其指向新的 `docs` 路徑。

## 2. 功能需求
### 2.1 內容遷移
- 將 `.wiki` 目錄下的所有檔案與資料夾移動或複製到 `docs` 目錄。
- 遷移過程中必須保持原始的目錄結構與檔案名稱。
- 如果 `docs` 目錄已存在，則併入；如果不存在，則建立。

### 2.2 連結更新
- 掃描整個專案（包括 `README.md`）。
- 尋找所有匹配 `https://github.com/softleader/slctl/wiki/` 前綴的連結。
- 將這些連結轉換為相對於檔案本身或專案根目錄的 `docs/` 路徑連結。
    - 範例：`https://github.com/softleader/slctl/wiki/Home` -> `docs/Home.md` (或對應的相對路徑)。
    - 注意：GitHub Wiki 連結通常不帶副檔名，遷移至 `docs` 後需確保連結能正確指向 `.md` 檔案。

## 3. 非功能需求
- **一致性**：確保所有文件間的交叉連結在遷移後依然有效。
- **清潔性**：完成後，應考慮是否保留舊的 `.wiki` 目錄（預設為保留，除非使用者要求刪除）。

## 4. 驗證準則 (Acceptance Criteria)
- [ ] `docs` 目錄包含 `.wiki` 中的所有內容。
- [ ] 執行全域搜尋不再發現指向 `https://github.com/softleader/slctl/wiki/` 的連結（或僅保留必要的外部參考）。
- [ ] `README.md` 中的連結已正確更新為指向本地 `docs` 目錄。
- [ ] 抽驗幾個連結，確保它們在本地 Markdown 預覽中能正常運作。

## 5. 範圍外 (Out of Scope)
- 重新編輯或重寫 Wiki 內容。
- 建立自動化的 Wiki 同步機制（這是一次性的遷移任務）。
