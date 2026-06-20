package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const dataFile = "data.json"

func LoadData() (*DataStore, error) {
	store := &DataStore{
		Gears:   []Gear{},
		Members: []Member{},
		Rents:   []RentRecord{},
	}
	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, fmt.Errorf("读取数据文件失败: %w", err)
	}
	if len(data) == 0 {
		return store, nil
	}
	if err := json.Unmarshal(data, store); err != nil {
		return nil, fmt.Errorf("解析数据文件失败: %w", err)
	}
	return store, nil
}

func SaveData(store *DataStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化数据失败: %w", err)
	}
	if err := os.WriteFile(dataFile, data, 0644); err != nil {
		return fmt.Errorf("写入数据文件失败: %w", err)
	}
	return nil
}

func (s *DataStore) FindGear(id string) *Gear {
	for i := range s.Gears {
		if s.Gears[i].ID == id {
			return &s.Gears[i]
		}
	}
	return nil
}

func (s *DataStore) FindMember(id string) *Member {
	for i := range s.Members {
		if s.Members[i].ID == id {
			return &s.Members[i]
		}
	}
	return nil
}

func (s *DataStore) FindActiveRent(gearID string) *RentRecord {
	for i := range s.Rents {
		if s.Rents[i].GearID == gearID && !s.Rents[i].Returned {
			return &s.Rents[i]
		}
	}
	return nil
}
