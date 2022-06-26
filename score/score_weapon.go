package score

var (
	weaponScores = make(map[int]*weaponScoreData)
)

type weaponScoreData struct {
	kills  int
	damage int
}

func Clear() {
	weaponScores = make(map[int]*weaponScoreData)
}

func weaponScore(weaponType int) *weaponScoreData {
	s, ok := weaponScores[weaponType]
	if !ok {
		s = &weaponScoreData{}
		weaponScores[weaponType] = s
	}
	return s
}

func AddKill(weaponType int) {
	weaponScore(weaponType).kills++
}

func AddDamage(weaponType int, damage int) {
	weaponScore(weaponType).damage += damage
}

func Damage(weaponType int) int {
	return weaponScore(weaponType).damage
}

func Kills(weaponType int) int {
	return weaponScore(weaponType).kills
}
