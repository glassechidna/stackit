FROM alpine
RUN apk add --update --no-cache ca-certificates
COPY stackit /usr/bin/
CMD ["/usr/bin/stackit"]
