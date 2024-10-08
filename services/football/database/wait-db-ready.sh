#!/bin/bash

ENV_FILE="../.env"

# Verifica se o arquivo .env existe
if [ -f "$ENV_FILE" ]; then
    echo "Carregando variáveis de ambiente de $ENV_FILE"
    source "$ENV_FILE"
else
    echo "Arquivo $ENV_FILE não encontrado!"
    exit 1
fi

# Variáveis de ambiente
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_NAME=${DB_NAME:-postgres}
DB_PASSWORD=${DB_PASSWORD}

if [[ -z "$DB_PASSWORD" ]]; then
  echo "A variável de ambiente DB_PASSWORD não está definida."
  exit 1
fi

MAX_ATTEMPTS=5
SLEEP_INTERVAL=5

check_postgres() {
  PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "\q" > /dev/null 2>&1
}

attempt=1
while [[ $attempt -le $MAX_ATTEMPTS ]]
do
  echo "Tentativa $attempt de $MAX_ATTEMPTS..."
  
  check_postgres

  if [[ $? -eq 0 ]]; then
    echo "Banco de dados PostgreSQL está funcionando!"
    exit 0
  else
    echo "Falha ao conectar ao banco de dados. Tentando novamente em $SLEEP_INTERVAL segundos..."
    sleep $SLEEP_INTERVAL
  fi

  attempt=$((attempt + 1))
done

echo "Não foi possível conectar ao banco de dados PostgreSQL após $MAX_ATTEMPTS tentativas."
exit 1
