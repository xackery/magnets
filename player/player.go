package player

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/rs/zerolog/log"
	"github.com/xackery/aseprite"
	"github.com/xackery/magnets/camera"
	"github.com/xackery/magnets/entity"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/input"
	"github.com/xackery/magnets/library"
	"github.com/xackery/magnets/score"
	"github.com/xackery/magnets/ui/bar"
	"github.com/xackery/magnets/ui/gameover"
	"github.com/xackery/magnets/ui/level"
	"github.com/xackery/magnets/ui/life"
	"github.com/xackery/magnets/weapon"
)

var (
	cheats = []*cheat{
		{key: ebiten.Key1, weapon: weapon.WeaponCrystal},
		{key: ebiten.Key2, weapon: weapon.WeaponBoomerang},
		{key: ebiten.Key3, weapon: weapon.WeaponBoot},
		{key: ebiten.Key4, weapon: weapon.WeaponShovel},
		{key: ebiten.Key5, weapon: weapon.WeaponMagneticGloves},
		{key: ebiten.Key6, weapon: weapon.WeaponHammer},
		{key: ebiten.Key7, weapon: weapon.WeaponShuriken},
	}
)

type cheat struct {
	key    ebiten.Key
	weapon int
}

type Player struct {
	hid                   string
	entityID              uint
	layer                 *aseprite.Layer
	image                 *ebiten.Image
	maxHP                 int
	hp                    int
	spriteName            string
	layerName             string
	x                     float64
	y                     float64
	exp                   int
	level                 int
	expBar                *bar.Bar
	isLastLeft            bool
	animation             animation
	direction             int
	weapons               map[int]*weapon.Weapon
	weaponUpgradeCooldown time.Time
	invulCooldown         time.Time
	invulBlink            bool
	invulBlinkCooldown    time.Time
}

type animation struct {
	tag              *aseprite.Tag
	delay            time.Time
	index            int16
	isPingPongToggle bool
}

func New(spriteName string, layerName string) (*Player, error) {
	name := "walk"
	layer, err := library.Layer(spriteName, layerName)
	if err != nil {
		return nil, fmt.Errorf("library.Layer: %w", err)
	}
	if len(layer.Cells) < 1 {
		return nil, fmt.Errorf("no cells found on layer %s", name)
	}

	n := &Player{
		spriteName: spriteName,
		layerName:  layerName,
		hp:         2,
		maxHP:      2,
		level:      1,
		direction:  global.DirectionRight,
		layer:      layer,
		entityID:   entity.NextEntityID(),
		image:      layer.Cells[0].EbitenImage,
		weapons:    make(map[int]*weapon.Weapon),
	}
	global.Player = n

	newPick := n.randomWeaponUpgrade([]int{weapon.WeaponBoot, weapon.WeaponMagneticGloves, weapon.WeaponHeart, weapon.WeaponMagnet})
	w, err := weapon.New(newPick.Bullet.SourceWeaponType)
	if err != nil {
		fmt.Println("weapon new:", err)
		os.Exit(1)
	}
	err = n.weaponAdd(w)
	if err != nil {
		log.Error().Err(err).Msgf("weaponAdd")
	}

	life.SetHP(n.hp)
	life.SetMaxHP(n.maxHP)
	tag := "walk"
	err = n.SetAnimation(tag)
	if err != nil {
		return nil, fmt.Errorf("SetAnimation %s: %w", tag, err)
	}

	players = append(players, n)
	err = entity.Register(n)
	if err != nil {
		return nil, fmt.Errorf("entity.Register: %w", err)
	}

	n.expBar, err = bar.New("bar", "exp")
	if err != nil {
		return nil, fmt.Errorf("bar.New: %w", err)
	}
	n.expBar.SetText(fmt.Sprintf("%d", n.level))
	n.expBar.SetWidth(float64(global.ScreenWidth()))
	n.expBar.SetHeight(16)
	return n, nil
}

func (n *Player) IsHit(x, y float64) bool {
	if n.IsDead() {
		return false
	}
	/*if n.x > float64(x) {
		return false
	}
	if n.x+float64(n.layer.SpriteWidth) < float64(x) {
		return false
	}
	if n.y > float64(y) {
		return false
	}
	if n.y+float64(n.layer.SpriteHeight) < float64(y) {
		return false
	}
	return true*/
	return n.image.At(int(x), int(y)).(color.RGBA).A > 0
	//return s.image.At(x-s.x, y-s.y).(color.RGBA).A > 0
}

