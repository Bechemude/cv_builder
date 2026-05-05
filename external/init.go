package external

import "cvbuilder/config"

type External struct {
	LLM *LLM
	Web *Web
}

func Init(c *config.Config) *External {
	return &External{
		LLM: InitLLM(c),
		Web: InitWeb(),
	}
}
