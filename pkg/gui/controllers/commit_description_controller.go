package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitDescriptionController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &CommitMessageController{}

func NewCommitDescriptionController(
	common *controllerCommon,
) *CommitDescriptionController {
	return &CommitDescriptionController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *CommitDescriptionController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:     opts.GetKey(opts.Config.Universal.TogglePanel),
			Handler: self.close,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.Return),
			Handler: self.close,
		},
	}

	return bindings
}

func (self *CommitDescriptionController) Context() types.Context {
	return self.context()
}

// this method is pointless in this context but I'm keeping it consistent
// with other contexts so that when generics arrive it's easier to refactor
func (self *CommitDescriptionController) context() types.Context {
	return self.contexts.CommitMessage
}

func (self *CommitDescriptionController) close() error {
	return self.c.PopContext()
}
