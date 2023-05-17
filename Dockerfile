FROM alpine:latest
COPY resmo-db-mapper /resmo-db-mapper
ENTRYPOINT ["/resmo-db-mapper"]