package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type file struct {
	name string
	size int64
}

type directory struct {
	name  string
	units []unit
}

type unit interface {
	String() string
}

func (dir directory) String() string {
	return dir.name
}

func (file file) String() string {
	if file.size == 0 {
		return file.name + " (empty)"
	}
	return file.name + " (" + strconv.FormatInt(file.size, 10) + "b)"
}

func readUnits(path string, units []unit, onlyDir bool) ([]unit, error) {

	rootDir, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	contents, err := rootDir.Readdir(0)
	sort.SliceStable(contents, func(i, j int) bool { return contents[i].Name() < contents[j].Name() })

	for _, info := range contents {

		if !info.IsDir() && onlyDir {
			continue
		}

		var newUnit unit
		if info.IsDir() {
			units, err := readUnits(filepath.Join(path, info.Name()), []unit{}, onlyDir)
			if err != nil {
				panic(err.Error())
			}
			newUnit = directory{info.Name(), units}
		} else {
			newUnit = file{info.Name(), info.Size()}
		}
		units = append(units, newUnit)
	}

	return units, err
}

func writeUnits(output io.Writer, units []unit, prefixes []string) {

	if len(units) == 0 {
		return
	}
	var elementSign string
	var newPrefix string
	for n, unit := range units {
		fmt.Fprintf(output, "%s", strings.Join(prefixes, ""))
		if n != len(units)-1 {
			elementSign = "├───"
			newPrefix = "│\t"

		} else {
			elementSign = "└───"
			newPrefix = "\t"
		}

		fmt.Fprintf(output, "%s%s\n", elementSign, unit)
		if p, ok := unit.(directory); ok {
			writeUnits(output, p.units, append(prefixes, newPrefix))
		}

	}
	// unit := units[0]

}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(output io.Writer, dir string, printFiles bool) error {
	// var b bytes.Buffer
	// var file os.File

	units, err := readUnits(dir, []unit{}, !printFiles)
	writeUnits(output, units, []string{})

	// _ = units

	// fmt.Println(files.Name())

	// b.Write([]byte("Hello "))
	// fmt.Fprintf(&b, "world!")
	// b.WriteTo(os.Stdout)
	return err
}
