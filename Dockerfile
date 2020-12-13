From golang:alpine

WORKDIR /build

COPY . .

RUN apk add --update make
RUN make build

EXPOSE 8080

CMD ["./main"]
