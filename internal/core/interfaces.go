package core

type LLMClient interface {
	Generate(prompt string) (string, error)
}

type Planner interface {
	CreatePlan(prompt string) (string, error)
}

type Coder interface {
	Implement(spec string) error
}
