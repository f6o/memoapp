# ベースイメージとしてGoを使用
FROM golang:1.23 as builder

RUN apt-get update && \
    apt-get install -y make protobuf-compiler

# 作業ディレクトリを設定
WORKDIR /app

# 必要なファイルをコピー
COPY go.mod go.sum ./
RUN go mod tidy

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# ソースコードをコピー
COPY . .

# プロトコルバッファのコード生成
RUN make gen-proto

# サーバーとクライアントをビルド
RUN mkdir -p dist && \
    go build -o dist/memoapp-server ./main.go

# 実行用の軽量イメージを使用
FROM debian:bullseye-slim

# 必要なディレクトリを作成
WORKDIR /app

# ビルド成果物をコピー
COPY --from=builder /app/dist/memoapp-server /app/memoapp-server

# 必要なポートを公開 (例: gRPCサーバー用)
EXPOSE 50051

# デフォルトでサーバーを実行
CMD ["/app/memoapp-server"]