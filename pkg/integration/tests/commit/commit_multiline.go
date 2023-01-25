package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitMultiline = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit with a multi-line commit message",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "myfile content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsEmpty()

		t.Views().Files().
			IsFocused().
			PressPrimaryAction().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Type("first line")

		t.Views().CommitMessage().Press(keys.Universal.TogglePanel)

		t.ExpectPopup().CommitDescriptionPanel().AddNewline().AddNewline().Type("third line")

		t.Views().CommitDescription().Press(keys.Universal.TogglePanel)

		t.ExpectPopup().CommitMessagePanel().Confirm()

		t.Views().Commits().
			Lines(
				Contains("first line"),
			)

		t.Views().Commits().Focus()
		t.Views().Main().Content(MatchesRegexp("first line\n\\s*\n\\s*third line"))
	},
})
