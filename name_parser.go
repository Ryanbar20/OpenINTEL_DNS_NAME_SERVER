package ns

import (
	"strings"
	"time"
)

type Name_data struct {
	tld    string
	domain string
	year   int
	month  int
	day    int
}

// returns the data in the domain name
func parseName(name string) (data Name_data, b bool) {
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
	if l < 5 ||
		labels[l-3] != "history" ||
		labels[l-2] != "openintel" ||
		labels[l-1] != "nl" {
		return Name_data{}, false
	}

	ymd, err := time.Parse("20060102", labels[0])
	if err != nil { // if not a valid YYYYMMDD
		return Name_data{}, false
	}

	y, m, d := ymd.Year(), int(ymd.Month()), ymd.Day()

	// add . to the end as openINTEL stores domains this way
	domain := strings.Join(labels[1:l-3], ".") + "."

	tld := labels[l-4]

	return Name_data{tld: tld, domain: domain, year: y, month: m, day: d}, true

}
