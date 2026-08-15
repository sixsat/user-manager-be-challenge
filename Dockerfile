FROM golang:1.26-alpine3.24 AS base

WORKDIR /app

RUN apk update && apk add --no-cache tzdata

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o api .

####################

FROM gcr.io/distroless/static-debian13:nonroot

ENV TZ=Asia/Bangkok

COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /app/api .
COPY --from=base /app/config/config.yml ./config/

USER nonroot

CMD ["./api"]