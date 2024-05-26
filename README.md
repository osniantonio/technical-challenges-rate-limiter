# Desafios técnicos: Rate Limiter

## Objetivo e descrição

Objetivo: Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

Descrição: O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:

Endereço IP: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
Token de Acesso: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:
API_KEY: <TOKEN>
As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.

## Requisitos

O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.
O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “dev.env” na pasta raiz.
Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.
O sistema deve responder adequadamente quando o limite é excedido:
Código HTTP: 429
Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame
Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.
Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.
A lógica do limiter deve estar separada do middleware.

# Exemplos

Limitação por IP: Suponha que o rate limiter esteja configurado para permitir no máximo 5 requisições por segundo por IP. Se o IP 192.168.1.1 enviar 6 requisições em um segundo, a sexta requisição deve ser bloqueada.
Limitação por Token: Se um token abc123 tiver um limite configurado de 10 requisições por segundo e enviar 11 requisições nesse intervalo, a décima primeira deve ser bloqueada.
Nos dois casos acima, as próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Ex: Se o tempo de expiração é de 5 minutos, determinado IP poderá realizar novas requisições somente após os 5 minutos.

# Dicas

Teste seu rate limiter sob diferentes condições de carga para garantir que ele funcione conforme esperado em situações de alto tráfego.

# Entrega

O código-fonte completo da implementação.
Documentação explicando como o rate limiter funciona e como ele pode ser configurado.
Testes automatizados demonstrando a eficácia e a robustez do rate limiter.
Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
O servidor web deve responder na porta 8080.

### Como subir o ambiente e o Rate Limiter em Golang

Para subir o ambiente via docker, execute:

```sh
docker-compose --env-file dev.env up -d --build
```

Ou, utilizando o make, execute:

```sh
make up
```

Para para o ambiente, execute:

```sh
make down
```

### Como testar

Primeiro execute os comandos descritos em:

1. Como subir o ambiente;
2. Como subir o Rate Limiter em Golang;
3. Por fim, execute os testes seguindo o descrito a seguir.

As requisições devem ser realizadas para a porta 8080 e o token de acesso vai no cabeçalho como uma API_KEY. Segue exemplo:

```sh
http GET http://localhost:8080 API_KEY:g8dXuf2MqNkqJ5tb47qw4m6thqYbrsK24SFZV4OiS83Lmbp8NCYulXtO3tyHJyZN
```

### Testar várias requisições simultâneas

Utilizando a ferramenta ApacheBench:

Instale o ApacheBench: sudo apt install apache2-utils (Linux);

Com o ApacheBench, execute:

```sh
ab -n 20 http://localhost:8080/
ab -n 101 -H "API_KEY: g8dXuf2MqNkqJ5tb47qw4m6thqYbrsK24SFZV4OiS83Lmbp8NCYulXtO3tyHJyZN" http://localhost:8080/
```

### Demais testes

```sh
go test -v ./...
```

Ou, ainda:

```sh
go test -v ./... -coverprofile=c.out
go tool cover -html=c.out
```
