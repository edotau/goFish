package simpleio

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
)

// BytesToInt a function that converts a byte slice and return a number, type int.
func BytesToInt(b []byte) int {
	answer, err := strconv.Atoi(string(b))
	ErrorHandle(err)
	return answer
}

// StringToInt is a function that converts a string and return a number, type int.
func StringToInt(s string) int {
	answer, err := strconv.Atoi(s)
	ErrorHandle(err)
	return answer
}

// StringToUint16 parses a string into a uint16 and will exit on error
func StringToUInt16(s string) uint16 {
	n, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		log.Panic(fmt.Sprintf("Error: trouble converting \"%s\" to a uint16\n", s))
	}
	return uint16(n)
}

// StringToInt is a function that converts a string and return a number, type int.
func StringToUInt32(s string) uint32 {
	answer, err := strconv.Atoi(s)
	ErrorHandle(err)
	return uint32(answer)
}

// StringToFloat is a function that converts a string to a type float64.
func StringToFloat(s string) float32 {
	answer, err := strconv.ParseFloat(s, 32)
	ErrorHandle(err)
	return float32(answer)
}

// ScientificNotation will convert a string with scientific notation into a float64
func ScientificNotation(s string) float64 {
	num, _, err := big.ParseFloat(s, 10, 0, big.ToNearestEven)
	ErrorHandle(err)
	ans, _ := num.Float64()
	return ans
}

// IntToString a function that converts a number of type int and return a string.
func IntToString(i int) string {
	return fmt.Sprintf("%d", i)
}

// StringToInts will process strings (usually from column data) and return a slice of []int
func StringToIntSlice(column string) []int {
	work := strings.Split(column, ",")
	sliceSize := len(work)
	if column[len(column)-1] == ',' {
		sliceSize--
	}
	var answer []int = make([]int, sliceSize)
	for i := 0; i < sliceSize; i++ {
		answer[i] = StringToInt(work[i])
	}
	return answer
}

// intListToString will process a slice of type int as an input and return a each value separated by a comma as a string.
func IntSliceToString(nums []int) string {
	ans := strings.Builder{}
	ans.Grow(2 * len(nums))
	for i := 0; i < len(nums); i++ {
		ans.WriteString(IntToString(nums[i]))
		ans.WriteByte(',')
	}
	return ans.String()
}

func Int16ToString(num int16) string {
	return strconv.FormatInt(int64(num), 10)
}
