package ns

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/netip"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/rdata"
	_ "github.com/duckdb/duckdb-go/v2"
)

func A_query(nd NameData, db *sql.DB) ([]*dns.A, error) {

	result := make([]*dns.A, 0)

	query := fmt.Sprintf(`
		SELECT ip4_address 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%04d/month=%02d/day=%02d/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'A';
	`, nd.tld, nd.year, nd.month, nd.day, dom)

	rows, err := db.Query(query)

	if err != nil {
		return nil, errors.New("rows could not be gotten")
	}
	defer rows.Close()
	var hdr = &dns.Header{Name: nd.domain + dom, Class: dns.ClassINET}

	for rows.Next() {
		var ip4addr string

		err := rows.Scan(&ip4addr)

		if err != nil {
			return nil, errors.New("row could not be parsed")

		}

		addr, err := netip.ParseAddr(ip4addr)

		if err != nil {
			return nil, errors.New("row could not be parsed")

		}
		result = append(result, &dns.A{Hdr: *hdr, A: rdata.A{Addr: addr}})

	}

	return result, nil
}

func AAAA_query(nd NameData, db *sql.DB) ([]*dns.RR, error) {
	query := fmt.Sprintf(`
		SELECT query_name, query_type, ip6_address 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%04d/month=%02d/day=%02d/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'AAAA';
	`, nd.tld, nd.year, nd.month, nd.day, dom)

	rows, err := db.Query(query)

	if err != nil {
		return nil, errors.New("rows could not be gotten")
	}
	defer rows.Close()
	return nil, nil
}

func TXT_query(nd NameData, db *sql.DB) ([]*dns.RR, error) {
	query := fmt.Sprintf(`
		SELECT query_name, query_type, txt_text 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%04d/month=%02d/day=%02d/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'TXT';
	`, nd.tld, nd.year, nd.month, nd.day, dom)

	rows, err := db.Query(query)

	if err != nil {
		return nil, errors.New("rows could not be gotten")
	}
	defer rows.Close()
	return nil, nil
}

func MX_query(nd NameData, db *sql.DB) ([]*dns.RR, error) {
	query := fmt.Sprintf(`
		SELECT query_name, query_type, mx_address, mx_preference 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%04d/month=%02d/day=%02d/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'MX';
	`, nd.tld, nd.year, nd.month, nd.day, dom)

	rows, err := db.Query(query)

	if err != nil {
		return nil, errors.New("rows could not be gotten")
	}
	defer rows.Close()
	return nil, nil
}

func NS_query(nd NameData, db *sql.DB) ([]*dns.RR, error) {
	query := fmt.Sprintf(`
		SELECT query_name, query_type, ns_address 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%04d/month=%02d/day=%02d/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'NS';
	`, nd.tld, nd.year, nd.month, nd.day, dom)

	rows, err := db.Query(query)

	if err != nil {
		return nil, errors.New("rows could not be gotten")
	}
	defer rows.Close()
	return nil, nil
}

func query_thread(query_queue *Queue, cache *Cache) {

	db, err := sql.Open("duckdb", "")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	setup := []string{
		"INSTALL httpfs;",
		"LOAD httpfs;",
		"SET s3_region='nl-utwente';",
		"SET s3_url_style='path';",
		"SET s3_endpoint='object.openintel.nl';",
		"SET s3_use_ssl=true;",
		"SET threads=1;",
	}

	for _, q := range setup {
		if _, err := db.Exec(q); err != nil {
			log.Fatal(err)
		}
	}

	for {
		question := query_queue.PeekBlocking() // get the question

		name := question.Header().Name
		data, success := parseName(name)
		if !success {
			// refuse
			continue
		}
		var rrs *[]dns.RR
		switch question.(type) {
		case *dns.A:
			A_query(data, db)
		case *dns.AAAA:
			A_query(data, db)
		case *dns.TXT:
			A_query(data, db)
		case *dns.MX:
			A_query(data, db)
		case *dns.NS:
			A_query(data, db)
		default:
			// refuse
		}
		cache.Put(question, rrs)

		query_queue.PopBlocking() // remove the question from the queue
	}
}
