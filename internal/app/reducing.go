package app

import (
	"fmt"
	"github.com/Aerok925/shortrurl/internal/entities"
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

func (s *Service) GetURL(id string) (entities.ShortUrl, error) {
	value, err := s.imMemory.GetValue(id)
	if err != nil {
		return entities.ShortUrl{}, err
	}
	res := entities.ShortUrl{
		ID:     id,
		URL:    value,
		Create: false,
	}
	return res, nil
}

func (s *Service) createURL(id string) string {
	return fmt.Sprintf("%s/%s", s.hostName, id)
}

func (s *Service) CreateOrUpdateNewURL(shortURL entities.UnprocessedURL) (entities.ShortUrl, error) {
	key := s.r.TruncateLine(shortURL.URL)
	create, err := s.imMemory.CreateOrUpdate(key, shortURL.URL)
	if err != nil {
		return entities.ShortUrl{}, err
	}
	res := entities.ShortUrl{
		ID:     key,
		Create: create,
		URL:    s.createURL(key),
	}
	return res, nil
}
