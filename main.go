package main

import (
	"fmt"
	"os"
	"strconv"
)

type Args map[string]string

func parseArgs(args []string) ([]string, Args) {
	positional := []string{}
	flags := Args{}
	i := 0
	for i < len(args) {
		a := args[i]
		if len(a) > 2 && a[:2] == "--" {
			key := a[2:]
			if i+1 < len(args) {
				flags[key] = args[i+1]
				i += 2
				continue
			}
		}
		positional = append(positional, a)
		i++
	}
	return positional, flags
}

func printUsage() {
	fmt.Println("用法:")
	fmt.Println("  add-gear <ID> <名称> --type 雪板|雪鞋|头盔|护具|雪仗 --size 尺寸 --deposit 押金(整数分)")
	fmt.Println("  add-member <ID> --name 姓名 --phone 电话 --type 日卡|季卡|年卡 --expire 有效期(YYYY-MM-DD)")
	fmt.Println("  rent <雪具ID> --member 会员ID --date 日期(YYYY-MM-DD)")
	fmt.Println("  return <雪具ID> --condition 完好|轻微磨损|有损坏 --date 日期(YYYY-MM-DD)")
	fmt.Println("  daily --date 日期(YYYY-MM-DD)")
	fmt.Println("  overdue")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	positional, flags := parseArgs(os.Args[2:])

	var err error
	switch cmd {
	case "add-gear":
		if len(positional) < 2 {
			err = fmt.Errorf("参数不足: 需要 ID 和 名称")
			break
		}
		id := positional[0]
		name := positional[1]
		gearType, ok := flags["type"]
		if !ok {
			err = fmt.Errorf("缺少 --type 参数")
			break
		}
		size, ok := flags["size"]
		if !ok {
			err = fmt.Errorf("缺少 --size 参数")
			break
		}
		depositStr, ok := flags["deposit"]
		if !ok {
			err = fmt.Errorf("缺少 --deposit 参数")
			break
		}
		deposit, cerr := strconv.Atoi(depositStr)
		if cerr != nil {
			err = fmt.Errorf("押金必须是整数: %w", cerr)
			break
		}
		err = CmdAddGear(id, name, gearType, size, deposit)

	case "add-member":
		if len(positional) < 1 {
			err = fmt.Errorf("参数不足: 需要 ID")
			break
		}
		id := positional[0]
		name, ok := flags["name"]
		if !ok {
			err = fmt.Errorf("缺少 --name 参数")
			break
		}
		phone, ok := flags["phone"]
		if !ok {
			err = fmt.Errorf("缺少 --phone 参数")
			break
		}
		memberType, ok := flags["type"]
		if !ok {
			err = fmt.Errorf("缺少 --type 参数")
			break
		}
		expire, ok := flags["expire"]
		if !ok {
			err = fmt.Errorf("缺少 --expire 参数")
			break
		}
		err = CmdAddMember(id, name, phone, memberType, expire)

	case "rent":
		if len(positional) < 1 {
			err = fmt.Errorf("参数不足: 需要 雪具ID")
			break
		}
		gearID := positional[0]
		memberID, ok := flags["member"]
		if !ok {
			err = fmt.Errorf("缺少 --member 参数")
			break
		}
		date, ok := flags["date"]
		if !ok {
			err = fmt.Errorf("缺少 --date 参数")
			break
		}
		err = CmdRent(gearID, memberID, date)

	case "return":
		if len(positional) < 1 {
			err = fmt.Errorf("参数不足: 需要 雪具ID")
			break
		}
		gearID := positional[0]
		condition, ok := flags["condition"]
		if !ok {
			err = fmt.Errorf("缺少 --condition 参数")
			break
		}
		date, ok := flags["date"]
		if !ok {
			err = fmt.Errorf("缺少 --date 参数")
			break
		}
		err = CmdReturn(gearID, condition, date)

	case "daily":
		date, ok := flags["date"]
		if !ok {
			err = fmt.Errorf("缺少 --date 参数")
			break
		}
		err = CmdDaily(date)

	case "overdue":
		err = CmdOverdue()

	case "-h", "--help", "help":
		printUsage()
		return

	default:
		err = fmt.Errorf("未知命令: %s", cmd)
		printUsage()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
