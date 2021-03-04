package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func countDirs(files []os.FileInfo) int {
	var i int
	for _, file := range files {
		if file.IsDir() {
			i++
		}
	}
	return i
}

func putTabs(out io.Writer, line bool, path string, printFiles bool) error {
	splittedPath := strings.Split(path, "/")

	if len(splittedPath) == 1 {
		return nil
	}
	for idx := range splittedPath {
		if idx != 0 {
			splittedPath[idx] = splittedPath[idx-1] + "/" + splittedPath[idx]
		}
	}

Exit:
	for idx, dir := range splittedPath {
		files, _ := ioutil.ReadDir(dir)
		dirCount := countDirs(files)
		for pos, file := range files {

			if (printFiles == false && dirCount == pos+1) || pos+1 == len(files) {
				if len(splittedPath) > idx+1 {
					pathName := strings.Split(splittedPath[idx+1], "/")
					if pathName[len(pathName)-1] == file.Name() {
						fmt.Fprint(out, "\t")
						continue Exit
					}
				}
			}
		}
		if idx+1 == len(splittedPath) {
			return nil
		}
		fmt.Fprint(out, "│\t")

	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}

	iteration := 0
	for pos, file := range files {
		if file.IsDir() == false && printFiles == false {
			continue
		}
		putTabs(out, true, path, printFiles)

		dirCount := countDirs(files)
		if (printFiles == false && dirCount == iteration+1) || pos+1 == len(files) {
			fmt.Fprint(out, "└───")
		} else {
			fmt.Fprint(out, "├───")
		}
		fmt.Fprint(out, file.Name())
		if file.IsDir() == false {
			if file.Size() != 0 {
				fmt.Fprint(out, " (", file.Size(), "b)")
			} else {
				fmt.Fprint(out, " (empty)")
			}
		}
		fmt.Fprintln(out)
		dirTree(out, path+"/"+file.Name(), printFiles)
		iteration++
	}
	return nil
}

func checkArgs() error {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		return fmt.Errorf("Invalid num of args")
	} else if len(os.Args) == 3 && os.Args[2] != "-f" {
		return fmt.Errorf("Invalide usage of -f")
	}
	return nil
}

func main() {

	if err := checkArgs(); err != nil {
		panic(err.Error())
	}

	out := os.Stdout
	path := os.Args[1]
	var printFiles bool = len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}

}
