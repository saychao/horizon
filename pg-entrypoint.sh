#!/usr/bin/dumb-init /bin/bash
set -eax

POSTGRES_PORT=${POSTGRES_PORT:-5432}
pg-entrypoint.sh "postgres -p $POSTGRES_PORT" &

while ! pg_isready -U$POSTGRES_USER -hlocalhost -p$POSTGRES_PORT > /dev/null 2> /dev/null; do
    echo "postgres not ready yet"
    sleep 3
done

horizon-entrypoint.sh $@
