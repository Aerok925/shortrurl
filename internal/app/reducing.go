package app

import (
	"fmt"
	"go.uber.org/zap"
)

type cache interface {
	GetValue(id string) (string, error)
	CreateOrUpdate(key, value string) (bool, error)
}

type reducing interface {
	TruncateLine(str string) string
}

type Service struct {
	imMemory cache
	r        reducing
	logger   *zap.Logger
	hostName string
}

func New(cache cache, r reducing, logger *zap.Logger, hostname string) *Service {
	return &Service{
		imMemory: cache,
		r:        r,
		logger:   logger,
		hostName: hostname,
	}
}

func (s *Service) GetUrl(id string) (string, error) {
	value, err := s.imMemory.GetValue(id)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (s *Service) createURL(id string) string {
	return fmt.Sprintf("http://%s/%s", s.hostName, id)
}

func (s *Service) CreateOrUpdateNewUrl(value string) (string, bool, error) {
	key := s.r.TruncateLine(value)
	create, err := s.imMemory.CreateOrUpdate(key, value)
	if err != nil {
		return "", false, err
	}
	if create {
		return s.createURL(key), true, nil
	}
	return s.createURL(key), false, nil
}
