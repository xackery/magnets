package global

import "github.com/xackery/magnets/entity"

var (
	Player Playerer
)

type Playerer interface {
	entity.Entiter
	AttractionRange() float64
}
