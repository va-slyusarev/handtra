package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type ExtCommander interface {
	SetExt(extOpts ExtensionOpts)
	Execute(args []string) error
}

type ExtensionOpts struct {
	AppRevision string
	InputFile   string
	OutputFile  string
	Expr        string
}

func (o *ExtensionOpts) SetExt(extOpts ExtensionOpts) {
	o.AppRevision = extOpts.AppRevision
	o.InputFile = extOpts.InputFile
	o.OutputFile = extOpts.OutputFile
	o.Expr = extOpts.Expr
}

func getFileLines(fileName string) ([]string, error) {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла %q: %v", fileName, err)
	}
	lines := make([]string, 0)
	scanner := bufio.NewScanner(bytes.NewBuffer(f))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обработки файла %q: %v", fileName, err)
	}
	return lines, nil
}

func prepareInput(fileName string, expr string) ([][]string, int, error) {
	lines, err := getFileLines(fileName)
	if err != nil {
		return nil, -1, fmt.Errorf("файл для перевода: %v", err)
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка обработки регулярного выражения %q: %v", expr, err)
	}

	i := 0
	result := make([][]string, 0)
	for _, line := range lines {
		r := re.FindStringSubmatch(line)
		el := make([]string, 2)
		el[0] = line
		if len(r) == 1 {
			el[1] = r[0]
		}
		if len(r) == 2 {
			el[1] = r[1]
		}
		if len(el[1]) > 0 {
			i++
		}
		result = append(result, el)
	}
	return result, i, nil
}

func process(inFile string, outFile string, trFile string, expr string, prepare bool) error {
	inLines, count, err := prepareInput(inFile, expr)
	if err != nil {
		return err
	}
	trLines := make([]string, 0)
	if !prepare {
		if trLines, err = getFileLines(trFile); err != nil {
			return fmt.Errorf("файл со списком фраз: %v", err)
		}
		if count != len(trLines) {
			return fmt.Errorf("количество строк перевода не соответствует ожидаемым: передано %d шт., ожидаем %d шт.", len(trLines), count)
		}
	}

	i := 0
	var buf bytes.Buffer
	for _, line := range inLines {

		// формируем строки для перевода
		if prepare {
			if line[1] != "" {
				buf.WriteString(fmt.Sprintf("%s\n", line[1]))
				i++
			}
			continue
		}

		if line[1] == "" {
			buf.WriteString(fmt.Sprintf("%s\n", line[0]))
			continue
		}

		buf.WriteString(fmt.Sprintf("%s\n", strings.ReplaceAll(line[0], line[1], trLines[i])))
		i++
	}

	out, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("ошибка создания исходящего файла %q: %v", outFile, err)
	}
	defer func() {
		_ = out.Close()
	}()
	if _, err := buf.WriteTo(out); err != nil {
		return fmt.Errorf("ошибка записи в исходящий файл %q: %v", outFile, err)
	}

	if prepare {
		fmt.Printf("Создан файл %q со списком фраз для перевода в количестве - %d шт.\n", outFile, count)
		return nil
	}
	fmt.Printf("Файл %q переведён в файл %q, переведено фраз - %d шт.\n", inFile, outFile, count)
	return nil
}
