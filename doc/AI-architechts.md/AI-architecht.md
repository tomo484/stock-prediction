プロジェクト設計・実装コンテキスト
1. プロジェクトの目的
日本株（当面はITセクターに限定）を対象に、「中期投資（数ヶ月〜半年）で株価上昇が見込める銘柄」 を発掘するAI分析システムを構築する。 人間のアナリストと同様に、「定性情報（ビジネスモデル・ニュース）」と「定量情報（財務・株価）」の両面から多角的に分析し、根拠のある投資判断を出力することをゴールとする。

2. システムアーキテクチャ概要
Go (Backend) と Python (AI Microservice) のマイクロサービス構成を採用し、間を gRPC で通信させる。

Go Service: データ収集、DB管理、APIゲートウェイ、バッチ処理の実行責任を持つ（"土管"兼"マネージャー"）。

Python Service: LangGraphを用いた高度な推論、分析ロジックの実行責任を持つ（"頭脳"）。

3. Go (Backend) の役割と実装範囲
Goは**「分析に必要な全ての素材（データ）を集めて整形し、Pythonに渡すこと」**に徹する。Python側で外部APIを叩くことは極力避ける。

A. データ収集 (J-Quants API & Tavily API)
以下のデータをAPIから取得し、Supabase (PostgreSQL) に保存する。

銘柄マスタ取得:

API: GET /listed/info

処理: 全銘柄を取得し、DBの companies テーブルに保存。分析時はここから「情報・通信業（SectorCode=10）」を抽出する。

株価データ取得:

API: GET /prices/daily_quotes

処理: 過去 6ヶ月分 の日足データを取得。トレンド分析用。

財務データ取得:

API: GET /fins/statements (詳細版ではなくStatementsを使用)

処理: 過去 5年分（四半期・通期全て） のデータを全量取得。

DB: ＃＃＃＃＃

定性情報・ニュース取得:

API: Tavily Search API

処理: 「{企業名} ビジネスモデル」「{企業名} 決算短信 要約」等のクエリで検索し、テキストデータを取得する。

B. Pythonへのデータ提供 (gRPC Client)
DBから取得したデータを .proto で定義された形式に詰め込み、PythonのgRPCサーバーへリクエストを送信する。

財務データなどの複雑な構造体は、あえてパースせず JSON文字列 として渡す。（要件等だけど、結構スキーマの数も多くなるから、これでいいかなとは思っている不便だったら言ってください）

4. Python (AI Agent) の役割と実装範囲
Goから渡されたデータを基に、LangGraph を用いて思考・分析を行う。

A. LangGraph 設計 (Supervisorパターン)
「3段階思考プロセス」 を採用し、動的かつ精度の高い分析を実現する。

Phase 1: 戦略策定 (Strategy)

入力: Goから渡された「Tavilyの検索結果（ビジネスモデル等）」

思考: 企業の成長フェーズを特定し、「この企業を評価するために見るべき重要なKPI（財務指標）」 を決定する。

例: 「SaaS企業なので、売上高成長率と営業CFを見るべき」

出力: ビジネスモデルの理解、見るべきKPIのリスト

Phase 2: 分析実行 (Execution)

入力: Phase 1の戦略 + Goから渡された「財務データ(JSON)」 + 「株価推移」

思考: JSONデータから実際の数値を参照し、Phase 1で定めたKPIが基準を満たしているか検証する。同時に株価トレンドと合わせて最終判断を下す。この段階で、株価データと財務データの要約情報もAIで生成する。

出力: 投資判断（Strong Buy / Buy / Hold / Sell）、その根拠、株価要約情報、財務要約情報

Phase 3: 比較分析 (Comparison)

入力: 全銘柄のPhase 2結果 + 各銘柄の株価要約情報 + 各銘柄の財務要約情報

思考: セクター内の全銘柄を比較し、投資価値の高いTop3銘柄を選出する。Phase 2で生成された要約情報を参照しながら、銘柄間の優劣を判断する。

出力: Top3銘柄の選出結果、選出理由、セクター全体のサマリー

B. 並列処理の方針
Phase 1とPhase 2は並列処理で実行する。これにより、処理時間を大幅に短縮できる。

