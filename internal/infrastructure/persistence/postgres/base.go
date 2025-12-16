package postgres

import "github.com/Masterminds/squirrel"

var Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Table struct {
	Name    string
	Columns []string
}

func (t Table) Select(columns ...string) squirrel.SelectBuilder {
	return Builder.Select(columns...).From(t.Name)
}

func (t Table) Insert() squirrel.InsertBuilder {
	return Builder.Insert(t.Name)
}

func (t Table) Update() squirrel.UpdateBuilder {
	return Builder.Update(t.Name)
}

func (t Table) Delete() squirrel.DeleteBuilder {
	return Builder.Delete(t.Name)
}
