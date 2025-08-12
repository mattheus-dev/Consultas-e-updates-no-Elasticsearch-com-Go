package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"consultas-updates-no-elasticsearch-com-go/internal/core"
)

type ElasticsearchRepository struct {
	es    *elasticsearch.Client
	index string
}

func NewElasticsearchRepository(es *elasticsearch.Client, index string) *ElasticsearchRepository {
	return &ElasticsearchRepository{es: es, index: index}
}

func (r *ElasticsearchRepository) GetByID(id string) (*core.POI, error) {
	res, err := r.es.Get(r.index, id)
	if err != nil {
		return nil, fmt.Errorf("erro de conexão: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		return nil, fmt.Errorf("documento não encontrado")
	}
	if res.IsError() {
		return nil, fmt.Errorf("erro ao buscar documento: %s", res.String())
	}
	var resp struct {
		Source core.POI `json:"_source"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}
	return &resp.Source, nil
}

func (r *ElasticsearchRepository) SearchByField(field, value string) ([]core.POI, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{field: value},
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := r.es.Search(
		r.es.Search.WithContext(context.Background()),
		r.es.Search.WithIndex(r.index),
		r.es.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("erro de conexão: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("erro na busca: %s", res.String())
	}
	var resp struct {
		Hits struct {
			Hits []struct {
				Source core.POI `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}
	pois := make([]core.POI, 0, len(resp.Hits.Hits))
	for _, h := range resp.Hits.Hits {
		pois = append(pois, h.Source)
	}
	return pois, nil
}

func (r *ElasticsearchRepository) Insert(poi *core.POI) (string, error) {
	body, err := json.Marshal(poi)
	if err != nil {
		return "", err
	}
	res, err := r.es.Index(r.index, strings.NewReader(string(body)))
	if err != nil {
		return "", fmt.Errorf("erro de conexão: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return "", fmt.Errorf("erro ao inserir: %s", res.String())
	}
	var resp struct {
		ID string `json:"_id"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (r *ElasticsearchRepository) Update(id string, poi *core.POI) error {
	body, err := json.Marshal(map[string]interface{}{"doc": poi})
	if err != nil {
		return err
	}
	res, err := r.es.Update(r.index, id, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("erro de conexão: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		return fmt.Errorf("documento não encontrado")
	}
	if res.IsError() {
		return fmt.Errorf("erro ao atualizar: %s", res.String())
	}
	return nil
}
