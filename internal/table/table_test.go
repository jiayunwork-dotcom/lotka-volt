package table

import (
	"testing"

	"lotka-volt/internal/lv"
)

func TestFromOrbit(t *testing.T) {
	orbit := lv.OrbitResult{Times: []float64{0, 1}, V: []float64{1, 2}, P: []float64{3, 4}, H: []float64{0, 0}}
	table := FromOrbit(orbit)
	if Rows(table) != 2 || Columns(table) != 4 {
		t.Fatalf("table=%+v", table)
	}
}

func TestDownsample(t *testing.T) {
	table := Table{Header: []string{"a", "b"}, Rows: [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}}
	out := Downsample(table, 2)
	if len(out.Rows) != 2 {
		t.Fatalf("len=%d", len(out.Rows))
	}
}

func TestFloatColumn(t *testing.T) {
	table := Table{Header: []string{"t"}, Rows: [][]string{{"1"}, {"2"}}}
	values, err := FloatColumn(table, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(values) != 2 || values[1] != 2 {
		t.Fatalf("values=%v", values)
	}
}

func TestSortByColumn(t *testing.T) {
	table := Table{Header: []string{"v"}, Rows: [][]string{{"3"}, {"1"}, {"2"}}}
	sorted, err := SortByColumn(table, 0)
	if err != nil {
		t.Fatal(err)
	}
	if sorted.Rows[0][0] != "1" || sorted.Rows[2][0] != "3" {
		t.Fatalf("rows=%v", sorted.Rows)
	}
}

func TestSummaryAndValidate(t *testing.T) {
	table := Table{Header: []string{"a"}, Rows: [][]string{{"1"}}}
	if Summary(table) == "" {
		t.Fatal("empty summary")
	}
	if err := Validate(table); err != nil {
		t.Fatal(err)
	}
}
