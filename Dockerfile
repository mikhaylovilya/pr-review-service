FROM golang:alpine

RUN apk update && apk upgrade && apk add git bash openssh

COPY . .

RUN go mod download

RUN go build -o main .

EXPOSE 3081

CMD ["./main"]