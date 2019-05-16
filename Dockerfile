FROM golang:1.12-alpine AS gobuild
WORKDIR /golang
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN apk add --update --no-cache git \
    && go mod download
COPY . .
RUN go build -o app .

FROM alpine:3.9
WORKDIR /golang
COPY --from=gobuild /golang/app /golang/app
EXPOSE 1323
CMD [ "./app" ]