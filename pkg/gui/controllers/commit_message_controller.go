package controllers

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitMessageController struct {
	baseController
	*controllerCommon
	contextsean *context.CommitMessageContext

	getCommitMessage func() string
	getCommitDescription func() string
	setCommitMessage func(message string)
	setCommitDescription func(message string)
	onCommitAttempt  func(message string)
	onCommitSuccess  func()
}

var _ types.IController = &CommitMessageController{}

func NewCommitMessageController(
	common *controllerCommon,
	getCommitMessage func() string,
	getCommitDescription func() string,
	onCommitAttempt func(message string),
	onCommitSuccess func(),
	setCommitMessage func(message string),
	setCommitDescription func(message string),
) *CommitMessageController {
	return &CommitMessageController{
		baseController:   baseController{},
		controllerCommon: common,
		contextsean: &context.CommitMessageContext{},

		getCommitMessage: getCommitMessage,
		getCommitDescription: getCommitDescription,
		onCommitAttempt:  onCommitAttempt,
		onCommitSuccess:  onCommitSuccess,
		setCommitMessage: setCommitMessage,
		setCommitDescription: setCommitDescription,
	}
}

func (self *CommitMessageController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:     opts.GetKey(opts.Config.Universal.SubmitEditorText),
			Handler: self.confirm,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.Return),
			Handler: self.close,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.PrevItem),
			Handler: self.handlePreviousCommit,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.NextItem),
			Handler: self.handleNextCommit,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.NextBlockAlt2),
			Handler: self.handleCommitDescriptionPress,
		},
	}

	return bindings
}

func (self *CommitMessageController) Context() types.Context {
	return self.context()
}

// this method is pointless in this context but I'm keeping it consistent
// with other contexts so that when generics arrive it's easier to refactor
func (self *CommitMessageController) context() types.Context {
	return self.contexts.CommitMessage
}

func (self *CommitMessageController) handlePreviousCommit() error {
	return self.handleCommitIndexChange(1)
}

func (self *CommitMessageController) handleNextCommit() error {
	if (self.contextsean.GetSelectedIndex() == 0) {
		return nil
	}
	return self.handleCommitIndexChange(-1)
}

func (self *CommitMessageController) handleCommitDescriptionPress() error {
	if err := self.c.PushContext(self.contexts.CommitDescription); err != nil {
		return err
	}
	return nil
}

func (self *CommitMessageController) handleCommitIndexChange(value int) error {
	self.contextsean.IncrementSelectedIndexBy(value)
	currentIndex:= self.contextsean.GetSelectedIndex()
	if (currentIndex == 0) {
		self.setCommitMessage("")
		return nil
	}
	return self.setCommitMessageAtIndex(currentIndex)
}

func splitMessageAndDescription(commitMessage string) (string, string) {
	// when saving the commit + description with the CommitCmdObj it creates two \n between the message & description
	parts := strings.Split(commitMessage, "\n\n")
    msg := parts[0]
    var description string
    if len(parts) > 1 {
        description = parts[1]
    }
    return msg, description
}

func (self *CommitMessageController) setCommitMessageAtIndex(index int) error {
	msg, err := self.git.Commit.GetMessageShawn(index)
	commitMessage, commitDescription := splitMessageAndDescription(msg)
	if (err != nil) {
		return self.c.ErrorMsg(self.c.Tr.CommitWithoutMessageErr)
	}
	self.setCommitMessage(commitMessage)
	self.setCommitDescription(commitDescription)
	return nil
}

func buildCommitMessage(message string, description string) string {
	if len(description) == 0 {
		return message
	}
	return message + "\n" + description
}

func (self *CommitMessageController) confirm() error {
	message :=  buildCommitMessage(self.getCommitMessage(), self.getCommitDescription())
	
	self.onCommitAttempt(message)

	if message == "" {
		return self.c.ErrorMsg(self.c.Tr.CommitWithoutMessageErr)
	}

	cmdObj := self.git.Commit.CommitCmdObj(message)
	self.c.LogAction(self.c.Tr.Actions.Commit)

	_ = self.c.PopContext()
	return self.helpers.GPG.WithGpgHandling(cmdObj, self.c.Tr.CommittingStatus, func() error {
		self.onCommitSuccess()
		return nil
	})
}

func (self *CommitMessageController) close() error {
	return self.c.PopContext()
}
