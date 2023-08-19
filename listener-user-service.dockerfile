FROM alpine:latest

RUN mkdir /app

COPY listenerUserApp /app

CMD [ "/app/listenerUserApp"]