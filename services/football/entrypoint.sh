echo "Aguardando banco de dados ficar pronto"
./wait-to-db-ready.sh
echo "Banco de dados pronto, iniciando a API..."
./api
