# Buildkit が docker build をどのように並列化しているのかを可視化する為のツールを作ったよ！

BuildKit が Docker に正式統合されて、ビルドが高速化しました。その要因の一つがビルドの並列化です。しかし、どこがどう並列化されているのかはよく分かっていない方も多いと思います。そこで今回、BuildKitがどのように並列処理しているのかをすぐに理解することのできるツール「llb2dot」を作ったので紹介します。こちら！！

## llb2dotは何をしてくれるのか

DockerfileもしくはLLBというビルドの過程で使用される中間言語を渡せば、グラフ構造をDOT言語に変換してくれます。DOT言語に変換することで、DOT言語をサポートする様々なツールが使用可能になります。例えばDOT言語にすることで下の図のように簡単にグラフ構造を図にできます。これでどこがどう並列処理されているのかが一眼でわかるようになります。

<img src="static/example.png">

## 使い方

GitHubからOSにあうバイナリをダウンロードしていただくか、Goの環境があれば下記でインストールできます。

```go
$ go get -u github.com/po3rin/llb2dot/cmd/llb2dot
```

Dockerfileのある場所で実行します(もしくはDockerfileのあるパスを指定する```-f```フラッグを使います)。

```console
$ llb2dot
```

そうするとDockerfileから生成されるLLBをDOT言語に変換してくれます。

```
go run cmd/llb2dot/main.go -f cmd/_testdockerllb/Dockerfile.test
strict digraph llb {
// Node definitions.
"[internal] helper image for file operations" [digest="sha256:e391a01c9f93647605cf734a2b0b39f844328cc22c46bde7535cec559138357b"];
"[internal] load build context" [digest="sha256:f8b67d93fc8b84291982a0262fcbfcc09194ac77e45bb56f8f1a2a474ff60730"];
"FROM golang:1.12" [digest="sha256:bd2b1f44d969cd5ac71886b9f842f94969ba3b2db129c1fb1da0ec5e9ba30e67"];
"ADD ./ /go" [digest="sha256:92a40b3e5547570a24293ed3023406c4ef8022c1bf4ed618f1d2773e6706d064"];
"RUN go build -o stage1_bin" [digest="sha256:f303bfaadb023a5106a78630b06c6f7bd59b042334e7a60db3a983c9a062e263"];
"RUN go build -o stage0_bin" [digest="sha256:1614fb7900bb8bfa5ba820832e96401e16535d6f58a29fc6f6fb5b9c58a6786f"];
"COPY --from=stage0 /go/stage0_bin /" [digest="sha256:1be70d79f63b8ff759cdd75f1de1ea3385ff0216ae5774b0ef469b212ae6ddc4"];
"COPY --from=stage1 /go/stage1_bin /" [digest="sha256:3fcd9cac9c924cf15e4ad7f811cab56f4362016076549057b03af3e53ac0b67c"];
"sha256:503..." [digest="sha256:5034fa69849c26ce76e1727501d65760883c554dfe0fa01df649f232c7cd19d1"];

// Edge definitions.
"[internal] helper image for file operations" -> "ADD ./ /go";
"[internal] helper image for file operations" -> "COPY --from=stage0 /go/stage0_bin /";
"[internal] helper image for file operations" -> "COPY --from=stage1 /go/stage1_bin /";
"[internal] load build context" -> "ADD ./ /go";
"FROM golang:1.12" -> "ADD ./ /go";
"FROM golang:1.12" -> "COPY --from=stage0 /go/stage0_bin /";
"ADD ./ /go" -> "RUN go build -o stage1_bin";
"ADD ./ /go" -> "RUN go build -o stage0_bin";
"RUN go build -o stage1_bin" -> "COPY --from=stage1 /go/stage1_bin /";
"RUN go build -o stage0_bin" -> "COPY --from=stage0 /go/stage0_bin /";
"COPY --from=stage0 /go/stage0_bin /" -> "COPY --from=stage1 /go/stage1_bin /";
"COPY --from=stage1 /go/stage1_bin /" -> "sha256:503...";
}
```

これはDOT言語と呼ばれるグラフ構造を書く為の言語です。これだけだと訳がわからないですが、DOT言語には様々な便利ツールがあります。例えば dot コマンドは DOT言語を受け取って、グラフ構造を図にしてくれます。

```
$ llb2dot | dot -T png -o out.png
```

これで下記のような画像ができます。Nodeの中はDOckerfileの場合はLLBのメタ情報が入ります。LLBから直で使った場合はDigest値を使って表示します。

<img src="static/example.png">

上の画像は下記のようなDockerfileを解析しました。3ステージのマルチステージビルドです。

```
FROM golang:1.12 AS stage0

WORKDIR /go
ADD ./ /go
RUN go build -o stage0_bin

FROM golang:1.12 AS stage1

WORKDIR /go
ADD ./ /go
RUN go build -o stage1_bin

FROM golang:1.12

COPY --from=stage0 /go/stage0_bin /
COPY --from=stage1 /go/stage1_bin /

```

このDockerfileとグラフ構造を見比べると、ステージの並列化だけでなく、ステージ同士の共通箇所も並列化できていることが分かります。LLBすごいですね。

## Go言語に処理を組み込む

パッケージはCLIだけでなく、コードとしても提供しています。独自のグラフ構造解析をしたいときに便利です。

```go
package main

import (
	"os"

	"github.com/po3rin/llb2dot"
)
func main(){
    // load llb (you want to load from Dockerfile? use LoadDockerfile)
    ops, _ := llb2dot.LoadLLB(os.Stdin)

    // convert graph
    g, _ := llb2dot.LLB2Graph(ops)

    // something with graph ...

    // write graph as dot language
    llb2dot.WriteDOT(os.Stdout, g)
}
```

## 作る際に使ったパッケージの紹介

ここからはGo言語の話ですが、内部で使っている便利パッケージも紹介します。

gonum.org/v1/gonum/graph

gonumは数学する際に便利なパッケージで有名ですが、その中にグラフ構造を扱う為のインターフェースが提供されています。このパッケージについての使い方の解説があまりにも少ないので今度記事にするかもです。

github.com/moby/buildkit

言わずもがな buildkit のパッケージです。Dockerfile から LLBに変換する部分や、LLBを直接受け取る処理をこれを使って書いています。

## おわりに

今回は Dockerfile が BuildKit でどのように並列化できるのかを簡単に掴めるツールを紹介しました。BuildKitそのものに関する記事も書いているので良かったらどうぞ！！