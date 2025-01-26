# gotradecypto

**!!注意!!** 本プロジェクトは学習目的です。実際の投資環境で使用する場合は、投資に伴うリスクを十分に理解し、本プログラムの仕組みを正しくご理解の上でご利用ください。本プログラムの使用により生じたいかなる損失や損害についても、当方は一切の責任を負いません。

## TL;DR
go を勉強したく、また最近仮想通貨は話題になり、本プロジェクトを始まりました。

## 使用 (開発中ので不安定可能性あり、安定版リリース予定v0.1.0)
1. git cloneでソースコードをダウンロードし、Maria DBをインストールし、 `scripts/db/init_mariadb.sql` のパスワードを設定してDB初期化
2. config.example.yamlをconfig.yamlリネームする
   - [bitflyer lighting](https://lightning.bitflyer.com/)をログインし、左側三本線を開き、APIをクリックする。新しいAPIを追加、「資産」と「トレード」は必要が、「入出金」では必要ない。API KEYとAPI SECRETをconfig.yamlに記入
   - invest_moneyは毎回投資の金額
   - cut_lossは保有する仮想通貨全部売ったあと、この金額より少ないなら取引停止（プログラムは停止しない）
   - safe_moneyは口座中最低限の金額
   - symbolsは取引したい取引所と仮想通貨
    
     ※現時点仮想通貨-日本円のペアだけサポートので、ETH-BTC, BCH-BTCは使わないてください
   - 1に設定されたパスワードをconnection_stringに
4. `go build -o gotradecrypto cmd/main.go` でプログラムをビルドし、ターミナルで実行する
   - macの場合、scripts/mac/launchd にlaunchdスクリプトを用意した。パスを変更してください

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
- [ ] v0.1.0 初期リリース
- [ ] v0.0.5 より豊富なデプロイ方法（docker, クラウドなど）
- [ ] v0.0.4 configやguiなど、簡単に取引ロジック（インジケーターで買売判断）を作成
- [x] v0.0.3 シミュレーター
- [x] v0.0.2 約定履歴をdbに書き込む、API節約・データ分析目的
- [x] v0.0.1 複数仮想通貨取扱

## Experiment
- [ ] LSTMを利用して、価格予測 
- [ ] LLMを利用して、ニュース解析により、価格予測
- [ ] RL/NNを使って、OrderBook・出来高を見てリスク回避
