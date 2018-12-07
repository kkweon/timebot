package timebot

import (
	"strings"
	"time"
)

var koreaTZ *time.Location
var californiaTZ *time.Location

var formats = []string{
	"2006-01-02 15:04",
}

const KST = "KST"
const PST = "PST"
const PDT = "PDT"

func init() {
	var err error

	koreaTZ, err = time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}

	californiaTZ, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
}

// ParseTime parse time using sane defaults
func ParseTime(t string) (time.Time, error) {
	var newTime string
	var loc *time.Location

	if strings.Contains(t, KST) {
		newTime = cleanText(t, KST)
		loc = koreaTZ
	} else if strings.Contains(t, PST) {
		newTime = cleanText(t, PST)
		loc = californiaTZ
	} else if strings.Contains(t, PDT) {
		newTime = cleanText(t, PDT)
		loc = californiaTZ
	}

	return time.ParseInLocation(formats[0], newTime, loc)
}

func cleanText(t string, timeZone string) string {
	return strings.Trim(strings.Replace(t, timeZone, "", -1), " ")
}

// ToKoreaTime returns using KST
func ToKoreaTime(t time.Time) string {
	return t.In(koreaTZ).Format(formats[0] + " " + KST)
}

// ToCaliforniaTime returns using PST/PDT
func ToCaliforniaTime(t time.Time) string {
	return t.In(californiaTZ).Format(formats[0] + " MST")
}
