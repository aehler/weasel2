package request

import "fmt"

type CharacterSkills struct{
	Path string
	CharacterID uint
	Token string
}

func (m CharacterSkills) IsValid() bool {
	return m.CharacterID != 0 && m.Token != ""
}

func (m CharacterSkills) Url() string {

	return fmt.Sprintf("characters/%d/skills/", m.CharacterID)
}

func (m CharacterSkills) RequiresAuth() bool {
	return true
}

func (m CharacterSkills) GetToken() string {
	return m.Token
}