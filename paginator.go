package utils

import (
	"encoding/json"
	"math"
)

const (
	PSQL_TOTAL_ROW_KEY = "total_row"
)

type Paginator struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_page"`
	TotalRows  int `json:"total_rows"`
}

func (p Paginator) String() string {
	ju, _ := json.Marshal(p)
	return string(ju)
}

func NewPaginator() Paginator {
	return Paginator{Page: 1, PerPage: 20}
}

func NewPaginatorWithConfig(page int, perPage int) Paginator {
	return Paginator{Page: page, PerPage: perPage}
}

func (p *Paginator) SetPaginatorByAllRows(allRows int) {
	p.setTotalEntrySizes(allRows)
	p.setTotalPages()
}

func (p *Paginator) setTotalEntrySizes(allRows int) {
	p.TotalRows = allRows
}

func (p *Paginator) setTotalPages() {
	p.TotalPages = int(math.Ceil(float64(p.TotalRows) / float64(p.PerPage)))
}
