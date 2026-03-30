package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type TreeEntry struct {
	os.DirEntry
}

func (entry TreeEntry) Size() int64 {
	info, err := entry.Info()
	if err != nil {
		log.Fatalf("Failed to get info for %s\n", entry.Name())
	}

	return info.Size()
}

func readDir(path string, readFiles bool) ([]TreeEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var dirs []TreeEntry
	for _, e := range entries {
		if readFiles || e.IsDir() {
			dirs = append(dirs, TreeEntry{e})
		}
	}

	return dirs, nil
}

func dfsDir(out io.Writer, path string, depth int, lastInTree []bool, printFiles bool) error {
	entries, err := readDir(path, printFiles)
	if err != nil {
		log.Fatalf("Failed to read directory: %s\n", path)
		return err
	}

	lastInTree = append(lastInTree, false)

	for i, e := range entries {
		for i := 0; i < depth; i++ {
			if lastInTree[i] != true {
				fmt.Fprint(out, "│\t")
			} else {
				fmt.Fprint(out, "\t")
			}
		}

		if i != len(entries)-1 {
			fmt.Fprint(out, "├───")
		} else {
			lastInTree[depth] = true
			fmt.Fprint(out, "└───")
		}

		if e.IsDir() {
			fmt.Fprintf(out, "%s\n", e.Name())
			err = dfsDir(out, filepath.Join(path, e.Name()), depth+1, lastInTree, printFiles)
			if err != nil {
				return err
			}
		} else {
			var size string
			if e.Size() != 0 {
				size = fmt.Sprintf("%db", e.Size())
			} else {
				size = "empty"
			}
			fmt.Fprintf(out, "%s (%s)\n", e.Name(), size)
		}
	}

	lastInTree = lastInTree[:len(lastInTree)-1]

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return dfsDir(out, path, 0, []bool{}, printFiles)
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

	fmt.Print(out)
}
