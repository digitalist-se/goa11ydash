FROM golang:1.13.4
WORKDIR /go/src/github.com/digitalist-se/goa11ydash
COPY . .
RUN go get github.com/labstack/echo/v4
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/digitalist-se/goa11ydash .
CMD ["./app"]  

