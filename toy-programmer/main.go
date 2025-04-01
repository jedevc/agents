package main

import (
	"dagger/toy-programmer/internal/dagger"
	"fmt"
)

type ToyProgrammer struct{}

// Write a Go program
func (m *ToyProgrammer) GoProgram(assignment string) *dagger.Container {
	env := dag.Env().
		WithToyWorkspaceInput("input", dag.ToyWorkspace(), "input workspace").
		WithToyWorkspaceOutput("output", "output workspace")

	env = dag.LLM().
		WithEnv(env).
		WithPrompt(fmt.Sprintf(`
You are an expert go programmer. You have access to a workspace.
Use the read, write, build tools to complete the following assignment.
Do not try to access the container directly.
Don't stop until your code builds.

Assignment: %s
`, assignment)).
		Env()

	return env.
		Output("output").
		AsToyWorkspace().
		Container()
}
