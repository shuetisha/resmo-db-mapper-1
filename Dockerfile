FROM alpine:latest
COPY resmo-db-mapper /resmo-db-mapper/resmo-db-mapper
WORKDIR /app
ENTRYPOINT ["/resmo-db-mapper/resmo-db-mapper"]