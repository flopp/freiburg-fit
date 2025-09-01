package app

import (
	"fmt"
	"html/template"

	"github.com/flopp/freiburg-fit/internal/utils"
)

type List struct {
	Name        string
	Description template.HTML
	Venues      []*Venue
	Type        string
}

func (l List) Slug() string {
	return fmt.Sprintf("%s/%s", l.Type, utils.SanitizeName(l.Name))
}
