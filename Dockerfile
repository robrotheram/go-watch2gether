FROM node:lts-alpine as UI_BUILDER
ARG VER
WORKDIR /ui
ADD /ui .
RUN sed -i "s/{WATCH2GETHER_VERSION}/$VER/g" index.html
RUN npm i; npm run build; 

FROM golang:1.20.1 as GO_BUILDER
ARG VER
WORKDIR /server
ADD /server .
RUN sed -i "s/{WATCH2GETHER_VERSION}/$VER/g"  pkg/datastore/version.go 
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine
WORKDIR /app
RUN apk upgrade -U \
 && apk add ca-certificates ffmpeg \
 && rm -rf /var/cache/*
RUN mkdir -p /app/ui
COPY server/app.sample.env /app/app.env
COPY --from=GO_BUILDER /server/watch2gether /app/watch2gether
COPY --from=UI_BUILDER /ui/dist /app/ui/build
EXPOSE 8080
ENTRYPOINT ["./watch2gether"]
