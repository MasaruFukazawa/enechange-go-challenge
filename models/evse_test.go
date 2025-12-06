package models

import (
	"testing"
)

func TestEVSEStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   EVSEStatus
		expected string
	}{
		{"Available", StatusAvailable, "AVAILABLE"},
		{"Blocked", StatusBlocked, "BLOCKED"},
		{"Charging", StatusCharging, "CHARGING"},
		{"Inoperative", StatusInoperative, "INOPERATIVE"},
		{"OutOfOrder", StatusOutOfOrder, "OUTOFORDER"},
		{"Planned", StatusPlanned, "PLANNED"},
		{"Removed", StatusRemoved, "REMOVED"},
		{"Reserved", StatusReserved, "RESERVED"},
		{"Unknown", StatusUnknown, "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("EVSEStatus.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEVSE_ToResponse(t *testing.T) {
	evse := EVSE{
		UID:        "evse-001",
		LocationID: "loc-001",
		Status:     int(StatusAvailable),
	}

	response := evse.ToResponse()

	if response.UID != "evse-001" {
		t.Errorf("ToResponse().UID = %v, want %v", response.UID, "evse-001")
	}
	if response.Status != "AVAILABLE" {
		t.Errorf("ToResponse().Status = %v, want %v", response.Status, "AVAILABLE")
	}
}

func TestEVSE_TableName(t *testing.T) {
	evse := EVSE{}
	if got := evse.TableName(); got != "evses" {
		t.Errorf("TableName() = %v, want %v", got, "evses")
	}
}
