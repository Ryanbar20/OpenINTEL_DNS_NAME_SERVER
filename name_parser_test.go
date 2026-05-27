package ns

import (
	"fmt"
	"testing"
)

func TestNameParser(t *testing.T) {
	type questionDataPair struct {
		question string
		data     Name_data
	}
	validQueries := []questionDataPair{
		{question: "20060202.www.example.com.history.openintel.nl", data: Name_data{tld: "com", domain: "www.example.com.", year: 2006, month: 2, day: 2}},
		{question: "20060302.example.com.history.openintel.nl", data: Name_data{tld: "com", domain: "example.com.", year: 2006, month: 3, day: 2}},
		{question: "www.20070202.www.example.com.history.openintel.nl", data: Name_data{tld: "com", domain: "www.example.com.", year: 2007, month: 2, day: 2}},
		{question: "20060202.www.test.abc.edf.com.history.openintel.nl.", data: Name_data{tld: "com", domain: "www.test.abc.edf.com.", year: 2006, month: 2, day: 2}},
	}
	for _, q := range validQueries {
		if n, b := parseName(q.question); b != true || n != q.data {
			t.Error(fmt.Sprintf("'%s' invalidly parsed as %v", q.question, n))
		}

	}

	invalidQueries := []string{
		"200600202.www.example.com.history.openintel.nl",
		"www.www.20060202.www.example.com.history.openintel.nl",
		"20060202.www.example.com.history.openintel.cl",
		"20060202.www.example.com.history.opedintel.nl",
		"20060202.www.example.com.hisrory.openintel.nl",
		"20060202.history.openintel.nl",
	}
	for _, q := range invalidQueries {
		if n, b := parseName(q); b == true {
			t.Error(fmt.Sprintf("'%s' should fail %v %v", q, n, b))
		}

	}

}
