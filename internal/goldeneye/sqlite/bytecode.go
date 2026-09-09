package sqlite

import (
	"strconv"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/analysis"
)

// SQLite has nothing to say about a parameter's type: a bound value is
// whatever it is, and the column a parameter stands in for appears nowhere
// in what the library reports about a prepared statement. What it does
// report is the program it compiled the statement into. EXPLAIN prints the
// bytecode, in which a parameter is a Variable loading a register, and the
// column it stands in for is the other operand of the comparison the
// register reaches, the row the seek it drives lands on, or the position
// in the record that Insert writes. So each parameter is followed from the
// register it loads, through copies and through the expressions it is an
// argument of, to the first opcode that uses it against something the
// catalog can name, the way the ClickHouse check follows a placeholder
// through the query tree to the column it is compared with.
//
// Registers hold what the last opcode to write them wrote, and remember
// which parameters have been loaded into them, which a later constant does
// not clear: coalesce(?, 'x') overwrites the parameter's register with the
// fallback on the branch where the parameter is NULL, and the comparison
// that follows is still the parameter's. A parameter with no such use is
// described by what the program does to it before using it, when it does
// anything: MustBeInt says it has to be an integer, as LIMIT's does, and
// Affinity says which affinity it is coerced to.

// instr is one line of EXPLAIN output.
type instr struct {
	Addr   int     `json:"addr"`
	Opcode string  `json:"opcode"`
	P1     int     `json:"p1"`
	P2     int     `json:"p2"`
	P3     int     `json:"p3"`
	P4     *string `json:"p4"`
	P5     int     `json:"p5"`
}

func (in instr) p4() string {
	if in.P4 == nil {
		return ""
	}
	return *in.P4
}

type valueKind int

const (
	vNone     valueKind = iota
	vParam              // a Variable
	vColumn             // Column: a stored column of a cursor
	vRowid              // Rowid, IdxRowid: the rowid of a cursor
	vConstant           // a literal, with its storage class
	vFunction           // the result of a function call, with its name
	vExpr               // the result of an operator
	vRecord             // MakeRecord: the registers packed into a record
)

// value is what a register holds.
type value struct {
	kind   valueKind
	cursor int
	index  int
	class  string // vConstant
	fn     string // vFunction
	regs   []int  // vRecord
}

type register struct {
	desc   value
	params []int // every parameter loaded into the register
}

// tracer follows the parameters of one statement through its bytecode.
type tracer struct {
	cat   *catalog
	names []string // the statement's result columns
	// cursors is what each open cursor is on; ephemeral holds, for the
	// cursors on ephemeral tables and sorters, the parameters stored in
	// each column, so that a parameter put into an IN list's table comes
	// back out with the Column that reads it.
	cursors   map[int]object
	ephemeral map[int]map[int][]int
	regs      map[int]*register
	// found is the column each parameter was found to stand in for, and
	// hint what the program said about a parameter's type.
	found map[int]analysis.Column
	hint  map[int]string
}

func newTracer(cat *catalog, names []string) *tracer {
	return &tracer{
		cat:       cat,
		names:     names,
		cursors:   map[int]object{},
		ephemeral: map[int]map[int][]int{},
		regs:      map[int]*register{},
		found:     map[int]analysis.Column{},
		hint:      map[int]string{},
	}
}

// run walks the program in the order it runs: the prologue Init jumps to,
// where the Variables of a query are loaded, and then the body.
func (t *tracer) run(prog []instr) {
	order := prog
	if len(prog) > 0 && prog[0].Opcode == "Init" && prog[0].P2 > 0 && prog[0].P2 < len(prog) {
		order = append(append([]instr{}, prog[prog[0].P2:]...), prog[1:prog[0].P2]...)
	}
	for _, in := range order {
		t.step(in)
	}
}

// param describes what the parameter was found to stand in for.
func (t *tracer) param(number int) analysis.Column {
	if ac, ok := t.found[number]; ok {
		return ac
	}
	if typ := t.hint[number]; typ != "" {
		return analysis.Column{Type: &analysis.TypeExpr{Name: typ}}
	}
	return analysis.Column{}
}

func (t *tracer) reg(n int) *register {
	r := t.regs[n]
	if r == nil {
		r = &register{}
		t.regs[n] = r
	}
	return r
}

func (t *tracer) params(n int) []int {
	if r := t.regs[n]; r != nil {
		return r.params
	}
	return nil
}

// set records a write to a register. A Variable is the parameter and
// nothing else; a column read replaces whatever the register held, since
// the compiler reuses registers for the columns it reads; a constant keeps
// the parameters, for the coalesce case above; anything else carries the
// parameters it was computed from.
func (t *tracer) set(n int, v value, params []int) {
	r := t.reg(n)
	r.desc = v
	switch v.kind {
	case vConstant:
	default:
		r.params = params
	}
}

