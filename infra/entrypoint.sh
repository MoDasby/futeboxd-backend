echo "Aguardando banco de dados ficar pronto..."
./wait-to-db-ready.sh
if [ $? -ne 0 ]; then
    echo "Falha ao conectar ao banco de dados, encerrando."
    exit 1
fi

echo "Banco de dados pronto, iniciando a API..."
./api
