FROM golang:latest as dev

ENV CGO_ENABLED 1

RUN apt-get update -qq
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev

RUN apt-get install -y -qq \
    tesseract-ocr-eng \
    tesseract-ocr-jpn

WORKDIR /go/src/
COPY src/go.sum src/go.mod ./

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
RUN go mod download

ARG PORT
EXPOSE ${PORT}

COPY ./src . 

CMD ["air", "-c", ".air.toml"]
