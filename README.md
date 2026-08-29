# Social-game (backend)
**[概要]** 
ソーシャルゲームのバックエンドAPI

**[構成]** 
ログイン・放置報酬・ガチャ

**[技術スタック]** 
* Go 1.25.7
* mysql v1.8.1 (Docker Compose管理)

**[ディレクトリ]** 
social-game/
├── docker-compose.yml
├── go.mod
├── go.sum
├── migrations/
│   ├── 0001_init.sql      # users / auth_tokens / user_resources
│   └── 0002_gacha.sql     # characters / user_characters
├── cmd/
│   └── server/
│       └── main.go        # エントリーポイント・ルーティング
└── internal/
    ├── model/              # DBテーブルに対応する構造体
    ├── repository/         # DBアクセス層
    ├── service/             # 放置計算・トークン生成・ガチャ抽選ロジック
    └── handler/              # HTTPハンドラ


**[API一覧]** 
メソッド	パス	　　　認証	　　　　概要
POST　/auth/login	　不要	　　　device_idでログイン。無ければ新規作成してトークンを発行
GET	　/user/state　　必要	　　　現在のコイン・生産レート・未受取の放置分を取得（表示専用）
POST　/user/claim	　必要	　　　放置分を確定してコインに加算
POST　/gacha/draw　必要	　　　100コイン消費してキャラを1体抽選（均等確率）
GET	　/healthz　　　不要	　　　死活確認用

**[設計]** 
* 放置報酬はバッチ処理ではなく都度計算方式
* 放置報酬には12時間の上限
* ガチャのコイン消費はSELECT ... FOR UPDATEで行ロックしてから減算
* ガチャは現状キャラのみ・均等確率