func (n *Player) Draw(screen *ebiten.Image) error {

	if len(n.layer.Cells) <= int(n.animation.index) {
		return fmt.Errorf("animationIndex %d is out of bounds for body cells %d", n.animation.index, len(n.layer.Cells))
	}
	c := n.layer.Cells[int(n.animation.index)]

	op := &ebiten.DrawImageOptions{}
	if n.isLastLeft {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(n.layer.SpriteWidth/2), 0)
	}
	op.GeoM.Translate(-float64(n.layer.SpriteWidth/2), -float64(n.layer.SpriteHeight/2))
	if n.isLastLeft {
		op.GeoM.Translate(float64(-c.PositionX), float64(c.PositionY))
		op.GeoM.Translate(float64(n.layer.SpriteWidth/2), 0)
	} else {
		op.GeoM.Translate(float64(c.PositionX), float64(c.PositionY))
	}
	op.GeoM.Translate(float64(camera.X), float64(camera.Y))
	op.GeoM.Translate(float64(n.x), float64(n.y))
	op.GeoM.Scale(global.ScreenScaleX(), global.ScreenScaleY())

	//op.GeoM.Translate(float64(global.ScreenWidth())/2, float64(global.ScreenHeight())/2)
	//op.GeoM.Translate(float64(n.x), float64(n.y))

	/*if n.IsDead() {
		op.ColorM.Scale(1, 1, 1, 0.6)
	} else if n.hp <= (n.maxHP / 2) {
		op.ColorM.Scale(1, 0.5, 0.5, 1)
	} else {
		op.ColorM.Scale(1, 1, 1, 1)
	}*/
	//duplication spell effect?
	/*for j := -128; j <= 128; j += 32 {
		//for i := -1; i <= 1; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(global.ScreenScaleX(), global.ScreenScaleY())
		op.GeoM.Translate(float64(int(c.PositionX)+n.x+n.dragX), float64(int(c.PositionY)+n.y+n.dragY+j))

		//op.GeoM.Translate(1, float64(j))
		// Alpha scale should be 1.0/49.0, but accumulating 1/49 49 times doesn't reach to 1, because
		// the final color is affected by the destination alpha when CompositeModeSourceOver is used.
		// This composite mode is the default mode. See how this is calculated at the doc:
		// https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#CompositeMode
		//
		// Use a higher value than 1.0/49.0. Here, 1.0/25.0 here to get a reasonable result.
		op.ColorM.Scale(1, 1, 1, 0.8)
		screen.DrawImage(c.EbitenImage, op)
		//}
	}*/

	if time.Now().Before(n.invulCooldown) {
		if time.Now().After(n.invulBlinkCooldown) {
			n.invulBlink = !n.invulBlink
			n.invulBlinkCooldown = time.Now().Add(100 * time.Millisecond)
		}
		if n.invulBlink {
			return nil
		}
	}

	screen.DrawImage(c.EbitenImage, op)

	if n.IsDead() {
		return nil
	}
	//text.Draw(screen, n.nameTag, font.NormalFont(), n.x-(len(n.nameTag)*2)+1, n.y+int(n.layer.SpriteHeight)+40+1, color.Black)
	//text.Draw(screen, n.nameTag, font.NormalFont(), n.x-(len(n.nameTag)*2), n.y+int(n.layer.SpriteHeight)+40, color.White)
	//x := n.x
	//y := n.y + 30
	//x -= float64(n.SWidth() / 2)
	//y += float64(n.SHeight() / 2)

	return nil
}

// SetAnimation sets the animation of the player
func (n *Player) SetAnimation(name string) error {
	name = strings.ToLower(name)
	tag, err := library.Tag(n.spriteName, name)
	if err != nil {
		return fmt.Errorf("tag: %w", err)
	}
	n.animation.tag = tag
	n.animation.index = tag.From
	return nil
}

func (n *Player) animationStep() {
	if n.IsDead() {
		return
	}
	if global.IsPaused() {
		return
	}
	if n.animation.delay.After(time.Now()) {
		return
	}

	if n.animation.tag.AnimationDirection == 2 && n.animation.isPingPongToggle {
		n.animation.index--

		if n.animation.index <= n.animation.tag.From {
			n.animation.isPingPongToggle = false
		}
		if n.animation.index < 0 {
			n.animation.index = 0
		}
	} else {
		n.animation.index++
		if n.animation.index > n.animation.tag.To {
			if n.animation.tag.AnimationDirection == 2 {
				n.animation.index -= 2
				n.animation.isPingPongToggle = true
			} else {
				n.animation.index = n.animation.tag.From
			}
		}
	}

	if n.layer == nil {
		return
	}

	if len(n.layer.Cells) <= int(n.animation.index) {
		return
	}

	c := n.layer.Cells[int(n.animation.index)]
	n.animation.delay = time.Now().Add(time.Duration(c.Duration) * time.Millisecond)
}

