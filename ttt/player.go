package ttt

type Player struct {
	id           int
	name         string
	char         string
	nextPlayerId int
	nextPlayer   *Player
}

func (p *Player) GetId() int {
	return p.id
}

func (p *Player) GetChar() string {
	return p.char
}

func (p *Player) GetName() string {
	return p.name
}

func (p *Player) GetNextPlayerId() int {
	return p.nextPlayerId
}
