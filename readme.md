
GUIA PASSO A PASSO - MONEY FLOW BACKEND -
Como rodar o backend da aplicação
================================================================================

PRÉ-REQUISITOS
================================================================================

1. Go instalado (versão 1.21 ou superior)
   - Baixe em: https://go.dev/dl/
   - Configure o PATH do sistema para incluir o Go
   - Verifique a instalação executando(no terminal): go version

2. SQLite instalado (opcional, geralmente já vem com o sistema)
   - O projeto usa SQLite como banco de dados
   - O banco será criado automaticamente no diretório do projeto

3. Terminal/Command Prompt configurado
   - Windows: PowerShell ou CMD
   - Linux/Mac: Terminal padrão


PASSO 1: VERIFICAR A INSTALAÇÃO DO GO
================================================================================

1. Abra o terminal na pasta do projeto

2. Execute o comando para verificar a versão:
   
   go version

3. Verifique se o Go está configurado corretamente:
   
   go env GOPATH
   go env GOROOT


PASSO 2: INSTALAR DEPENDÊNCIAS DO PROJETO
================================================================================

1. No terminal, navegue até a pasta do projeto:
   
   cd MoneyFlow-Backend

2. Execute o comando para baixar as dependências:
   
   go mod download

3. Isso irá baixar todas as dependências listadas no go.mod:
   - gin-gonic/gin (framework web)
   - gorm.io/gorm (ORM para banco de dados)
   - gorm.io/driver/sqlite (driver SQLite)
   - dgrijalva/jwt-go (autenticação JWT)
   - golang.org/x/crypto/bcrypt (hash de senhas)
   - golang.org/x/net/html (parsing HTML)


PASSO 3: CONFIGURAR VARIÁVEIS DE AMBIENTE (OPCIONAL)
================================================================================

1. Crie um arquivo .env na raiz do projeto (opcional):
   
   PORT=8000
   JWT_SECRET=sua_chave_secreta_aqui

2. Ou defina a variável PORT diretamente no terminal antes de rodar:
   
   Windows (PowerShell):
   $env:PORT="8000"
   
   Windows (CMD):
   set PORT=8000
   
   Linux/Mac:
   export PORT=8000

3. Se não definir PORT, o servidor usará a porta padrão: 8000


PASSO 4: EXECUTAR O BACKEND
================================================================================

1. No terminal, certifique-se de estar na pasta do projeto:
   
   cd MoneyFlow-Backend

2. Execute o comando para rodar o servidor:
   
   go run .

3. Aguarde a mensagem indicando que o servidor está rodando:
   
   [GIN-debug] Listening and serving HTTP on :8000

4. O servidor estará disponível em:
   - http://localhost:8000
   - http://127.0.0.1:8000


PASSO 5: VERIFICAR SE O BACKEND ESTÁ FUNCIONANDO
================================================================================

1. Abra um navegador ou use ferramentas como Postman/Insomnia

2. Teste as rotas públicas:
   
   POST http://localhost:8000/signup
   {
     "email": "teste@example.com",
     "senha": "senha123"
   }
   
   POST http://localhost:8000/login
   {
     "email": "teste@example.com",
     "senha": "senha123"
   }

3. Se tudo estiver funcionando, você receberá respostas JSON


CONFIGURAÇÃO PARA ACESSO EXTERNO (DISPOSITIVO FÍSICO/FRONTEND)
================================================================================

IMPORTANTE: Para que o frontend Flutter ou outros clientes acessem o backend:

1. Para EMULADOR ANDROID:
   - O frontend deve usar: http://10.0.2.2:8000
   - Isso é um alias do Android Emulator para localhost da máquina
   - O backend deve estar rodando em localhost:8000

2. Para DISPOSITIVO FÍSICO:
   - Descubra o IP da sua máquina na rede local:
     
     Windows: ipconfig (procure por IPv4)
     Linux/Mac: ifconfig ou ip addr
   
   - Exemplo: 192.168.1.100
   - O frontend deve usar: http://192.168.1.100:8000
   - O backend deve estar acessível nesse IP

3. Para permitir acesso externo:
   - Verifique o firewall do sistema
   - Permita conexões na porta 8000
   - Windows: Firewall do Windows > Permitir um aplicativo
   - Linux: sudo ufw allow 8000


