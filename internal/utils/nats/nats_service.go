package nats

import (
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"encoding/json"
	"github.com/nats-io/nats.go"
)

type NatsService interface {
	RequestImageDeletion(clientID string, images []string) error
	RequestImageUsage(images []assets.ImageDeleteRequest) error
}

type natsService struct {
	nats string
}

func NewNatsService(nats string) NatsService {
	return &natsService{
		nats: nats,
	}
}

func (cs natsService) RequestImageDeletion(clientID string, images []string) error {
	nc, err := nats.Connect(cs.nats)
	if err != nil {
		return err
	}
	defer nc.Close()

	data, _ := json.Marshal(assets.ImageDeleteRequest{ClientID: clientID, Images: images})

	return nc.Publish(utils.NatsAssetImageDelete, data)
}

func (cs natsService) RequestImageUsage(images []assets.ImageDeleteRequest) error {
	nc, err := nats.Connect(cs.nats)
	if err != nil {
		return err
	}
	defer nc.Close()

	data, _ := json.Marshal(images)

	return nc.Publish(utils.NatsAssetImageUsage, data)
}
