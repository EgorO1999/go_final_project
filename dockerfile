FROM golang:1.24

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY ./web ./web
COPY ./pkg ./pkg 
COPY *.go ./  

ENV TODO_PORT=7540
ENV TODO_DBFILE="/app/scheduler.db"
ENV TODO_PASSWORD="12345"

EXPOSE 7540

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go_final_project

RUN chmod +x ./go_final_project

CMD ["./go_final_project"]    