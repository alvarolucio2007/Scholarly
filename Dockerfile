FROM golang:1.27-alpine as builder
WORKDIR /app
COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scholarly cmd/api/*.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/scholarly .
EXPOSE 8080
CMD ["./scholarly"]
