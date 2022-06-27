package models

type (
	ChainUpdate struct {
		ID    string
		Field UpdateField
		Data  interface{}
	}
)
