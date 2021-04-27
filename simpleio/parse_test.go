package simpleio

import (
	"bytes"
	"strings"
	"testing"
)

var strs = []string{"a", "b", "c", "d", "aaa", "bbb", "ccc", "ddd", "aaaa", "bbbb", "cccc", "dddd"}

func BenchmarkConcatA(b *testing.B) {
	var buf bytes.Buffer
	for n := 0; n < b.N; n++ {
		buf.WriteString("START")
		buf.WriteString(strs[0])
		for _, v := range strs {
			buf.WriteString("-")
			buf.WriteString(v)
		}
		buf.WriteString("END")
		b.Logf("%s\n", buf.Bytes())
		buf.Reset()
	}
}

func BenchmarkConcatB(b *testing.B) {
	var buf bytes.Buffer
	for n := 0; n < b.N; n++ {
		buf.WriteString("START")
		buf.WriteString(strings.Join(strs, "-"))
		buf.WriteString("END")
		b.Logf("%s\n", buf.Bytes())
		buf.Reset()
	}
}

func BenchmarkConcatD(b *testing.B) {
	n := len("START") + len("END")

	for t := 0; t < b.N; t++ {
		for _, s := range strs {
			n += len(s) + 1
		}
		buf := make([]byte, n)
		n = copy(buf, "START")

		for _, s := range strs {
			n += copy(buf[n:], s)
			buf[n] = '-'
			n++
		}
		copy(buf[n:], "END")
	}
}