- Phase 1: 各企業のnews取得と推論を並列実行（goroutine使用）
- Phase 2: 各企業の財務・株価データ取得と推論を並列実行（goroutine使用）
- Phase 3: 全Phase 2完了後に、全銘柄の結果を統合して実行（順次処理）

5. gRPC インターフェース設計 (.proto)
GoとPythonの境界線となるデータ構造の定義方針。

Request (Go -> Python)
Goが集めた「素材」を全て渡すスタイル。

ticker (string): 銘柄コード

company_info (string): 企業名、セクター

financial_statements_json (string): 重要。 J-Quantsから取得した過去5年分の財務データ配列を、そのままJSON文字列化したもの。Python側で pandas.read_json 等で解析する。

stock_prices (repeated Struct): 日付、始値、終値などの構造体リスト（過去6ヶ月分）。

qualitative_info (repeated string): Tavilyで取得したニュースや記事のテキストリスト。

Response (Python -> Go)
分析結果とログ。

sentiment (Enum): 投資判断結果。

summary_reasoning (string): 最終的な分析レポート（Markdown形式）。

thought_log (string): Phase 1/Phase 2でAIがどう考えたかの思考ログ（デバッグ・表示用）。

## 開発の進め方
1. DB定義（日本株用のテーブル: companies, financial_statements, stock_prices, news等）
2. J-Quants API認証の仕組みを実装（リフレッシュトークン管理）
3. J-Quants APIクライアントの実装（fmp.goを参考に）
   - 銘柄マスタ取得 (/listed/info)
   - 株価データ取得 (/prices/daily_quotes)
   - 財務データ取得 (/fins/statements)
4. Tavily APIクライアントの実装（news/tavily_search.go の拡張or新規）
5. SyncData()のような統合メソッド構築
6. protoファイルの定義とコード生成
7. Go側のgRPCクライアント実装
8. Python側でLangGraph + gRPCサーバー実装
9. gRPCでの疎通確認

## 設計フロー

### 全体フロー図

```
┌─────────────────────────────────────────────────────────────┐
│ 1. セクター企業リスト取得                                    │
│    companiesテーブルからSectorCode=xの企業を取得              │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 【第一段階】定性分析（Newsベース）【並列処理】                │
│                                                              │
│ 2. 各企業のnews取得&保存（並列）                             │
│    ┌──────────┐ ┌──────────┐ ┌──────────┐                │
│    │ 企業A    │ │ 企業B    │ │ 企業C    │ ...            │
│    │ news取得 │ │ news取得 │ │ news取得 │                │
│    └────┬─────┘ └────┬─────┘ └────┬─────┘                │
│         │            │            │                       │
│ 3. Phase 1推論（並列）                                       │
│    ┌────▼─────┐ ┌────▼─────┐ ┌────▼─────┐                │
│    │ Phase 1  │ │ Phase 1  │ │ Phase 1  │                │
│    │ (企業A)  │ │ (企業B)  │ │ (企業C)  │                │
│    └────┬─────┘ └────┬─────┘ └────┬─────┘                │
│         │            │            │                       │
│ 4. 結果収集・DB保存                                          │
│    AnalysisResultテーブルにPhase1結果を保存                 │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 【第二段階】定量分析（財務・株価ベース）【並列処理】          │
│                                                              │
│ 5. stock_data, fundamentals取得（並列）                      │
│    ┌──────────┐ ┌──────────┐ ┌──────────┐                │
│    │ 企業A    │ │ 企業B    │ │ 企業C    │                │
│    │ データ取得│ │ データ取得│ │ データ取得│                │
│    └────┬─────┘ └────┬─────┘ └────┬─────┘                │
│         │            │            │                       │
│ 6. Phase 2推論（並列）                                       │
│    ┌────▼─────┐ ┌────▼─────┐ ┌────▼─────┐                │
│    │ Phase 2  │ │ Phase 2  │ │ Phase 2  │                │
│    │ (企業A)  │ │ (企業B)  │ │ (企業C)  │                │
│    │ +要約生成│ │ +要約生成│ │ +要約生成│                │
│    └────┬─────┘ └────┬─────┘ └────┬─────┘                │
│         │            │            │                       │
│ 7. 結果収集・DB保存                                          │
│    - AnalysisResultテーブルにPhase2結果を更新                │
│    - StockSummaryテーブルに株価要約を保存                   │
│    - FinancialSummaryテーブルに財務要約を保存               │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 【第三段階】比較分析（銘柄間比較）【順次処理】                │
│                                                              │
│ 8. 全銘柄の分析結果と要約情報を取得                         │
│    - AnalysisResultテーブルから全銘柄の結果                 │
│    - StockSummaryテーブルから全銘柄の株価要約               │
│    - FinancialSummaryテーブルから全銘柄の財務要約           │
│                                                              │
│ 9. protoに入れてPythonに送る                                │
│    - 各銘柄のCode, Sentiment, SummaryReasoning              │
│    - 各銘柄のStockSummary                                    │
│    - 各銘柄のFinancialSummary                                │
│                                                              │
│ 10. Python側で第三段階推論                                   │
│     - 銘柄間比較（要約情報を参照）                          │
│     - Top3銘柄の選出                                         │
│     - 選出理由の生成                                         │
│                                                              │
│ 11. 返り値をprotoに入れてGoに送る                           │
│                                                              │
│ 12. SectorAnalysisResultテーブルに保存                      │
└─────────────────────────────────────────────────────────────┘
```

