package main

import (
	"bytes"
	"io"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment)

	for _, file := range files {
		// имя переменной не должно содержать символа "="
		if strings.Contains(file.Name(), "=") {
			continue
		}

		// если файл полностью пустой, помечаю переменную на удаление
		stat, err := file.Info()
		if err != nil {
			return nil, err
		}
		if stat.Size() == 0 {
			env[file.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		value, err := readEnvFromFile(dir + "/" + file.Name())
		if err != nil {
			return nil, err
		}
		env[file.Name()] = value
	}

	return env, nil
}

func readEnvFromFile(name string) (EnvValue, error) {
	f, err := os.Open(name)
	if err != nil {
		return EnvValue{}, err
	}
	defer f.Close()
	buf := make([]byte, 30, 120)
	offset := 0

	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			return readEnvFromBytes(buf[:offset])
		}
		if err != nil {
			return EnvValue{}, err
		}

		if i := bytes.IndexByte(buf[:offset+n], '\n'); i >= 0 {
			return readEnvFromBytes(buf[:offset+i])
		}
		offset += n
	}
}

func readEnvFromBytes(data []byte) (EnvValue, error) {
	// пробелы и символы табуляции в конце значения переменной должны быть удалены
	data = bytes.TrimRight(data, " \t")

	// терминальные нули должны быть трансформированы в символ перевода строки
	data = bytes.ReplaceAll(data, []byte{0}, []byte{'\n'})

	return EnvValue{Value: string(data)}, nil
}