func (t *tracer) copy(from, to int) {
	src := t.reg(from)
	dst := t.reg(to)
	dst.desc = src.desc
	dst.params = append([]int{}, src.params...)
}

// assoc records the column the parameters in a register stand in for,
// keeping the first found for each.
func (t *tracer) assoc(params []int, ac analysis.Column) {
	for _, p := range params {
		if _, ok := t.found[p]; !ok {
			t.found[p] = ac
		}
	}
}

// typeHint records what the program requires of a register's parameters.
func (t *tracer) typeHint(n int, typ string) {
	if typ == "" {
		return
	}
	for _, p := range t.params(n) {
		if _, ok := t.hint[p]; !ok {
			t.hint[p] = typ
		}
	}
}

// describe turns what a register holds into the column a parameter
// compared with it stands in for, when it can be named: a stored column
// or rowid, a constant's storage class, or a function's name.
func (t *tracer) describe(v value) (analysis.Column, bool) {
	switch v.kind {
	case vColumn:
		return t.stored(v.cursor, v.index)
	case vRowid:
		if owner := t.cursors[v.cursor].owner(); owner != nil {
			return owner.rowid(), true
		}
	case vConstant:
		if v.class != "null" && v.class != "" {
			return analysis.Column{Type: &analysis.TypeExpr{Name: v.class}}, true
		}
	case vFunction:
		return analysis.Column{Name: v.fn}, true
	}
	return analysis.Column{}, false
}

// stored describes the ith stored column of a cursor.
func (t *tracer) stored(cursor, i int) (analysis.Column, bool) {
	col, owner, ok := t.cursors[cursor].column(i)
	switch {
	case !ok:
		return analysis.Column{}, false
	case col != nil:
		return col.describe(), true
	default:
		return owner.rowid(), true
	}
}

// rank orders what a parameter's fellow operands may be described by: a
// column first, then a function, then a constant.
func rank(v value) int {
	switch v.kind {
	case vColumn, vRowid:
		return 0
	case vFunction:
		return 1
	case vConstant:
		if v.class != "null" && v.class != "" {
			return 2
		}
	}
	return -1
}

// expression handles an opcode computing a result from operand registers:
// each parameter among the operands is described by the best other
// operand, and one with no other operand to be described by is carried
// into the result, for whatever the result is then compared with.
func (t *tracer) expression(args []int, dest int, result value) {
	var carried []int
	for _, a := range args {
		params := t.params(a)
		if len(params) == 0 {
			continue
		}
		best := -1
		for _, b := range args {
			if b == a {
				continue
			}
			if r := rank(t.reg(b).desc); r >= 0 && (best < 0 || r < rank(t.reg(best).desc)) {
				best = b
			}
		}
		if best >= 0 {
			if ac, ok := t.describe(t.reg(best).desc); ok {
				t.assoc(params, ac)
				continue
			}
		}
		carried = append(carried, params...)
	}
	t.set(dest, result, carried)
}

// compare handles a comparison of two registers: the parameters of each
// stand in for what the other holds.
func (t *tracer) compare(a, b int) {
	if ac, ok := t.describe(t.reg(b).desc); ok {
		t.assoc(t.params(a), ac)
	}
	if ac, ok := t.describe(t.reg(a).desc); ok {
		t.assoc(t.params(b), ac)
	}
}

// store handles a register written as the ith column of a cursor's row:
// an ephemeral table remembers the parameters for the reads to come, and
// a real one is the column the parameters stand in for.
func (t *tracer) store(cursor, i, reg int) {
	params := t.params(reg)
	if len(params) == 0 {
		return
	}
	if cols, ok := t.ephemeral[cursor]; ok {
		cols[i] = append(cols[i], params...)
		return
	}
	if ac, ok := t.stored(cursor, i); ok {
		t.assoc(params, ac)
	}
}

// key handles registers used as a key into a cursor: n registers from
// start, or the record in start when n is zero.
func (t *tracer) key(cursor, start, n int) {
	regs := t.span(start, n)
	if n == 0 {
		if r := t.reg(start); r.desc.kind == vRecord {
			regs = r.desc.regs
		}
	}
	for i, reg := range regs {
		t.store(cursor, i, reg)
	}
}

func (t *tracer) span(start, n int) []int {
	regs := make([]int, 0, max(n, 0))
	for i := 0; i < n; i++ {
		regs = append(regs, start+i)
	}
	return regs
}

// affinity names the type an affinity character coerces to.
func affinity(c byte) string {
	switch c {
	case 'A':
		return "blob"
	case 'B':
		return "text"
	case 'C':
		return "numeric"
	case 'D':
		return "integer"
	case 'E':
		return "real"
	}
	return ""
}

