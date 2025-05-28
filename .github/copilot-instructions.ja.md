# Copilot Instruction for gogh Project

このドキュメントは、goghプロジェクトの設計と構造を理解するためのガイドです。特に、プログラムのアーキテクチャと依存関係のルールに焦点を当てています。

## Copilotの前提

Goのプロとして適切なコードを記述してください。

### 出力の言語

説明は主に日本語を使ってください。
コード内のコメント、エラーメッセージ、テストケース名などは常に英語を用いてください。

### プロジェクトの在処

gogh は、github.com/kyoh86/gogh というリポジトリで管理されているGo言語のプロジェクトです。
Goのパッケージ名はバージョン名を接尾辞に持つ可能性があるので、go.mod の `module` 行を確認してください。

### テストの書き方

privateな関数やメソッドに対してテストを書く際は、privateなテストを書いてください。
privateなテストは、`*_privavte_test.go` というファイル名で、同じパッケージ内に配置してください。

exportedな関数やメソッドに対してテストを書く際は、`*_test.go` というファイル名で、パッケージ名に `_test`を付けて、同じディレクトリに配置してください。
exportedなテストパッケージの対象をimportするときは、 `testtarget` という名前を付けてimportしてください。

例えば `github.com/kyoh86/gogh/vx/core/foo` というパッケージのテストのexportedな関数やメソッドに対してテストを書く場合は、以下のようにimportしてください。

```go
import (
    testtarget "github.com/kyoh86/gogh/vx/core/foo"
)
```

テストにおいてはtestify/assertやtestify/requireは使用せず、標準の`testing`パッケージを使用してください。

### Mockの置き方

Mockは基本的にmockgenでI/F定義したパッケージと同じ階層の別パッケージに生成してください。
mockgenは、依存先として "go.uber.org/mock/gomock" を要求します。
gen.goを置いて、`go:generate`でmockgenを実行するようにしてください。

例(ディレクトリ構造):

```
gogh/
├── core/                # コア機能 - プログラムの本質
│   ├── gen.go           # Mock生成用のファイル
│   ├── repository/      # リポジトリ関連の定義
│   ├── repository_mock/ # リポジトリ関連の定義
│   └── ...
...
```

例(gen.go):

```go
package core

//go:generate go tool mockgen -source ./repository/default_name_service.go -destination ./repository_mock/gen_default_name_service_mock.go -package repository_mock
//go:generate go tool mockgen -source ./repository/location_format.go      -destination ./repository_mock/gen_location_format_mock.go      -package repository_mock
```

## プロジェクトの目的

主にGitHubを対象としたリポジトリ管理ツールであり、リポジトリのクローン、作成、削除、リスト表示などの機能を提供します。
goghは、リポジトリの参照をローカルで構造化して配置し、一貫した方法で処理し、ユーザーが簡単に操作できるように設計されています。


## プログラム構造の全体像

### 1. 4つの主要な関心事（レイヤー）

#### コア機能層 (core)

場所: `core/`

- 役割
    - プログラムが実行する本質的な操作と概念の定義
- このプロジェクトでの例
    - リポジトリ参照の処理、リポジトリ操作の定義
- 定義するもの
    - ドメインエンティティ（具体的実装：`Reference`, `BaseRepository`など）
    - 外部システムとのやり取りに必要なインターフェース（`RepositoryService`など）
- 依存するもの
    - なし（最内層）
- 実装するもの
    - ドメインエンティティと純粋なロジック
    - 外部システムやUIに依存しない純粋な機能

#### アプリケーション層 (app)

場所: `app/`

- 役割
    - ユースケースの実現とコア機能の調整
- このプロジェクトでの例
    - リポジトリのクローン処理、リスト取得ロジック
- 依存するもの
    - 
    - コア層のインターフェースのみ
- 実装するもの
    - コア機能を組み合わせて実際の処理を行う
    - ユースケース全体のフローの実装（`CloneService`など）
    - サービスオブジェクト、ユースケースハンドラ

#### 外部システム連携層 (infra)

- 役割
    - プログラムの外部との連携や技術的実装
- このプロジェクトでの例
    - GitHub API呼び出し、トークン保存、ファイル操作
- 依存するもの
    - コア層のインターフェースと型
- 実装するもの
    - APIクライアント、データベース接続、ファイル操作など
    - コア層で定義されたインターフェース（`RepositoryService` → `GitHubRepositoryService`）の実装

#### ユーザーインターフェース層 (ui)

- 役割
    - ユーザーとのインタラクション方法を定義する
- このプロジェクトでの例
    - CLIコマンド、出力フォーマット
- 特徴
    - ユーザーの要求をアプリケーション層に伝える
- 依存するもの
    - アプリケーション層のサービス
- 実装するもの
    - コマンドやビュー
    - コマンド実装、表示ロジック、入力処理

### 2. 実践的な依存関係ルール

1. コア機能層の依存ルール
   - コア機能層は他のどの層もインポートしない
   - (他の層がコア層をインポートすることは可能)
   - 理由: コアはアプリケーション特有の処理や外部連携に依存すべきでない

2. アプリケーション層の依存ルール
   - アプリケーション層はコア機能層をインポートできる
   - UIや外部システム層をインポートしない
   - 理由: ユースケースはコア機能を使うが、具体的な実装方法やUIには依存しない

3. 外部システム連携層の依存ルール
   - 外部システム連携層はコア機能層をインポートできる
   - UIやアプリケーション層をインポートしない
   - 理由: 外部連携はコアのインターフェースを実装するが、UIやユースケースの詳細を知る必要はない

4. UI層の依存ルール
   - アプリケーション層をインポートできる
   - 理由: UIはユースケースを呼び出して実行する

### 3. インターフェースによる連携

