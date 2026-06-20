package main

import "time"

const (
	GearTypeSki    = "雪板"
	GearTypeBoot   = "雪鞋"
	GearTypeHelmet = "头盔"
	GearTypePad    = "护具"
	GearTypePole   = "雪仗"
)

var ValidGearTypes = map[string]bool{
	GearTypeSki:    true,
	GearTypeBoot:   true,
	GearTypeHelmet: true,
	GearTypePad:    true,
	GearTypePole:   true,
}

const (
	ConditionGood    = "完好"
	ConditionMinor   = "轻微磨损"
	ConditionDamaged = "有损坏"
)

var ValidConditions = map[string]bool{
	ConditionGood:    true,
	ConditionMinor:   true,
	ConditionDamaged: true,
}

const (
	MemberTypeDay    = "日卡"
	MemberTypeSeason = "季卡"
	MemberTypeYear   = "年卡"
)

var ValidMemberTypes = map[string]bool{
	MemberTypeDay:    true,
	MemberTypeSeason: true,
	MemberTypeYear:   true,
}

const (
	GearStatusAvailable = "在库"
	GearStatusRented    = "借出中"
)

type Gear struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Size     string `json:"size"`
	Deposit  int    `json:"deposit"`
	Status   string `json:"status"`
}

type Member struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Type      string `json:"type"`
	ExpireStr string `json:"expire"`
}

type RentRecord struct {
	GearID       string `json:"gear_id"`
	MemberID     string `json:"member_id"`
	RentDateStr  string `json:"rent_date"`
	Returned     bool   `json:"returned"`
	ReturnDateStr string `json:"return_date,omitempty"`
	Condition    string `json:"condition,omitempty"`
}

type DataStore struct {
	Gears   []Gear       `json:"gears"`
	Members []Member     `json:"members"`
	Rents   []RentRecord `json:"rents"`
}

func (m Member) ParseExpire() (time.Time, error) {
	return time.Parse("2006-01-02", m.ExpireStr)
}

func (r RentRecord) ParseRentDate() (time.Time, error) {
	return time.Parse("2006-01-02", r.RentDateStr)
}

func (r RentRecord) ParseReturnDate() (time.Time, error) {
	if r.ReturnDateStr == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", r.ReturnDateStr)
}
