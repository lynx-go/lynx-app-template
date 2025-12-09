package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type RuntimeVar struct {
	ent.Schema
}

// Fields of the Identity.
func (RuntimeVar) Fields() []ent.Field {

	return []ent.Field{

		field.Int64("id").Unique(),

		field.String("key"),

		field.String("value"),

		field.Bool("locked").Default(false),

		field.Bool("enabled").Default(true),

		field.String("desc").Default(""),

		field.Time("created_at"),

		field.Time("updated_at"),

		field.Int64("created_by"),

		field.Int64("updated_by"),
	}
}

// Edges of the Identity.
func (RuntimeVar) Edges() []ent.Edge {
	return nil
}

func (RuntimeVar) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "runtime_vars",
		},
	}
}
