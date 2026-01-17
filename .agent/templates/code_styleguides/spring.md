# Spring Boot 開發指南

## 一般準則

* 在審查程式碼變更時，僅提出具有高度把握的建議。
* 撰寫具備良好可維護性的程式碼，並包含解釋「為何」做出特定設計決策的註解。
* 妥善處理邊界情況（Edge Cases）並撰寫清晰的例外處理機制。
* 對於函式庫或外部依賴，請在註解中說明其使用方式與目的。

## Spring Boot 準則

### 相依性注入 (Dependency Injection)

* 對所有必要的依賴項目使用建構子注入（Constructor Injection）。
* 將依賴欄位宣告為 `private final`。

### 設定 (Configuration)

* 使用 YAML 檔案 (`application.yml`) 進行外部化設定。
* 環境 Profiles：針對不同環境（開發 dev、測試 test、生產 prod）使用 Spring Profiles。
* 設定屬性：使用 `@ConfigurationProperties` 進行型別安全（Type-safe）的設定綁定。
* 機密管理：使用環境變數或密鑰管理系統將機密資訊（Secrets）外部化。

### 程式碼組織

* 套件結構：依據功能/領域（Feature/Domain）而非分層（Layer）來組織程式碼。
* 關注點分離：保持 Controller 精簡、Service 專注於業務邏輯、Repository 簡單明確。
* 工具類別：將工具類別設為 `final` 並具備私有建構子（private constructors）。

### 服務層 (Service Layer)

* 將業務邏輯放在標註 `@Service` 的類別中。
* Service 應為無狀態（Stateless）且可被測試的。
* 透過建構子注入 Repository。
* Service 的方法簽章應使用領域 ID 或 DTO，除非必要，否則不應直接暴露 Repository 的實體（Entity）。

### 日誌記錄 (Logging)

* 所有日誌皆使用 SLF4J (`private static final Logger logger = LoggerFactory.getLogger(MyClass.class);`)。
* 請勿直接使用具體實作（如 Logback, Log4j2）或 `System.out.println()`。
* 使用參數化日誌：`logger.info("User {} logged in", userId);`。

### 安全性與輸入處理

* 使用參數化查詢 | 務必使用 Spring Data JPA 或 `NamedParameterJdbcTemplate` 以防止 SQL 注入攻擊（SQL Injection）。
* 使用 JSR-380（`@NotNull`, `@Size` 等）註解與 `BindingResult` 來驗證請求主體（Request Body）與參數。

## 建置與驗證

* 新增或修改程式碼後，請驗證專案是否仍能成功建置。
* 若專案使用 Maven，請執行 `mvn clean package`。
* 若專案使用 Gradle，請執行 `./gradlew build`（Windows 上則為 `gradlew.bat build`）。
* 確保所有測試皆通過，這應是建置過程的一部分。

## 常用指令

| Gradle 指令         | Maven 指令                     | 說明                      |
| ------------------- | ------------------------------ | ------------------------- |
| `./gradlew bootRun` | `./mvnw spring-boot:run`       | 執行應用程式。            |
| `./gradlew build`   | `./mvnw package`               | 建置應用程式。            |
| `./gradlew test`    | `./mvnw test`                  | 執行測試。                |
| `./gradlew bootJar` | `./mvnw spring-boot:repackage` | 將應用程式打包為 JAR 檔。 |