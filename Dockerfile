# FROM golang:latest

# RUN mkdir /build
# WORKDIR /build

# COPY . /build

# RUN cd /build && go build -o bin/PerfectPick_Likes_ms .

# EXPOSE 3000

# ENTRYPOINT ["/build/bin/PerfectPick_Likes_ms"]

FROM golang:latest

RUN mkdir /build
WORKDIR /build

COPY . /build

EXPOSE 3000

# Install Dockerize
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

RUN cd /build && go build -o bin/PerfectPick_Likes_ms .


ENTRYPOINT ["dockerize", "-wait", "tcp://neo4j:7687", "-timeout", "60s", "/build/bin/PerfectPick_Likes_ms"]
