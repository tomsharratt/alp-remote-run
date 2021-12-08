FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /alp-remote-run

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /alp-remote-run /alp-remote-run

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/alp-remote-run"]
