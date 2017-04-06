FROM golang:1.8-onbuild

RUN mkdir -p /go/src/github.com/lanmj2004/task
WORKDIR /go/src/github.com/lanmj2004/task
COPY . /go/src/github.com/lanmj2004/task

VOLUME $PWD:/go/src/github.com/lanmj2004/task


RUN go-wrapper download
RUN go-wrapper install

ENV PORT 24115
# 给主机暴露 8080 端口，这样外部网络可以访问你的应用
EXPOSE 24115
# 告诉 Docker 启动容器运行的命令
CMD ["go-wrapper", "run"]