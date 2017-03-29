package simpletable

import (
	"fmt"
	"strings"

	//"github.com/davecgh/go-spew/spew"
)

type Table struct {
	Header   *Header
	Body     *Body
	Footer   *Footer
	style    *Style
	rows     []*Row
	columns  []*Column
	spanned  []*TextCell
	dividers []*Divider
}

func (t *Table) SetStyle(style *Style) {
	t.style = style
}

func (t *Table) String() string {
	t.prepareRows()
	t.prepareColumns()
	t.resizeColumns()

	s := []string{}

	for _, r := range t.rows {
		s = append(s, r.String())
	}

	return strings.Join(s, "\n")
}

func (t *Table) Print() {
	fmt.Print(t.String())
}

func (t *Table) Println() {
	fmt.Println(t.String())
}

func (t *Table) prepareRows() {
	hlen := len(t.Header.Cells)
	if hlen > 0 {
		t.rows = append(t.rows, &Row{
			Cells: t.Header.Cells,
		})

		t.rows = append(t.rows, &Row{
			Cells: []Cell{
				&Divider{
					Span: hlen,
				},
			},
		})
	}

	for _, r := range t.Body.Cells {
		t.rows = append(t.rows, &Row{
			Cells: r,
		})
	}

	flen := len(t.Footer.Cells)
	if flen > 0 {
		t.rows = append(t.rows, &Row{
			Cells: []Cell{
				&Divider{
					Span: hlen,
				},
			},
		})

		t.rows = append(t.rows, &Row{
			Cells: t.Footer.Cells,
		})
	}
}

func (t *Table) prepareColumns() {
	m := [][]Cell{}

	for _, r := range t.rows {
		row := []Cell{}

		for _, c := range r.Cells {
			row = append(row, c)
			span := 0
			var p Cell
			var tc *TextCell

			switch v := c.(type) {
			case *TextCell:
				span = v.Span
				p = v
				tc = v
			case *Divider:
				span = v.Span
				p = v
			}

			if span > 1 {
				empty := []*EmptyCell{}

				for i := 1; i < span; i++ {
					empty = append(empty, &EmptyCell{
						parent: p,
					})
				}

				for _, c := range empty {
					row = append(row, c)
				}

				if tc != nil {
					t.spanned = append(t.spanned, tc)
				}

				switch v := c.(type) {
				case *TextCell:
					v.children = empty
				case *Divider:
					v.children = empty
				}
			}
		}

		m = append(m, row)
	}

	m = t.transposeCells(m)
	for _, r := range m {
		c := &Column{
			Cells: r,
		}

		for _, cell := range c.Cells {
			cell.SetColumn(c)
		}

		t.columns = append(t.columns, c)
	}
}

func (t *Table) transposeCells(i [][]Cell) [][]Cell {
	r := [][]Cell{}

	for x := 0; x < len(i[0]); x++ {
		r = append(r, make([]Cell, len(i)))
	}

	for x, row := range i {
		for y, c := range row {
			r[y][x] = c
		}
	}

	return r
}

func (t *Table) resizeColumns() {
	for _, c := range t.columns {
		c.Resize()
	}

	for _, c := range t.spanned {
		c.Resize()
	}

	for _, d := range t.dividers {
		s := t.size()
		d.SetWidth(s)
	}
}

func (t *Table) size() int {
	return 0
}

func New() *Table {
	return &Table{
		style: StyleDefault,
		Header: &Header{
			Cells: []Cell{},
		},
		Body: &Body{
			Cells: [][]Cell{},
		},
		Footer: &Footer{
			Cells: []Cell{},
		},
		rows:     []*Row{},
		columns:  []*Column{},
		spanned:  []*TextCell{},
		dividers: []*Divider{},
	}
}
