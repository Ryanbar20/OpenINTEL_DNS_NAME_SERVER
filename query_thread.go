package ns

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/netip"
	"strings"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/rdata"
	_ "github.com/duckdb/duckdb-go/v2"
)

func getQueryString(nd NameData, columns []string, qtype string) string {
	cols := strings.Join(columns, ", ")
	return fmt.Sprintf(`
		SELECT %s
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%04d/month=%02d/day=%02d/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = '%s';
	`, cols, nd.tld, nd.year, nd.month, nd.day, nd.domain, qtype)
}

func rrQuery(db *sql.DB, query string, scan func(*sql.Rows) (dns.RR, error)) (*[]dns.RR, error) {
	result := make([]dns.RR, 0)

	rows, err := db.Query(query)

	if err != nil {
		return nil, errors.New("rows could not be gotten")
	}
	defer rows.Close()

	for rows.Next() {
		rr, err := scan(rows)
		if err != nil {
			return nil, errors.New("rows could not be parsed")
		}
		result = append(result, rr)
	}
	return &result, nil
}

func aQuery(nd NameData, db *sql.DB) (*[]dns.RR, error) {
	query := getQueryString(nd, []string{"ip4_address"}, "A")
	return rrQuery(db, query, func(rows *sql.Rows) (dns.RR, error) {
		var ip4addr string

		err := rows.Scan(&ip4addr)

		if err != nil {
			return nil, errors.New("row could not be parsed")

		}

		addr, err := netip.ParseAddr(ip4addr)

		if err != nil {
			return nil, errors.New("row could not be parsed")

		}
		return &dns.A{Hdr: dns.Header{Name: nd.domain + dom, Class: dns.ClassINET}, A: rdata.A{Addr: addr}}, nil
	})
}

func aaaaQuery(nd NameData, db *sql.DB) (*[]dns.RR, error) {
	query := getQueryString(nd, []string{"ip6_address"}, "AAAA")
	return rrQuery(db, query, func(rows *sql.Rows) (dns.RR, error) {
		var ip6addr string

		err := rows.Scan(&ip6addr)

		if err != nil {
			return nil, errors.New("row could not be parsed")

		}

		addr, err := netip.ParseAddr(ip6addr)

		if err != nil {
			return nil, errors.New("row could not be parsed")

		}
		return &dns.AAAA{Hdr: dns.Header{Name: nd.domain + dom, Class: dns.ClassINET}, AAAA: rdata.AAAA{Addr: addr}}, nil
	})
}

func txtQuery(nd NameData, db *sql.DB) (*[]dns.RR, error) {
	query := getQueryString(nd, []string{"txt_text"}, "TXT")

	return rrQuery(db, query, func(rows *sql.Rows) (dns.RR, error) {
		var txt string

		err := rows.Scan(&txt)

		if err != nil {
			return nil, errors.New("row could not be parsed")

		}
		return &dns.TXT{Hdr: dns.Header{Name: nd.domain + dom, Class: dns.ClassINET}, TXT: rdata.TXT{Txt: []string{txt}}}, nil
	})
}

func mxQuery(nd NameData, db *sql.DB) (*[]dns.RR, error) {
	query := getQueryString(nd, []string{"mx_address", "mx_preference"}, "MX")

	return rrQuery(db, query, func(rows *sql.Rows) (dns.RR, error) {
		var mx_addr string
		var mx_pref uint16
		err := rows.Scan(&mx_addr, &mx_pref)

		if err != nil {
			return nil, errors.New("row could not be parsed")

		}
		return &dns.MX{Hdr: dns.Header{Name: nd.domain + dom, Class: dns.ClassINET},
			MX: rdata.MX{Mx: mx_addr, Preference: mx_pref}}, nil

	})
}

func nsQuery(nd NameData, db *sql.DB) (*[]dns.RR, error) {
	query := getQueryString(nd, []string{"ns_address"}, "NS")

	return rrQuery(db, query, func(rows *sql.Rows) (dns.RR, error) {
		var ns_addr string
		err := rows.Scan(&ns_addr)

		if err != nil {
			return nil, errors.New("row could not be parsed")

		}
		return &dns.NS{Hdr: dns.Header{Name: nd.domain + dom, Class: dns.ClassINET}, NS: rdata.NS{Ns: ns_addr}}, nil
	})
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
		rrs, err := handleQuestion(question, db, cache)
		if err == nil {
			cache.Put(question, rrs)
		} else {
			cache.Put(question, nil) // if there is no answer, put that as an answer
		}
		query_queue.PopBlocking() // remove the question from the queue
	}
}

func handleQuestion(question dns.RR, db *sql.DB, cache *Cache) (*[]dns.RR, error) {
	name := question.Header().Name
	data, success := parseName(name)
	if !success {
		return nil, errors.New("Invalid name format")
	}
	var rrs *[]dns.RR = nil
	var err error = nil
	switch question.(type) {
	case *dns.A:
		rrs, err = aQuery(data, db)
	case *dns.AAAA:
		rrs, err = aaaaQuery(data, db)
	case *dns.TXT:
		rrs, err = txtQuery(data, db)
	case *dns.MX:
		rrs, err = mxQuery(data, db)
	case *dns.NS:
		rrs, err = nsQuery(data, db)
	default:
		err = errors.New("Unsupported Type")
	}
	return rrs, err
}
