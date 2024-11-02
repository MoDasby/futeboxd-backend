check_db() {
  PGPASSWORD=$POSTGRES_PASSWORD psql -h $DB_HOST -p $PGPORT -U $DB_USER -d $POSTGRES_DB -c "\q" >/dev/null 2>&1
}

echo $POSTGRES_PASSWORD $DB_HOST $PGPORT $DB_USER $POSTGRES_DB

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
