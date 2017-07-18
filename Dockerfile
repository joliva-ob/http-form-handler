#
# Dockerfile for Form Handler project
#

FROM golang

MAINTAINER Joan Oliva

RUN mkdir /formhandler
RUN mkdir /formhandler/bin
RUN mkdir /formhandler/cfg
RUN mkdir /formhandler/logs

ADD *.yml /formhandler/cfg/

ENV CONF_PATH /formhandler/cfg
ENV ENV pro

CMD GOOS=linux GOARCH=386 go build -o formhandler

ADD formhandler /formhandler/bin/

ENTRYPOINT /formhandler/bin/formhandler

EXPOSE 8001