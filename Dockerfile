FROM golang:1.18
WORKDIR /dist
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o fampay .
EXPOSE 3000
CMD  ["./fampay"]