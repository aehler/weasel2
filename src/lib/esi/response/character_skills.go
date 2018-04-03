package response

type CharacterSkills struct {
	Skills []struct{
		SkillID uint `json:"skill_id"`
		Sp uint `json:"skillpoints_in_skill"`
		Trained uint `json:"trained_skill_level"`
		Active uint `json:"active_skill_level"`
	}
	TotalSP uint `json:"total_sp"`
}