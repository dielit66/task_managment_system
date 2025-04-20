FROM golang:1.24

WORKDIR /workspace/

COPY go.mod go.sum ./
RUN go mod download 

COPY . .
RUN go build -v -o ./app ./cmd

CMD [ "app" ]