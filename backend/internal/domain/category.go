package domain

type Category struct {
	ID   int    `db:"category_id"`
	Name string `db:"category_name"`
	Slug string `db:"category_slug"`
}
