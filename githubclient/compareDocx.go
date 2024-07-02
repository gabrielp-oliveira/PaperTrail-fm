package githubclient

import (
	"fmt"
	"strings"

	"baliance.com/gooxml/document"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// readDocxText lê o conteúdo de texto de um arquivo DOCX e retorna o texto junto com a página
func readDocxText(filePath string) ([]string, []int, error) {
	doc, err := document.Open(filePath)
	if err != nil {
		return nil, nil, err
	}

	var content []string
	var pages []int
	pageNum := 1

	for _, para := range doc.Paragraphs() {
		var paraText strings.Builder
		for _, run := range para.Runs() {
			paraText.WriteString(run.Text())
		}
		content = append(content, paraText.String())
		pages = append(pages, pageNum)

		// Supondo que cada parágrafo seja em uma página distinta para fins de exemplo
		pageNum++
	}
	return content, pages, nil
}

// highlightDiffs destaca as diferenças dentro de uma linha
func highlightDiffs(line1, line2 string) (string, string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(line1, line2, false)

	var highlightedLine1, highlightedLine2 strings.Builder

	for _, diff := range diffs {
		text := diff.Text
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			highlightedLine2.WriteString(fmt.Sprintf("<span class=\"diff-insert\">%s</span>", text))
		case diffmatchpatch.DiffDelete:
			highlightedLine1.WriteString(fmt.Sprintf("<span class=\"diff-delete\">%s</span>", text))
		case diffmatchpatch.DiffEqual:
			highlightedLine1.WriteString(text)
			highlightedLine2.WriteString(text)
		}
	}

	return highlightedLine1.String(), highlightedLine2.String()
}

// compareTexts compara o texto de dois documentos e retorna as diferenças como uma lista de strings
func compareTexts(lines1 []string, pages1 []int, lines2 []string, pages2 []int) []string {
	var differences []string
	maxLines := len(lines1)
	if len(lines2) > maxLines {
		maxLines = len(lines2)
	}

	for i := 0; i < maxLines; i++ {
		line1 := ""
		line2 := ""
		page1 := 0

		if i < len(lines1) {
			line1 = lines1[i]
			page1 = pages1[i]
		}
		if i < len(lines2) {
			line2 = lines2[i]

		}

		if line1 != line2 {
			contextStart := i - 2
			if contextStart < 0 {
				contextStart = 0
			}

			differences = append(differences, fmt.Sprintf("<div class=\"diff-section\" data-page=\"%d\">", page1))

			for j := contextStart; j < i; j++ {
				if j < len(lines1) {
					differences = append(differences, fmt.Sprintf("<pre class=\"context-line\">%s</pre>", lines1[j]))
				}
			}

			differences = append(differences, fmt.Sprintf("<pre class=\"page-info\">Página %d</pre>", page1))

			if line1 == "" {
				// Linha adicionada
				differences = append(differences, fmt.Sprintf("<pre class=\"line-insert\">+ %s</pre>", line2))
			} else if line2 == "" {
				// Linha deletada
				differences = append(differences, fmt.Sprintf("<pre class=\"line-delete\">- %s</pre>", line1))
			} else {
				// Linha alterada
				highlightedLine1, highlightedLine2 := highlightDiffs(line1, line2)
				differences = append(differences, fmt.Sprintf("<pre class=\"line-delete\">- %s</pre>", highlightedLine1))
				differences = append(differences, fmt.Sprintf("<pre class=\"line-insert\">+ %s</pre>", highlightedLine2))
			}

			differences = append(differences, fmt.Sprintf("</div>"))
		}
	}

	return differences
}

func GetDocxDiff(filePath1 string, filePath2 string) (*[]string, error) {
	// Caminhos dos arquivos DOCX

	// Ler o conteúdo de texto dos arquivos
	text1, pages1, err := readDocxText(filePath1)
	if err != nil {
		return nil, fmt.Errorf("Erro ao ler o arquivo %s: %v", filePath1, err)

	}

	text2, pages2, err := readDocxText(filePath2)
	if err != nil {
		return nil, fmt.Errorf("Erro ao ler o arquivo %s: %v", filePath1, err)
	}

	// Comparar os textos e obter as diferenças como uma lista de strings
	differences := compareTexts(text1, pages1, text2, pages2)
	return &differences, nil
}
