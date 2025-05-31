package notebookparser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Cell represents a Jupyter notebook cell
type Cell struct {
	CellType string   `json:"cell_type"`
	Source   []string `json:"source"`
}

// Notebook represents a Jupyter notebook structure
type Notebook struct {
	Cells []Cell `json:"cells"`
}

// ParseNotebook reads and parses a Jupyter notebook file, extracting text content
func ParseNotebook(filePath string) (string, error) {
	// Read the notebook file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read notebook file: %w", err)
	}

	// Parse JSON
	var notebook Notebook
	err = json.Unmarshal(data, &notebook)
	if err != nil {
		return "", fmt.Errorf("failed to parse notebook JSON: %w", err)
	}

	var content strings.Builder
	content.WriteString("# Jupyter Notebook Content\n\n")

	// Process each cell
	for i, cell := range notebook.Cells {
		// Only process code, markdown, and raw cells
		if cell.CellType == "code" || cell.CellType == "markdown" || cell.CellType == "raw" {
			content.WriteString(fmt.Sprintf("## Cell %d (%s)\n\n", i+1, cell.CellType))

			// Join source lines
			if len(cell.Source) > 0 {
				cellContent := strings.Join(cell.Source, "")
				content.WriteString(cellContent)

				// Ensure proper spacing between cells
				if !strings.HasSuffix(cellContent, "\n") {
					content.WriteString("\n")
				}
				content.WriteString("\n")
			}
		}
	}

	return content.String(), nil
}
