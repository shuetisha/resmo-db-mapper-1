FROM alpine:latest
COPY resmo-db-mapper /resmo-db-mapper/resmo-db-mapper
WORKDIR "/resmo-db-mapper"
ENV SCHEDULE=""
ENV CONTEXT_TIME=""
ENV DATASOURCE_NAME=""
ENV INGEST_KEY=""
ENTRYPOINT ["/resmo-db-mapper/resmo-db-mapper"]