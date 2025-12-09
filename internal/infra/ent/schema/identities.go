package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Identity holds the schema definition for the Identity entity.
type Identity struct {
	ent.Schema
}

// Fields of the Identity.
func (Identity) Fields() []ent.Field {

	return []ent.Field{

		field.Int64("id").Unique(),

		field.Int64("user_id"),

		field.String("provider"),

		field.String("provider_id"),

		field.JSON("identity_data", map[string]any{}).Optional(),

		field.Time("last_sign_in_at"),

		field.Time("created_at"),

		field.Time("updated_at"),

		field.Int64("created_by"),

		field.Int64("updated_by"),
	}
}

// Edges of the Identity.
func (Identity) Edges() []ent.Edge {
	return nil
}

func (Identity) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "identities",
		},
	}
}
