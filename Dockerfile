FROM node:lts-alpine as UI_BUILDER
ARG VER
WORKDIR /ui
ADD /ui .
RUN npm i; npm run build; 

FROM golang:1.21.4 as GO_BUILDER
ARG VER
WORKDIR /server
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine
LABEL org.opencontainers.image.source="https://github.com/robrotheram/go-watch2gether"
WORKDIR /app
RUN apk upgrade -U \
 && apk add ca-certificates ffmpeg \
 && rm -rf /var/cache/*
RUN mkdir -p /app/ui
ADD app.sample.env /app/app.env
COPY --from=GO_BUILDER /server/w2g /app/w2g
COPY --from=UI_BUILDER /ui/dist /app/ui/dist
EXPOSE 8080
ENV GOMAXPROCS=100
ENTRYPOINT ["./w2g"]