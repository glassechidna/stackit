FROM gcr.io/distroless/static
COPY stackit /usr/bin/
CMD ["/usr/bin/stackit"]
