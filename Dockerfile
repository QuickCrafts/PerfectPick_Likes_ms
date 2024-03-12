FROM golang:latest

RUN mkdir /build
WORKDIR /build

COPY . /build

RUN cd /build && go build -o bin/PerfectPick_Likes_ms .

EXPOSE 3000

ENTRYPOINT ["/build/bin/PerfectPick_Likes_ms"]