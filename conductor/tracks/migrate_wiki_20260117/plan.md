# 實作計畫：遷移 Wiki 內容至 docs 目錄並更新連結

## 第一階段：環境準備與內容遷移
- [ ] Task: 建立基礎目錄結構
    - [ ] 確認 `docs` 目錄是否存在，若不存在則建立 `/Users/samwang/Dev/slctl/docs`
- [ ] Task: 內容遷移
    - [ ] 將 `.wiki` 目錄下的所有檔案與子目錄同步至 `docs`
    - [ ] 驗證檔案同步後的完整性
- [ ] Task: Conductor - User Manual Verification '環境準備與內容遷移' (Protocol in workflow.md)

## 第二階段：全域連結更新
- [ ] Task: 更新 README.md
    - [ ] 替換 `README.md` 中所有指向 `https://github.com/softleader/slctl/wiki/` 的連結為相對路徑 `./docs/`
- [ ] Task: 全域連結掃描與更新
    - [ ] 使用 grep 掃描專案內所有檔案，尋找舊 Wiki 連結
    - [ ] 批次更新這些連結，並加上 `.md` 副檔名（如果原本沒有）
- [ ] Task: Conductor - User Manual Verification '全域連結更新' (Protocol in workflow.md)

## 第三階段：驗證與清理
- [ ] Task: 驗證連結有效性
    - [ ] 隨機抽選 5-10 個更新後的連結，確認其指向的檔案確實存在於 `docs`
    - [ ] 檢查是否有遺漏的連結未被轉換
- [ ] Task: 專案清理 (選用)
    - [ ] 詢問使用者是否刪除舊的 `.wiki` 目錄
- [ ] Task: Conductor - User Manual Verification '驗證與清理' (Protocol in workflow.md)
