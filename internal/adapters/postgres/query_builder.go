package postgres

import (
	"fmt"
	"strings"
)

type whereBuilder struct {
	conditions []string
	args       []any
	argPos     int
}

func newWhereBuilder() *whereBuilder {
	return &whereBuilder{argPos: 1}
}

func (b *whereBuilder) add(column string, value any) {
	b.conditions = append(b.conditions, fmt.Sprintf("%s = $%d", column, b.argPos))
	b.args = append(b.args, value)
	b.argPos++
}

func (b *whereBuilder) addIf(column string, value any) {
	if value == nil {
		return
	}
	b.add(column, value)
}

func (b *whereBuilder) build() (string, []any) {
	if len(b.conditions) == 0 {
		return "", b.args
	}
	return " WHERE " + strings.Join(b.conditions, " AND "), b.args
}

func (b *whereBuilder) addILike(column string, value string) {
	b.conditions = append(b.conditions, fmt.Sprintf("%s ILIKE '%%' || $%d || '%%'", column, b.argPos))
	b.args = append(b.args, value)
	b.argPos++
}
