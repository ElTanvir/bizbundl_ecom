package elastic

import (
	"bizbundl/internal/config"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
)

func NewElasticClient(cfg *config.Config) (*elasticsearch.Client, error) {
	if cfg.ElasticURL == "" {
		log.Warn().Msg("Elasticsearch URL not provided. Search will degrade to Database fallback.")
		return nil, nil
	}

	conf := elasticsearch.Config{
		Addresses: []string{cfg.ElasticURL},
	}

	if cfg.ElasticUsername != "" && cfg.ElasticPassword != "" {
		conf.Username = cfg.ElasticUsername
		conf.Password = cfg.ElasticPassword
	}

	client, err := elasticsearch.NewClient(conf)
	if err != nil {
		return nil, fmt.Errorf("error creating elastic client: %w", err)
	}

	// Verify connection
	res, err := client.Info()
	if err != nil {
		// If we can't connect, should we fail hard?
		// "In some Cases we Would utilize elastic Search... If not Provided the We fall back"
		// If URL IS provided but fails, maybe warn?
		// For now, let's error out if explicit config fails.
		return nil, fmt.Errorf("error connecting to elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error connecting to elasticsearch: %s", res.String())
	}

	log.Info().Str("url", cfg.ElasticURL).Msg("Connected to Elasticsearch")
	return client, nil
}
