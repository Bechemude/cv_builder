package external

import "cvbuilder/config"

type External struct {
	LLM *LLM
}

func Init(c *config.Config) *External {
	llm := InitLLM(c)

	return &External{
		LLM: llm,
	}
}
