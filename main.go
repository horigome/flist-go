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

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var version = "0.0.0.0"

// Params command line option values
type Params struct {
	matchString    string // match filename (regular expression)
	skipString     string // match filename (regular expression)
	matchDirString string // match filename (regular expression)
	skipDirString  string // match filename (regular expression)
	noCSV          bool   // false: no Output CSV
	noDetail       bool   // false: no detail
	specCSVFile    string // csv filename specification
	toSJIS         bool   // true: SJIS false:UTF8(def)
}

// CommandUsage commandline usage
var CommandUsage = `
Usage : flist <options>

options:

`

// newParams get commandline params
func newParams() *Params {

	var v bool

	var p Params
	flag.StringVar(&p.matchString, "m", "", " File Match String(Regular expression)")
	flag.StringVar(&p.skipString, "s", "", " File Skip String(Regular expression)")
	flag.StringVar(&p.matchDirString, "md", "", " Directory Match String(Regular expression)")
	flag.StringVar(&p.skipDirString, "sd", "", " Directory Skip String(Regular expression)")
	flag.StringVar(&p.specCSVFile, "f", "", " Specify CSV filename")
	flag.BoolVar(&p.noCSV, "no", false, " CSV Not Output")
	flag.BoolVar(&p.noDetail, "nd", false, " Print list only")
	flag.BoolVar(&p.toSJIS, "sjis", false, " SJIS encoding")
	flag.BoolVar(&v, "version", false, " Show version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, CommandUsage)
		flag.PrintDefaults()
	}
	flag.Parse()

	if v {
		fmt.Println("flist-go version ", version)
		os.Exit(0)
	}
	return &p
}

func isMatch(matchString string, s string) bool {
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
			"basepath", "path", "filename", "mod-date", "mod-time", "size",
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
			fmt.Sprintf("%v", info.Size()),
		})
	}
	writer.Flush()
}

// openCSV Open CSV File
func openCSV(filename string, toSJIS bool) (*os.File, *csv.Writer, error) {

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
	file.Truncate(0)

	if toSJIS == true {
		writer := csv.NewWriter(transform.NewWriter(file, japanese.ShiftJIS.NewEncoder()))
		return file, writer, nil
	}

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
		file, writer, err = openCSV(p.specCSVFile, p.toSJIS)
		if err == nil {
			defer file.Close()
			writeCSV(writer, "", "", nil, true)
		}
	}

	// * Parse dir
	num := 0
	sumSize := int64(0)
	root, _ := filepath.Abs(filepath.Dir("."))

	filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil // skip directory
			}
			rel, _ := filepath.Rel(root, path)

			// match dir name
			if p.matchDirString != "" && !isMatch(p.matchDirString, filepath.Dir(rel)) {
				return nil
			}
			// skip dir
			if p.skipDirString != "" && isMatch(p.skipDirString, filepath.Dir(rel)) {
				return nil
			}

			// match lename
			if p.matchString != "" && !isMatch(p.matchString, info.Name()) {
				return nil
			}
			// skip filename
			if p.skipString != "" && isMatch(p.skipString, info.Name()) {
				return nil // skip pattern
			}

			fmt.Println(rel)
			writeCSV(writer, root, rel, info, false)
			num++
			sumSize = sumSize + info.Size()

			return nil
		})

	if !p.noDetail {
		fmt.Println("\n---------------------------------------------------------")
		fmt.Println("Exec Date  : ", time.Now())
		fmt.Println("Root       : ", root)
		fmt.Println("NumFiles   : ", num)
		fmt.Println("Total Size : ", sumSize)
		fmt.Println("")
	}
}
