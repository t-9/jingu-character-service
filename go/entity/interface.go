package entity

// Entity represents a Entity Interface.
type Entity interface {
	Create() error
	Destroy() error
	Update() error
	SetFieldByName(name string, value string)
	GetID() uint
}
