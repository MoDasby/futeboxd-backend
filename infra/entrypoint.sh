#!/bin/zsh

ENV_FILE=$1

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

# Verifica se a variável DB_PASSWORD está definida
if [[ -z "$DB_PASSWORD" ]]; then
  echo "A variável de ambiente DB_PASSWORD não está definida."
  exit 1
fi

MAX_ATTEMPTS=5
SLEEP_INTERVAL=2

# Função para verificar a conexão com o PostgreSQL
check_postgres() {
  PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "\q" > /dev/null 2>&1
}

# Tentativas de conexão
attempt=1
while [[ $attempt -le $MAX_ATTEMPTS ]]
do
  echo "Tentativa $attempt de $MAX_ATTEMPTS..."

  check_postgres

  if [[ $? -eq 0 ]]; then
    echo "Banco de dados PostgreSQL está funcionando!"
    break
  else
    echo "Falha ao conectar ao banco de dados. Tentando novamente em $SLEEP_INTERVAL segundos..."
    sleep $SLEEP_INTERVAL
  fi

  attempt=$((attempt + 1))
done

# Se não conseguir conectar após todas as tentativas, aborta
if [[ $attempt -gt $MAX_ATTEMPTS ]]; then
  echo "Não foi possível conectar ao banco de dados PostgreSQL após $MAX_ATTEMPTS tentativas."
  exit 1
fi

# Banco de dados está pronto, então executa a API
echo "Banco de dados pronto, iniciando a API..."
exec ./api
