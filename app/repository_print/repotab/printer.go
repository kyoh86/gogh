package repotab

import (
	"fmt"
	"io"
	"os"

	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/mattn/go-runewidth"
	"github.com/morikuni/aec"
	"golang.org/x/term"
)

type Printer struct {
	w       io.Writer
	c       *runewidth.Condition
	columns []Column
	indices []int // column indices sorted by priority
	rows    [][]cell
	width   int
	styled  bool
}

type Align int

const (
	AlignLeft  Align = iota
	AlignRight Align = iota
)

type Option func(*Printer)

// Styled sets the printer to use styled output if the terminal supports it.
func Styled(force bool) Option {
	if force || term.IsTerminal(int(os.Stdout.Fd())) {
		return func(p *Printer) {
			p.styled = true
		}
	}
	return func(*Printer) {}
}

// TermWidth sets the terminal width to the printer if it is available.
func TermWidth() Option {
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		return func(p *Printer) {
			p.width = width
		}
	}
	return func(*Printer) {}
}

func Width(w int) Option {
	return func(p *Printer) {
		p.width = w
	}
}

func Columns(columns ...Column) Option {
	// get column indices sorted by priority
	insertPriority := func(indices, priors []int, newIndex, newPrior int) ([]int, []int) {
		for i, old := range priors {
			if old > newPrior {
				return append(indices[:i], append([]int{newIndex}, indices[i:]...)...),
					append(priors[:i], append([]int{newPrior}, priors[i:]...)...)
			}
		}
		return append(indices, newIndex),
			append(priors, newPrior)
	}

	var priors []int
	var indices []int
	for index, column := range columns {
		columns[index].width = 0
		indices, priors = insertPriority(indices, priors, index, column.Priority)
	}

	return func(p *Printer) {
		p.indices = indices
		p.columns = columns
	}
}

type Column struct {
	CellBuilder CellBuilder
	Elipsis     string
	MinWidth    int
	Priority    int
	Align       Align
	Truncatable bool
	width       int
}

type cell struct {
	style   aec.ANSI
	content string
}

var DefaultColumns = []Column{{
	Priority:    0,
	CellBuilder: RepoRefCell,
}, {
	Truncatable: true,
	MinWidth:    20,
	Elipsis:     "...",
	Priority:    3,
	CellBuilder: DescriptionCell,
}, {
	Priority:    1,
	CellBuilder: AttributesCell,
}, {
	Align:       AlignRight,
	Priority:    2,
	CellBuilder: UpdatedAtCell,
}}

func NewPrinter(w io.Writer, option ...Option) *Printer {
	c := runewidth.NewCondition()
	p := &Printer{
		w: w,
		c: c,
	}
	for _, o := range option {
		o(p)
	}
	if len(p.columns) == 0 {
		Columns(DefaultColumns...)(p)
	}
	return p
}

func (p *Printer) Print(r hosting.Repository) error {
	cells := make([]cell, len(p.columns))
	for i, column := range p.columns {
		content, style := column.CellBuilder.Build(r)
		width := p.c.StringWidth(content)
		if p.columns[i].width < width {
			p.columns[i].width = width
		}
		cells[i] = cell{
			content: content,
			style:   style,
		}
	}
	p.rows = append(p.rows, cells)
	return nil
}

const (
	sep       = "  "
	sepLength = 2
)

type cellPicker func(cell) *cell

func (p *Printer) skipCell(cell) *cell { return nil }

func (p *Printer) alignCell(column Column) cellPicker {
	align := p.c.FillRight
	if column.Align == AlignRight {
		align = p.c.FillLeft
	}
	return func(c cell) *cell {
		c.content = align(c.content, column.width)
		return &c
	}
}

func (p *Printer) truncateCell(rest int, column Column) (int, cellPicker) {
	if rest < 0 {
		return rest, p.skipCell
	}
	align := p.c.FillRight
	if column.Align == AlignRight {
		align = p.c.FillLeft
	}
	length := rest - sepLength
	if column.width > length {
		if column.Truncatable && column.MinWidth < length {
			return 0, func(c cell) *cell {
				c.content = align(p.c.Truncate(c.content, length, column.Elipsis), length)
				return &c
			}
		}
		return rest, p.skipCell
	}
	return length - column.width, p.alignCell(column)
}

func (p *Printer) Close() error {
	pickers := make([]cellPicker, len(p.columns))
	if rest := p.width; rest > 0 {
		for _, index := range p.indices {
			rest, pickers[index] = p.truncateCell(rest, p.columns[index])
		}
	} else {
		for i, column := range p.columns {
			pickers[i] = p.alignCell(column)
		}
	}
	if p.styled {
		for i, picker := range pickers {
			pickers[i] = func(c cell) *cell {
				newc := picker(c)
				newc.content = newc.style.Apply(newc.content)
				return newc
			}
		}
	}
	for _, row := range p.rows {
		s := ""
		for i, pick := range pickers {
			cell := pick(row[i])
			if cell == nil {
				continue
			}
			fmt.Fprint(p.w, s)
			fmt.Fprint(p.w, cell.content)
			s = sep
		}
		fmt.Fprintln(p.w)
	}
	return nil
}
