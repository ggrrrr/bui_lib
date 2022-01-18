package cassandra

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

func readCqlFile(f string) (string, error) {
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return "n/a", err
	}
	cql := string(body)
	err = Session.ExecStmt(cql)
	return cql, err
}

func CreateSchema(subdir string) error {
	if Session == nil {
		fmt.Printf("please use db.Connect first\n")
		return fmt.Errorf("please use db.Connect first")
	}
	var errors []error
	filter := fmt.Sprintf("%s/%s/*.cql", cqlBaseDir, subdir)
	log.Printf("CreateSchema: %s", filter)
	files, err := filepath.Glob(filter)
	if err != nil {
		return err
	}
	for i, file := range files {
		cql, err := readCqlFile(file)
		log.Printf("cql[%d/%d] %s", len(files), i+1, file)
		if err != nil {
			errors = append(errors, err)
			log.Printf("error during parsing file: %s %v CQL: %v", file, err, cql)
			continue
		}
	}
	return nil
}
