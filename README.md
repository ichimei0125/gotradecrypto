# gotradecypto

**!!注意!!** 本プロジェクトは学習目的です。実際の投資環境で使用する場合は、投資に伴うリスクを十分に理解し、本プログラムの仕組みを正しくご理解の上でご利用ください。本プログラムの使用により生じたいかなる損失や損害についても、当方は一切の責任を負いません。

## TL;DR
go を勉強したく、また最近仮想通貨は話題になり、本プロジェクトを始まりました。

## 使用
1. git cloneでソースコードをダウンロード
2. config.example.yamlをconfig.yamlリネームする
  - [bitflyer lighting](https://lightning.bitflyer.com/)をログインし、左側三本線を開き、APIをクリックする。新しいAPIを追加、「資産」と「トレード」は必要が、「入出金」では必要ない。API KEYとAPI SECRETをconfig.yamlに記入
  - invest_moneyは毎回投資の金額
  - cut_lossは保有する仮想通貨全部売ったあと、この金額より少ないなら取引停止（プログラムは停止しない）
  - ~~safe_moneyは口座中最低限の金額~~（現在バグ中）
  - dry_runがtrueの場合、取引しない、ただ価格を取得する。falseの場合取引する
3. 取引しない時、data/local.dbに約定履歴を書き込み、log/{取引所_仮想通貨}.logに価格や指標などを出力
4. configを設定したら、 `go build -o gotradecrypto cmd/main.go` でプログラムをビルドし、ターミナルで実行する
  - macの場合、scripts/mac/launchd にlaunchdスクリプトを用意した。パスを変更してください
5. BTCとXRPを例とする。他の仮想通貨で取引したい場合、cmd/main.go 中41行目　exchange.XRPJPYを他の仮想通貨に変更し、config.yamlのdry_runに仮想通貨を記入（定義はinternal/exchange/exchange.go）

> 複数の仮想通貨同時取引はプログラム上サポートしているが、取引所のAPI制限でエラー出る可能性ある。解決中が、しばらく1~2種類の仮想通貨を取引


## 構成概要

```
|-- cmd
|   `-- main.go                # エントリポイント
|-- internal
|   |-- common
|   |   |-- constants.go       # 定数
|   |-- exchange               # 取引所
|   |   |-- bitflyer           # 現在はbitflyerのみ実装
|   |   |-- exchange.go        # ※1
|   |-- indicator              # 指標の計算
|   `-- trade
|       |-- klinestrategies.go # 買い/売り判断:
|       |-- trade.go           # 取引 
|       `-- tradestrategies.go # 具体的な取引:
```

※1 他の取引所を適用したい場合、インターフェース「Exchange」を実装すれば動けるか。configやdbなど「取引所_仮想通貨」で一意の名前を区別するが、テストしたことない


## 取引ロジック
買い - 下記三つ条件全部達成した場合買い注文を発注

- 前時刻の終値がボリンジャーバンド(BBands)の $-2\sigma$ 以下
- 現在の価格はボリンジャーバンド(BBands)の $-2\sigma$ 以上
- ストキャスティクス（STOCK）の %D または %K が25以下

売り - 下記三つ条件全部達成した場合売り注文を発注

- 前時刻の終値がボリンジャーバンド(BBands)の $+2\sigma$ 以上
- 現在の価格はボリンジャーバンド(BBands)の $+2\sigma$ 以下
- ストキャスティクス（STOCK）の %D または %K が75以上

## milestone
- [x] v0.1 複数仮想通貨取扱
- [x] v0.2 約定履歴をdbに書き込む、API節約・データ分析目的
- [ ] v0.3 シミュレーター
- [ ] v0.5 AIを利用して、価格予測 
- [ ] v0.9 より豊富なデプロイ方法（docker, クラウドなど）