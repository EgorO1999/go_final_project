FROM golang:1.24

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./web ./web
COPY todoapp ./todoapp

ENV TODO_PORT=7540
ENV TODO_DBFILE="/app/scheduler.db"
ENV TODO_PASSWORD="12345"

EXPOSE 7540

CMD ["./todoapp"]    