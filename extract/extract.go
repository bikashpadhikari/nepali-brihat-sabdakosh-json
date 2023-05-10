// extract.go extracts word definitions from the JSON export of the realm
// database, decrypts the defintions, and writes the resulting JSON to stdout.
//
// Can be used like so:
//   ./extract -input db.json > raw.json
package main

import (
	"encoding/json"
	"fmt"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
)

var (
	input = flag.String("input", "", "Path to JSON exported from Realm Studio")
)

var opensslArgs = []string{
	"enc", "-aes-256-cbc",
	"-d", "-a", "-A",
	"-md", "md5",
	"-pass", "pass:058aa5325d7d2e7",
}

func decrypt(ciphertext string) (string, error) {
	cmd := exec.Command("openssl", opensslArgs...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, ciphertext)
	}()
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

type rawDefinition struct {
	Word string `json:"word"`
	Meaning string `json:"meaning"`
}

func parseIndex(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalf("parsing index: %v", err)
		return 0
	}
	return i
}

func extractWords(data []byte) ([]rawDefinition, error) {
	elements := []interface{}{}
	if err := json.Unmarshal(data, &elements); err != nil {
		return nil, err
	}
	indices, ok := elements[1].([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected list of indices, got %T", elements[1])
	}
	definitions := []rawDefinition{}
	for idx, d := range indices {
		indexStr := d.(string)
		i := parseIndex(indexStr)
		definition := elements[i].(map[string]interface{})
		word := elements[parseIndex(definition["word"].(string))].(string)
		meaning := elements[parseIndex(definition["meaning"].(string))].(string)
		decryptedMeaning, err := decrypt(meaning)
		if err != nil {
			return nil, err
		}
		d := rawDefinition{
			Word: word,
			Meaning: decryptedMeaning,
		}
		definitions = append(definitions, d)
		if num := idx + 1; (num % 1_000 == 0) {
			fmt.Fprintf(os.Stderr, "Decrypted %v/%v definitions\n", num, len(indices))
		}
	}
	return definitions, nil
}

func main() {
	flag.Parse()
	data, err := os.ReadFile(*input)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	rawDefinitions, err := extractWords(data)
	if err != nil {
		log.Fatalf("Error extracting words: %v", err)
	}
	marshalled, err := json.Marshal(rawDefinitions)
	if err != nil {
		log.Fatalf("Error marshalling extracted data: %v", err)
	}
	os.Stdout.Write(marshalled)
	fmt.Fprintln(os.Stderr, "Done")
}
