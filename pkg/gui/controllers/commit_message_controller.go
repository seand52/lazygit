package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitMessageController struct {
	baseController
	*controllerCommon
	contextsean *context.CommitMessageContext

	getCommitMessage func() string
	setCommitMessage func(message string)
	onCommitAttempt  func(message string)
	onCommitSuccess  func()
}

var _ types.IController = &CommitMessageController{}

func NewCommitMessageController(
	common *controllerCommon,
	getCommitMessage func() string,
	onCommitAttempt func(message string),
	onCommitSuccess func(),
	setCommitMessage func(message string),
) *CommitMessageController {
	return &CommitMessageController{
		baseController:   baseController{},
		controllerCommon: common,
		contextsean: &context.CommitMessageContext{},

		getCommitMessage: getCommitMessage,
		onCommitAttempt:  onCommitAttempt,
		onCommitSuccess:  onCommitSuccess,
		setCommitMessage: setCommitMessage,
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
	return self.handleCommitIndexChange(-1)
}

func (self *CommitMessageController) handleCommitIndexChange(value int) error {
	self.contextsean.IncrementSelectedIndexBy(value)
	currentIndex:= self.contextsean.GetSelectedIndex()
	if (currentIndex <= 0) {
		self.setCommitMessage("")
		return nil
	}
	return self.setCommitMessageAtIndex(currentIndex)
}

func (self *CommitMessageController) setCommitMessageAtIndex(index int) error {
	msg, err := self.git.Commit.GetMessageShawn(index)
	if (err != nil) {
		return self.c.ErrorMsg(self.c.Tr.CommitWithoutMessageErr)
	}
	self.setCommitMessage(msg)
	return nil
}

func (self *CommitMessageController) confirm() error {
	message := self.getCommitMessage()
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
