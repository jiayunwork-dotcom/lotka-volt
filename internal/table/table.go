package table

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"lotka-volt/internal/lv"
)

type Table struct {
	Header []string   `json:"header"`
	Rows   [][]string `json:"rows"`
}

func FromOrbit(orbit lv.OrbitResult) Table {
	table := Table{Header: []string{"t", "V", "P", "H"}}
	for i := range orbit.Times {
		table.Rows = append(table.Rows, []string{
			strconv.FormatFloat(orbit.Times[i], 'g', -1, 64),
			strconv.FormatFloat(orbit.V[i], 'g', -1, 64),
			strconv.FormatFloat(orbit.P[i], 'g', -1, 64),
			strconv.FormatFloat(orbit.H[i], 'g', -1, 64),
		})
	}
	return table
}

func SaveCSV(path string, table Table) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	if err := writer.Write(table.Header); err != nil {
		return err
	}
	for _, row := range table.Rows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func LoadCSV(path string) (Table, error) {
	file, err := os.Open(path)
	if err != nil {
		return Table{}, err
	}
	defer file.Close()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return Table{}, err
	}
	table := Table{}
	for i, record := range records {
		if i == 0 {
			table.Header = record
			continue
		}
		table.Rows = append(table.Rows, record)
	}
	return table, nil
}

func Text(table Table) string {
	var b strings.Builder
	for i, header := range table.Header {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(header)
	}
	b.WriteString("\n")
	for _, row := range table.Rows {
		for i, cell := range row {
			if i > 0 {
				b.WriteString(" ")
			}
			b.WriteString(cell)
		}
		b.WriteString("\n")
	}
	return b.String()
}

func Downsample(table Table, maxRows int) Table {
	if len(table.Rows) <= maxRows {
		return table
	}
	out := Table{Header: append([]string(nil), table.Header...)}
	step := float64(len(table.Rows)) / float64(maxRows)
	for i := 0; i < maxRows; i++ {
		index := int(float64(i) * step)
		out.Rows = append(out.Rows, table.Rows[index])
	}
	return out
}

func Columns(table Table) int {
	return len(table.Header)
}

func Rows(table Table) int {
	return len(table.Rows)
}

func Column(table Table, index int) ([]string, error) {
	if index < 0 || index >= len(table.Header) {
		return nil, fmt.Errorf("column out of range")
	}
	out := make([]string, 0, len(table.Rows))
	for _, row := range table.Rows {
		if index >= len(row) {
			return nil, fmt.Errorf("row too short")
		}
		out = append(out, row[index])
	}
	return out, nil
}

func FloatColumn(table Table, index int) ([]float64, error) {
	values, err := Column(table, index)
	if err != nil {
		return nil, err
	}
	out := make([]float64, 0, len(values))
	for _, value := range values {
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		out = append(out, parsed)
	}
	return out, nil
}

func MaxFloatColumn(table Table, index int) (float64, error) {
	values, err := FloatColumn(table, index)
	if err != nil {
		return 0, err
	}
	if len(values) == 0 {
		return 0, fmt.Errorf("empty column")
	}
	max := values[0]
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max, nil
}

func MinFloatColumn(table Table, index int) (float64, error) {
	values, err := FloatColumn(table, index)
	if err != nil {
		return 0, err
	}
	if len(values) == 0 {
		return 0, fmt.Errorf("empty column")
	}
	min := values[0]
	for _, value := range values {
		if value < min {
			min = value
		}
	}
	return min, nil
}

func MeanFloatColumn(table Table, index int) (float64, error) {
	values, err := FloatColumn(table, index)
	if err != nil {
		return 0, err
	}
	if len(values) == 0 {
		return 0, fmt.Errorf("empty column")
	}
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values)), nil
}

func SortByColumn(table Table, index int) (Table, error) {
	if index < 0 || index >= len(table.Header) {
		return Table{}, fmt.Errorf("column out of range")
	}
	keys := make([]float64, len(table.Rows))
	for i, row := range table.Rows {
		value, err := strconv.ParseFloat(row[index], 64)
		if err != nil {
			return Table{}, err
		}
		keys[i] = value
	}
	order := make([]int, len(table.Rows))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(i, j int) bool {
		return keys[order[i]] < keys[order[j]]
	})
	out := Table{Header: append([]string(nil), table.Header...)}
	for _, i := range order {
		out.Rows = append(out.Rows, table.Rows[i])
	}
	return out, nil
}

func Summary(table Table) string {
	return fmt.Sprintf("%d rows x %d columns", len(table.Rows), len(table.Header))
}

func IsEmpty(table Table) bool {
	return len(table.Rows) == 0
}

func Validate(table Table) error {
	if len(table.Header) == 0 {
		return fmt.Errorf("empty header")
	}
	for i, row := range table.Rows {
		if len(row) != len(table.Header) {
			return fmt.Errorf("row %d width mismatch", i)
		}
	}
	return nil
}
