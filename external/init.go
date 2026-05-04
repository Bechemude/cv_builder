package external

type External struct {
	LLM *LLM
}

func Init() (*External, error) {
	llm := InitLLM()

	return &External{
		LLM: llm,
	}, nil
}
