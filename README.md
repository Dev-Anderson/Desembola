# Desembola API

API construída em Golang (Gin) para upload e importação de planilhas CSV de tickets corporativos, e módulo nativo de gestão de usuários autenticados (com senhas criptografadas em Bcrypt via PostgreSQL).

## ⚙️ Como Executar

1. Configure o arquivo `.env` na raiz do projeto:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=seu_usuario
DB_PASSWORD=sua_senha
DB_NAME=seu_banco
DB_SSLMODE=disable
API_PORT=8080
```

2. Instale as dependências e inicie o servidor (As *Migrations* farão auto-healing e construirão qualquer tabela faltante na hora do Start!):
```bash
go mod tidy
go run cmd/main.go
```

---

## 🔗 Rotas Abertas (Usuários Comuns)

### `POST /api/users/register`
Registra um novo usuário na plataforma traduzindo imediatamente a senha bruta em Hash Irreversível.
**Ação via cURL:**
```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"nome": "Anderson Silva", "email": "anderson@desembola.com", "senha": "SenhaSecreta123"}'
```

### `POST /api/tickets/upload`
Faz o upload do arquivo CSV massivo de sistema e salva/atualiza todos os tickets no banco de dados, reconstruindo os acentos originais (de UTF-8 para ISO/Windows) perdidos por planilhas.
**Ação via cURL:**
```bash
curl -X POST http://localhost:8080/api/tickets/upload \
  -F "file=@./lista-tickets-BJ-csv.csv"
```

### `GET /api/tickets/all`
Retorna todos os tickets mesclados (Left Join) com as suas respectivas `Descrições de Status` do banco em JSON puro e limpo.
**Ação via cURL:**
```bash
curl http://localhost:8080/api/tickets/all
```

---

## 🔒 Rotas de Administração Restritas

### `GET /api/admin/users/all`
Lista todo o corpo de funcionários integrados omitindo intencionalmente os bytes do hash das senhas.
**Ação via cURL:**
```bash
curl http://localhost:8080/api/admin/users/all
```

### `POST /api/admin/tickets/clean`
Apaga de forma massiva a base de dados de tickets e status locais esvaziando as memórias para refazer testes (Apenas acessado sob parâmetro de liberação).
**Ação via cURL:**
```bash
curl -X POST http://localhost:8080/api/admin/tickets/clean \
  -H "Content-Type: application/json" \
  -d '{"limpar": true}'
```
