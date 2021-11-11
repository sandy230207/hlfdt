package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func WriteYaml(conf interface{}, filename string) error {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	yamlEncoder.Encode(&conf)

	if err := ioutil.WriteFile(filename, b.Bytes(), 0666); err != nil {
		return fmt.Errorf("error occur in writting %v: %v", filename, err)
	}
	return nil
}

func WriteSh(dataString string, filename string) error {
	data := []byte(dataString)
	if err := ioutil.WriteFile(filename, data, 0666); err != nil {
		return fmt.Errorf("error occur in writting %v: %v", filename, err)
	}
	return nil
}

func ExtractHost(endpoint string, index int) (string, error) {
	delimiter := func(r rune) bool {
		return r == '.' || r == ':'
	}
	orghost := ""
	endpointDelimited := strings.FieldsFunc(endpoint, delimiter)
	for i, v := range endpointDelimited {
		if i == len(endpointDelimited)-index {
			orghost += v
			break
		} else if i > 0 {
			orghost += v + "."
		}
	}
	return orghost, nil
}

func ConvertNet(filename string, start string, end string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	isSection := false
	content := ""
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		row := scanner.Text()
		if row == start && !isSection {
			isSection = true
		} else if row == end && isSection {
			isSection = false
		} else if isSection {
			target := strings.Split(row, ":")
			row = target[0] + ":"
		}
		content = content + row + "\n"
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err = ioutil.WriteFile(filename, []byte(content), 0666); err != nil {
		return err
	}
	return nil
}

func ConvertConfigtx(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	isSection1, isSection2, isSection3 := false, false, false
	content := ""
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		row := scanner.Text()
		rowWithoutIndent := strings.TrimSpace(row)
		if rowWithoutIndent == "Type: ImplicitMeta" && !isSection1 {
			isSection1 = true
		} else if isSection1 {
			target1 := strings.Split(row, ":")
			target2 := strings.Split(rowWithoutIndent, " ")
			row = target1[0] + ": \"" + target2[1] + " " + target2[2] + "\""
			isSection1 = false
		}
		if rowWithoutIndent == "Type: Signature" && !isSection2 {
			isSection2 = true
		} else if isSection2 {
			target1 := strings.Split(row, ":")
			target2 := strings.Split(rowWithoutIndent, " ")
			row = target1[0] + ": \""
			for i := 1; i < len(target2)-1; i++ {
				row = row + target2[i] + " "
			}
			row = row + target2[len(target2)-1] + "\""
			isSection2 = false
		}

		if row == "  Organizations:" && !isSection3 {
			isSection3 = true
		} else if row == "  Policies:" && isSection3 {
			isSection3 = false
		} else if isSection3 {
			row = ""
		}
		content = content + row + "\n"
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err = ioutil.WriteFile(filename, []byte(content), 0666); err != nil {
		return err
	}
	return nil
}
