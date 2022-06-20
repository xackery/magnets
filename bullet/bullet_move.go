package bullet

import (
	"github.com/xackery/magnets/global"
)

func (n *Bullet) bulletMove() {
	if n.isDead {
		return
	}
	n.bulletLinear()
	n.bulletBoomerang()
}

func (n *Bullet) bulletLinear() {
	if n.behaviorType != BehaviorLinear {
		return
	}
	if global.IsDirectionLeft(n.direction) {
		n.x -= n.moveSpeed
	}
	if global.IsDirectionRight(n.direction) {
		n.x += n.moveSpeed
	}

	if global.IsDirectionDown(n.direction) {
		n.y += n.moveSpeed
	}

	if global.IsDirectionUp(n.direction) {
		n.y -= n.moveSpeed
	}
}

func (n *Bullet) bulletBoomerang() {
	if n.behaviorType != BehaviorBoomerang {
		return
	}
	if n.isReturning {
		if n.player.X() > n.x {
			n.x += n.moveSpeed
		} else if n.player.X() != n.x {
			n.x -= n.moveSpeed
		}
		if n.player.Y() > n.y {
			n.y += n.moveSpeed
		} else if n.player.Y() != n.y {
			n.y -= n.moveSpeed
		}
	} else {
		if global.IsDirectionLeft(n.direction) {
			n.y -= 0.1
			n.x -= n.moveSpeed
		}
		if global.IsDirectionRight(n.direction) {
			n.x += n.moveSpeed
		}

		if global.IsDirectionDown(n.direction) {
			n.y += n.moveSpeed
		}

		if global.IsDirectionUp(n.direction) {
			n.y -= n.moveSpeed
		}
	}

	if !n.isReturning && global.Distance(n.x, n.y, n.spawnX, n.spawnY) >= n.distance {
		n.isReturning = true
	}

	if n.isReturning && global.Distance(n.x, n.y, n.player.X(), n.player.Y()) < 5 {
		n.isDead = true
	}
}
