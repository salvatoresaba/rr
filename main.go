package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	var (
		recursive   = flag.Bool("r", false, "Search recursively in directories")
		includeDirs = flag.Bool("d", false, "Include directory names")
		excludeExt  = flag.Bool("e", false, "Exclude file extension from the text matching")
		force       = flag.Bool("f", false, "Force replacement if file already exists")
		noConfirm   = flag.Bool("q", false, "Perform action without confirmation")
		help        = flag.Bool("h", false, "Print this help ")
	)

	flag.Parse()

	if *help || flag.NArg() < 3 {
		printHelp()
		return
	}

	rootDir := flag.Arg(0)
	matchPattern := flag.Arg(1)
	replacePattern := flag.Arg(2)

	regex, err := regexp.Compile(matchPattern)
	if err != nil {
		fmt.Printf("Regular expression compilation error: %v\n", err)
		return
	}

	//create the target files list
	var filesToRename []string

	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Esclude the root directory name
		if path == rootDir {
			return nil
		}

		if info.IsDir() { //check the directory
			if *includeDirs { //include or not the directory name
				//check the directory name
				dirName := filepath.Base(path)
				//add the file name to the list if it maches the given rule
				if regex.MatchString(dirName) {
					filesToRename = append(filesToRename, path)
				}
			}
			if !*recursive {
				return filepath.SkipDir
			}
			return nil
		} else { //check the files
			//check all the file name or the name without the extension
			oldName := filepath.Base(path)
			ext := filepath.Ext(oldName)
			nameToMatch := oldName
			if *excludeExt && !info.IsDir() {
				nameToMatch = strings.TrimSuffix(oldName, ext)
			}
			//add the file name to the list if it matches the given rule
			if regex.MatchString(nameToMatch) {
				filesToRename = append(filesToRename, path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		return
	}

	if len(filesToRename) == 0 {
		fmt.Println("No files found for renaming.")
		return
	}

	var newPaths []string
	fmt.Printf("Found %d files to rename:\n", len(filesToRename))
	for _, path := range filesToRename {
		relativePath, _ := filepath.Rel(rootDir, path)
		oldName := filepath.Base(relativePath)
		ext := filepath.Ext(oldName)

		nameToReplace := oldName
		if *excludeExt && (filepath.Ext(path) != "") {
			nameToReplace = strings.TrimSuffix(oldName, ext)
		}

		newName := regex.ReplaceAllString(nameToReplace, replacePattern)
		if *excludeExt && (filepath.Ext(path) != "") {
			newName += ext
		}
		newPath := filepath.Join(filepath.Dir(path), newName)

		newPaths = append(newPaths, newPath) //create the list of the new names
		newRelativePath, _ := filepath.Rel(rootDir, newPath)
		fmt.Printf("'%s' -> '%s'\n", relativePath, newRelativePath)
	}

	if !*noConfirm {
		fmt.Printf("Procede to rename %d files? (y/n):\n ", len(filesToRename))
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))

		if response != "y" {
			fmt.Println("Operation cancelled")
			return
		}
	}

	n_files := len(filesToRename)
	for i := 0; i < n_files; i++ {
		rename := true
		path := filesToRename[i]
		newPath := newPaths[i]
		if !*force {
			if _, err := os.Stat(newPath); err == nil {
				for {
					fmt.Printf("The file '%s' already exist. Do you wanto to replace it? (Yes/No/All): ", newPath)
					reader := bufio.NewReader(os.Stdin)
					response, _ := reader.ReadString('\n')
					response = strings.ToLower(strings.TrimSpace(response))

					switch response {
					case "y":
						break
					case "a":
						*force = true
						break
					case "n":
						rename = false
						fmt.Println("Operation cancelled")
						break
					default:
						fmt.Println("Answer not valid")
					}
				}
			}
		}
		if rename {
			err := os.Rename(path, newPath)
			if err != nil {
				fmt.Printf("Error renaming '%s' -> '%s': %v\n", path, newPath, err)
				return
			}
		}
	}

	fmt.Println("Renaming completed successfully")
	fmt.Print("Press enter to close the program...")
	if !*noConfirm {
		bufio.NewReader(os.Stdin).ReadString('\n')
	}
}

func printHelp() {
	fmt.Print(`
Version 1.0
Usage: rr [options] <match_rule> <replace_rule>
-r = Search recursively in directories
-d = Include directory names
-e = Exclude file extension from the text matching
-f = Force replacement if file already exists
-q = Perform action without confirmation
-h = Print this help

Example:

rr -r ./test "^t" "r"
Search recursively for all files in the test folder and rename those starting with “t” to “r”

rr ./test '(\d+)' '$1$1'
Search for all files in the test folder and if they contain a number it is written twice (‘test1.txt’ -> ‘test11.txt’)

rr -r -f ./test '\.JPG$' '.jpg'
Edit file extension overwriting existing files with the same name

`)
}
