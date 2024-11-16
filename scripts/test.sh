test_failed() {
    echo "\n\nTestes Falharam"

    exit 1
}

echo "Testando football"

make -C football test || test_failed

echo "Testando core"

make -C core test || test_failed

echo "\n\nTodos os testes passaram com sucesso"