### 詳細フロー

#### 【第一段階】定性分析（Newsベース）【並列処理】

1. **企業リスト取得**
   - companyテーブルからSectorCode=xの企業を取得

2. **News取得（並列実行）**
   - 各企業に対してgoroutineで並列実行
   - クエリ1: "{企業名} ビジネスモデル"
   - クエリ2: "{企業名} 決算短信 要約"
   - news_searchesテーブルとnews_itemsテーブルに保存

3. **Phase 1推論（並列実行）**
   - 各企業のnewsデータをprotoに入れてPythonに送る
   - Python側で第一段階推論（ビジネスモデル理解、見るべきKPI決定）
   - 返り値: ビジネスモデルの理解、見るべきKPIリスト

4. **結果保存**
   - 返り値をprotoから受け取り、AnalysisResultテーブルに保存
   - Phase1の結果として、BusinessModelとKPIフィールドに保存

#### 【第二段階】定量分析（財務・株価ベース）【並列処理】

5. **データ取得（並列実行）**
   - 各企業に対してgoroutineで並列実行
   - stock_dataテーブルから過去6ヶ月分の日足データ取得
   - fundamentalsテーブルから過去5年分の財務データ取得

6. **Phase 2推論（並列実行）**
   - protoに入れてPythonに送る
     - Phase1の結果（ビジネスモデル、KPI）
     - 財務データ(JSON)
     - 株価推移データ
   - Python側で第二段階推論
     - 財務データと株価で検証
     - 最終判断（Sentiment）
     - **株価データの要約情報をAI生成**
     - **財務データの要約情報をAI生成**

7. **結果保存**
   - 返り値をprotoから受け取り
   - AnalysisResultテーブルにPhase2の結果を更新
   - **StockSummaryテーブルに株価要約情報を保存**
   - **FinancialSummaryテーブルに財務要約情報を保存**

#### 【第三段階】比較分析（銘柄間比較）【順次処理】

8. **データ収集**
   - 全Phase 2推論が完了するまで待機
   - AnalysisResultテーブルから全銘柄の分析結果を取得
   - StockSummaryテーブルから全銘柄の株価要約情報を取得
   - FinancialSummaryテーブルから全銘柄の財務要約情報を取得

9. **Phase 3推論**
   - protoに入れてPythonに送る
     - 各銘柄のCode, Sentiment, SummaryReasoning
     - 各銘柄のStockSummary（株価要約）
     - 各銘柄のFinancialSummary（財務要約）
   - Python側で第三段階推論
     - 全銘柄の分析結果を比較
     - 株価要約と財務要約を参照しながら比較
     - Top3銘柄の選出
     - 各銘柄の選出理由を生成
     - セクター全体のサマリーを生成

10. **結果保存**
    - 返り値をprotoから受け取り
    - SectorAnalysisResultテーブルに保存

## データベース設計の詳細

### 既存テーブルの拡張

#### AnalysisResultテーブルに追加するカラム

Phase 1の推論結果を保存するためのフィールドを追加する。

