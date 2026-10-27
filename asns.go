package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	_ "github.com/duckdb/duckdb-go/v2"
)

var ASN_DIRECTORY = "data/as"
var OUTPUT_DIRECTORY = "data/db"

type MetaData struct {
	Handle      string
	Description string
	CountryCode string
	Country     string
	Origin      string
	Category    string
	NetworkRole string
}

type Prefixes struct {
	IPv4 []string
	IPv6 []string
}

type ASNMetaData struct {
	ASN      int
	MetaData MetaData
	Prefixes Prefixes
}

type TargetASNData struct {
	ASN         int
	Handle      string
	Description string
	CountryCode string
	Country     string
	Origin      string
	Category    string
	NetworkRole string
	IPv4        []string
	IPv6        []string
}

var ASN ASNMetaData

var query string = `COPY (
    WITH unnested AS (
        SELECT ASN, HANDLE, Description, CountryCode, Country, Origin, Category, NetworkRole,
               'v4' as Family,
               unnest(IPv4)::INET as CIDR,
        FROM read_ndjson_auto('data/db/asn_collected.json')
        UNION ALL
        SELECT ASN, Handle, Description, CountryCode, Country, Origin, Category, NetworkRole,
               'v6' as Family
               unnest(IPv6)::INET as CIDR
        FROM read_ndjson_auto('data/db/asn_collected.json')
    ), with_gateway_broadcast AS (
        SELECT *,
               host(network(CIDR)) as Gateway,
               host(broadcast(CIDR)) as Broadcast
        FROM unnested
    )
    SELECT *
    FROM with_gateway_broadcast
) TO 'data/db/asns.parquet' (FORMAT 'parquet', COMPRESSION zstd);`

func getAsns() {
	dir, err := os.ReadDir(ASN_DIRECTORY)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.OpenFile(filepath.Join(OUTPUT_DIRECTORY, "asn_collected.json"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	for _, d := range dir {
		var ASN *ASNMetaData
		var TargetASN TargetASNData
		f, err := os.Open(filepath.Join(ASN_DIRECTORY, d.Name(), "aggregated.json"))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		decoder := json.NewDecoder(f)
		if err := decoder.Decode(&ASN); err != nil {
			log.Fatal(err)
		}

		TargetASN.ASN = ASN.ASN
		TargetASN.Handle = ASN.MetaData.Handle
		TargetASN.Description = ASN.MetaData.Description
		TargetASN.CountryCode = ASN.MetaData.CountryCode
		TargetASN.Country = ASN.MetaData.Country
		TargetASN.Origin = ASN.MetaData.Origin
		TargetASN.Category = ASN.MetaData.Category
		TargetASN.NetworkRole = ASN.MetaData.NetworkRole
		TargetASN.IPv4 = ASN.Prefixes.IPv4
		TargetASN.IPv6 = ASN.Prefixes.IPv6

		encoder := json.NewEncoder(out)
		if err := encoder.Encode(TargetASN); err != nil {
			log.Fatal(err)
		}
	}

	db, err := sql.Open("duckdb", "")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
