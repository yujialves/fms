# fms
LINEのFlex Message Simulatorをプログラム上で実行可能にするライブラリ

## 導入手順

## 環境変数の設定
以下のように、LINEのユーザー情報を予め環境変数として読み込んでおく必要があります。
```env:.env
EMAIL=example@gmail.com
PASSWORD=12345678
```

## 使用例
```go:main.go
// jsonの読み込み
json, err := os.Open("./example.json")
if err != nil {
	log.Fatal(err)
}

// jsonからdomの文字列を生成
dom, err := fms.Generate(json)
if err != nil {
	log.Fatal(err)
}
```
