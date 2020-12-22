# Test - read_csv

Esse projeto compreende um serviço Go que recebe um arquivo, lê e em seguida armazena em banco de dados e o sanitiza.

## Pré requisitos:

1) Instalar o docker, docker-compose e go versão 1.15

## Execução passo a passo:

1) No terminal, navegar até o diretório raiz do projeto;
2) Executar o comando `docker-compose up -d` - comando para início de dois containers: Adminer e Postgres;
3) Valide se os containers estão em execução com o comando `docker ps -a`;
4) Se tudo estiver certo, em seu navegador, navegue até o endereço: `http://localhost:8080/?pgsql=db`;
5) Será aberta a página inicial do Adminer onde você pode gerenciar facilmente nosso database;
6) Escolha o PostgreSQL e nos formulários usar: `server`:`db`, `username`:`luccasman`, `password`:`admin` (esses dados estão nos arquivos `.env` e `.yml`);
7) Após Login, em seu Menu Lateral Esquerdo, clicar `Import`;
8) Executar um Upload do arquivo que contém no diretório do projeto `migrations`/`sqls`/`arquivo.sql`;
9) Selecionar o arquivo `20201219004151-dataset-input-up.sql` e clicar `Execute`;
10) Está pronto o ambiente, hora de testar, no terminal navegar até a pasta raiz do projeto e executar o comando `go run main.go`


Qualquer dúvida, estou disponível no email: lucas@itguru.com.br
Aceito PRs e sugestões de melhorias.