- コア層がインターフェースを定義
- 外部システム層がそれを実装
- アプリケーション層がインターフェースを使用

#### コア層でのインターフェース定義例

```go
// core/repository/service.go
type RepositoryService interface {
   Get(name string) (Repository, error)
   List() ([]Repository, error)
}
```

#### アプリケーション層でのユースケース例

```go
// app/clone/service.go
type CloneService struct {
   repoService core.RepositoryService
}

func (s *CloneService) CloneRepository(name string) error {
   repo, err := s.repoService.Get(name)
   if err != nil {
       return err
   }
   // クローン処理を実装
}
```

#### 外部連携層での実装例

```go
// infra/github/client.go
type githubClient struct {...}

// コア層のインターフェースを実装
func (c *githubClient) Get(name string) (core.Repository, error) {
   // GitHub APIを使った実装
}
```

#### メイン関数（または依存性注入コンテナ）での依存性注入

```go
func main() {
    // 1. インフラ層の実装を作成
    tokenStore := config.NewTokenStore()
    githubService := github.NewRepositoryService(tokenStore)

    // 2. アプリケーション層のサービスを作成し、インフラ実装を注入
    cloneService := clone.NewService(githubService)

    // 3. UI層のコマンドを作成し、アプリケーションサービスを注入
    cloneCommand := commands.NewCloneCommand(cloneService)

    // 4. アプリケーション実行
    rootCmd.AddCommand(cloneCommand)
    rootCmd.Execute()
}
```

#### ファクトリによる依存性注入

```go
// app/clone/factory.go
func NewCloneServiceWithGitHub(tokenStore core.TokenStore) *Service {
    githubService := github.NewRepositoryService(tokenStore)
    return NewService(githubService)
}
```

#### 注入フローの全体マップ

| 層               | 依存対象                     | 注入される実装        | 注入される場所          |
|------------------|------------------------------|-----------------------|-------------------------|
| UI               | アプリケーション層のサービス | `app.CloneService`    | main関数/コマンド初期化 |
| アプリケーション | コア層のインターフェース     | `infra.GitHubService` | main関数/ファクトリ     |
| 外部システム     | コア層のインターフェース     | - (自身が実装側)      | -                       |
| コア             | -                            | -                     | -                       |

### 4. 推奨ディレクトリ構造

```
gogh/
├── core/                # コア機能 - プログラムの本質
│   ├── repository/      # リポジトリ関連の定義
│   ├── repository_mock/ # リポジトリ関連の定義
│   ├── ...
│   ├── auth_mock/       # 認証関連のMock
│   └── auth/            # 認証関連の定義
│
├── app/                 # アプリケーション - ユースケース
│   ├── clone/           # クローン機能のユースケース
│   └── list/            # リスト機能のユースケース
│
├── infra/               # 外部システム連携 - 外の世界との接続
│   ├── github/          # GitHub API実装
│   └── filesystem/      # ファイルシステムに依存する実装
│
└── ui/                  # ユーザーインターフェース - 使い方
    └── cli/             # コマンドライン実装
```

### 5. 実際のケーススタディ

クローン操作を例にすると

1. UI層: `clone` コマンドがユーザーからのリポジトリ指定を受け取る
2. アプリケーション層: `CloneService` がクローン操作全体を調整
3. 外部システム層: GitHub APIクライアントがリポジトリ情報を取得
4. コア層: リポジトリインターフェースがクローンに必要な情報を定義

### 6. コア層の2種類の要素

コア層には2種類の要素がある

#### 6-1. インターフェース定義

- 外部システム層が実装するもの
- 例：`RepositoryService`, `RemoteService`など
- これらは技術的実装（GitHub API、ファイルシステムなど）に依存する操作

#### 6-2. 具体的な実装

- コア層自身が実装するもの
- 例：`Reference`型の実装、パース処理、バリデーションロジックなど
- 外部システムに依存しない純粋なドメインロジック

#### 6-3. 実例

```go
// core/repository/reference.go - コア層で実装される
type reference struct { // 具象型
    host string
    owner string
    name string
}

func (r *reference) Host() string { return r.host } // 実装

// core/repository/parser.go - コア層で実装される
func Parse(s string) (Reference, error) { // 実装
    // パース処理...
    return &reference{...}, nil
}

// core/repository/service.go - 外部システム層が実装する
type RepositoryService interface { // インターフェースのみ
    Get(ref Reference) (Repository, error)
    // ...
}
```

#### 6-4. 判断基準

コア層で実装するか、インターフェースだけにするかの基準：

1. 外部システム（API、DB、ファイルなど）に依存するか？
   - Yes → インターフェースのみ定義（実装は外部システム層）
   - No → コア層で実装可能

2. 純粋なドメインロジックか？
   - Yes → コア層で実装
   - No → 適切なレイヤーで実装

## goghのコアルールの設計

goghにおける「コアルール」は以下の通り

### 1. リポジトリの概念モデル

- リポジトリ参照の定義: `github.com/owner/repo`形式の解釈ルール
- リポジトリの構成要素: 名前、オーナー、URL形式などの基本属性
- リポジトリ間の関係: フォークとオリジン、親子関係など

### 2. 構造化・組織化のルール

- リポジトリのグループ化: 複数リポジトリの管理方法
- ルート下の構造: `<root>/<host>/<owner>/<name>`というパス構成
- 名前空間の管理: 同名リポジトリの衝突回避ルール

### 3. 操作の基本概念

- クローンの定義: リモートからローカルへの複製とは何か
- 検索の抽象概念: リポジトリをどう検索するか

### 4. ユーザーとリポジトリの関係

- ユーザー所有とアクセス権: ユーザーとリポジトリの関係性
- 権限モデル: 読み取り/書き込み権限の概念
