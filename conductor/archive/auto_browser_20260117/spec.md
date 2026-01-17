# Specification - 自動開啟瀏覽器與複製驗證碼

## Overview
優化 `slctl init` 的使用者體驗。當進行 GitHub Device Flow 驗證時，自動開啟瀏覽器至驗證頁面，並將 User Code 自動複製到系統剪貼簿，減少使用者的手動操作步驟。

## Functional Requirements
1.  **自動開啟瀏覽器**：
    -   當取得 Device Code 後，自動開啟預設瀏覽器並導向至 `verification_uri`。
    -   使用專案現有的 `github.com/skratchdot/open-golang` 或同等功能的函式庫。
2.  **自動複製驗證碼**：
    -   將 `user_code` 自動寫入系統剪貼簿。
    -   需引入新的依賴庫（建議：`github.com/atotto/clipboard`）以支援跨平台剪貼簿操作。
3.  **終端機輸出**：
    -   **必須**保留原本的提示訊息：`Please go to https://github.com/login/device and enter the code: D896-85C0`。
    -   (建議) 當自動操作成功時，可額外顯示 Log 提示（如 "Browser opened", "Code copied"），以便使用者感知。

## Non-Functional Requirements
1.  **容錯處理**：
    -   若無法開啟瀏覽器或寫入剪貼簿（例如無 GUI 環境、SSH 連線），程式 **嚴禁崩潰 (Panic)**。
    -   應僅記錄錯誤或警告，確保使用者仍可依照終端機的文字提示手動完成驗證。
2.  **跨平台支援**：
    -   支援 macOS, Windows, Linux。

## Acceptance Criteria
1.  執行 `slctl init` (需觸發 Device Flow)：
    -   瀏覽器自動開啟至正確的 GitHub 驗證頁面。
    -   使用者在輸入框貼上 (Paste)，內容應為正確的 User Code。
    -   終端機仍完整顯示原始的 `Please go to...` 訊息。
2.  異常測試：
    -   在無瀏覽器或無剪貼簿支援的環境執行時，程式不會崩潰，且流程可繼續。

## Out of Scope
-   自動點擊網頁上的按鈕或自動登入。
