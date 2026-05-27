package ns

import "fmt"

func A_query(tld string, year string, month string, day string, dom string) string {
	return fmt.Sprintf(`
		SELECT query_name, query_type, ip4_address 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%s/month=%s/day=%s/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'A';
	`, tld, year, month, day, dom)
}

func AAAA_query(tld string, year string, month string, day string, dom string) string {
	return fmt.Sprintf(`
		SELECT query_name, query_type, ip6_address 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%s/month=%s/day=%s/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'AAAA';
	`, tld, year, month, day, dom)
}

func TXT_query(tld string, year string, month string, day string, dom string) string {
	return fmt.Sprintf(`
		SELECT query_name, query_type, txt_text 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%s/month=%s/day=%s/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'TXT';
	`, tld, year, month, day, dom)
}

func MX_query(tld string, year string, month string, day string, dom string) string {
	return fmt.Sprintf(`
		SELECT query_name, query_type, mx_address, mx_preference 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%s/month=%s/day=%s/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'MX';
	`, tld, year, month, day, dom)
}

func NS_query(tld string, year string, month string, day string, dom string) string {
	return fmt.Sprintf(`
		SELECT query_name, query_type, ns_address 
		FROM read_parquet('s3://openintel-public/fdns/basis=zonefile/source=%s/year=%s/month=%s/day=%s/*.gz.parquet') 
		WHERE query_name = '%s' AND query_type = 'NS';
	`, tld, year, month, day, dom)
}

func query_thread() {

}
