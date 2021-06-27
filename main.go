package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

var (
	err error
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	//fmt.Println(os.Args[1])
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, filePath string, printFiles bool) error {
	var fileFolder []string

	home := func(path string, info os.FileInfo, err error) error {
		fileFolder = append(fileFolder, path)
		return nil
	}

	err := filepath.Walk(filePath, home)
	if err != nil {
		//fmt.Println("No")
		log.Fatal(err)
	}

	//fmt.Println(string(os.PathSeparator) + string(os.PathSeparator))

	sort.Strings(fileFolder)
	var sortDataF, sortData []string
	for k := range fileFolder {
		if fileFolder[k] == "hw1.md" || fileFolder[k] == "dockerfile" || fileFolder[k] == "./" || fileFolder[k] == filePath {
		} else {
			sortDataF = append(sortDataF, fileFolder[k])
		}
	}
	for i := range sortDataF {
		files, err := ioutil.ReadDir(filepath.Dir(sortDataF[i]))
		if err != nil {
			log.Fatal(err)
		}
		for k, file := range files {
			if file.IsDir() && getCompare(files[k].Name(), filepath.Base(sortDataF[i])) {
				sortData = append(sortData, sortDataF[i])
				//fmt.Println(file.Name())
			}
		}
	}
	//fmt.Println(sortData)
	var slice []string
	var result string
	if printFiles {
		for i := range sortDataF {
			if printFiles {
				result, slice = getTreeF(sortDataF[i], slice)
				fmt.Fprintln(out, result+getFileSize(sortDataF[i]))
				//fmt.Println(result + getFileSize(sortDataF[i]))
			}
		}
	} else {
		for i := range sortData {
			result, slice = getTree(sortData[i], slice)
			if result != "" {
				fmt.Fprintln(out, result)
				//fmt.Println(result)
			}
		}
	}
	return err
}

func getFileSize(path string) string {
	var fileSize string
	fileInfo, _ := os.Stat(path)
	if !fileInfo.IsDir() {
		size := fileInfo.Size()
		if size == 0 {
			fileSize = " (empty)"
		} else {
			fileSize = fmt.Sprintf(" (%vb)", size)
		}
	}
	return fileSize
}

func getTree(path string, sl []string) (string, []string) {
	var tabResult string
	files, err := ioutil.ReadDir(filepath.Dir(path))
	if err != nil {
		log.Fatal(err)
	}
	var blaklist []int
	for j, file := range files {
		if !file.IsDir() {
			blaklist = append(blaklist, j)
		}
	}

	for i, file := range files {
		if getResulutionI(blaklist, i) {
			if getCompare(files[i].Name(), filepath.Base(path)) {
				if file.IsDir() { //file.IsDir()
					if i == len(files)-1-len(blaklist) {
						tabResult = setSlash(sl, filepath.Dir(path)) + `└───` + filepath.Base(path)
						if getResulution(sl, filepath.Base(path)) {
							sl = append(sl, filepath.Base(path))
						}
						break
					} else {
						tabResult = setSlash(sl, filepath.Dir(path)) + `├───` + filepath.Base(path)
						break
					}
				}
			}
		} else {
			blaklist = blaklist[1:]
		}
	}
	return tabResult, sl
}

func getTreeF(path string, sl []string) (string, []string) {
	var tabResult string
	files, err := ioutil.ReadDir(filepath.Dir(path))
	if err != nil {
		log.Fatal(err)
	}
	for i, file := range files {
		if getCompare(files[i].Name(), filepath.Base(path)) {
			if file.IsDir() {
				if i == len(files)-1 {
					tabResult = setSlash(sl, filepath.Dir(path)) + `└───` + filepath.Base(path)
					if getResulution(sl, filepath.Base(path)) {
						sl = append(sl, filepath.Base(path))
					}
					break
				} else {
					tabResult = setSlash(sl, filepath.Dir(path)) + `├───` + filepath.Base(path)
					break
				}
			} else {
				if i == len(files)-1 {
					tabResult = setSlash(sl, filepath.Dir(path)) + `└───` + filepath.Base(path)
					break
				} else {
					tabResult = setSlash(sl, filepath.Dir(path)) + `├───` + filepath.Base(path)
					break
				}
			}
		}
	}
	return tabResult, sl
}

func getSlash(text string) int {
	re := regexp.MustCompile(string(os.PathSeparator))
	return len(re.FindAll([]byte(text), -1))
}

func getCompare(text1, text2 string) bool {
	re := regexp.MustCompile(string(os.PathSeparator))
	var buf1 []string
	var buf2 []string
	buf1 = re.Split(text1, -1)
	buf2 = append(buf2, text2)
	if buf1[len(buf1)-1] == buf2[0] {
		return true
	}
	return false
}

func setSlash(slice []string, path string) string {
	var i int = getSlash(path)
	re := regexp.MustCompile(string(os.PathSeparator))
	var buf1 []string
	buf1 = re.Split(path, -1)
	buf1 = buf1[1:]
	var otab string = ""
	//var k int
	var b bool = false
	//var b2 bool = false
	for j := 0; j < i; j++ {
		b = false

		for k := 0; k < len(slice); k++ {
			if buf1[j] == slice[k] {
				var mtab string = "\t"
				otab = otab + mtab
				//slice = slice[1:]
				b = true
				//break
			}
		}
		if b == false {
			var mtab string = "│\t"
			otab = otab + mtab
		}
		//var mtab string = "│\t"
		//otab = otab + mtab
	}
	return otab
}

func getResulutionI(sl []int, dbase int) bool {
	for i := 0; i < len(sl); i++ {
		if sl[i] == dbase {
			return false
		}
	}
	return true
}

func getResulution(sl []string, dbase string) bool {
	for i := 0; i < len(sl); i++ {
		if sl[i] == dbase {
			return false
		}
	}
	return true
}
