FROM alpine:latest

RUN mkdir /

COPY binary_file/gymondoApp /
COPY docker_env/.env /

CMD [ "/gymondoApp" ]
