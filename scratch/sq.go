package scratch

import (
"fmt"
sq "github.com/Masterminds/squirrel"
)

func DoIt() {
	s, args, err := sq.Select("data").
		From("policies").
		Where(sq.Eq{"data->>'label'": "lmao"}).
		ToSql()

	if err != nil {
		panic(err)
	}

	fmt.Println("Yeet", s, args)

	s, args, err = sq.Insert("policies").
		Columns("data").
		Values("dasdsa").
		ToSql()

	fmt.Println("again", s, args)
}

