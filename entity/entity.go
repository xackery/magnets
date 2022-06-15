package entity

import "fmt"

var (
	lastEID  uint
	entities map[uint]Entiter = make(map[uint]Entiter)
)

func Entity(entityID uint) Entiter {
	return entities[entityID]
}

func NextEntityID() uint {
	lastEID++
	if lastEID == 0 {
		lastEID = 1
	}
	return lastEID
}

func Clear() {
	lastEID = 1
	entities = make(map[uint]Entiter)
}

func Register(entity Entiter) error {
	_, ok := entities[entity.EntityID()]
	if ok {
		return fmt.Errorf("entity %d already exists", entity.EntityID())
	}
	entities[entity.EntityID()] = entity
	return nil
}
