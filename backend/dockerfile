FROM golang:1.22

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /lorapanel-backend

EXPOSE 8080

CMD [ "/lorapanel-backend" ]