# For build run: podman build -t rpn_image -f Dockerfile
FROM docker.io/library/golang:1.23-alpine AS go_builder
ARG work_dir=/rpn

ADD https://github.com/AlexS25/rpn/archive/refs/heads/master.zip ./
RUN unzip master.zip && mv rpn-master $work_dir; rm master.zip
WORKDIR $work_dir
RUN go build ./cmd/main.go

