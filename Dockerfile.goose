FROM gomicro/goose:3.26.0

COPY migrations/*.sql .

ENTRYPOINT ["goose", "postgres", "up"]