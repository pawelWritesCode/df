// Package df holds stuff related to data format.
package df

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/goccy/go-yaml"
)

// Pre-defined well known data formats.
const (
	// JSON describes JSON data format.
	JSON DataFormat = "json"

	// YAML describes Yaml data format.
	YAML DataFormat = "yaml"

	// XML describes XML data format.
	XML DataFormat = "xml"

	// HTML describes HTML data format.
	HTML DataFormat = "html"

	// PlainText describes plan text data format.
	PlainText DataFormat = "plain text"
)

// DataFormat describes format of data.
type DataFormat string

type input interface {
	~string | ~[]byte
}

// IsJSON checks whether bytes are in JSON format.
func IsJSON[i input](b i) bool {
	var js json.RawMessage
	err := json.Unmarshal([]byte(b), &js)

	return err == nil
}

// IsYAML checks whether bytes are in YAML format.
func IsYAML[i input](b i) bool {
	if IsJSON(b) {
		return false
	}

	if IsXML(b) {
		return false
	}

	// yaml.UnmarshalWithOptions parses successfully any plain text,
	// to detect text that is not in yaml format, we assume, there must be,
	// at least one key: value pair in yaml
	if !bytes.Contains([]byte(b), []byte(":")) {
		return false
	}

	var y any
	return yaml.UnmarshalWithOptions([]byte(b), &y, yaml.Strict()) == nil
}

// IsXML checks whether bytes are in XML format.
// Function does not guarantee that standard xml.Unmarshal will work, instead
// it only looks for characteristics of XML formatted data.
func IsXML[i input](b i) bool {
	str := string(b)
	idx := strings.Index(strings.TrimSpace(str), "<?xml version=")
	if idx == 0 || idx == 1 {
		return true
	}

	if !(strings.Contains(str, ">") && strings.Contains(str, "<")) {
		return false
	}

	if strings.Count(str, "<") >= (strings.Count(str, "</") + strings.Count(str, "/>")) {
		return true
	}

	return false
}

// IsHTML checks whether bytes are in HTML format.
func IsHTML[i input](b i) bool {
	var points, confidenceLevel = 0, 3
	shouldContain := []string{"<!doctype html>", "</head>", "</html>", "</body>", "</title>", "</a>", "</div>"}

	for _, s := range shouldContain {
		if strings.Contains(strings.ToLower(string(b)), s) {
			points++
		}
	}

	return points >= confidenceLevel
}

// IsPlainText checks whether bytes are in plain text format.
func IsPlainText[i input](b i) bool {
	if IsJSON(b) {
		return false
	}

	return len(b) > 0
}
