package output

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type Table struct {
	out     io.Writer
	headers []string
	rows    [][]string
}

func NewTable(out io.Writer, headers ...string) *Table {
	return &Table{out: out, headers: headers}
}

func (t *Table) Row(values ...string) {
	row := make([]string, len(t.headers))
	copy(row, values)
	t.rows = append(t.rows, row)
}

func (t *Table) Render() error {
	if len(t.headers) == 0 {
		return nil
	}

	widths := t.widths()
	border := renderBorder(widths)

	if _, err := fmt.Fprintln(t.out, border); err != nil {
		return err
	}
	if err := writeTableRow(t.out, t.headers, widths); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(t.out, border); err != nil {
		return err
	}

	for _, row := range t.rows {
		if err := writeTableRow(t.out, row, widths); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(t.out, border); err != nil {
			return err
		}
	}

	return nil
}

func (t *Table) widths() []int {
	widths := make([]int, len(t.headers))
	for i, header := range t.headers {
		widths[i] = displayLen(header)
	}

	for _, row := range t.rows {
		for i := range widths {
			if i >= len(row) {
				continue
			}
			if n := displayLen(row[i]); n > widths[i] {
				widths[i] = n
			}
		}
	}
	return widths
}

func renderBorder(widths []int) string {
	var b strings.Builder
	for _, width := range widths {
		b.WriteByte('+')
		b.WriteString(strings.Repeat("-", width+2))
	}
	b.WriteByte('+')
	return b.String()
}

func writeTableRow(out io.Writer, values []string, widths []int) error {
	var b strings.Builder
	for i, width := range widths {
		value := ""
		if i < len(values) {
			value = values[i]
		}

		b.WriteString("| ")
		b.WriteString(value)
		b.WriteString(strings.Repeat(" ", width-displayLen(value)))
		b.WriteByte(' ')
	}
	b.WriteString("|\n")

	_, err := io.WriteString(out, b.String())
	return err
}

func displayLen(value string) int {
	return utf8.RuneCountInString(value)
}
