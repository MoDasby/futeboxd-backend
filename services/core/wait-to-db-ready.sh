# Carrega as variáveis do arquivo .env
export $(grep -v '^#' .env | xargs)

# Função para verificar se o banco de dados está pronto
check_db() {
  PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "\q" >/dev/null 2>&1
}

# Número de tentativas
max_attempts=5
attempt=1

# Loop de tentativas
while [ $attempt -le $max_attempts ]; do

  if check_db; then
    echo "Conexão com o banco de dados estabelecida com sucesso!"
    exit 0
  else
    echo "Banco de dados não está pronto, tentando novamente em 2 segundos..."
    attempt=$((attempt + 1))
    sleep 2
  fi
done

echo "Falha ao conectar ao banco de dados."
exit 1
