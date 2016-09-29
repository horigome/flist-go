// flist-go
//   : parse dir and output csv
// ----------------------------------------------------------------------
// copyright 2016-09-29 M.Horigome
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// Params command line option values
type Params struct {
	matchString string // match filename (regular expression)
	skipString  string // match filename (regular expression)
	noCSV       bool   // false: no Output CSV
	specCSVFile string // csv filename specification
}

// CommandUsage commandline usage
var CommandUsage = `
flist go
  ver 0.1.0.0
  copyright M.Horigome

Usage : flist <options>

options:

`

// newParams get commandline params
func newParams() *Params {

	var p Params
	flag.StringVar(&p.matchString, "m", "", " File Match String(Regular expression)")
	flag.StringVar(&p.skipString, "s", "", " File Skip String(Regular expression)")
	flag.StringVar(&p.specCSVFile, "f", "", " Specify CSV filename")
	flag.BoolVar(&p.noCSV, "no", false, " CSV Not Output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, CommandUsage)
		flag.PrintDefaults()
	}
	flag.Parse()

	return &p
}

func isFileMatch(matchString string, s string) bool {
	r := regexp.MustCompile(matchString)
	return r.MatchString(s)
}

// writeCSV Write to CSV record
func writeCSV(writer *csv.Writer, basepath, path string, info os.FileInfo, head bool) {
	if writer == nil {
		return
	}

	if head {
		// Write to CSV Header
		writer.Write([]string{
			"basepath", "path", "filename", "mod-date", "mod-time",
		})
	} else {
		// Write to CSV Record
		writer.Write([]string{
			basepath,
			path,
			info.Name(),
			fmt.Sprintf("%04d-%02d-%02d",
				info.ModTime().Year(), info.ModTime().Month(), info.ModTime().Day()),
			fmt.Sprintf("%02d:%02d:%d",
				info.ModTime().Hour(), info.ModTime().Minute(), info.ModTime().Second()),
		})
	}
	writer.Flush()
}

// openCSV Open CSV File
func openCSV(filename string) (*os.File, *csv.Writer, error) {

	csvfile := filename
	if csvfile == "" {
		// default csv name
		t := time.Now()
		csvfile = fmt.Sprintf("filelist-%04d%02d%02d-%02d%02d%02d.csv",
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	} else {
		// make subdirectory
		d := filepath.Dir(filename)
		_, err := os.Stat(d)
		if !os.IsExist(err) {
			os.MkdirAll(d, 0777)
		}
	}

	file, err := os.OpenFile(csvfile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, nil, err
	}
	err = file.Truncate(0)
	writer := csv.NewWriter(file)

	return file, writer, nil
}

// main
func main() {

	// * Getting options
	p := newParams()

	// * CSV Open setting
	var file *os.File
	var writer *csv.Writer
	var err error

	if !p.noCSV {
		file, writer, err = openCSV(p.specCSVFile)
		if err == nil {
			defer file.Close()
			writeCSV(writer, "", "", nil, true)
		}
	}

	// * Parse dir
	root, _ := filepath.Abs(filepath.Dir("."))
	fmt.Println("\n(root)---> ", root)
	num := 0

	filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil // skip directory
			}
			if p.matchString != "" && !isFileMatch(p.matchString, info.Name()) {
				return nil
			}
			if p.skipString != "" && isFileMatch(p.skipString, info.Name()) {
				return nil // skip pattern
			}

			rel, _ := filepath.Rel(root, path)
			fmt.Println(rel)
			num++

			writeCSV(writer, root, rel, info, false)
			return nil
		})

	fmt.Printf("--------->  %v files\n", num)
}
