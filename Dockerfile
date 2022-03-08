FROM golang:1.16.2

WORKDIR /app

RUN apt update && \
    apt install -y reflex

CMD ["reflex", "-r", "(\\.go$|go\\.mod)", "-s", "go", "run", "main/main.go"]
