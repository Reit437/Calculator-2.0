package Calc

import (
	"fmt"
	"strconv"
	"strings"
)

func Calc(expression string) map[string]string {
	id := 0
	exp := ""
	bid, eid := 0, 0
	mapid := make(map[string]string)
	for strings.Index(expression, "(") != -1 {
		bs := strings.Index(expression, "(") + 1
		es := strings.Index(expression, ")")
		exp = expression[bs:es]
		exp2 := "(" + exp + ")"
		exp = exp + " "
		for strings.Index(exp, "*") != -1 || strings.Index(exp, "/") != -1 {
			mult := strings.Index(exp, "*")
			div := strings.Index(exp, "/")
			fmt.Println(mult, div)
			if (mult < div && mult != -1) || div == -1 {
				for i := mult - 2; i >= 0; i-- {
					if string(exp[i]) == string(" ") {
						bid = i + 1
						break
					}
				}
				for i := mult + 2; i <= len(exp); i++ {
					if string(exp[i]) == string(" ") {
						eid = i + 1
						break
					}
				}
			} else if (mult > div && div != -1) || mult == -1 {
				for i := div - 2; i >= 0; i-- {
					if string(exp[i]) == string(" ") {
						bid = i + 1
						break
					}
				}
				for i := div + 2; i <= len(exp); i++ {
					if string(exp[i]) == string(" ") {
						eid = i + 1
						break
					}
				}
			}
			id++
			sid := "id" + strconv.Itoa(id)
			fmt.Println(bid, eid)
			mapid[sid] = exp[bid:eid]
			exp = strings.Replace(exp, exp[bid:eid], sid, 1)
			fmt.Println(exp, mapid)
		}
		for strings.Index(exp, "+") != -1 || strings.Index(exp, " - ") != -1 {
			add := strings.Index(exp, "+")
			sub := strings.Index(exp, " - ")
			fmt.Println(add, sub)
			if (add < sub && add != -1) || sub == -1 {
				for i := add - 2; i >= 0; i-- {
					if string(exp[i]) == string(" ") {
						bid = i + 1
						break
					}
				}
				for i := add + 2; i <= len(exp); i++ {
					if string(exp[i]) == string(" ") {
						eid = i + 1
						break
					}
				}
			} else if (add > sub && sub != -1) || add == -1 {
				for i := sub + 1 - 2; i >= 0; i-- {
					if string(exp[i]) == string(" ") {
						bid = i + 1
						break
					}
				}
				for i := sub + 1 + 2; i < len(exp); i++ {
					if string(exp[i]) == string(" ") {
						eid = i + 1
						break
					}
				}
			}
			id++
			sid := "id" + strconv.Itoa(id)
			fmt.Println(bid, eid)
			mapid[sid] = exp[bid:eid]
			exp = strings.Replace(exp, exp[bid:eid], sid, 1)
			fmt.Println(exp, mapid)
		}
		lk := "id" + strconv.Itoa(id)
		expression = strings.Replace(expression, exp2, lk, 1)
		fmt.Println(expression)
	}

	exp = " " + expression + " "
	for strings.Index(exp, "*") != -1 || strings.Index(exp, "/") != -1 {
		mult := strings.Index(exp, "*")
		div := strings.Index(exp, "/")
		fmt.Println(mult, div)
		if (mult < div && mult != -1) || div == -1 {
			for i := mult - 2; i >= 0; i-- {
				if string(exp[i]) == string(" ") {
					bid = i + 1
					break
				}
			}
			for i := mult + 2; i <= len(exp); i++ {
				if string(exp[i]) == string(" ") {
					eid = i + 1
					break
				}
			}
		} else if (mult > div && div != -1) || mult == -1 {
			for i := div - 2; i >= 0; i-- {
				if string(exp[i]) == string(" ") {
					bid = i + 1
					break
				}
			}
			for i := div + 2; i < len(exp); i++ {
				if string(exp[i]) == string(" ") {
					eid = i + 1
					break
				}
			}
		}
		id++
		sid := "id" + strconv.Itoa(id)
		fmt.Println(bid, eid)
		mapid[sid] = exp[bid:eid]
		exp = strings.Replace(exp, exp[bid:eid], sid+" ", 1)
		fmt.Println(exp, mapid)
	}
	for strings.Index(exp, "+") != -1 || strings.Index(exp, " - ") != -1 {
		add := strings.Index(exp, "+")
		sub := strings.Index(exp, " - ")
		fmt.Println(add, sub)
		if (add < sub && add != -1) || sub == -1 {
			for i := add - 2; i >= 0; i-- {
				if string(exp[i]) == string(" ") {
					bid = i + 1
					break
				}
			}
			for i := add + 2; i <= len(exp); i++ {
				if string(exp[i]) == string(" ") {
					eid = i + 1
					break
				}
			}
		} else if (add > sub && sub != -1) || add == -1 {
			for i := sub + 1 - 2; i >= 0; i-- {
				if string(exp[i]) == string(" ") {
					bid = i + 1
					break
				}
			}
			for i := sub + 1 + 2; i <= len(exp); i++ {
				if string(exp[i]) == string(" ") {
					eid = i + 1
					break
				}
			}
		}
		id++
		sid := "id" + strconv.Itoa(id)
		fmt.Println(bid, eid)
		mapid[sid] = exp[bid:eid]
		exp = strings.Replace(exp, exp[bid:eid], sid+" ", 1)
		fmt.Println(exp, mapid)
	}
	lk := "id" + strconv.Itoa(id)
	fmt.Println(lk)
	return mapid
}
