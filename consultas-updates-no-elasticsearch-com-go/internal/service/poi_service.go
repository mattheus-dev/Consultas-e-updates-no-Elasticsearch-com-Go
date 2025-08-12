package service

import (
	"consultas-updates-no-elasticsearch-com-go/internal/core"
)

type POIService struct {
	repo core.POIRepository
}

func NewPOIService(repo core.POIRepository) *POIService {
	return &POIService{repo: repo}
}

func (s *POIService) GetByID(id string) (*core.POI, error) {
	return s.repo.GetByID(id)
}

func (s *POIService) SearchByField(field, value string) ([]core.POI, error) {
	return s.repo.SearchByField(field, value)
}

func (s *POIService) Insert(poi *core.POI) (string, error) {
	return s.repo.Insert(poi)
}

func (s *POIService) Update(id string, poi *core.POI) error {
	return s.repo.Update(id, poi)
}
