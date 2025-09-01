package app

import (
	"fmt"
	"html/template"

	"github.com/flopp/freiburg-fit/internal/utils"
)

type Label struct {
	Type        string
	ID          string
	Name        string
	Names       string
	Description template.HTML
	Venues      []*Venue
}

func (l Label) Slug() string {
	return fmt.Sprintf("%s/%s", utils.SanitizeName(l.Type), utils.SanitizeName(l.ID))
}
