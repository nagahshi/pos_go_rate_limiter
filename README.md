# Rate limiter

Gerencie requisições por segundo de sua aplicação por token e por endereço de IP

### Instalação
```sh
download 
ou
git clone git@github.com/nagahshi/pos_go_rate_limiter
cd pos_go_rate_limiter
```
### Configuração
Crie um arquivo `.env` com base no exemplo `.env-example` se for usar via docker configure o HOST conforme o nome do service em `Dockerfile`, por padrão, usamos `redis`

```sh
PORT=8080 #porta da aplicação 

REDIS_HOST=redis # endereço do serviço de redis 
REDIS_PORT=6379 # porta de acesso ao serviço redis
REDIS_PASSWORD= # senha de acesso ao serviço redis
REDIS_DATABASE_INDEX=0 # index do database redis

RATE_LIMITER_IP=0 # seta o limite de requisições por IP
RATE_LIMIT_TOKEN=0 # seta o limite de requisições por por token
RATE_LIMITER_TIMEOUT=0 # seta o tempo de timeout após o bloqueio
RATE_LIMITER_WINDOW_TIME=1 # tempo para o bloqueio padrão 1 sec
```

### Como usar
local
```sh
go run cmd/main.go
```

via docker a aplicação estará escutando na porta 8080
```sh
docker-compose up -d
```

Veja os exemplos de uso na pasta `api`, os bloqueios por token sempre será avaliado o header `API_KEY` da aplicação.

Implementação do middleware está em `cmd/main.go`
```sh
# cria-se uma istância do useCase limiter
# injeta no middleware
middleware := middlewareLimiter.NewMiddleware(limiter)
# func Run desse middleware é um HandlerFunc 
# que pode ser utilizado em servidores WEB
```