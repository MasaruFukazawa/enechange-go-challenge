package models

// EVSEStatus はEVSEのステータスを表す列挙型
type EVSEStatus int

const (
	StatusAvailable   EVSEStatus = 1
	StatusBlocked     EVSEStatus = 2
	StatusCharging    EVSEStatus = 3
	StatusInoperative EVSEStatus = 4
	StatusOutOfOrder  EVSEStatus = 5
	StatusPlanned     EVSEStatus = 6
	StatusRemoved     EVSEStatus = 7
	StatusReserved    EVSEStatus = 8
	StatusUnknown     EVSEStatus = 9
)

var statusStrings = map[EVSEStatus]string{
	StatusAvailable:   "AVAILABLE",
	StatusBlocked:     "BLOCKED",
	StatusCharging:    "CHARGING",
	StatusInoperative: "INOPERATIVE",
	StatusOutOfOrder:  "OUTOFORDER",
	StatusPlanned:     "PLANNED",
	StatusRemoved:     "REMOVED",
	StatusReserved:    "RESERVED",
	StatusUnknown:     "UNKNOWN",
}

// String はEVSEStatusの文字列表現を返す
func (s EVSEStatus) String() string {
	if str, ok := statusStrings[s]; ok {
		return str
	}
	return "UNKNOWN"
}

// EVSE はEVSE（Electric Vehicle Supply Equipment）エンティティ
type EVSE struct {
	UID        string `gorm:"primaryKey;size:36" json:"uid"`
	LocationID string `gorm:"size:36;not null" json:"-"`
	Status     int    `gorm:"not null" json:"-"`
}

// TableName はGORMのテーブル名を返す
func (EVSE) TableName() string {
	return "evses"
}

// EVSEResponse はJSON出力用の構造体
type EVSEResponse struct {
	UID    string `json:"uid"`
	Status string `json:"status"`
}

// ToResponse はEVSEをJSON出力用構造体に変換する
func (e EVSE) ToResponse() EVSEResponse {
	return EVSEResponse{
		UID:    e.UID,
		Status: EVSEStatus(e.Status).String(),
	}
}
