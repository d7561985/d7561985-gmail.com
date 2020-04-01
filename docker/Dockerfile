FROM golang:1 as build

MAINTAINER Dzmitry Harupa <dzmirty.harupa@1pt.com>
ARG BUILD_DATE
ARG VCS_REF

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="question-service" \
      org.label-schema.description="Short description of application" \
      org.label-schema.url="https://github.com/d7561985/questions/" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/d7561985/questions" \
      org.label-schema.vcs-type="Git" \
      org.label-schema.docker.cmd="/question-service"

COPY . /build
WORKDIR /build
RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -o /app /build/main.go

# multi-staging
FROM golang:1-alpine

ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/ca-certificates.crt
ADD data/ /data

COPY --from=build /app /app

RUN addgroup -g 1000 -S username && \
    adduser -u 1000 -S username -G username -h /data

# ToDo: current version require permission to folder data
# It's unsefe! So please not use in production till solve that.
#USER 1000
EXPOSE 8080/tcp

ENV DELIVERY_HTTP_ADDR :8080
ENV PROMETHEUS_HTTP_ADDR :8089

WORKDIR /
ENTRYPOINT ["/app"]
CMD ["runserver"]
