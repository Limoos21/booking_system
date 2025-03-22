FROM ubuntu:latest
LABEL authors="nikita"

ENTRYPOINT ["top", "-b"]