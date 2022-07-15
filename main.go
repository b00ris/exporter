package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	origin := flag.String("json_path", "", "path to json format")
	outputDir := flag.String("output_path", "", "path to output dir")
	flag.Parse()

	if *origin == "" || *outputDir == "" {
		fmt.Println("Empty param json_path or output_path")
		return
	}

	originFileName := strings.TrimSuffix(filepath.Base(*origin), ".json")

	f, err := os.Open(*origin)
	if err != nil {
		fmt.Println(err)
		return
	}

	res := struct {
		ABI      []interface{} `json:"abi"`
		Bytecode string        `json:"bytecode"`
		UserDoc  interface{}   `json:"userdoc"`
		DevDoc   interface{}   `json:"devdoc"`
	}{}

	err = json.NewDecoder(f).Decode(&res)
	if err != nil {
		fmt.Println(err)
		return
	}

	abiFilePath := filepath.Join(*outputDir, originFileName+".abi")

	err = encode(abiFilePath, res.ABI)
	if err != nil {
		fmt.Println(err)
		return
	}

	bytecodeFilePath := filepath.Join(*outputDir, originFileName+".bin")

	err = writeString(bytecodeFilePath, res.Bytecode)
	if err != nil {
		fmt.Println(err)
		return
	}

	devDocFilePath := filepath.Join(*outputDir, originFileName+".dev")

	err = encode(devDocFilePath, res.DevDoc)
	if err != nil {
		fmt.Println(err)
		return
	}

	userDocFilePath := filepath.Join(*outputDir, originFileName+".user")

	err = encode(userDocFilePath, res.UserDoc)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Successful written abi: %s, bytecode: %s, userDoc: %s, devDoc: %s \n", abiFilePath, bytecodeFilePath, userDocFilePath, devDocFilePath)
}

func encode(filepath string, content any) error {
	abiFile, err := os.Create(filepath)
	if err != nil {
		return err
	}

	return json.NewEncoder(abiFile).Encode(content)
}

func writeString(filepath string, content string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	_, err = io.WriteString(file, content)
	if err != nil {
		return err
	}

	return nil
}