func (n *Player) Update() {

	if global.Countdown <= 0 {
		gameover.Clear()

		isStats := false
		for i := 0; i < weapon.WeaponMax; i++ {
			dmg := score.Damage(i)
			if dmg > 0 {
				gameover.AddText(fmt.Sprintf("%s damage: %d", weapon.Name(i), dmg))
				gameover.AddText(fmt.Sprintf("%s kills: %d", weapon.Name(i), score.Kills(i)))
				isStats = true
			}
		}
		if !isStats {
			gameover.AddText("No damage was done")
		}

		gameover.SetTitle("Victory!")
		global.SetIsGameOver(true)
		global.SetIsPaused(true)
	}

	if global.IsLevelUp() {
		var button *level.Level
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			button = level.IsHit(float64(x), float64(y))
		}

		/*if button == nil && input.IsPressed(ebiten.Key1) {
			button = level.ByIndex(0)
		}

		if button == nil && input.IsPressed(ebiten.Key2) {
			button = level.ByIndex(1)
		}

		if button == nil && input.IsPressed(ebiten.Key3) {
			button = level.ByIndex(2)
		}*/
		if button == nil {
			return
		}
		if n.hasWeapon(button.Data.Bullet.SourceWeaponType) {
			n.weaponUpgrade(button.Data.Bullet.SourceWeaponType)
		} else {
			w, err := weapon.New(button.Data.Bullet.SourceWeaponType)
			if err != nil {
				fmt.Println("weapon new:", err)
				os.Exit(1)
			}
			err = n.weaponAdd(w)
			if err != nil {
				log.Error().Err(err).Msgf("weaponAdd")
			}
		}

		global.SetIsLevelUp(false)
		global.SetIsPaused(false)
		level.Clear()
		return
	}

	for _, cheat := range cheats {
		if time.Now().Before(n.weaponUpgradeCooldown) {
			break
		}

		if input.IsPressed(cheat.key) {
			n.addExp(50)
			if n.hasWeapon(cheat.weapon) {
				n.weaponUpgrade(cheat.weapon)
			} else {
				w, err := weapon.New(cheat.weapon)
				if err != nil {
					fmt.Println("weapon new:", err)
					os.Exit(1)
				}
				err = n.weaponAdd(w)
				if err != nil {
					log.Error().Err(err).Msgf("weaponAdd")
				}
			}
			n.weaponUpgradeCooldown = time.Now().Add(1 * time.Second)
			break
		}

	}

	if global.IsPaused() {
		return
	}
	isMoving := input.IsPlayerMoving()
	moveSpeed := float64(0.5)
	bootLvl := n.weaponLevel(weapon.WeaponBoot)
	if bootLvl > 0 {
		moveSpeed += float64(bootLvl) * 0.2
	}
	if isMoving {
		n.direction = input.PlayerDirection
		n.animationStep()
		delta := float64(moveSpeed)
		if global.IsDirectionLeft(n.direction) {
			/*delta = 0
			for ; delta < moveSpeed; delta += 0.1 {
				if world.IsCollision(n.x-delta, n.y) {
					delta -= 0.1
					break
				}
			}*/

			camera.X += delta
			n.x -= delta
			n.isLastLeft = true
		}
		if global.IsDirectionRight(n.direction) {
			/*delta = 0
			for ; delta < moveSpeed; delta += 0.1 {
				if world.IsCollision(n.x+delta, n.y) {
					delta -= 0.1
					break
				}
			}*/
			camera.X -= delta
			n.x += delta
			n.isLastLeft = false
		}

		if global.IsDirectionDown(n.direction) {
			/*delta = 0
			for ; delta < moveSpeed; delta += 0.1 {
				if world.IsCollision(n.x, n.y+delta) {
					delta -= 0.1
					break
				}
			}*/
			camera.Y -= delta
			n.y += delta
		}

		if global.IsDirectionUp(n.direction) {
			/*delta = 0
			for ; delta < moveSpeed; delta += 0.1 {
				if world.IsCollision(n.x, n.y-delta) {
					delta -= 0.1
					break
				}
			}*/
			camera.Y += delta
			n.y -= delta
		}
	} else {
		n.animation.index = 0
	}

	n.weaponUpdate()
}