ESTRUTURA DO PROJETO
================================================================================

controllers/       - Handlers HTTP para cada recurso
  - usuario.go
  - transacao.go
  - categoria.go
  - estabelecimento.go
  - formaPagamento.go
  - dadosCompra.go
  - metaFinanceira.go

models/            - Estruturas de dados (structs)
  - usuario.go
  - transacao.go
  - categoria.go
  - estabelecimento.go
  - formaPagamento.go
  - metaFinanceira.go

repositories/      - Camada de acesso a dados
  - usuario.go
  - transacao.go
  - categoria.go
  - estabelecimento.go
  - formaPagamento.go
  - metaFinanceira.go
  - tokenBlacklist.go

services/          - Lógica de negócio
  - usuario.go
  - dadosCompra.go
  - cnae.go
  - cnaeMap.go
  - bootstrap.go

middlewares/       - Middlewares HTTP (autenticação)
  - auth.go

router/            - Configuração de rotas
  - router.go

database/          - Configuração do banco de dados
  - database.go

utils/             - Utilitários
  - context.go

main.go            - Arquivo principal, ponto de entrada
go.mod             - Gerenciamento de dependências
go.sum             - Checksums das dependências
MoneyFlow.db       - Banco de dados SQLite (criado automaticamente)


ROTAS DISPONÍVEIS
================================================================================

ROTAS PÚBLICAS (sem autenticação):
-----------------------------------
POST   /login                    - Login de usuário
POST   /signup                   - Registro de novo usuário


ROTAS PROTEGIDAS (requer token JWT):
-----------------------------------
POST   /usuarios/logout          - Logout do usuário
GET    /usuarios/usuario         - Obter ID do usuário logado

POST   /transacoes               - Criar transação
GET    /transacoes               - Listar todas as transações
GET    /transacoes/:id           - Obter transação por ID
GET    /transacoes/usuario       - Transações do usuário logado
GET    /transacoes/qtd/:qtd      - Últimas N transações
GET    /transacoes/periodo       - Transações por período
GET    /transacoes/tipo/:tipo    - Transações por tipo
GET    /transacoes/gastos-categorias/ultimo-mes - Gastos por categoria (últimos 30 dias)
PUT    /transacoes/:id           - Atualizar transação
DELETE /transacoes/:id           - Deletar transação

POST   /categorias               - Criar categoria
GET    /categorias               - Listar todas as categorias
GET    /categorias/:id           - Obter categoria por ID
PUT    /categorias/:id           - Atualizar categoria
DELETE /categorias/:id           - Deletar categoria

POST   /estabelecimentos         - Criar estabelecimento
GET    /estabelecimentos         - Listar todos os estabelecimentos
GET    /estabelecimentos/:id     - Obter estabelecimento por ID
PUT    /estabelecimentos/:id     - Atualizar estabelecimento
DELETE /estabelecimentos/:id     - Deletar estabelecimento

POST   /formas-pagamento         - Criar forma de pagamento
GET    /formas-pagamento         - Listar todas as formas de pagamento
GET    /formas-pagamento/:id     - Obter forma de pagamento por ID
GET    /formas-pagamento/qtd/:qtd - Primeiras N formas de pagamento
PUT    /formas-pagamento/:id     - Atualizar forma de pagamento
DELETE /formas-pagamento/:id     - Deletar forma de pagamento

POST   /metas-financeiras        - Criar meta financeira
GET    /metas-financeiras        - Listar todas as metas
GET    /metas-financeiras/usuario - Metas do usuário logado
GET    /metas-financeiras/ativa  - Meta financeira ativa
GET    /metas-financeiras/:id    - Obter meta por ID
GET    /metas-financeiras/:id/transacoes - Transações de uma meta
PUT    /metas-financeiras/:id    - Atualizar meta
DELETE /metas-financeiras/:id    - Deletar meta

POST   /scan                     - Extrair dados de compra via URL


AUTENTICAÇÃO
================================================================================

1. Para acessar rotas protegidas, você precisa:
   
   a) Fazer login primeiro:
      POST http://localhost:8000/login
      Body: {
        "email": "usuario@example.com",
        "senha": "senha123"
      }
   
   b) Usar o token retornado no header:
      Authorization: Bearer <token_aqui>

