package main

// Monster represents a basic enemy NPC for battles.
type Monster struct {
	Name  string `json:"name"`
	Class string `json:"class"`
	Level int    `json:"level"`
	HP    int    `json:"hp"`
	AP    int    `json:"ap"`
}

func NewMonster(name, cls string, level, hp, ap int) *Monster {
	return &Monster{
		Name:  name,
		Class: cls,
		Level: level,
		HP:    hp,
		AP:    ap,
	}
}