func (n *Player) SetPosition(x, y float64) {
	n.x, n.y = x, y
}

func (n *Player) Position() (float64, float64) {
	return n.x, n.y
}

func (n *Player) HID() string {
	return n.hid
}

func (n *Player) Damage(damage int) bool {
	if time.Now().Before(n.invulCooldown) {
		return false
	}
	if n.hp < 1 {
		return false
	}
	n.hp -= damage
	err := library.AudioPlay("player-hurt")
	if err != nil {
		log.Error().Err(err).Msgf("audioplay playerHurt")
	}
	if n.hp <= 0 {
		log.Debug().Msgf("player got killed")
		gameover.Clear()

		isStats := false
		for i := 0; i < weapon.WeaponMax; i++ {
			dmg := score.Damage(i)
			if dmg > 0 {
				gameover.AddText(fmt.Sprintf("%s damage: %d", weapon.Name(i), dmg))
				gameover.AddText(fmt.Sprintf("%s kills: %d", weapon.Name(i), score.Kills(i)))
				isStats = true
			}
		}
		if !isStats {
			gameover.AddText("No damage was done")
		}

		gameover.SetTitle("Game Over")
		n.hp = 0
		life.SetHP(0)
		global.SetIsGameOver(true)
		global.SetIsPaused(true)
		return true
	}
	n.invulCooldown = time.Now().Add(1 * time.Second)
	life.SetHP(n.hp)
	return false
}

func (n *Player) IsDead() bool {
	return n.hp < 1
}

func (n *Player) EntityID() uint {
	return n.entityID
}

func (n *Player) SWidth() int {
	return int(n.layer.SpriteWidth * uint16(global.ScreenScaleX()))
}

func (n *Player) SHeight() int {
	return int(n.layer.SpriteHeight * uint16(global.ScreenScaleY()))
}

func (n *Player) AnimationAttack() error {
	return nil
}

func (n *Player) AnimationGotHit() error {
	return nil
}

func (n *Player) X() float64 {
	return n.x
}

func (n *Player) Y() float64 {
	return n.y
}

func (n *Player) maxExp() int {
	return n.level * 5
}

func (n *Player) addExp(value int) {
	n.exp += value
	if n.exp >= n.maxExp() {
		n.exp = 0
		n.level++
		n.levelup()
	}
	log.Debug().Msgf("%0.2f %d/%d", float64(n.exp)/float64(n.maxExp()), n.exp, n.maxExp())
	n.expBar.SetPercent(float64(n.exp) / float64(n.maxExp()))
}

func (n *Player) AttractionRange() float64 {
	dist := float64(20)
	w, ok := n.weapons[weapon.WeaponMagnet]
	if !ok {
		return dist
	}
	return float64(dist + float64(w.Level*10))
}

func (n *Player) levelup() {
	var err error
	level.Clear()
	lastPicks := []int{}

	for i := 0; i < 3; i++ {
		newPick := n.randomWeaponUpgrade(lastPicks)
		_, err = level.New(newPick)
		if err != nil {
			log.Error().Err(err).Msgf("level new")
		}
		lastPicks = append(lastPicks, newPick.Bullet.SourceWeaponType)
	}
	global.SetIsLevelUp(true)
	global.SetIsPaused(true)
	err = library.AudioPlay("levelup")
	if err != nil {
		log.Error().Err(err).Msgf("play levelup")
	}
	n.expBar.SetText(fmt.Sprintf("%d", n.level))
	log.Debug().Msgf("player is now level %d", n.level)
}

func (n *Player) randomWeaponUpgrade(lastPicks []int) *weapon.WeaponData {
	for i := 0; i < 30; i++ {
		weaponType := rand.Intn(weapon.WeaponMax-1) + 1
		isNew := true
		for _, last := range lastPicks {
			if weaponType != last {
				continue
			}
			isNew = false
			break
		}
		if !isNew {
			continue
		}
		level := n.weaponLevel(weaponType)
		if level >= 8 {
			continue
		}
		if level > 0 {
			return n.weapons[weaponType].Data
		}
		return weapon.WeaponInfo(weaponType)
	}
	return nil
}
