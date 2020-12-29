# Test - read_csv

Esse projeto compreende um serviço Go que recebe um arquivo, lê e em seguida armazena em banco de dados e o sanitiza.

## Pré requisitos:

1) Instalar o docker, docker-compose e go versão 1.15;

## Execução passo a passo:

1) No terminal, navegar até o diretório raiz do projeto;
2) Executar o comando `docker-compose up -d db` - comando para início do postgres;
3) Valide se os containers estão em execução com o comando `docker ps -a`;
4) Está pronto o ambiente, hora de testar, no terminal navegar até a pasta raiz do projeto e executar o comando `go run main.go "/dir/path/to/the/file"`;

## Informação adicional:

O arquivo modelo para execução está localizado no diretório <read_csv/external/base_teste.txt>
Você pode copiar para um diretório local antes da execução ou mover para uma pasta de sua ;preferência; 

Qualquer dúvida, estou disponível no email: lucas@itguru.com.br
Aceito PRs e sugestões de melhorias.