2. Exemplo de requisição autenticada:
   
   GET http://localhost:8000/transacoes
   Headers:
     Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...


PASSO 7: RESOLVER PROBLEMAS COMUNS
================================================================================

PROBLEMA: "cannot find package" ou erros de import
SOLUÇÃO:
- Execute: go mod download
- Execute: go mod tidy
- Verifique se está na pasta correta do projeto

PROBLEMA: "port already in use"
SOLUÇÃO:
- Altere a porta usando variável de ambiente: PORT=8001
- Ou feche o processo que está usando a porta 8000
- Windows: netstat -ano | findstr :8000
- Linux/Mac: lsof -i :8000

PROBLEMA: "database is locked"
SOLUÇÃO:
- Feche outras conexões com o banco de dados
- Reinicie o servidor
- Verifique se há processos usando o arquivo MoneyFlow.db

PROBLEMA: "no such table"
SOLUÇÃO:
- O banco será criado automaticamente na primeira execução
- Delete o arquivo MoneyFlow.db e execute novamente
- Verifique se database.SetupDatabase() está sendo chamado

PROBLEMA: Erro de compilação do Go
SOLUÇÃO:
- Verifique a versão do Go: go version (deve ser 1.21+)
- Atualize o Go se necessário
- Execute: go clean -cache
- Execute: go mod tidy

PROBLEMA: Conexão recusada do frontend
SOLUÇÃO:
- Verifique se o backend está rodando
- Confirme que está usando a URL correta:
  * Emulador: http://10.0.2.2:8000
  * Dispositivo físico: http://IP_DA_MAQUINA:8000
- Verifique o firewall


COMANDOS ÚTEIS
================================================================================

go run main.go                    - Executa o servidor
go build                          - Compila o projeto
go mod download                   - Baixa dependências
go mod tidy                       - Organiza e atualiza go.mod
go mod verify                     - Verifica dependências
go clean                          - Remove arquivos de compilação


FUNCIONALIDADES AUTOMÁTICAS
================================================================================

1. Auto-migração do banco de dados:
   - Na primeira execução, cria todas as tabelas automaticamente
   - Inclui: usuarios, transacoes, categorias, estabelecimentos, formas_pagamento, metas_financeiras

2. Seed inicial de categorias:
   - Popula categorias com base em CNAE na primeira execução
   - Inclui categorias como: Supermercados, Padarias, Farmácias, etc.
   - Categoria padrão "Outros" (ID: 9999999)

3. Backfill de categorias:
   - Na inicialização, atribui categorias automaticamente a estabelecimentos sem categoria
   - Consulta CNAE via BrasilAPI baseado no CNPJ
   - Estabelecimentos sem CNPJ recebem categoria "Outros"

4. Categorização automática:
   - Ao criar/atualizar estabelecimento, categoria é atribuída automaticamente via CNPJ/CNAE
   - Estabelecimentos sem CNPJ recebem categoria "Outros"


OBSERVAÇÕES IMPORTANTES
================================================================================

1. Banco de dados SQLite:
   - Arquivo: MoneyFlow.db na raiz do projeto
   - Criado automaticamente na primeira execução
   - Backup: copie o arquivo .db antes de atualizações importantes

2. Segurança:
   - Senhas são hasheadas com bcrypt antes de salvar
   - Tokens JWT expiram em 72 horas
   - Sistema de blacklist para tokens revogados no logout

3. API Externa:
   - O backend consulta BrasilAPI para obter CNAE de CNPJs
   - URL: https://brasilapi.com.br/api/cnpj/v1/{cnpj}
   - Requer conexão com internet para funcionar

4. CORS:
   - Se necessário, configure CORS no router.go para permitir requisições do frontend

5. Logs:
   - Logs são exibidos no console durante a execução
   - Inclui informações de requisições, erros e operações do banco


CONTATO E SUPORTE
================================================================================

Para problemas ou dúvidas sobre o projeto:
- Documentação Go: https://go.dev/doc/
- Documentação Gin: https://gin-gonic.com/docs/
- Documentação GORM: https://gorm.io/docs/

