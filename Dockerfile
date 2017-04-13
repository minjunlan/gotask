FROM golang:1.8-onbuild

COPY . /go/src/app
RUN go get -d -v
RUN go install -v
EXPOSE 24115

CMD ["app"]