- **BusinessModel** (string, type:text): Phase 1で理解したビジネスモデルの説明
  - 例: "SaaS企業で、クラウドベースのソフトウェア提供を主事業とする"
  
- **KPI** (string, type:text): Phase 1で決定した見るべきKPIのリスト（JSON形式またはカンマ区切り）
  - 例: "売上高成長率, 営業CF, 顧客獲得コスト"

### 新規テーブル

#### StockSummaryテーブル（株価要約情報）

Phase 2推論時にAIで生成される株価データの要約情報を保存する。

**目的**: Phase 3の比較分析で、生の株価データではなく要約情報を参照することで、メモリ効率を向上させる。

**主要フィールド**:
- Code: 銘柄コード
- AnalysisDate: 分析日時
- DataFrom/DataTo: データ期間
- AveragePrice: 平均株価
- PriceChangeRate: 期間内の変化率（%）
- Trend: トレンド（"upward", "downward", "sideways"）
- Volatility: ボラティリティ
- MaxPrice/MinPrice: 最高値/最安値
- VolumeTrend: 出来高トレンド
- SummaryText: AI生成の要約テキスト（自然言語での説明）
- AnalysisResultID: 関連するAnalysisResultへの外部キー

**ユニーク制約**: (Code, AnalysisDate) の組み合わせ

#### FinancialSummaryテーブル（財務要約情報）

Phase 2推論時にAIで生成される財務データの要約情報を保存する。

**目的**: Phase 3の比較分析で、生の財務データではなく要約情報を参照することで、メモリ効率を向上させる。

**主要フィールド**:
- Code: 銘柄コード
- AnalysisDate: 分析日時
- DataFrom/DataTo: データ期間
- GrowthStage: 成長段階（"growth", "mature", "decline"）
- RevenueGrowthRate: 売上高成長率（%）
- OperatingCF: 営業CF（百万円）
- NetIncome: 純利益（百万円）
- ROE: ROE（%）
- DebtRatio: 負債比率（%）
- KPIMetrics: Phase 1で決定したKPIの値（JSON形式）
- SummaryText: AI生成の要約テキスト（自然言語での説明）
- AnalysisResultID: 関連するAnalysisResultへの外部キー

**ユニーク制約**: (Code, AnalysisDate) の組み合わせ

#### SectorAnalysisResultテーブル（セクター分析結果）

第三段階の推論結果を保存する。Top3銘柄の選出結果を含む。

**主要フィールド**:
- SectorCode: セクターコード
- AnalyzedAt: 分析日時
- Top1Code/Top2Code/Top3Code: Top3銘柄のコード
- Top1Reasoning/Top2Reasoning/Top3Reasoning: 各銘柄の選出理由
- ComparisonLog: 比較分析の思考ログ
- OverallSummary: セクター全体のサマリー
- AnalysisResultIDs: 使用したAnalysisResultのIDリスト（カンマ区切り）

## 並列処理の詳細

### Phase 1の並列処理

- **実装方法**: Goのgoroutineを使用
- **同時実行数**: セマフォで制限（例: 最大10社同時）
- **エラーハンドリング**: 1社の処理が失敗しても、他の企業の処理は継続
- **結果収集**: channelを使用して結果を収集

### Phase 2の並列処理

- **実装方法**: Goのgoroutineを使用
- **同時実行数**: セマフォで制限（例: 最大10社同時）
- **エラーハンドリング**: 1社の処理が失敗しても、他の企業の処理は継続
- **結果収集**: channelを使用して結果を収集
- **要約情報生成**: Phase 2推論時に、株価と財務の要約情報も同時に生成

### Phase 3の実行タイミング

- **実行条件**: 全Phase 2推論が完了した時点で実行
- **処理方式**: 順次処理（全銘柄の結果を統合して比較するため）
- **データ参照**: StockSummaryとFinancialSummaryテーブルから要約情報を一括取得

## メモリ使用量の見積もり

### Phase 3でのメモリ使用量

- **Phase 2結果（50社）**: 約2.5MB
- **StockSummary（50社）**: 約100KB
- **FinancialSummary（50社）**: 約150KB
- **合計**: 約2.75MB

