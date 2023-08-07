#!/bin/bash

# Pacotes para testar
declare -a test_dirs=("modules/book") # ("modules/book" "modules/other" "modules/another")

# Loop através dos diretórios e executa os testes
for dir in "${test_dirs[@]}"; do
    echo "Running tests in $dir..."
    ENVIRONMENT=test go test ./"$dir"/ -v -cover
done