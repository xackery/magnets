package bullet

import (
	"math"

	"github.com/xackery/magnets/global"
)

func (n *Bullet) bulletMove() {
	if n.isDead {
		return
	}
	n.bulletLinear()
	n.bulletBoomerang()
	n.bulletCircle()
	n.bulletLasso()
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

func (n *Bullet) bulletCircle() {
	if n.behaviorType != BehaviorCircle {
		return
	}
	theta := n.rotation * math.Pi
	n.x = n.player.X() + math.Sin(theta)*n.distance
	n.y = n.player.Y() + math.Cos(theta)*n.distance
	n.rotation += n.moveSpeed * 0.01
	if n.rotation >= 2 {
		n.rotation = 0
	}
}

func (n *Bullet) bulletLasso() {
	if n.behaviorType != BehaviorLasso {
		return
	}

	if n.rotation == 0 {
		//n.rotation = -math.Sqrt(3)
	}
	t := n.rotation
	n.x = n.player.X() + 3*t - math.Pow(t, 3)
	n.y = n.player.Y() + 3*math.Pow(t, 2)
	//rotation := 9 * t / (3 - 3*math.Pow(t, 2))
	n.rotation += 0.2

	/*theta := n.rotation * math.Pi
	n.x = n.spawnX + math.Sin(theta)*n.distance
	n.y = n.spawnY + math.Cos(theta)*n.distance

	n.rotation += n.moveSpeed * 0.01*/
}
