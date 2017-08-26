FROM      gliderlabs/alpine

WORKDIR /myapp

COPY ./picresize /myapp

RUN pwd

CMD ["./picresize"]

