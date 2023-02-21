FROM debian:stretch-slim

COPY --from=horizonbuild /usr/local/bin/horizon /usr/local/bin/horizon
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

EXPOSE 8000

ENTRYPOINT ["entrypoint.sh"]
