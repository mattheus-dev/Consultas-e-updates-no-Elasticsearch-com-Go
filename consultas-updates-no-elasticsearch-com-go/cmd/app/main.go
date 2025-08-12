package main

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"

	"consultas-updates-no-elasticsearch-com-go/internal/adapter/repository"
	"consultas-updates-no-elasticsearch-com-go/internal/core"
	"consultas-updates-no-elasticsearch-com-go/internal/service"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Erro ao criar cliente Elasticsearch: %v", err)
	}

	repo := repository.NewElasticsearchRepository(es, "pois")
	poiService := service.NewPOIService(repo)

	// Inserir novo POI
	poi := &core.POI{
		Name:    "Cristo Redentor",
		Type:    "monument",
		Lat:     -22.9519,
		Lon:     -43.2105,
		Address: "Parque Nacional da Tijuca, Rio de Janeiro",
	}
	id, err := poiService.Insert(poi)
	if err != nil {
		log.Printf("Erro ao inserir POI: %v", err)
	} else {
		fmt.Printf("POI inserido com id: %s\n", id)
	}

	// Buscar por campo
	results, err := poiService.SearchByField("type", "monument")
	if err != nil {
		log.Printf("Erro na busca: %v", err)
	} else {
		fmt.Printf("Resultados da busca por type=monument: %+v\n", results)
	}

	// Buscar por ID (usando o id retornado acima)
	poiFound, err := poiService.GetByID(id)
	if err != nil {
		log.Printf("Erro ao buscar por ID: %v", err)
	} else {
		fmt.Printf("POI encontrado por ID: %+v\n", poiFound)
	}

	// Atualizar POI
	poi.Address = "Novo endereço atualizado"
	err = poiService.Update(id, poi)
	if err != nil {
		log.Printf("Erro ao atualizar POI: %v", err)
	} else {
		fmt.Println("POI atualizado com sucesso")
	}
}
