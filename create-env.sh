#!/usr/bin/env /bin/sh

SECRET=$(openssl rand -base64 32)

echo  SESSION_SECRET="${SECRET}" > .env
{
  echo  DB_USERNAME="root"
  echo  DB_PASSWORD="password"
  # DB_NAME is the name of the database
  echo  DB_NAME="integra_db"
  # Docker uses this host to connect it should be the same name as docker db service name
  echo  HOST="integra_db"

  echo  DB_PORT="5432"
} >> .env

echo "${SECRET}"