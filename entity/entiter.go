package entity

type Entiter interface {
	EntityID() uint
	X() float64
	Y() float64
	SWidth() int
	SHeight() int
}
