FROM golang:1.21 as build

WORKDIR /app
COPY . ./

RUN go mod download
RUN go build -o app.so ./cmd/server/main.go

FROM debian

EXPOSE 8080

WORKDIR /
COPY --from=build /app/app.so .

ENTRYPOINT [ "./app.so" ]
