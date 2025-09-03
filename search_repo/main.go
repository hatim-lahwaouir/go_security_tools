package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

type Commit struct {
	Msg    string
	Author string
	Hash   string
	Branche string
	FilesPath []string
}

type Branche struct {
	Name string
}

// commits that we need to search for

// files that we should ignore 
var ignored_files []string = []string{".git"}

func IgnoreFile(path string) bool {

	for _, fn := range(ignored_files){
		if strings.Contains(path, fn) {
			return true
		}
	}
	return false
}


func GetBranchesName() []string {

	var (
		cmd *exec.Cmd
		out bytes.Buffer
		brs []string
	)

	cmd = exec.Command("git", "branch", "--format=%(refname:short)")
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	for {
		line, err := out.ReadBytes('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		line = line[:len(line)-1]

		brs = append(brs, string(line))
	}

	fmt.Println("----branches currenty searching ----")
	for _, br := range brs {
		fmt.Printf("-- %s\n", br)
	}
	return brs
}

func GetCommits(br string, commits map[string]Commit) {
	var (
		cmd  *exec.Cmd
		out  bytes.Buffer
	)
	// switch to the branche
	cmd = exec.Command("git", "switch", br)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command("git", "log", `--pretty=format:%h|%an|%s`)
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	for {
		line, err := out.ReadBytes('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		line = line[:len(line)-1]
		arr := strings.SplitN(string(line), "|", 3)
		commits[arr[0]] = Commit{Hash: arr[0], Author: arr[1], Msg: arr[2], Branche : br}
	}
}



func Reset(){
	cmd := exec.Command("git", "switch", "main")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}


func ListAllFiles(cmt Commit) {
	var (
		cmd *exec.Cmd
		filesPath []string
	)

	cmd = exec.Command("git", "checkout", cmt.Hash)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("-- files in this commit %s --\n", cmt.Hash)
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {

		if !info.IsDir() && IgnoreFile(path) == false{
			fmt.Printf("file to search for : %q\n", path)
			filesPath = append(filesPath, path)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	cmt.FilesPath = filesPath
	Reset()
}


func Search(filePath string, pattern string) {


}

func main() {
	commits := make(map[string]Commit)
	brs := GetBranchesName()
	for _, br := range brs {
		GetCommits(br, commits)
	}

	for _, cmt := range commits {
		ListAllFiles(cmt)
	}

	Reset()
}