**結論**: GPT-4などのモデルは32K以上のコンテキストを扱えるため、この程度のデータ量は問題なく処理可能。

# 今後のタスク集
1. Repository層のメソッド追加（部分実装）
CreateOrUpdateStockSummary
CreateOrUpdateFinancialSummary
CreateOrUpdateSectorAnalysisResult
FindStockSummariesByCodes（Phase 3用）
FindFinancialSummariesByCodes（Phase 3用）
2. データ取得メソッドの実装（部分実装）
FindNewsByCode（newsテーブルから取得）
FindDailyQuotesByCode（stock_dataテーブルから取得）
FindFinancialStatementsByCode（fundamentalsテーブルから取得）
3. 個別のSyncメソッド実装（部分実装）
SyncNewsForCompany（1社分のnews取得→保存）
SyncStockDataForCompany（1社分のstock_data取得→保存）
SyncFundamentalsForCompany（1社分のfundamentals取得→保存）
4. 並列化処理の追加（部分実装）
SyncNewsForSector（セクター全体を並列で処理）
SyncStockDataForSector（セクター全体を並列で処理）
SyncFundamentalsForSector（セクター全体を並列で処理）

【レベル1】基盤実装（データ収集・保存）
├─ 1. DBテーブル作成（新規テーブル追加）
├─ 2. Repository層のメソッド追加
│   ├─ CreateOrUpdateStockSummary
│   ├─ CreateOrUpdateFinancialSummary
│   ├─ CreateOrUpdateSectorAnalysisResult
│   └─ Find系メソッド（Phase 3用）
│
├─ 3. データ取得メソッドの実装
│   ├─ FindNewsByCode
│   ├─ FindDailyQuotesByCode
│   └─ FindFinancialStatementsByCode
│
└─ 4. 個別のSyncメソッド実装
    ├─ SyncNewsForCompany（1社分）
    ├─ SyncStockDataForCompany（1社分）
    └─ SyncFundamentalsForCompany（1社分）

【レベル2】並列化と統合制御
├─ 5. 並列化処理の追加
│   ├─ SyncNewsForSector（セクター全体を並列）
│   ├─ SyncStockDataForSector（セクター全体を並列）
│   └─ SyncFundamentalsForSector（セクター全体を並列）
│
└─ 6. 段階制御ロジック（重要！）
    ├─ Phase 1→Phase 2の制御（全Phase 1完了後にPhase 2開始）
    ├─ Phase 2→Phase 3の制御（全Phase 2完了後にPhase 3開始）
    └─ エラーハンドリング（一部失敗時の処理）

【レベル3】gRPC連携
├─ 7. proto定義の作成
│   ├─ Phase 1用のRequest/Response
│   ├─ Phase 2用のRequest/Response
│   └─ Phase 3用のRequest/Response
│
├─ 8. gRPCクライアント実装（Go側）
│   ├─ Phase 1推論の呼び出し
│   ├─ Phase 2推論の呼び出し
│   └─ Phase 3推論の呼び出し
│
└─ 9. Python側からのレスポンス保存メソッド
    ├─ SavePhase1Result（AnalysisResultに保存）
    ├─ SavePhase2Result（AnalysisResult更新 + StockSummary/FinancialSummary保存）
    └─ SavePhase3Result（SectorAnalysisResultに保存）

【レベル4】Python側実装（別プロジェクト）
└─ 10. Python側のAI実装
    ├─ LangGraphの実装
    ├─ Phase 1推論ロジック
    ├─ Phase 2推論ロジック
    └─ Phase 3推論ロジック

推奨するタスクの順序
フェーズ1: 基盤実装（データ収集・保存）
DBテーブル作成 ← 現在ここ
Repository層のメソッド追加
データ取得メソッドの実装
個別のSyncメソッド実装
フェーズ2: 並列化と統合
並列化処理の追加
段階制御ロジック（Phase 1→2→3の制御）
フェーズ3: gRPC連携準備
proto定義の作成
gRPCクライアント実装（Go側）
Python側からのレスポンス保存メソッド
フェーズ4: Python側実装（別プロジェクト）
Python側のAI実装

「並列化処理の追加」には以下が含まれます：
並列化処理の実装（goroutine使用）
段階制御ロジック（Phase 1→2→3の遷移制御）