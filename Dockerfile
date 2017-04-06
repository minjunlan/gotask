FROM golang:latest

RUN mkdir -p /go/src/web-app
WORKDIR /go/src/web-app
COPY . /go/src/web-app

RUN go get github.com/codegangsta/gin
RUN go-wrapper download
RUN go-wrapper install

ENV PORT 8080
# 给主机暴露 8080 端口，这样外部网络可以访问你的应用
EXPOSE 8080 3030
CMD gin run
# 告诉 Docker 启动容器运行的命令
CMD ["go-wrapper", "run"]