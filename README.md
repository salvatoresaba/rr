# Recursive Rename tool
The purpose of this command line program is to change the name of the given files using regular expressions.
The syntax of the regular expressions accepted is the same general syntax used by Perl, Python, and other languages. More precisely, it is the syntax accepted by RE2 and described at [https://golang.org/s/re2syntax](https://golang.org/s/re2syntax), except for \C. For an overview of the syntax, see the [regexp/syntax](https://pkg.go.dev/regexp/syntax) package.

Version 1.0 
Usage: rr [options] <path> <match_rule> <replace_rule> 
**-r** = Search recursively in directories 
**-d** = Include directory names 
**-e** = Exclude file extension from the text matching 
**-f** = Force replacement if file already exists 
**-q** = Perform action without confirmation 
**-h** = Print this help 

Example: 
- `rr -r ./test "^t" "r"` 
Search recursively for all files in the test folder and rename those starting with "t" to "r" 
- `rr ./test '(\d+)' '$1$1'` 
Search for all files in the test folder and if they contain a number it is written twice ('test1.txt' -> 'test11.txt') 
- `rr -r -f ./test '\.JPG$' '.jpg'` 
Edit  file extension overwriting existing files with the same name