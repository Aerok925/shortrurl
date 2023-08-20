package app

import (
	"github.com/Aerok925/shortrurl/internal/inmemory"
	reducing2 "github.com/Aerok925/shortrurl/internal/reducing"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func Test(t *testing.T) {
	// setup all struct
	hostName := "http://localhost:8080"
	wantHostName := "http://localhost:8080/"
	cache := inmemory.New()
	r := reducing2.New()
	logger, _ := zap.NewDevelopment()
	s := New(cache, r, logger, hostName)

	type get struct {
		id  string
		url string
		err error
	}

	type set struct {
		url    string
		create bool
		err    error
	}

	tests := []struct {
		name    string
		getWant get
		setWant set
		url     string
	}{
		{
			name: "default",
			getWant: get{
				id:  "d05ece8a",
				err: nil,
				url: "https://arenda.yandex.ru/lk/tenant/flats-search/a12d1f9460194783a57d60a2e874d875/",
			},
			setWant: set{
				create: true,
				err:    nil,
				url:    wantHostName + "d05ece8a",
			},
			url: "https://arenda.yandex.ru/lk/tenant/flats-search/a12d1f9460194783a57d60a2e874d875/",
		},
		{
			name: "default v2",
			getWant: get{
				id:  "4c60d71b",
				err: nil,
				url: "https://career.avito.com/vacancies/razrabotka/762/",
			},
			setWant: set{
				create: true,
				err:    nil,
				url:    wantHostName + "4c60d71b",
			},
			url: "https://career.avito.com/vacancies/razrabotka/762/",
		},
		{
			name: "dublicate default v2",
			getWant: get{
				id:  "4c60d71b",
				err: nil,
				url: "https://career.avito.com/vacancies/razrabotka/762/",
			},
			setWant: set{
				create: false,
				err:    nil,
				url:    wantHostName + "4c60d71b",
			},
			url: "https://career.avito.com/vacancies/razrabotka/762/",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url, create, err := s.CreateOrUpdateNewURL(test.url)
			assert.NoError(t, err)
			assert.Equal(t, test.setWant.create, create)
			assert.Equal(t, test.setWant.url, url)

			getURL, err := s.GetURL(test.getWant.id)
			assert.NoError(t, err)
			assert.Equal(t, test.getWant.url, getURL)

		})
	}
}
