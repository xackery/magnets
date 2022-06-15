package library

import "fmt"

const (
	tileFlippedHorizontal = 0x80000000
	tileFlippedVertical   = 0x40000000
	tileFlippedDiagonal   = 0x20000000
	tileFlipped           = tileFlippedHorizontal | tileFlippedVertical | tileFlippedDiagonal
)

// GID represents a global ID, uses rotation padding on the last 3 bytes
type GID struct {
	gid uint32
}

// NewGID returns a GID
func NewGID(gid uint32) (g *GID) {
	g = &GID{
		gid: gid,
	}
	return
}

// Index returns an index without flipping data
func (g *GID) Index() uint32 {
	return uint32(g.gid &^ tileFlipped)
}

// ValueRead returns the true GID value packed with rotation data
func (g *GID) ValueRead() uint32 {
	return g.gid
}

// ValueUpdate sets the true GID value packed with rotation data
func (g *GID) ValueUpdate() uint32 {
	return g.gid
}

// RotationUpdate sets the rotation of a tile clockwise, supports 0, 90, 180, 270
func (g *GID) RotationUpdate(rotation int) (err error) {
	//771        //0 fff
	//2684355437 // 270 ccw tft
	//1073742701 // 180 ftf
	//536871683  // 90 ccw fft
	switch rotation {
	case 0:
		g.HUpdate(false)
		g.VUpdate(false)
		g.DUpdate(false)
	case 90:
		g.HUpdate(true)
		g.VUpdate(false)
		g.DUpdate(true)
	case 180:
		g.HUpdate(false)
		g.VUpdate(true)
		g.DUpdate(false)
	case 270:
		g.HUpdate(false)
		g.VUpdate(false)
		g.DUpdate(true)
	default:
		err = fmt.Errorf("invalid rotation, must be 0, 90, 180, or 270")
	}
	return
}

// RotationRead returns the rotation of the object. Can only be 0, 90, 180, or 270
func (g *GID) RotationRead() (rotation int) {
	if !g.HRead() && !g.VRead() && !g.DRead() {
		return 0
	}
	if g.HRead() && !g.VRead() && g.DRead() {
		return 90
	}
	if !g.HRead() && g.VRead() && !g.DRead() {
		return 180
	}
	if !g.HRead() && !g.VRead() && g.DRead() {
		return 270
	}
	return 0
}

// HUpdate changes the horizontal flag
func (g *GID) HUpdate(val bool) {
	if val {
		g.gid |= tileFlippedHorizontal
	} else {
		g.gid &^= tileFlippedHorizontal
	}
}

// HRead reads if the tile should be changed horizontal
func (g *GID) HRead() bool {
	return g.gid&tileFlippedHorizontal != 0
}

// VUpdate changes the horizontal flag
func (g *GID) VUpdate(val bool) {
	if val {
		g.gid |= tileFlippedVertical
	} else {
		g.gid &^= tileFlippedVertical
	}
}

// VRead returns true if a vertical rotation is in place
func (g *GID) VRead() bool {
	return g.gid&tileFlippedVertical != 0
}

// DUpdate changes the horizontal flag
func (g *GID) DUpdate(val bool) {
	if val {
		g.gid |= tileFlippedDiagonal
	} else {
		g.gid &^= tileFlippedDiagonal
	}
}

// DRead returns true if diagonal flagged
func (g *GID) DRead() bool {
	return g.gid&tileFlippedDiagonal != 0
}

// Index strips a gid of rotation data
func Index(gid uint32) uint32 {
	return uint32(gid &^ tileFlipped)
}
