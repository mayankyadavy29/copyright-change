package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func UpdateCopyright(filename string) bool {
	ext := filepath.Ext(filename)
	if ext != ".go" {
		return false
	}
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return false
	}
	var firstLine string
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		firstLine = scanner.Text()
		newLine, err := checkCopyright(firstLine)
		if err != nil {
			//fmt.Println(err)
			return false
		}
		err = file.Truncate(0)
		if err != nil {
			//fmt.Println(err)
			return false
		}
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			//fmt.Println(err)
			return false
		}
		_, err = file.WriteString(newLine + "\n")
		if err != nil {
			//fmt.Println(err)
			return false
		}
		for scanner.Scan() {
			_, err = file.WriteString(scanner.Text() + "\n")
			if err != nil {
				//fmt.Println(err)
				return false
			}
		}
	}
	if err := scanner.Err(); err != nil {
		//fmt.Println(err)
		return false
	}
	return true
}

/*
	checkCopyright checks for copyright conditions

Possible conditions:
1. No copyright statement. Add a new copyright statement.
2. Copyright contains current year. No change required.
3. Copyright contains only old year. Updated it to oldYear-curYear.
4. Copyright contains oldYear-curYear. No change required.
5. Copyright contains old1Year-old2Year. Update it to oldYear-curYear.
*/
func checkCopyright(line string) (string, error) {
	copyrightRegex := `// Copyright (\d{4})(-(\d{4}))? Dell Inc. or its subsidiaries\. All Rights Reserved\.`
	re := regexp.MustCompile(copyrightRegex)
	matches := re.FindStringSubmatch(line)

	if len(matches) == 0 {
		copyrightStatement := fmt.Sprintf("// Copyright %d Dell Inc. or its subsidiaries. All Rights Reserved.\n\n%s", CUR_YEAR, line)
		return copyrightStatement, nil
	}

	oldYear, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", err
	}
	if oldYear == CUR_YEAR {
		return line, fmt.Errorf("no change required")
	}
	newYear, err := strconv.Atoi(matches[3])
	if err != nil {
		return fmt.Sprintf("// Copyright %d-%d Dell Inc. or its subsidiaries. All Rights Reserved.", oldYear, CUR_YEAR), nil
	}
	if newYear == CUR_YEAR {
		return line, fmt.Errorf("no change required")
	}
	return fmt.Sprintf("// Copyright %d-%d Dell Inc. or its subsidiaries. All Rights Reserved.", oldYear, CUR_YEAR), nil
}
