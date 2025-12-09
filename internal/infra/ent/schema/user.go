package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {

	return []ent.Field{

		field.Int64("id").Unique(),

		field.String("username").Unique(),

		field.String("display_name"),

		field.String("password_hash").Optional(),

		field.String("avatar_url").Optional(),

		field.String("phone").Optional(),

		field.Time("phone_confirmed_at").Optional(),

		field.String("email").Optional(),

		field.Time("email_confirmed_at").Optional(),

		field.Int8("status").Default(0),

		field.Int8("gender").Default(0),

		field.Time("confirmed_at").Optional(),

		field.String("confirmation_token").Optional(),

		field.Time("confirmation_sent_at").Optional(),

		field.String("role"),

		field.String("recovery_token").Optional(),

		field.Time("recovery_sent_at").Optional(),

		field.JSON("app_metadata", map[string]any{}).Optional(),

		field.JSON("user_metadata", map[string]any{}).Optional(),

		field.Time("last_sign_in_at").Default(time.Now),

		field.Bool("is_super_admin").Default(false),

		field.Int64("created_by"),

		field.Int64("updated_by"),

		field.Time("banned_until").Optional(),

		field.Time("created_at").Default(time.Now),

		field.Time("updated_at").Default(time.Now),
	}

}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "users",
		},
	}
}
