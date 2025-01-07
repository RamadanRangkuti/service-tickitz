FROM golang

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o backend main.go

expose 8888

ENTRYPOINT ["/app/backend"]