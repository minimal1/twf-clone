package state

type SortType int

const (
	SortByName SortType = iota
	SortBySize
	SortByDate
)

type ViewMode int

const (
	ViewModeNormal ViewMode = iota
	ViewModeSearch
	ViewModeHelp
)

const sortTypeCount = 3

type ViewState struct {
	scrollOffset int
	sortBy       SortType
	mode         ViewMode
	filterText   string
	showHidden   bool
}

func NewViewState() *ViewState {
	return &ViewState{
		scrollOffset: 0,
		sortBy:       SortByName,
		mode:         ViewModeNormal,
		filterText:   "",
		showHidden:   false,
	}
}

// 스크롤 관리
func (vs *ViewState) GetScrollOffset() int {
	return vs.scrollOffset
}
func (vs *ViewState) SetScrollOffset(offset int) {
	vs.scrollOffset = offset
}
func (vs *ViewState) ScrollUp(amount int) {
	vs.scrollOffset -= amount

	if vs.scrollOffset < 0 {
		vs.scrollOffset = 0
	}
}
func (vs *ViewState) ScrollDown(amount int) {
	vs.scrollOffset += amount
}

// 표시 옵션
func (vs *ViewState) ToggleHidden() {
	vs.showHidden = !vs.showHidden
}
func (vs *ViewState) ShowHidden() bool {
	return vs.showHidden
}

// 정렬
func (vs *ViewState) GetSortType() SortType {
	return vs.sortBy
}
func (vs *ViewState) SetSortType(sortType SortType) {
	vs.sortBy = sortType
}
func (vs *ViewState) CycleSortType() {
	vs.sortBy = (vs.sortBy + 1) % sortTypeCount
}

// 필터, 검색
func (vs *ViewState) SetFilter(text string) {
	vs.filterText = text
}
func (vs *ViewState) GetFilter() string {
	return vs.filterText
}
func (vs *ViewState) ClearFilter() {
	vs.filterText = ""
}

// 뷰 모드
func (vs *ViewState) GetMode() ViewMode {
	return vs.mode
}
func (vs *ViewState) SetMode(mode ViewMode) {
	vs.mode = mode
}
