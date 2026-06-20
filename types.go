package main

import (
	"fmt"
	"time"
)

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

var ValidGearStatuses = map[string]bool{
	GearStatusAvailable: true,
	GearStatusRented:    true,
}

type Gear struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Size    string `json:"size"`
	Deposit int    `json:"deposit"`
	Status  string `json:"status"`
}

func (g Gear) Validate() error {
	if !ValidGearTypes[g.Type] {
		return fmt.Errorf("雪具类型非法 %q，只能是：雪板/雪鞋/头盔/护具/雪仗", g.Type)
	}
	if !ValidGearStatuses[g.Status] {
		return fmt.Errorf("雪具状态非法 %q，只能是：在库/借出中", g.Status)
	}
	return nil
}

func (g Gear) IsAvailable() bool {
	return g.Status == GearStatusAvailable
}

type Member struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Type      string `json:"type"`
	ExpireStr string `json:"expire"`
}

type RentRecord struct {
	GearID        string `json:"gear_id"`
	MemberID      string `json:"member_id"`
	RentDateStr   string `json:"rent_date"`
	Returned      bool   `json:"returned"`
	ReturnDateStr string `json:"return_date,omitempty"`
	Condition     string `json:"condition,omitempty"`
}

func (r RentRecord) Validate() error {
	if r.Condition != "" && !ValidConditions[r.Condition] {
		return fmt.Errorf("归还状况非法 %q，只能是：完好/轻微磨损/有损坏", r.Condition)
	}
	return nil
}

type DataStore struct {
	Gears   []Gear       `json:"gears"`
	Members []Member     `json:"members"`
	Rents   []RentRecord `json:"rents"`
}

func parseDateUTC(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func (m Member) ParseExpire() (time.Time, error) {
	return parseDateUTC(m.ExpireStr)
}

func (r RentRecord) ParseRentDate() (time.Time, error) {
	return parseDateUTC(r.RentDateStr)
}

func (r RentRecord) ParseReturnDate() (time.Time, error) {
	if r.ReturnDateStr == "" {
		return time.Time{}, nil
	}
	return parseDateUTC(r.ReturnDateStr)
}
