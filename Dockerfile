FROM golang:1.14
RUN apt-get update && apt-get install -y libjpeg-dev
WORKDIR /root
COPY . /app
RUN cd /app && go build -o gluedd-cli
ENTRYPOINT ["/app/gluedd-cli"]
