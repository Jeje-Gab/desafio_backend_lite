**Desafio Backend**

Este projeto fornece uma API em Go defafio backend. A seguir, uma orientação passo a passo para que você possa rodar tudo na sua máquina.

---

## 📂 Estrutura do Projeto

```
desafio_backend/
├── bruno                  # Cliente CLI opcional para testes
├── cmd
│   ├── api                # Aplicação principal da API
│   └── importer           # Importador de CSVs para o banco
├── config
│   └── docker-compose.yml # Definições do Docker Compose
├── data
│   └── data-csv           # Arquivos CSV de origem
├── dist                   # Binários compilados (se houver)
├── dist-importer          # Binário do importador (se houver)
├── dump
│   └── trading_localhost-2025_06_09_03_31_43-dump.sql # Dump completo pronto para restauração
├── internal
│   └── migrations         # Scripts SQL de migração
│       ├── 001_create_negociacoes.sql
│       └── 002_aggregator_negociacoes.sql
├── pkg                    # Pacotes compartilhados
├── public                 # Recursos estáticos (se houver)
└── go.mod                 # Dependências do Go
```

---

## ⚙️  Pré-requisitos

* [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/).
* [Go](https://go.dev/doc/install) (versão 1.18 ou superior).
* Cliente de testes como [Postman](https://www.postman.com/), [Bruno](https://github.com/Jeje-Gab/bruno) ou `curl`.

---

## 🚀 Passo a Passo

1. **Subir o ambiente de banco de dados**

   ```bash
   cd config
   docker-compose up -d
   ```

   Isso iniciará um container PostgreSQL pronto para receber suas tabelas.

2. **Populando as tabelas**

   Você tem duas opções:

   ### Opção A: Executar migrações + importação manual

    1. Aplique o script de criação de estruturas:

       ```bash
       psql -h localhost -U postgres -d trading -f internal/migrations/001_create_negociacoes.sql
       ```

    2. Importe os dados dos CSVs:

       ```bash
       go run cmd/importer/main.go
       ```

    3. Execute o agregador de valores:

       ```bash
       psql -h localhost -U postgres -d trading -f internal/migrations/002_aggregator_negociacoes.sql
       ```

   ### Opção B: Restaurar dump completo

   Se preferir pular as etapas acima, restaure o dump:

   ```bash
   psql -h localhost -U postgres -d trading -f dump/trading_localhost-2025_06_09_03_31_43-dump.sql
   ```

3. **Iniciar a API**

   No diretório raiz do projeto:

   ```bash
   go run cmd/api/main.go
   ```

   A API estará disponível em `http://localhost:8080`.

4. **Testar o endpoint de estatísticas**

    * **Via navegador/Postman/Bruno**

      ```
      ```

GET /api/negociacoes/stats?ticker=WINM25\&from=2025-06-06

````

   - **Via curl**

     ```bash
     curl -G http://localhost:8080/api/negociacoes/stats \
       --data-urlencode "ticker=WINM25" \
       --data-urlencode "from=2025-06-06"
     ```
---

## 🔧 Ajustes Opcionais

- Se quiser alterar usuário, senha ou porta do PostgreSQL, edite o `config/docker-compose.yml` e o `.env`.
- Para compilar binários:

  ```bash
  go build -o dist/importer cmd/importer/main.go
  go build -o dist/api     cmd/api/main.go
````


OBS: apontar local de build para dist e deixar .env no mesmo local para pegar as variaveis de ambiente!
---