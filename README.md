# gorm のリレーション制御サンプル

[Gorm v2](https://gorm.io/) の内容。

## モデル定義

モデル定義は下記のように行う。

### 構造体の作成とタグ設定

`gorm` ではテーブルと構造体を関連づけて更新・参照を行う。

- `gorm:"column:fruit_id"` みたいなカラム名のタグ設定

  - 必須ではない。
  - 構造体のフィールド名を `gorm` がよしなにスネークケースとして認識してくれる。
  - 構造体のフィールド名と DB のカラム名が異なる場合は設定する。

- `gorm:"foreignKey:CustomerID;references:CustomerID;"` みたいな外部キー設定
  - あくまで `gorm` がリレーションを認識するための設定であり、DB に反映されるわけではない。
  - この構造体定義から `gorm` 経由でテーブル作成を行う場合は DB 側の設定として定義される？（未検証なので不明）
  - `foreignKey` には参照先の構造体のキーを設定し、 `references` には自身のメンバのキーを設定する。
    - 構造体のフィールド名を設定する。

```
type Fruit struct {
	FruitID    uint     `gorm:"column:fruit_id"`
	MarketID   uint     `gorm:"column:market_id"`
	Fruit      string   `gorm:"column:fruit_name"`
	CustomerID uint     `gorm:"column:customer_id"`
	Customer   Customer `gorm:"foreignKey:CustomerID;references:CustomerID;"`
}

type Customer struct {
	CustomerID   uint   `gorm:"column:customer_id"`
	CustomerName string `gorm:"column:customer_name"`
}

type Market struct {
	MarketID   uint    `gorm:"column:market_id"`
	MarketName string  `gorm:"column:market_name"`
	Fruits     []Fruit `gorm:"foreignKey:MarketID;references:MarketID;"`
}
```

- （おまけ）リレーションではなく構造体を埋め込みたいだけの場合、 `embedded` タグを使う。
  - JOIN の結果を Scan したりするときはリレーションではなく構造体の埋め込みで結果取得を行う
  ```
  type Fruit struct {
      FruitID    uint     `gorm:"column:fruit_id"`
      MarketID   uint     `gorm:"column:market_id"`
      Fruit      string   `gorm:"column:fruit_name"`
      CustomerID uint     `gorm:"column:customer_id"`
      Customer   Customer `gorm:"embedded"` // ただの埋め込み
  }
  ```

### 構造体名とテーブル名

構造体のフィールド名とテーブルのカラム名はタグでよしなにリンクさせることができるが、
構造体名とテーブル名は `gorm` 側で用意されているインタフェースを実装させる必要がある。
参考: https://gorm.io/ja_JP/docs/conventions.html#%E3%83%86%E3%83%BC%E3%83%96%E3%83%AB%E5%90%8D

```
type Fruit struct {
	FruitID    uint     `gorm:"column:fruit_id"`
	MarketID   uint     `gorm:"column:market_id"`
	Fruit      string   `gorm:"column:fruit_name"`
	CustomerID uint     `gorm:"column:customer_id"`
	Customer   Customer `gorm:"foreignKey:CustomerID;references:CustomerID;"`
}

// Fruit 構造体と fruit テーブルを関連づける
func (f *Fruit) TableName() string {
	return "fruit"
}
```

- `gorm` 側のインタフェース

  ```
  // https://github.com/go-gorm/gorm/blob/5daa413f418d8b745d5e7178b07405b0a215f5f2/schema/schema.go#L70-L72
  type Tabler interface {
      TableName() string
  }
  ```

- 用途
  - `gorm` がデフォルトで参照しようとするテーブル名は「構造体名+"s"」であるため、そのようになっていない場合は利用者側で定義する必要がある
    ```
    // TableName で独自設定をしてない場合は fruits テーブルを参照しようとする。
    db.Model(&Fruit{})
    ```

## `has many` と `has one`

### `has many`

- 下記のように構造体のメンバに別の構造体の配列が定義されている状態は `has many` を表す。

  ```
  type Market struct {
      MarketID   uint    `gorm:"column:market_id"`
      MarketName string  `gorm:"column:market_name"`
      Fruits     []Fruit `gorm:"foreignKey:MarketID;references:MarketID;"` // has many
  }

  type Fruit struct {
      FruitID    uint     `gorm:"column:fruit_id"`
      MarketID   uint     `gorm:"column:market_id"`
      Fruit      string   `gorm:"column:fruit_name"`
      CustomerID uint     `gorm:"column:customer_id"`
  }
  ```

  - 上記場合、下記のようにレコード参照を行うことができる。
    ```
    var data []Market
    db.Table("market").
            Preload("Fruits"). // Preload で Fluits (= fluit テーブル)をキャシングしている。必ず構造体名+"s"の表記にすること。
            Find(&data)
    ```

### `has one`

- 下記のように構造体のメンバに別の構造体が定義され、外部キー設定が定義されている状態は `has one` を表す。

  ```
  type Fruit struct {
      FruitID    uint     `gorm:"column:fruit_id"`
      MarketID   uint     `gorm:"column:market_id"`
      Fruit      string   `gorm:"column:fruit_name"`
      CustomerID uint     `gorm:"column:customer_id"`
      Customer   Customer `gorm:"foreignKey:CustomerID;references:CustomerID;"` // has one
  }

  type Customer struct {
      CustomerID   uint   `gorm:"column:customer_id"`
      CustomerName string `gorm:"column:customer_name"`
  }
  ```

  - 上記場合、下記のようにレコード参照を行うことができる。
    ```
    var data []Fruit
    db.Table("fluit").
            Preload("Customer"). // Preload で Customer (= customer テーブル)をキャシングしている。`has many` の場合とは異なり、構造体名の末尾に "s" はつけない。
            Find(&data)
    ```

### ネストされている場合

- `has many` の中で `has many` / `has one` をしている場合などは下記のようにレコード参照を行うことができる。

  ```
  type Fruit struct {
      FruitID    uint     `gorm:"column:fruit_id"`
      MarketID   uint     `gorm:"column:market_id"`
      Fruit      string   `gorm:"column:fruit_name"`
      CustomerID uint     `gorm:"column:customer_id"`
      Customer   Customer `gorm:"foreignKey:CustomerID;references:CustomerID;"` // （Market から見て） has many の中で has one
  }

  type Customer struct {
      CustomerID   uint   `gorm:"column:customer_id"`
      CustomerName string `gorm:"column:customer_name"`
  }

  type Market struct {
      MarketID   uint    `gorm:"column:market_id"`
      MarketName string  `gorm:"column:market_name"`
      Fruits     []Fruit `gorm:"foreignKey:MarketID;references:MarketID;"` // has many
  }

  var data []Market
  db.Table("market").
          Preload("Fruits.Customer"). // Fruit 構造体内の Customer メンバを表している
          Preload("Fruits"). // Fruit も忘れずに Preload
          Find(&data)
  ```
