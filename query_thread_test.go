package ns

import "testing"

func TestGetQueryString(t *testing.T) {
	nameData := NameData{tld: "com", domain: "www.example.com.", year: 2006, month: 2, day: 2}
	query := getQueryString(nameData, []string{"test", "1", "2", "3"}, "TESTQTYPE")
	if query != `
		SELECT test, 1, 2, 3
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=com/year=2006/month=02/day=02/*.gz.parquet') 
		WHERE query_name = 'www.example.com.' AND query_type = 'TESTQTYPE';
	` {
		t.Error("Query string generation failed")
	}
}
