# searcher

検索アプリのAPIです。

## 構成図

![](./docs/architecture.png)

## OpenSearch について

### 検索

- query
  - 検索条件。かなり複雑な指定ができる。
- filter
  - 検索結果の絞り込み条件。スコアに影響しない。
- source
  - 検索結果のフィールド指定。
- highlight
  - 検索結果のハイライト指定。
- size
  - 検索結果の取得件数。
- from
  - 検索結果の取得開始位置。

### template

#### search template の登録および更新

```
POST _scripts/{テンプレート名}
{
  "script": {
    "lang": "mustache",
    "source": """
    {template.mustache}
    """,
    "params": {}
  }
}
```

#### search template のレンダリング

```
POST _render/template
{
  "id": "{テンプレート名}",
  "params": {
    "{パラメータ}": "{値}"
  }
}
```

#### 未登録の search template を使った検索

```
GET _search/template
{
  "source": {テンプレート},
  "params": {
    "{パラメータ}": "{値}"
  }
}
```

#### search template を使った検索

```
GET /{インデックス名}/_search/template
{
  "id": "{テンプレート名}",
  "params": {
    "{パラメータ}": "{値}"
  }
}
```