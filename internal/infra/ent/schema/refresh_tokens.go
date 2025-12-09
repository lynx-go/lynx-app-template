package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type RefreshToken struct {
	ent.Schema
}

func (RefreshToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.Int64("user_id"),
		field.String("token"),
		field.Bool("revoked"),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

func (RefreshToken) Edges() []ent.Edge {
	return nil
}

func (RefreshToken) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "refresh_tokens",
		},
	}
}