func (t *tracer) step(in instr) {
	switch in.Opcode {
	case "OpenRead", "OpenWrite", "ReopenIdx":
		t.cursors[in.P1] = t.cat.root(in.P3, in.P2)
	case "OpenDup":
		t.cursors[in.P1] = t.cursors[in.P2]
		if cols, ok := t.ephemeral[in.P2]; ok {
			t.ephemeral[in.P1] = cols
		}
	case "OpenEphemeral", "OpenAutoindex", "SorterOpen", "OpenPseudo":
		t.cursors[in.P1] = object{}
		t.ephemeral[in.P1] = map[int][]int{}
	case "Variable":
		t.set(in.P2, value{kind: vParam}, []int{in.P1})
	case "Column", "VColumn":
		t.set(in.P3, value{kind: vColumn, cursor: in.P1, index: in.P2}, t.ephemeral[in.P1][in.P2])
	case "Rowid", "IdxRowid":
		t.set(in.P2, value{kind: vRowid, cursor: in.P1}, nil)
	case "Copy":
		for i := 0; i <= in.P3; i++ {
			t.copy(in.P1+i, in.P2+i)
		}
	case "SCopy", "IntCopy":
		t.copy(in.P1, in.P2)
	case "Null", "SoftNull":
		end := max(in.P2, in.P3)
		if in.Opcode == "SoftNull" {
			end = in.P1
			in.P2 = in.P1
		}
		for r := in.P2; r <= end; r++ {
			t.set(r, value{kind: vConstant, class: "null"}, nil)
		}
	case "Integer", "Int64":
		t.set(in.P2, value{kind: vConstant, class: "integer"}, nil)
	case "Real":
		t.set(in.P2, value{kind: vConstant, class: "real"}, nil)
	case "String8", "String":
		t.set(in.P2, value{kind: vConstant, class: "text"}, nil)
	case "Blob":
		t.set(in.P2, value{kind: vConstant, class: "blob"}, nil)
	case "Function", "PureFunc":
		// P4 prints the function's arity, which for a variadic function
		// is negative and says nothing about this call. Its arguments
		// are the registers from P2 written so far, allocated together
		// for the call.
		name, n := parseCall(in.p4())
		if n < 0 {
			n = 0
			for t.regs[in.P2+n] != nil {
				n++
			}
		}
		t.expression(t.span(in.P2, n), in.P3, value{kind: vFunction, fn: name})
	case "Add", "Subtract", "Multiply", "Divide", "Remainder", "Concat", "BitAnd", "BitOr", "ShiftLeft", "ShiftRight":
		t.expression([]int{in.P2, in.P1}, in.P3, value{kind: vExpr})
	case "MakeRecord":
		t.set(in.P3, value{kind: vRecord, regs: t.span(in.P1, in.P2)}, nil)
	case "Insert":
		if r := t.reg(in.P2); r.desc.kind == vRecord {
			for i, reg := range r.desc.regs {
				t.store(in.P1, i, reg)
			}
		}
		if owner := t.cursors[in.P1].owner(); owner != nil {
			t.assoc(t.params(in.P3), owner.rowid())
		}
	case "IdxInsert", "SorterInsert":
		n, _ := strconv.Atoi(in.p4())
		if n > 0 {
			t.key(in.P1, in.P3, n)
		} else {
			t.key(in.P1, in.P2, 0)
		}
	case "Eq", "Ne", "Lt", "Le", "Gt", "Ge":
		t.compare(in.P1, in.P3)
	case "SeekRowid", "NotExists":
		if owner := t.cursors[in.P1].owner(); owner != nil {
			t.assoc(t.params(in.P3), owner.rowid())
		}
	case "SeekGE", "SeekGT", "SeekLE", "SeekLT":
		// On a table with a rowid the seek is by the rowid in P3; on an
		// index it is by the P4 registers from P3.
		if o := t.cursors[in.P1]; o.table != nil && !o.table.withoutRowid {
			t.assoc(t.params(in.P3), o.table.rowid())
			break
		}
		n, _ := strconv.Atoi(in.p4())
		t.key(in.P1, in.P3, n)
	case "IdxGE", "IdxGT", "IdxLE", "IdxLT", "Found", "NotFound", "NoConflict":
		n, _ := strconv.Atoi(in.p4())
		t.key(in.P1, in.P3, n)
	case "MustBeInt":
		t.typeHint(in.P1, "integer")
	case "Affinity":
		for i, c := range []byte(in.p4()) {
			if i < in.P2 {
				t.typeHint(in.P1+i, affinity(c))
			}
		}
	case "ResultRow":
		for i := 0; i < in.P2; i++ {
			if i < len(t.names) {
				t.assoc(t.params(in.P1+i), analysis.Column{Name: t.names[i]})
			}
		}
	}
}

// parseCall splits the `name(n)` a Function's P4 prints into the function's
// name and the number of arguments it is called with.
func parseCall(p4 string) (string, int) {
	open := strings.LastIndexByte(p4, '(')
	if open < 0 || !strings.HasSuffix(p4, ")") {
		return p4, 0
	}
	n, _ := strconv.Atoi(p4[open+1 : len(p4)-1])
	return p4[:open], n
}
