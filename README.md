# Spike: Consultas e updates no Elasticsearch com Go

## Subindo o Elasticsearch localmente

```sh
docker run -d --name elasticsearch -p 9200:9200 \
  -e "discovery.type=single-node" \
  -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
  -e "xpack.security.enabled=false" \
  -e "xpack.security.http.ssl.enabled=false" \
  docker.elastic.co/elasticsearch/elasticsearch:8.13.4
```

Acesse http://localhost:9200 para validar.

## Criando índice e inserindo documento manualmente

```sh
curl -X PUT "http://localhost:9200/pois" -H 'Content-Type: application/json' -d '{
  "mappings": {
    "properties": {
      "name":    { "type": "text" },
      "type":    { "type": "keyword" },
      "lat":     { "type": "float" },
      "lon":     { "type": "float" },
      "address": { "type": "text" }
    }
  }
}'

curl -X POST "http://localhost:9200/pois/_doc/" -H 'Content-Type: application/json' -d '{
  "name": "Museu do Amanhã",
  "type": "museum",
  "lat": -22.8968,
  "lon": -43.1807,
  "address": "Praça Mauá, Rio de Janeiro"
}'
```

## Estrutura do Projeto (Hexagonal)

```
consultas-updates-no-elasticsearch-com-go/
├── cmd/app/main.go           # Exemplo de uso
├── internal/core/poi.go      # Entidade e porta
├── internal/adapter/repository/elasticsearch_repository.go # Adapter
├── internal/service/poi_service.go     # Casos de uso
├── go.mod
└── README.md
```

## Rodando o exemplo Go

```sh
cd consultas-updates-no-elasticsearch-com-go
# Instale as dependências
go mod tidy
# Execute o main
cd cmd/app
go run main.go
```

## Exemplo de uso com curl

Buscar por campo:
```sh
curl -X GET "http://localhost:9200/pois/_search" -H 'Content-Type: application/json' -d '{
  "query": { "match": { "type": "museum" } }
}'
```

Buscar por ID:
```sh
curl -X GET "http://localhost:9200/pois/_doc/<ID>"
```

Atualizar documento:
```sh
curl -X POST "http://localhost:9200/pois/_update/<ID>" -H 'Content-Type: application/json' -d '{
  "doc": { "address": "Novo endereço" }
}'
```

---

- O código trata erros de índice, documento não encontrado e conexão.
- Use a struct POI para inserir/atualizar dados.
