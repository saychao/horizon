package corer

import (
	"net/http"

	"gitlab.com/distributed_lab/corer"
	"gitlab.com/distributed_lab/logan"
)

// Connector - wrapper for corer.Connector
type Connector struct {
	corer.Connector
}

func NewConnector(httpClient *http.Client, coreURL string) (*Connector, error) {
	connector, err := corer.NewConnector(httpClient, coreURL)
	if err != nil {
		return nil, logan.Wrap(err, "Failed to create corer connector")
	}

	return &Connector{
		Connector: connector,
	}, nil
}

// GetCoreInfo returns response for core /info request
func (c *Connector) GetCoreInfo() (*Info, error) {
	var response infoResponse
	err := c.Connector.GetCoreInfo(&response)
	if err != nil {
		return nil, logan.Wrap(err, "Failed to get core info")
	}

	err = response.Info.validate()
	if err != nil {
		return nil, logan.Wrap(err, "Invalid info response from core")
	}

	return &response.Info, nil
}
