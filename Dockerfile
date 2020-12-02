FROM alpine
RUN mkdir -p /app/ui
ADD watch2gether /app/watch2gether
ADD ui/build /app/ui/build
EXPOSE 8080
WORKDIR /app
ENTRYPOINT ["./watch2gether"]
