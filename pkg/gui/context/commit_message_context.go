package context

type CommitMessageContext struct {
	selectedIndex int
}

func (self *CommitMessageContext) IncrementSelectedIndexBy(delta int) {
	self.selectedIndex = self.selectedIndex + delta
}

func (self *CommitMessageContext) GetSelectedIndex() int {
	return self.selectedIndex
}
