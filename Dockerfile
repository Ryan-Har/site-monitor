FROM golang:1.22.5

WORKDIR /app
COPY src ./
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o /site-monitor
EXPOSE 3000

CMD ["/site-monitor"]