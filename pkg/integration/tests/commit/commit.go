package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Commit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a couple files and committing",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "myfile content")
		shell.CreateFile("myfile2", "myfile2 content")
		shell.CreateFile("myfile3", "myfile3 content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsEmpty()

		t.Views().Files().
			IsFocused().
			PressPrimaryAction(). // stage file
			SelectNextItem().
			PressPrimaryAction(). // stage other file
			Press(keys.Files.CommitChanges)

		commitMessage := "my commit message"

		t.ExpectPopup().CommitMessagePanel().Type(commitMessage).Confirm()

		t.Views().Commits().
			Lines(
				Contains(commitMessage),
			)
		t.Views().Files().
		IsFocused().
		PressPrimaryAction().
		Press(keys.Files.CommitChanges).
		Press(keys.Universal.PrevItem)

		additionalText := "number 2"

		t.ExpectPopup().CommitMessagePanel().Type(additionalText).Confirm()
		t.Views().Commits().
			Lines(
				Contains(commitMessage + additionalText),
			)

	},
})
