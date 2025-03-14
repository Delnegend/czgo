package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
)

var (
	typeOfChange  string
	scopeOfChange string
	shortDesc     string
	longDesc      string

	confirm bool
)

func theWholeThing() (bool, string, error) {

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select the type of change that you're committing").
				Options(
					huh.NewOption("fix:      A bug fix", "fix"),
					huh.NewOption("feat:     A new feature", "feat"),
					huh.NewOption("docs:     Documentation only changes", "docs"),
					huh.NewOption("style:    Changes that do not affect the meaning of the code", "style"),
					huh.NewOption("refactor: A code change that neither fixes a bug nor adds a feature", "refactor"),
					huh.NewOption("perf:     A code change that improves performance", "perf"),
					huh.NewOption("test:     Adding missing tests or correcting existing tests", "test"),
					huh.NewOption("build:    Changes that affect the build system or external dependencies", "build"),
					huh.NewOption("ci:       Changes to our CI configuration files and scripts", "ci"),
					huh.NewOption("chore:    Other changes that don't modify src or test files", "chore"),
					huh.NewOption("revert:   Reverts a previous commit", "revert"),
				).
				Height(7).
				Filtering(true).
				Value(&typeOfChange), // store the chosen option in the "burger" variable

			huh.NewInput().
				Title("Denote the SCOPE of this change (optional)").
				Value(&scopeOfChange),

			huh.NewInput().
				Title("Write a SHORT, IMPERATIVE tense description of the change").
				Value(&shortDesc),

			huh.NewText().
				Title("Provide a LONGER description of the change (optional)").
				Value(&longDesc),
		),
	)
	if err := form.Run(); err != nil {
		return false, "", err
	}

	if shortDesc == "" {
		return false, "", fmt.Errorf("You must provide a short description")
	}

	commitMessage := func() string {
		var sb strings.Builder
		sb.WriteString(typeOfChange)
		if scopeOfChange != "" {
			sb.WriteString("(")
			sb.WriteString(scopeOfChange)
			sb.WriteString(")")
		}
		sb.WriteString(": ")
		sb.WriteString(shortDesc)
		if longDesc != "" {
			sb.WriteString("\n\n")
			sb.WriteString(longDesc)
		}
		return sb.String()
	}()

	confirmForm := huh.NewConfirm().
		Title("Commit with this message?").
		Description(commitMessage).
		Affirmative("Yes!").
		Negative("No.").
		Value(&confirm)

	if err := confirmForm.Run(); err != nil {
		return false, "", err
	}

	return confirm, commitMessage, nil
}

// isGitRepo checks if the current directory is a git repository
func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

// hasStaged checks if there are any changes staged for commit
func hasStaged() bool {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	return err != nil // Exit status 1 means there are changes, which is what we want
}

func main() {
	// Check if current directory is a git repository
	if !isGitRepo() {
		fmt.Println("Not a git repository")
		os.Exit(1)
	}

	// Check if there are staged changes
	if !hasStaged() {
		fmt.Println("No changes staged for commit. Use 'git add' to stage changes.")
		os.Exit(1)
	}

	confirm, commitMessage, err := theWholeThing()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !confirm {
		confirm, commitMessage, err = theWholeThing()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd := exec.Command("git", "commit", "-m", commitMessage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
