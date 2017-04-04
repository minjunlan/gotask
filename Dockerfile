FROM golang:1.8-onbuild

RUN mkdir -p /go/src/mytask
WORKDIR /go/src/mytask
COPY . /go/src/mytask
CMD ["go-wrapper", "run"]


