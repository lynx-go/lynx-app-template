package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --template ./template --feature entql,sql/upsert,sql/modifier ./schema --target ./gen
