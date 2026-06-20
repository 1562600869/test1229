package main

import (
	"fmt"
	"time"
)

func CmdAddGear(id, name, gearType, size string, deposit int) error {
	gear := Gear{
		ID:      id,
		Name:    name,
		Type:    gearType,
		Size:    size,
		Deposit: deposit,
		Status:  GearStatusAvailable,
	}
	if err := gear.Validate(); err != nil {
		return err
	}
	if deposit < 0 {
		return fmt.Errorf("押金不能为负数")
	}
	store, err := LoadData()
	if err != nil {
		return err
	}
	if store.FindGear(id) != nil {
		return fmt.Errorf("雪具ID %s 已存在", id)
	}
	store.Gears = append(store.Gears, gear)
	if err := SaveData(store); err != nil {
		return err
	}
	fmt.Printf("添加雪具成功: %s %s\n", gear.ID, gear.Name)
	return nil
}

func CmdAddMember(id, name, phone, memberType, expire string) error {
	if !ValidMemberTypes[memberType] {
		return fmt.Errorf("会员类型非法，只能是：日卡/季卡/年卡")
	}
	if _, err := parseDateUTC(expire); err != nil {
		return fmt.Errorf("有效期格式错误，应为 2006-01-02: %w", err)
	}
	store, err := LoadData()
	if err != nil {
		return err
	}
	if store.FindMember(id) != nil {
		return fmt.Errorf("会员ID %s 已存在", id)
	}
	member := Member{
		ID:        id,
		Name:      name,
		Phone:     phone,
		Type:      memberType,
		ExpireStr: expire,
	}
	store.Members = append(store.Members, member)
	if err := SaveData(store); err != nil {
		return err
	}
	fmt.Printf("添加会员成功: %s %s\n", member.ID, member.Name)
	return nil
}

func CmdRent(gearID, memberID, date string) error {
	if _, err := parseDateUTC(date); err != nil {
		return fmt.Errorf("日期格式错误，应为 2006-01-02: %w", err)
	}
	store, err := LoadData()
	if err != nil {
		return err
	}
	gear := store.FindGear(gearID)
	if gear == nil {
		return fmt.Errorf("雪具 %s 不存在", gearID)
	}
	if !gear.IsAvailable() {
		return fmt.Errorf("雪具 %s 当前状态为 %s，无法借出", gearID, gear.Status)
	}
	member := store.FindMember(memberID)
	if member == nil {
		return fmt.Errorf("会员 %s 不存在", memberID)
	}
	expire, err := member.ParseExpire()
	if err != nil {
		return fmt.Errorf("会员有效期解析失败: %w", err)
	}
	rentDate, _ := parseDateUTC(date)
	if rentDate.After(expire) {
		return fmt.Errorf("会员 %s 的有效期至 %s，已过期", memberID, member.ExpireStr)
	}
	gear.Status = GearStatusRented
	record := RentRecord{
		GearID:      gearID,
		MemberID:    memberID,
		RentDateStr: date,
		Returned:    false,
	}
	store.Rents = append(store.Rents, record)
	if err := SaveData(store); err != nil {
		return err
	}
	fmt.Printf("借出成功: %s 借给 %s，日期 %s\n", gear.Name, member.Name, date)
	return nil
}

func CmdReturn(gearID, condition, date string) error {
	record := RentRecord{Condition: condition}
	if err := record.Validate(); err != nil {
		return err
	}
	if _, err := parseDateUTC(date); err != nil {
		return fmt.Errorf("日期格式错误，应为 2006-01-02: %w", err)
	}
	store, err := LoadData()
	if err != nil {
		return err
	}
	gear := store.FindGear(gearID)
	if gear == nil {
		return fmt.Errorf("雪具 %s 不存在", gearID)
	}
	if gear.Status != GearStatusRented {
		return fmt.Errorf("雪具 %s 当前状态为 %s，无需归还", gearID, gear.Status)
	}
	rent := store.FindActiveRent(gearID)
	if rent == nil {
		return fmt.Errorf("未找到雪具 %s 的借出记录", gearID)
	}
	member := store.FindMember(rent.MemberID)
	memberName := rent.MemberID
	if member != nil {
		memberName = member.Name
	}
	gear.Status = GearStatusAvailable
	rent.Returned = true
	rent.ReturnDateStr = date
	rent.Condition = condition
	if err := SaveData(store); err != nil {
		return err
	}
	fmt.Printf("归还成功: %s 由 %s 归还，状况 %s\n", gear.Name, memberName, condition)
	return nil
}

func CmdDaily(date string) error {
	if _, err := parseDateUTC(date); err != nil {
		return fmt.Errorf("日期格式错误，应为 2006-01-02: %w", err)
	}
	store, err := LoadData()
	if err != nil {
		return err
	}
	rentCount := 0
	for _, r := range store.Rents {
		if r.RentDateStr == date {
			rentCount++
		}
	}
	availableCount := 0
	for _, g := range store.Gears {
		if g.Status == GearStatusAvailable {
			availableCount++
		}
	}
	fmt.Printf("日期: %s\n", date)
	fmt.Printf("  借出次数: %d\n", rentCount)
	fmt.Printf("  在库雪具数量: %d\n", availableCount)
	return nil
}

func CmdOverdue() error {
	store, err := LoadData()
	if err != nil {
		return err
	}
	today := time.Now().UTC()
	found := false
	for _, r := range store.Rents {
		if r.Returned {
			continue
		}
		rentDate, err := r.ParseRentDate()
		if err != nil {
			continue
		}
		days := int(today.Sub(rentDate).Hours() / 24)
		if days > 3 {
			gear := store.FindGear(r.GearID)
			member := store.FindMember(r.MemberID)
			gearName := r.GearID
			if gear != nil {
				gearName = gear.Name
			}
			memberName := r.MemberID
			if member != nil {
				memberName = member.Name
			}
			if !found {
				fmt.Println("逾期未还雪具（借出超过3天）:")
				found = true
			}
			fmt.Printf("  雪具: %s, 借用人: %s, 借出日期: %s, 已逾期 %d 天\n",
				gearName, memberName, r.RentDateStr, days-3)
		}
	}
	if !found {
		fmt.Println("暂无逾期未还的雪具")
	}
	return nil
}
