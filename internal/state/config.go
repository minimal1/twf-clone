package state

type ConfigState struct {
	defaultPath string
	maxHistory  int

	colorScheme     string
	showLineNumbers bool

	confirmDelete  bool
	followSymlinks bool
}

func NewConfigState() *ConfigState {
	return &ConfigState{
		defaultPath:     ".",
		maxHistory:      50,
		colorScheme:     "dark",
		showLineNumbers: false,
		confirmDelete:   true,
		followSymlinks:  false,
	}
}

func (cs *ConfigState) GetDefaultPath() string {
	return cs.defaultPath
}
func (cs *ConfigState) SetDefaultPath(defaultPath string) {
	cs.defaultPath = defaultPath
}

func (cs *ConfigState) GetMaxHistory() int {
	return cs.maxHistory
}
func (cs *ConfigState) SetMaxHistory(value int) {
	cs.maxHistory = value
}

func (cs *ConfigState) GetColorScheme() string {
	return cs.colorScheme
}
func (cs *ConfigState) SetColorScheme(value string) {
	cs.colorScheme = value
}

func (cs *ConfigState) GetShowLineNumbers() bool {
	return cs.showLineNumbers
}
func (cs *ConfigState) SetShowLineNumbers(value bool) {
	cs.showLineNumbers = value
}

func (cs *ConfigState) GetConfirmDelete() bool {
	return cs.confirmDelete
}
func (cs *ConfigState) SetConfirmDelete(value bool) {
	cs.confirmDelete = value
}

func (cs *ConfigState) GetFollowSymlinks() bool {
	return cs.followSymlinks
}
func (cs *ConfigState) SetFollowSymlinks(value bool) {
	cs.followSymlinks = value
}
