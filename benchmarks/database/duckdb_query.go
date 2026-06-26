package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Provide a benchmark number (1-6) to execute")
		return
	}

	benchmark, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Provide a valid integer for the benchmark number")
		return
	}

	// setup
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
		"SET memory_limit='4GB';",
	}

	for _, q := range setup {
		if _, err := db.Exec(q); err != nil {
			log.Fatal(err)
		}
	}

	// benchmarks
	switch benchmark {
	case 1:
		bench1(db)
	case 2:
		bench2(db)
	case 3:
		bench3(db)
	case 4:
		bench4(db)
	case 5:
		bench5(db)
	case 6:
		bench6(db)
	default:
		fmt.Println("Benchmark number should be 1,2,3,4,5 or 6")

	}
}

// query .nu on 01-october-2023 for the A record(s) of google.nu.
func bench1(db *sql.DB) {
	fmt.Println("Bench1")
	start := time.Now()
	query := `
		SELECT ip4_address
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=nu/year=2023/month=10/day=01/*.gz.parquet') 
		WHERE query_name = 'google.nu.' AND query_type = 'A';
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	end := time.Now()

	for rows.Next() {
		var ipAddr string
		err := rows.Scan(&ipAddr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", ipAddr)
	}
	fmt.Printf("Start : \t %v\nEnd : \t %v\nDuration : \t %v\n", start, end, end.Sub(start))
}

// query .fr on 01-october-2023 for the A record(s) of google.fr.
func bench2(db *sql.DB) {
	fmt.Println("Bench2")
	start := time.Now()
	query := `
		SELECT ip4_address
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=fr/year=2023/month=10/day=01/*.gz.parquet') 
		WHERE query_name = 'google.fr.' AND query_type = 'A';
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	end := time.Now()

	for rows.Next() {
		var ipAddr string
		err := rows.Scan(&ipAddr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", ipAddr)
	}
	fmt.Printf("Start : \t %v\nEnd : \t %v\nDuration : \t %v\n", start, end, end.Sub(start))
}

// query .nu on 01-october-2023 for the AAAA record(s) of google.nu.
func bench3(db *sql.DB) {
	fmt.Println("Bench3")
	start := time.Now()
	query := `
		SELECT ip6_address
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=nu/year=2023/month=10/day=01/*.gz.parquet') 
		WHERE query_name = 'google.nu.' AND query_type = 'AAAA';
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	end := time.Now()

	for rows.Next() {
		var ipAddr string
		err := rows.Scan(&ipAddr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", ipAddr)
	}
	fmt.Printf("Start : \t %v\nEnd : \t %v\nDuration : \t %v\n", start, end, end.Sub(start))
}

// query .nu on 01-october-2023 for the MX record(s) of google.nu.
func bench4(db *sql.DB) {
	fmt.Println("Bench4")
	start := time.Now()
	query := `
		SELECT mx_address, mx_preference
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=nu/year=2023/month=10/day=01/*.gz.parquet') 
		WHERE query_name = 'google.nu.' AND query_type = 'MX';
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	end := time.Now()

	for rows.Next() {
		var mx_addr string
		var mx_pref string
		err := rows.Scan(&mx_addr, &mx_pref)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s \t| %s\n", mx_addr, mx_pref)
	}
	fmt.Printf("Start : \t %v\nEnd : \t %v\nDuration : \t %v\n", start, end, end.Sub(start))

}

// query .nu on 01-october-2023 for the NS record(s) of google.nu.
func bench5(db *sql.DB) {
	fmt.Println("Bench5")
	start := time.Now()
	query := `
		SELECT ns_address
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=nu/year=2023/month=10/day=01/*.gz.parquet') 
		WHERE query_name = 'google.nu.' AND query_type = 'NS';
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	end := time.Now()

	for rows.Next() {
		var ns_addr string
		err := rows.Scan(&ns_addr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", ns_addr)
	}
	fmt.Printf("Start : \t %v\nEnd : \t %v\nDuration : \t %v\n", start, end, end.Sub(start))

}

// query .nu on 01-october-2023 for the TXT record(s) of google.nu.
func bench6(db *sql.DB) {
	fmt.Println("Bench6")
	start := time.Now()
	query := `
		SELECT txt_text
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=nu/year=2023/month=10/day=01/*.gz.parquet') 
		WHERE query_name = 'google.nu.' AND query_type = 'TXT';
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	end := time.Now()

	for rows.Next() {
		var text string
		err := rows.Scan(&text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", text)
	}
	fmt.Printf("Start : \t %v\nEnd : \t %v\nDuration : \t %v\n", start, end, end.Sub(start))

}
