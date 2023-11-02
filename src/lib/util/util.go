package util

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"html"
	"io"
	"math"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Buf struct {
	data     []string
	delList  []int
	isDelete bool
	str      string
}

func NewBuf() *Buf {
	rv := new(Buf)
	rv.data = make([]string, 0, 20)
	rv.delList = make([]int, 0, 20)
	rv.isDelete = true

	return rv
}

func (buf *Buf) Clear() {
	buf.data = make([]string, 0, 20)
	buf.delList = make([]int, 0, 20)
	buf.isDelete = true
}

func (buf *Buf) Add(format string, args ...interface{}) {
	if format == "" {
		return
	}

	str := format
	if len(args) > 0 {
		str = fmt.Sprintf(format, args...)
	}

	buf.data = append(buf.data, str)
}

func (buf *Buf) AddLater(format string, args ...interface{}) {
	buf.Add(format, args...)
	buf.delList = append(buf.delList, len(buf.data)-1)
}

func (buf *Buf) Forge() {
	buf.isDelete = false
}

func (buf *Buf) clearDelList() {
	for _, v := range buf.delList {
		buf.data[v] = ""
	}
}

func (buf *Buf) Copy(srcBuf *Buf) {
	for _, v := range srcBuf.data {
		buf.Add(v)
	}
}

func (buf *Buf) String() *string {
	if buf.isDelete {
		buf.clearDelList()
	}

	buf.str = strings.Join(buf.data, "\n")
	return &buf.str
}

func (buf *Buf) StringSep(sep string) *string {
	if buf.isDelete {
		buf.clearDelList()
	}

	buf.str = strings.Join(buf.data, sep)
	return &buf.str
}

func (buf *Buf) Len() int {
	return len(buf.data)
}

func (buf *Buf) IsEmpty() bool {
	return len(buf.data) == 0
}

func NewUUID() string {
	id := uuid.New()
	return id.String()
}

func ScrStr(str string) string {
	return html.EscapeString(str)
}

func DbStr(v string) string {
	v = strings.ReplaceAll(v, "%", "")
	v = strings.ReplaceAll(v, "?", "")

	buf := make([]byte, len(v)*2)
	pos := 0

	for i := 0; i < len(v); i++ {
		c := v[i]
		switch c {
		case '\x00':
			buf[pos] = '\\'
			buf[pos+1] = '0'
			pos += 2
		case '\n':
			buf[pos] = '\\'
			buf[pos+1] = 'n'
			pos += 2
		case '\r':
			buf[pos] = '\\'
			buf[pos+1] = 'r'
			pos += 2
		case '\x1a':
			buf[pos] = '\\'
			buf[pos+1] = 'Z'
			pos += 2
		case '\'':
			buf[pos] = '\\'
			buf[pos+1] = '\''
			pos += 2
		case '"':
			buf[pos] = '\\'
			buf[pos+1] = '"'
			pos += 2
		case '\\':
			buf[pos] = '\\'
			buf[pos+1] = '\\'
			pos += 2
		default:
			buf[pos] = c
			pos += 1
		}
	}

	return string(buf[:pos])
}

func PasswordHash(str string) string {
	h := sha1.New()
	io.WriteString(h, str+"987$#sd")
	return fmt.Sprintf("%X", h.Sum(nil))
}

var intRe *regexp.Regexp = regexp.MustCompile("(\\d+)(\\d{3})")

func FormatInt(num int64) string {
	numStr := strconv.FormatInt(num, 10)

	for {
		rv := intRe.ReplaceAllString(numStr, "$1.$2")
		if rv == numStr {
			return rv
		}

		numStr = rv
	}
}

func PadLeft(str string, padStr string, pLen int) string {
	p := pLen - len(str)
	if p <= 0 {
		return str
	}

	return strings.Repeat(padStr, p) + str
}

func PadRight(str string, padStr string, pLen int) string {
	p := pLen - len(str)
	if p <= 0 {
		return str
	}

	return str + strings.Repeat(padStr, p)
}

func ShortText(str string, length int) (string, bool) {
	wordList := strings.Fields(str)
	if len(wordList) > length {
		wordList = wordList[0:length]
		return strings.Join(wordList, " "), true
	}

	return str, false
}

func LimitStr(str string, length int) string {
	if len(str) < length {
		return str
	}

	return str[0:length]
}

func NullToString(val sql.NullString) string {
	if !val.Valid {
		return ""
	}

	return val.String
}

func NullToInt64(val sql.NullInt64) int64 {
	if !val.Valid {
		return 0
	}

	return val.Int64
}

func PrintStack(err error) {
	list := strings.Split(string(debug.Stack()[:]), "\n")
	println(strings.Join(list[6:], "\n"))
}

func IsFloatEqual(a, b float64) bool {
	tolerance := 0.000001
	if diff := math.Abs(a - b); diff < tolerance {
		return true
	}

	return false
}
