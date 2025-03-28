// A module for GitReview functions

package main

import (
	"context"
	"dagger/git-review/internal/dagger"
	"fmt"
)

type GitReview struct{}

func (m *GitReview) Review(
	ctx context.Context,

	base *dagger.GitRef,
	ref *dagger.GitRef,

	// +optional
	additionalInstructions string,
) (string, error) {
	baseCommit, err := base.Commit(ctx)
	if err != nil {
		return "", err
	}
	refCommit, err := ref.Commit(ctx)
	if err != nil {
		return "", err
	}

	ctr := dag.
		Alpine(dagger.AlpineOpts{
			Packages: []string{"git"},
		}).
		Container().
		WithWorkdir("/src").
		WithMountedDirectory(".", ref.Tree()).
		// HACK: until fix for dagger/dagger#7637 is merged
		WithExec([]string{"git", "fetch", "--unshallow"})

	diff, err := ctr.
		WithExec([]string{"git", "log", "-p", baseCommit + ".." + refCommit}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	if additionalInstructions == "" {
		additionalInstructions = fmt.Sprintf("\n\nAdditional Instructions: %s\n", additionalInstructions)
	}

	llm := dag.LLM().
		WithPromptVar("diff", diff).
		WithPromptVar("additionalInstructions", additionalInstructions).
		WithPrompt(`Review the following git commit log.

Git log:
<diff>
$diff
</diff>

Generate a succinct review of the Pull Request. Include the following information:
- The changes made to the code
- The rationale for the changes
- Any potential risks, considerations or security issues
- Any other relevant details
$additionalInstructions

In the review, make a recommendation for merging the PR or requesting changes,
but do not repeat the PR title or body, or summarizing the changes, focus on the
merge recommendation and assessment of the changes.

At the very end of the message, mentions if you recommends merging the PR or requesting changes, in bold, with a corresponding emoji.

Only output the review, nothing else.`)

	review, err := llm.LastReply(ctx)
	if err != nil {
		return "", err
	}

	return review, nil
}
