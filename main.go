package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
)

func theWholeThing() (bool, string, error) {
	var (
		typeOfChange  string
		scopeOfChange string
		shortDesc     string
		longDesc      string

		confirm bool
	)

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

func main() {
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
