package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

// Format represents the output format
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

// Print outputs data in the specified format
func Print(data interface{}, format Format, headers []string, rowFunc func(interface{}) []string) error {
	switch format {
	case FormatJSON:
		return printJSON(data)
	case FormatYAML:
		return printYAML(data)
	case FormatTable:
		return printTable(data, headers, rowFunc)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// printJSON outputs data as JSON
func printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// printYAML outputs data as YAML
func printYAML(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	defer encoder.Close()
	return encoder.Encode(data)
}

// printTable outputs data as a table
func printTable(data interface{}, headers []string, rowFunc func(interface{}) []string) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	// Handle slice types
	switch v := data.(type) {
	case []interface{}:
		for _, item := range v {
			table.Append(rowFunc(item))
		}
	default:
		// Try to convert to slice
		if slice, ok := toSlice(data); ok {
			for _, item := range slice {
				table.Append(rowFunc(item))
			}
		} else {
			// Single item
			table.Append(rowFunc(data))
		}
	}

	table.Render()
	return nil
}

// toSlice attempts to convert data to []interface{}
func toSlice(data interface{}) ([]interface{}, bool) {
	switch v := data.(type) {
	case []interface{}:
		return v, true
	default:
		// Use reflection if needed in the future
		return nil, false
	}
}

// PrintSuccess prints a success message
func PrintSuccess(message string) {
	fmt.Printf("✓ %s\n", message)
}

// PrintError prints an error message
func PrintError(message string) {
	fmt.Fprintf(os.Stderr, "✗ %s\n", message)
}

// PrintInfo prints an info message
func PrintInfo(message string) {
	fmt.Printf("ℹ %s\n", message)
}
