package ns

import (
	"strings"
	"time"
)

type NameData struct {
	tld    string
	domain string
	year   int
	month  int
	day    int
}

// parses a domain name and returns the data from it
func parseName(name string) (data NameData, b bool) {
	labels := strings.Split(name, ".")
	if labels[0] == "www" { // ignore www label
		labels = labels[1:]
	}
	l := len(labels)
	if labels[l-1] == "" { // drop the trailing . if it is there
		labels = labels[:l-1]
		l -= 1
	}

	// ignore if name does not end in history.openintel.nl or if it has no domain
	// (<5 labels means that there is no space to put a domain)
	if l < 5 ||
		labels[l-3] != "history" ||
		labels[l-2] != "openintel" ||
		labels[l-1] != "nl" {
		return NameData{}, false
	}

	ymd, err := time.Parse("20060102", labels[0])
	if err != nil { // if not a valid YYYYMMDD
		return NameData{}, false
	}

	y, m, d := ymd.Year(), int(ymd.Month()), ymd.Day()

	// add . to the end as openINTEL stores domains this way
	domain := strings.Join(labels[1:l-3], ".") + "."

	tld := labels[l-4]

	return NameData{tld: tld, domain: domain, year: y, month: m, day: d}, true

}
