package players

type PlayerInfo struct {
	Kills  uint8
	Deaths uint8
}

type PlayerEvents string

const (
	INCREMENT_KILL_COUNT  PlayerEvents = "INCREMENT_KILL_COUNT"
	INCREMENT_DEATH_COUNT PlayerEvents = "INCREMENT_DEATH_COUNT"
)
