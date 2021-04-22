package simpleio

import (
	"log"
)

// ToNumber to covert roman numeral to decimal
func ToNumber(n string) int {
	out := 0
	ln := len(n)
	var c, cnext string
	var vc, vcnext int
	for i := 0; i < ln; i++ {
		c = string(n[i])
		vc = num[c]
		if i < ln-1 {
			cnext = string(n[i+1])
			vcnext = LookUpRoman(cnext)
			if vc < vcnext {
				out += vcnext - vc
				i++
			} else {
				out += vc
			}
		} else {
			out += vc
		}
	}
	return out
}

// ToRoman is to convert decimal number to roman numeral
func ToRoman(n int) string {
	out := ""
	var v int
	for n > 0 {
		v = highestDecimal(n)
		out += InvMap(v)
		n -= v
	}
	return out
}

func highestDecimal(n int) int {
	for _, v := range maxTable {
		if v <= n {
			return v
		}
	}
	return 1
}

var maxTable = []int{
	1000,
	900,
	500,
	400,
	100,
	90,
	50,
	40,
	10,
	9,
	5,
	4,
	1,
}

func LookUpRoman(s string) int {
	switch s {
	case "I":
		return 1
	case "V":
		return 5
	case "X":
		return 10
	case "L":
		return 50
	case "C":
		return 100
	case "M":
		return 1000
	}
	log.Fatalf("Error: did not find roman string...\n")
	return 0
}

var num = map[string]int{
	"I": 1,
	"V": 5,
	"X": 10,
	"L": 50,
	"C": 100,
	"D": 500,
	"M": 1000,
}

func InvMap(num int) string {
	switch num {
	case 1000:
		return "M"
	case 900:
		return "CM"
	case 500:
		return "D"
	case 400:
		return "CD"
	case 100:
		return "C"
	case 90:
		return "XC"
	case 50:
		return "L"
	case 40:
		return "XL"
	case 10:
		return "X"
	case 9:
		return "IX"
	case 5:
		return "V"
	case 4:
		return "IV"
	case 1:
		return "I"
	}
	log.Fatalf("Error: roman string does not exist for the given number")
	return ""
}
