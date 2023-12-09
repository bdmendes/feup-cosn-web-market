FROM golang:alpine

COPY . /app
WORKDIR /app

RUN go mod download
RUN go build -o main .

EXPOSE 8001
CMD /app/main
