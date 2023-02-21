FROM tokend/postgres-ubuntu:9.6

COPY --from=horizonbuild /usr/local/bin/horizon /usr/local/bin/horizon

COPY ./entrypoint.sh /usr/local/bin/horizon-entrypoint.sh
COPY ./pg-entrypoint.sh /usr/local/bin/entrypoint.sh
RUN mv /usr/bin/entrypoint.sh /usr/local/bin/pg-entrypoint.sh

RUN true \
 && apt-get update \
 && apt-get install -y --no-install-recommends dumb-init ca-certificates \
 && chmod +x /usr/local/bin/horizon-entrypoint.sh \
 && chmod +x /usr/local/bin/entrypoint.sh \
 && chmod +x /usr/local/bin/pg-entrypoint.sh

EXPOSE 8000
EXPOSE 5432

ENTRYPOINT ["entrypoint.sh"]
