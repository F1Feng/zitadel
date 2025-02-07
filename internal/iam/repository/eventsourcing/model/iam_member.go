package model

import (
	"encoding/json"

	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

type IAMMember struct {
	es_models.ObjectRoot
	UserID string   `json:"userId,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

func (m *IAMMember) SetData(event *es_models.Event) error {
	m.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, m); err != nil {
		logging.Log("EVEN-e4dkp").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
