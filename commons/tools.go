package commons

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"io"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ReadLinesOfFile(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {

	}
	lines := strings.Split(string(content), "\n")
	return lines
}

func TimeInSeconds() int64 {
	now := time.Now()
	return now.Unix()
}

func TimeInMilliseconds() int64 {
	now := time.Now()
	return int64(now.UnixNano() / 1000000)
}

func TimeInNanoseconds() int64 {
	now := time.Now()
	return now.UnixNano()
}

func YYYYMMDDHH() string {
	now := time.Now()
	return fmt.Sprintf("%d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour())
}

func YYYYMMDD() string {
	now := time.Now()
	return fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
}

func Try(body func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	body()
}

func eb64(s string) string {
	return base64.URLEncoding.EncodeToString([]byte(s))
}

func db64(s string) string {
	output, _ := base64.URLEncoding.DecodeString(s)
	return string(output)
}

func stringTimestamp() string {
	now := time.Now().UnixNano()
	return strconv.FormatInt(now, 10)
}

func String2Milliseconds(t string) int {
	last := t[len(t)-1]
	digits := t[0 : len(t)-1]
	num, _ := strconv.Atoi(digits)
	switch last {
	case 's':
		return num * 1000
	case 'm':
		return num * 60 * 1000
	case 'h':
		return num * 3600 * 1000
	case 'd':
		return num * 24 * 3600 * 1000
	default:
		return num * 1000
	}
}

func UUID() string {
	u := uuid.New()
	return u.String()
}

func SHA1(str string) string {
	hasher := sha1.New()
	hasher.Write([]byte(str))
	sha := eb64(string(hasher.Sum(nil)))
	return sha
}

var charsets = map[string]encoding.Encoding{
	"big5":         traditionalchinese.Big5,
	"euc-jp":       japanese.EUCJP,
	"gb2312":       simplifiedchinese.GBK,
	"iso-2022-jp":  japanese.ISO2022JP,
	"iso-8859-1":   charmap.ISO8859_1,
	"iso-8859-2":   charmap.ISO8859_2,
	"iso-8859-3":   charmap.ISO8859_3,
	"iso-8859-4":   charmap.ISO8859_4,
	"iso-8859-10":  charmap.ISO8859_10,
	"iso-8859-13":  charmap.ISO8859_13,
	"iso-8859-14":  charmap.ISO8859_14,
	"iso-8859-15":  charmap.ISO8859_15,
	"iso-8859-16":  charmap.ISO8859_16,
	"koi8-r":       charmap.KOI8R,
	"shift_jis":    japanese.ShiftJIS,
	"windows-1250": charmap.Windows1250,
	"windows-1251": charmap.Windows1251,
	"windows-1252": charmap.Windows1252,
}

func DetectContentCharset(body io.Reader) string {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		if _, name, ok := charset.DetermineEncoding(data, ""); ok {
			return name
		}
	}
	return "utf-8"
}

// Reader returns an io.Reader that converts the provided charset to UTF-8.
func Reader(charset string, input io.Reader) (io.Reader, error) {
	charset = strings.ToLower(charset)
	if charset == "utf-8" || charset == "us-ascii" {
		return input, nil
	}
	if enc, ok := charsets[charset]; ok {
		return enc.NewDecoder().Reader(input), nil
	}
	return nil, fmt.Errorf("unhandled charset %q", charset)
}

func Shuffle(vals []string) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
		n := len(vals)
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
		vals = vals[:n-1]
	}
}
