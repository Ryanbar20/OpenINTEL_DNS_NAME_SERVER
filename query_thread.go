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

func A_query(nd NameData, db *sql.DB) (*[]dns.RR, error) {

	result := make([]dns.RR, 0)

	query := fmt.Sprintf(`
		SELECT ip4_address 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%04d/month=%02d/day=%02d/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'A';
	`, nd.tld, nd.year, nd.month, nd.day, nd.domain)

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

	return &result, nil
}

func AAAA_query(nd NameData, db *sql.DB) (*[]dns.RR, error) {
	result := make([]dns.RR, 0)

	query := fmt.Sprintf(`
		SELECT ip6_address 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%04d/month=%02d/day=%02d/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'AAAA';
	`, nd.tld, nd.year, nd.month, nd.day, nd.domain)

	rows, err := db.Query(query)

	if err != nil {
		return nil, errors.New("rows could not be gotten")
	}
	defer rows.Close()
	var hdr = &dns.Header{Name: nd.domain + dom, Class: dns.ClassINET}

	for rows.Next() {
		var ip6addr string

		err := rows.Scan(&ip6addr)

		if err != nil {
			return nil, errors.New("row could not be parsed")

		}

		addr, err := netip.ParseAddr(ip6addr)

		if err != nil {
			return nil, errors.New("row could not be parsed")

		}
		result = append(result, &dns.AAAA{Hdr: *hdr, AAAA: rdata.AAAA{Addr: addr}})

	}

	return &result, nil
}

func TXT_query(nd NameData, db *sql.DB) (*[]dns.RR, error) {
	result := make([]dns.RR, 0)

	query := fmt.Sprintf(`
		SELECT txt_text 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%04d/month=%02d/day=%02d/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'TXT';
	`, nd.tld, nd.year, nd.month, nd.day, nd.domain)

	rows, err := db.Query(query)

	if err != nil {
		return nil, errors.New("rows could not be gotten")
	}
	defer rows.Close()
	var hdr = &dns.Header{Name: nd.domain + dom, Class: dns.ClassINET}

	for rows.Next() {
		var txt string

		err := rows.Scan(&txt)

		if err != nil {
			return nil, errors.New("row could not be parsed")

		}
		result = append(result, &dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{txt}}})

	}

	return &result, nil
}

func MX_query(nd NameData, db *sql.DB) (*[]dns.RR, error) {
	result := make([]dns.RR, 0)
	query := fmt.Sprintf(`
		SELECT mx_address, mx_preference 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%04d/month=%02d/day=%02d/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'MX';
	`, nd.tld, nd.year, nd.month, nd.day, nd.domain)

	rows, err := db.Query(query)

	if err != nil {
		return nil, errors.New("rows could not be gotten")
	}
	defer rows.Close()
	var hdr = &dns.Header{Name: nd.domain + dom, Class: dns.ClassINET}

	for rows.Next() {
		var mx_addr string
		var mx_pref uint16
		err := rows.Scan(&mx_addr, &mx_pref)

		if err != nil {
			return nil, errors.New("row could not be parsed")

		}
		result = append(result, &dns.MX{Hdr: *hdr, MX: rdata.MX{Mx: mx_addr, Preference: mx_pref}})

	}

	return &result, nil
}

func NS_query(nd NameData, db *sql.DB) (*[]dns.RR, error) {
	result := make([]dns.RR, 0)
	query := fmt.Sprintf(`
		SELECT ns_address 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%04d/month=%02d/day=%02d/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'NS';
	`, nd.tld, nd.year, nd.month, nd.day, nd.domain)

	rows, err := db.Query(query)

	if err != nil {
		return nil, errors.New("rows could not be gotten")
	}
	defer rows.Close()
	var hdr = &dns.Header{Name: nd.domain + dom, Class: dns.ClassINET}

	for rows.Next() {
		var ns_addr string

		err := rows.Scan(&ns_addr)

		if err != nil {
			return nil, errors.New("row could not be parsed")

		}
		result = append(result, &dns.NS{Hdr: *hdr, NS: rdata.NS{Ns: ns_addr}})

	}

	return &result, nil
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

	// message refusing is handled in name_server
	// the checks on data parsing and message type are here for added security
	for {
		question := query_queue.PeekBlocking() // get the question

		name := question.Header().Name
		data, success := parseName(name)
		if !success {
			query_queue.PopBlocking() // remove the question from the queue
			continue
		}
		switch question.(type) {
		case *dns.A:
			rrs, err := A_query(data, db)
			if err == nil {
				cache.Put(question, rrs)
			}
		case *dns.AAAA:
			rrs, err := AAAA_query(data, db)
			if err == nil {
				cache.Put(question, rrs)
			}
		case *dns.TXT:
			rrs, err := TXT_query(data, db)
			if err == nil {
				cache.Put(question, rrs)
			}
		case *dns.MX:
			rrs, err := MX_query(data, db)
			if err == nil {
				cache.Put(question, rrs)
			}
		case *dns.NS:
			rrs, err := NS_query(data, db)
			if err == nil {
				cache.Put(question, rrs)
			}
		default:
		}

		query_queue.PopBlocking() // remove the question from the queue
	}
}
