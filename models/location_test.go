package models

import (
	"testing"
	"time"
)

func TestLocation_ToResponse(t *testing.T) {
	name := "東京駅充電ステーション"
	now := time.Now()

	location := Location{
		ID:          "loc-001",
		Name:        &name,
		Address:     "東京都千代田区丸の内1-9-1",
		Latitude:    35.681236,
		Longitude:   139.767125,
		LastUpdated: &now,
		EVSEs: []EVSE{
			{UID: "evse-001", LocationID: "loc-001", Status: int(StatusAvailable)},
			{UID: "evse-002", LocationID: "loc-001", Status: int(StatusCharging)},
		},
	}

	response := location.ToResponse()

	if response.ID != "loc-001" {
		t.Errorf("ToResponse().ID = %v, want %v", response.ID, "loc-001")
	}
	if *response.Name != name {
		t.Errorf("ToResponse().Name = %v, want %v", *response.Name, name)
	}
	if response.Address != "東京都千代田区丸の内1-9-1" {
		t.Errorf("ToResponse().Address = %v, want %v", response.Address, "東京都千代田区丸の内1-9-1")
	}
	if response.Coordinates.Latitude != "35.681236" {
		t.Errorf("ToResponse().Coordinates.Latitude = %v, want %v", response.Coordinates.Latitude, "35.681236")
	}
	if response.Coordinates.Longitude != "139.767125" {
		t.Errorf("ToResponse().Coordinates.Longitude = %v, want %v", response.Coordinates.Longitude, "139.767125")
	}
	if len(response.EVSEs) != 2 {
		t.Errorf("len(ToResponse().EVSEs) = %v, want %v", len(response.EVSEs), 2)
	}
	if response.EVSEs[0].Status != "AVAILABLE" {
		t.Errorf("ToResponse().EVSEs[0].Status = %v, want %v", response.EVSEs[0].Status, "AVAILABLE")
	}
}

func TestLocation_ToResponse_NilName(t *testing.T) {
	location := Location{
		ID:        "loc-002",
		Name:      nil,
		Address:   "大阪市北区",
		Latitude:  34.6937,
		Longitude: 135.5023,
		EVSEs:     []EVSE{},
	}

	response := location.ToResponse()

	if response.Name != nil {
		t.Errorf("ToResponse().Name = %v, want nil", response.Name)
	}
}

func TestLocation_TableName(t *testing.T) {
	location := Location{}
	if got := location.TableName(); got != "locations" {
		t.Errorf("TableName() = %v, want %v", got, "locations")
	}
}
