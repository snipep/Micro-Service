FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY frontEndApp /app
    
COPY ./cmd/web/templates /app/cmd/web/templates

CMD [ "/app/frontEndApp"]