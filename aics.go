package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/duckdb/duckdb-go/v2"
)

type AIC struct {
	AICName string
	Range   []string
}

var AIC_DIRECTORY = "data/aic"

var query2 string = `COPY (
SELECT AICName,
       unnest(Range)::INET as CIDR,
       host(network(unnest(Range)::INET)) as Gateway,
       host(broadcast(unnest(Range)::INET)) as Broadcast
FROM read_ndjson_auto('data/db/aic_collected.json')
ORDER BY unnest(Range)::INET
) TO 'data/db/aics.parquet' (FORMAT parquet, COMPRESSION zstd);`

func getAics() {

	dir, err := os.ReadDir(filepath.Join(AIC_DIRECTORY))
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.OpenFile(filepath.Join(OUTPUT_DIRECTORY, "aic_collected.json"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	for _, d := range dir {
		var AICinstance AIC
		fn := d.Name()
		replacer := strings.NewReplacer("ip_", "", ".txt", "")
		aicName := replacer.Replace(fn)
		AICinstance.AICName = aicName
		f, err := os.Open(filepath.Join(AIC_DIRECTORY, fn))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			AICinstance.Range = append(AICinstance.Range, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		encoder := json.NewEncoder(out)
		if err := encoder.Encode(AICinstance); err != nil {
			log.Fatal(err)
		}
	}
	db, err := sql.Open("duckdb", "")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(query2)
	if err != nil {
		log.Fatal(err)
	}
}
