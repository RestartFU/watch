package tokenizer

type Position struct {
	offset int
	line   int
	column int
}

func (p Position) Line() int {
	return p.line
}

func (p Position) Column() int {
	return p.column
}
