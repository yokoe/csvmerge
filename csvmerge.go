package csvmerge

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Result struct {
	File    string
	Header  []string
	NumRows int
}

// Returns a merge result contains path of generated file, header columns or an error
func MergeCSVFiles(files []string) (*Result, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files")
	}

	tmpFile, err := ioutil.TempFile("", "mergedcsv")
	if err != nil {
		return nil, fmt.Errorf("temp file error: %w", err)
	}
	outPath := tmpFile.Name()

	w := csv.NewWriter(tmpFile)

	header := []string{}
	rowsCount := 0
	for _, file := range files {
		if filepath.Ext(file) != ".csv" {
			continue
		}
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("file open error: %w", err)
		}
		defer f.Close()

		r := csv.NewReader(f)

		// skip header
		firstRow, err := r.Read()
		if err != nil {
			return nil, fmt.Errorf("header read error: %w", err)
		}
		if len(header) == 0 {
			header = firstRow
			w.Write(header)
		} else {
			// check header column orders
			if len(header) != len(firstRow) {
				return nil, fmt.Errorf("different header length at %s. %d - %d", file, len(header), len(firstRow))
			}

			for j, col := range header {
				if firstRow[j] != col {
					return nil, fmt.Errorf("different column %s found in %s", firstRow[j], file)
				}
			}
		}

		for {
			row, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("csv file read error: %w", err)
			}
			w.Write(row)
			rowsCount++
		}
		w.Flush()
	}
	return &Result{File: outPath, Header: header, NumRows: rowsCount}, nil
}
