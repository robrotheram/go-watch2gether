FROM alpine
RUN apk upgrade -U \
 && apk add ca-certificates ffmpeg \
 && rm -rf /var/cache/*
RUN mkdir -p /app/ui
ADD watch2gether /app/watch2gether
ADD ui/build /app/ui/build
ADD server/app.env /app/app.env
EXPOSE 8080
WORKDIR /app
ENTRYPOINT ["./watch2gether"]
