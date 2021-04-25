package repotab

import (
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v2"
	"github.com/mattn/go-runewidth"
	"github.com/morikuni/aec"
)

type Printer struct {
	w          io.Writer
	c          *runewidth.Condition
	columns    []builtColumn
	priorities []columnPriority
	rows       [][]Cell
	width      int
	styled     bool
}

type columnPriority struct {
	index    int
	priority int
}

type builtColumn struct {
	Column
	width int
}

type Cell struct {
	style   aec.ANSI
	content string
	width   int
}

type TableOption func(*Printer)

func OptionStyled() TableOption {
	return func(p *Printer) {
		p.styled = true
	}
}

func OptionWidth(w int) TableOption {
	return func(p *Printer) {
		p.width = w
	}
}

func insertPriority(priorities []columnPriority, newbee columnPriority) []columnPriority {
	for i, exist := range priorities {
		if exist.priority > newbee.priority {
			return append(priorities[:i], append([]columnPriority{newbee}, priorities[i:]...)...)
		}
	}
	return append(priorities, newbee)
}

func OptionColumns(columns ...Column) TableOption {
	return func(p *Printer) {
		for i, column := range columns {
			p.columns = append(p.columns, builtColumn{
				Column: column,
			})
			p.priorities = insertPriority(p.priorities, columnPriority{
				index:    i,
				priority: column.Priority,
			})
		}
	}
}

func NewPrinter(w io.Writer, option ...TableOption) *Printer {
	c := runewidth.NewCondition()
	p := &Printer{
		w: w,
		c: c,
	}
	for _, o := range option {
		o(p)
	}
	if len(p.columns) == 0 {
		OptionColumns(
			Column{
				Priority:    0,
				CellBuilder: SpecCell,
			}, Column{
				Truncatable: true,
				MinWidth:    20,
				Elipsis:     GenericElipsis,
				Priority:    3,
				CellBuilder: DescriptionCell,
			}, Column{
				Priority:    1,
				CellBuilder: AttributesCell,
			}, Column{
				Priority:    2,
				CellBuilder: PushedAtCell,
			})(p)
	}
	return p
}

func (p *Printer) Print(r gogh.Repository) error {
	cells := make([]Cell, len(p.columns))
	for i, column := range p.columns {
		content, style := column.CellBuilder.Build(r)
		width := p.c.StringWidth(content)
		if p.columns[i].width < width {
			p.columns[i].width = width
		}
		cells[i] = Cell{
			content: content,
			style:   style,
			width:   width,
		}
	}
	p.rows = append(p.rows, cells)
	return nil
}

const (
	separator = "  "
	sepLength = 2
)

func (p *Printer) truncater(rest int, column builtColumn) (int, func(Cell) *Cell) {
	if rest < 0 {
		return rest, func(Cell) *Cell { return nil }
	}
	length := rest - sepLength
	if column.width > length {
		if column.Truncatable && column.MinWidth < length {
			return 0, func(c Cell) *Cell {
				c.content = p.c.FillRight(p.c.Truncate(c.content, length, column.Elipsis), length)
				return &c
			}
		}
		return rest, func(Cell) *Cell { return nil }
	}
	return length - column.width, func(c Cell) *Cell {
		c.content = p.c.FillRight(c.content, column.width)
		return &c
	}
}

func (p *Printer) Close() error {
	picker := make([]func(Cell) *Cell, len(p.columns))
	if rest := p.width; rest > 0 {
		for _, priority := range p.priorities {
			rest, picker[priority.index] = p.truncater(rest, p.columns[priority.index])
		}
	} else {
		for i, column := range p.columns {
			column := column
			picker[i] = func(c Cell) *Cell {
				c.content = p.c.FillRight(c.content, column.width)
				return &c
			}
		}
	}
	for _, row := range p.rows {
		sep := ""
		for i, pick := range picker {
			cell := pick(row[i])
			if cell == nil {
				break
			}
			fmt.Fprint(p.w, sep)
			fmt.Fprint(p.w, cell.style.Apply(cell.content))
			sep = separator
		}
		fmt.Fprintln(p.w)
	}
	return nil
}
