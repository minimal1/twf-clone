# TWF Clone 다음 단계 학습 가이드

> **목적**: 프로젝트 완료 후 추가 학습 및 확장을 위한 로드맵
> **대상**: 기본 기능을 완성한 후 더 나아가고 싶은 개발자
> **방법**: 난이도별로 체계적으로 학습 및 구현

---

## 📋 학습 로드맵 개요

```
현재 위치: 5.5단계 완료 ✅
    ↓
[빠른 개선] → [중급 확장] → [고급 기능] → [전문가 레벨]
  (1-2시간)    (4-8시간)    (8-16시간)   (16+시간)
```

---

## 🚀 Level 1: 빠른 개선 (1-2시간)

### 1.1 코드 품질 개선

#### 목표: golangci-lint 도입

**학습 내용**:
- Go 린터 도구 이해
- 코드 품질 기준 설정

**실습 단계**:
```bash
# 1. golangci-lint 설치
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 2. 실행
golangci-lint run

# 3. 자동 수정
golangci-lint run --fix
```

**설정 파일 작성**:
```yaml
# .golangci.yml
linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
```

**예상 결과**:
- 코드 스타일 일관성 확보
- 잠재적 버그 발견
- 사용하지 않는 코드 제거

---

#### 목표: 주석 및 문서화 보강

**실습 단계**:
1. 모든 public 함수에 주석 추가
2. 패키지 문서 작성
3. 예제 코드 추가

**예시**:
```go
// Package terminal provides low-level terminal control functionality
// for TUI applications. It handles raw mode, ANSI escape sequences,
// and keyboard input parsing.
package terminal

// Terminal represents a terminal instance with raw mode capabilities.
// It manages the terminal state and provides methods for input/output.
type Terminal struct {
    // ...
}

// NewTerminal creates a new Terminal instance by opening /dev/tty.
// It returns an error if /dev/tty cannot be accessed.
//
// Example:
//     term, err := terminal.NewTerminal()
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer term.Cleanup()
func NewTerminal() (*Terminal, error) {
    // ...
}
```

---

### 1.2 사용자 경험 개선

#### 목표: 도움말 화면 추가 (? 키)

**난이도**: ⭐ (쉬움)
**소요 시간**: 1-2시간

**구현 계획**:

1. **새 파일 생성**: `internal/views/help-view.go`

```go
package views

type HelpView struct{}

func (h *HelpView) Render(term *Terminal, state *AppState, area Rect) error {
    helpText := []string{
        "TWF Clone - Keyboard Shortcuts",
        "",
        "Navigation:",
        "  j/↓       - Move down",
        "  k/↑       - Move up",
        "  h/←       - Parent directory",
        "  l/→       - Child directory",
        "",
        "Actions:",
        "  Enter     - Expand/collapse",
        "  Space     - Select",
        "  m + [key] - Set bookmark",
        "  ' + [key] - Jump to bookmark",
        "",
        "Other:",
        "  ?         - Show this help",
        "  q/ESC     - Quit",
    }

    // 화면 중앙에 렌더링
    for i, line := range helpText {
        term.WriteAt(area.X+2, area.Y+i+2, line)
    }
    return nil
}
```

2. **ViewState에 Help 모드 추가**:
```go
// internal/state/view.go
const (
    Normal ViewMode = iota
    Help
    // ...
)
```

3. **메인 루프에서 처리**:
```go
// cmd/twf/main.go
case terminal.KeyRune:
    if e.Rune == '?' {
        app.appState.View().SetMode(state.Help)
    }
```

**학습 포인트**:
- 새로운 뷰 추가 방법
- 모달 UI 구현
- 상태 전환 처리

---

#### 목표: 상태바 정보 확장

**난이도**: ⭐ (쉬움)
**소요 시간**: 30분

**추가할 정보**:
- 현재 디렉토리의 파일/폴더 개수
- 선택된 항목의 총 크기
- 현재 시간

**구현**:
```go
// internal/views/status-view.go
func (sv *StatusView) Render(...) error {
    // 기존: "Path: /home/user | Selected: 3"
    // 개선: "Path: /home/user (15 files, 3 dirs) | Selected: 3 (1.2MB) | 14:30"

    currentNode := state.Cursor().GetCurrentNode()
    fileCount, dirCount := countItems(currentNode)
    selectedSize := calculateSize(state.Selection().GetSelected())
    currentTime := time.Now().Format("15:04")

    leftInfo := fmt.Sprintf("%s (%d files, %d dirs)",
        currentNode.Path, fileCount, dirCount)
    rightInfo := fmt.Sprintf("Selected: %d (%s) | %s",
        len(state.Selection().GetSelected()),
        formatSize(selectedSize),
        currentTime)

    // ...
}
```

---

## 🎯 Level 2: 중급 확장 (4-8시간)

### 2.1 파일 미리보기 기능

**난이도**: ⭐⭐⭐ (중간)
**소요 시간**: 4-6시간

#### 설계

**화면 레이아웃**:
```
┌─────────────────────────────────────┐
│ File Tree     │ Preview             │
│ /home/user/   │ package main        │
│ ├─ docs/      │                     │
│ ├─ src/       │ import (            │
│ │  ├─ main.go │     "fmt"           │
│ │  └─ util.go │ )                   │
│ └─ README.md  │                     │
│               │ func main() {       │
│               │     fmt.Println()   │
│               │ }                   │
└─────────────────────────────────────┘
```

#### 구현 단계

**1단계: PreviewView 구현**

```go
// internal/views/preview-view.go
package views

type PreviewView struct {
    maxLines int
}

func NewPreviewView() *PreviewView {
    return &PreviewView{maxLines: 50}
}

func (pv *PreviewView) Render(term *Terminal, state *AppState, area Rect) error {
    currentNode := state.Cursor().GetCurrentNode()

    if currentNode == nil || currentNode.IsDir {
        return pv.renderDirPreview(term, area, currentNode)
    }

    return pv.renderFilePreview(term, area, currentNode)
}

func (pv *PreviewView) renderFilePreview(term *Terminal, area Rect, node *TreeNode) error {
    // 파일 읽기
    content, err := os.ReadFile(node.Path)
    if err != nil {
        return pv.renderError(term, area, err)
    }

    // 바이너리 파일 확인
    if !isTextFile(content) {
        return pv.renderBinaryInfo(term, area, node)
    }

    // 텍스트 파일 렌더링
    lines := strings.Split(string(content), "\n")
    for i := 0; i < min(len(lines), area.Height); i++ {
        term.WriteAt(area.X, area.Y+i, truncate(lines[i], area.Width))
    }

    return nil
}

func (pv *PreviewView) renderDirPreview(term *Terminal, area Rect, node *TreeNode) error {
    // 디렉토리 정보 표시
    entries, _ := os.ReadDir(node.Path)

    fileCount := 0
    dirCount := 0
    totalSize := int64(0)

    for _, entry := range entries {
        if entry.IsDir() {
            dirCount++
        } else {
            fileCount++
            info, _ := entry.Info()
            totalSize += info.Size()
        }
    }

    summary := []string{
        fmt.Sprintf("Directory: %s", node.Name),
        "",
        fmt.Sprintf("Files: %d", fileCount),
        fmt.Sprintf("Directories: %d", dirCount),
        fmt.Sprintf("Total size: %s", formatSize(totalSize)),
    }

    for i, line := range summary {
        term.WriteAt(area.X+2, area.Y+i+2, line)
    }

    return nil
}

func isTextFile(content []byte) bool {
    // 간단한 텍스트 파일 감지
    if len(content) == 0 {
        return true
    }

    // NULL 바이트가 있으면 바이너리
    for i := 0; i < min(512, len(content)); i++ {
        if content[i] == 0 {
            return false
        }
    }

    return true
}
```

**2단계: Layout 수정**

```go
// internal/views/layout.go
type Layout struct {
    treeView    View
    statusView  View
    previewView View  // 추가

    treeArea    Rect
    previewArea Rect  // 추가
    statusArea  Rect

    showPreview bool  // 토글
}

func (l *Layout) SetSize(width, height int) {
    if l.showPreview {
        // 반반 분할
        treeWidth := width / 2
        previewWidth := width - treeWidth

        l.treeArea = Rect{0, 0, treeWidth, height - 1}
        l.previewArea = Rect{treeWidth, 0, previewWidth, height - 1}
        l.statusArea = Rect{0, height - 1, width, 1}
    } else {
        // 전체 화면
        l.treeArea = Rect{0, 0, width, height - 1}
        l.statusArea = Rect{0, height - 1, width, 1}
    }
}

func (l *Layout) Render(term *Terminal, state *AppState) error {
    l.treeView.Render(term, state, l.treeArea)

    if l.showPreview {
        l.previewView.Render(term, state, l.previewArea)

        // 구분선 그리기
        for y := 0; y < l.treeArea.Height; y++ {
            term.WriteAt(l.treeArea.Width, y, "│")
        }
    }

    l.statusView.Render(term, state, l.statusArea)
    return nil
}
```

**3단계: 키바인딩 추가**

```go
// cmd/twf/main.go
case 'p':
    // 미리보기 토글
    app.layout.TogglePreview()
```

**학습 포인트**:
- 파일 I/O
- 텍스트 vs 바이너리 구분
- 화면 분할 레이아웃
- 효율적인 렌더링

**추가 도전**:
- 구문 강조 (syntax highlighting)
- 이미지 파일 정보 표시
- 스크롤 가능한 미리보기

---

### 2.2 검색 기능

**난이도**: ⭐⭐⭐ (중간)
**소요 시간**: 4-6시간

#### 기능 명세

- `/` 키로 검색 모드 진입
- 검색어 입력 (실시간 필터링)
- `Enter`로 검색 확정
- `n/N`으로 다음/이전 결과 이동
- `ESC`로 검색 취소

#### 구현 단계

**1단계: SearchMode 추가**

```go
// internal/state/view.go
const (
    Normal ViewMode = iota
    Search
    // ...
)

type ViewState struct {
    mode         ViewMode
    searchQuery  string
    searchResults []*TreeNode
    searchIndex   int
}

func (vs *ViewState) EnterSearchMode() {
    vs.mode = Search
    vs.searchQuery = ""
    vs.searchResults = nil
    vs.searchIndex = 0
}

func (vs *ViewState) UpdateSearchQuery(query string) {
    vs.searchQuery = query
}

func (vs *ViewState) NextSearchResult() {
    if len(vs.searchResults) == 0 {
        return
    }
    vs.searchIndex = (vs.searchIndex + 1) % len(vs.searchResults)
}

func (vs *ViewState) PrevSearchResult() {
    if len(vs.searchResults) == 0 {
        return
    }
    vs.searchIndex = (vs.searchIndex - 1 + len(vs.searchResults)) % len(vs.searchResults)
}
```

**2단계: 검색 로직 구현**

```go
// internal/filetree/walker.go
func (w *Walker) Search(root *TreeNode, query string) []*TreeNode {
    var results []*TreeNode

    // 대소문자 무시 검색
    lowerQuery := strings.ToLower(query)

    w.Walk(root, func(node *TreeNode) error {
        lowerName := strings.ToLower(node.Name)
        if strings.Contains(lowerName, lowerQuery) {
            results = append(results, node)
        }
        return nil
    })

    return results
}

// 정규표현식 지원
func (w *Walker) SearchRegex(root *TreeNode, pattern string) ([]*TreeNode, error) {
    re, err := regexp.Compile(pattern)
    if err != nil {
        return nil, err
    }

    var results []*TreeNode

    w.Walk(root, func(node *TreeNode) error {
        if re.MatchString(node.Name) {
            results = append(results, node)
        }
        return nil
    })

    return results, nil
}
```

**3단계: UI 업데이트**

```go
// internal/views/status-view.go
func (sv *StatusView) Render(...) error {
    viewState := state.View()

    if viewState.GetMode() == state.Search {
        // 검색 모드
        searchInfo := fmt.Sprintf("Search: %s (%d results)",
            viewState.GetSearchQuery(),
            len(viewState.GetSearchResults()))
        term.WriteAt(0, area.Y, searchInfo)

        if viewState.GetSearchIndex() >= 0 {
            resultInfo := fmt.Sprintf("[%d/%d]",
                viewState.GetSearchIndex()+1,
                len(viewState.GetSearchResults()))
            term.WriteAt(area.Width-len(resultInfo), area.Y, resultInfo)
        }
    } else {
        // 일반 모드
        // ...
    }
}
```

**4단계: 키 입력 처리**

```go
// cmd/twf/main.go
func (app *App) handleKeyPress(e terminal.KeyPressEvent) {
    viewState := app.appState.View()

    if viewState.GetMode() == state.Search {
        switch e.Key {
        case terminal.KeyEsc:
            // 검색 취소
            viewState.ExitSearchMode()

        case terminal.KeyEnter:
            // 검색 확정
            results := app.walker.Search(app.tree.GetRoot(), viewState.GetSearchQuery())
            viewState.SetSearchResults(results)
            if len(results) > 0 {
                app.appState.Cursor().SetCurrentNode(results[0])
            }

        case terminal.KeyBackspace:
            // 문자 삭제
            query := viewState.GetSearchQuery()
            if len(query) > 0 {
                viewState.UpdateSearchQuery(query[:len(query)-1])
            }

        case terminal.KeyRune:
            // 문자 추가
            query := viewState.GetSearchQuery()
            viewState.UpdateSearchQuery(query + string(e.Rune))

            // 실시간 검색
            results := app.walker.Search(app.tree.GetRoot(), viewState.GetSearchQuery())
            viewState.SetSearchResults(results)
        }
    } else {
        switch e.Key {
        case terminal.KeyRune:
            switch e.Rune {
            case '/':
                viewState.EnterSearchMode()
            case 'n':
                viewState.NextSearchResult()
                if results := viewState.GetSearchResults(); len(results) > 0 {
                    app.appState.Cursor().SetCurrentNode(results[viewState.GetSearchIndex()])
                }
            case 'N':
                viewState.PrevSearchResult()
                if results := viewState.GetSearchResults(); len(results) > 0 {
                    app.appState.Cursor().SetCurrentNode(results[viewState.GetSearchIndex()])
                }
            // ...
            }
        }
    }
}
```

**학습 포인트**:
- 실시간 검색 구현
- 정규표현식 사용
- 모드 기반 입력 처리
- 검색 결과 하이라이팅

---

### 2.3 설정 파일 시스템

**난이도**: ⭐⭐⭐ (중간)
**소요 시간**: 4-6시간

#### 설정 파일 형식

```yaml
# ~/.twfrc
ui:
  show_hidden_files: true
  color_scheme: "default"
  show_line_numbers: false
  tree_indent: 2

keybindings:
  quit: ["q", "Q"]
  toggle_select: [" "]
  search: ["/"]

behavior:
  confirm_delete: true
  follow_symlinks: false
  auto_save_bookmarks: true

bookmarks:
  h: "/home/user"
  d: "/home/user/Documents"
  p: "/home/user/Projects"
```

#### 구현

```go
// internal/config/config.go
package config

import (
    "gopkg.in/yaml.v3"
    "os"
    "path/filepath"
)

type Config struct {
    UI          UIConfig          `yaml:"ui"`
    Keybindings KeybindingsConfig `yaml:"keybindings"`
    Behavior    BehaviorConfig    `yaml:"behavior"`
    Bookmarks   map[string]string `yaml:"bookmarks"`
}

type UIConfig struct {
    ShowHiddenFiles bool   `yaml:"show_hidden_files"`
    ColorScheme     string `yaml:"color_scheme"`
    ShowLineNumbers bool   `yaml:"show_line_numbers"`
    TreeIndent      int    `yaml:"tree_indent"`
}

func LoadConfig() (*Config, error) {
    home, _ := os.UserHomeDir()
    configPath := filepath.Join(home, ".twfrc")

    data, err := os.ReadFile(configPath)
    if err != nil {
        // 기본 설정 반환
        return DefaultConfig(), nil
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

func (c *Config) Save() error {
    home, _ := os.UserHomeDir()
    configPath := filepath.Join(home, ".twfrc")

    data, err := yaml.Marshal(c)
    if err != nil {
        return err
    }

    return os.WriteFile(configPath, data, 0644)
}
```

**학습 포인트**:
- YAML 파싱
- 홈 디렉토리 처리
- 기본값 설정
- 런타임 설정 변경

---

## 🎓 Level 3: 고급 기능 (8-16시간)

### 3.1 증분 렌더링 (Incremental Rendering)

**난이도**: ⭐⭐⭐⭐ (높음)
**소요 시간**: 8-12시간

#### 문제점
현재는 매 프레임마다 전체 화면을 재렌더링합니다.
- CPU 사용량 높음
- 터미널 깜빡임 가능성

#### 해결 방안
변경된 부분만 업데이트

#### 구현 전략

**1단계: 화면 버퍼링**

```go
// internal/terminal/buffer.go
package terminal

type Cell struct {
    Rune  rune
    FG    Color
    BG    Color
    Style Style
}

type Buffer struct {
    width  int
    height int
    cells  [][]Cell
}

func NewBuffer(width, height int) *Buffer {
    cells := make([][]Cell, height)
    for i := range cells {
        cells[i] = make([]Cell, width)
    }
    return &Buffer{
        width:  width,
        height: height,
        cells:  cells,
    }
}

func (b *Buffer) Set(x, y int, r rune, fg, bg Color) {
    if x >= 0 && x < b.width && y >= 0 && y < b.height {
        b.cells[y][x] = Cell{
            Rune: r,
            FG:   fg,
            BG:   bg,
        }
    }
}

func (b *Buffer) Get(x, y int) Cell {
    if x >= 0 && x < b.width && y >= 0 && y < b.height {
        return b.cells[y][x]
    }
    return Cell{}
}
```

**2단계: Diff 계산**

```go
type Diff struct {
    X, Y  int
    Cells []Cell
}

func (b *Buffer) Diff(other *Buffer) []Diff {
    var diffs []Diff

    for y := 0; y < b.height; y++ {
        startX := -1
        var cells []Cell

        for x := 0; x < b.width; x++ {
            if b.cells[y][x] != other.cells[y][x] {
                if startX == -1 {
                    startX = x
                }
                cells = append(cells, b.cells[y][x])
            } else {
                if startX != -1 {
                    diffs = append(diffs, Diff{
                        X:     startX,
                        Y:     y,
                        Cells: cells,
                    })
                    startX = -1
                    cells = nil
                }
            }
        }

        if startX != -1 {
            diffs = append(diffs, Diff{
                X:     startX,
                Y:     y,
                Cells: cells,
            })
        }
    }

    return diffs
}
```

**3단계: 부분 렌더링**

```go
func (t *Terminal) RenderDiff(diffs []Diff) {
    for _, diff := range diffs {
        t.MoveCursorTo(diff.X, diff.Y)

        for _, cell := range diff.Cells {
            t.SetColors(cell.FG, cell.BG)
            fmt.Fprint(t.output, string(cell.Rune))
        }
    }
}
```

**학습 포인트**:
- 더블 버퍼링
- Diff 알고리즘
- 성능 최적화
- 메모리 효율성

**성능 개선 예상**:
- CPU 사용량: 50-70% 감소
- 깜빡임: 완전 제거

---

### 3.2 비동기 디렉토리 로딩

**난이도**: ⭐⭐⭐⭐ (높음)
**소요 시간**: 6-10시간

#### 문제점
대용량 디렉토리 로딩 시 UI가 멈춤

#### 해결 방안
고루틴을 사용한 백그라운드 로딩

#### 구현

```go
// internal/filetree/async.go
package filetree

type LoadRequest struct {
    Node *TreeNode
    Done chan error
}

type AsyncLoader struct {
    requests chan LoadRequest
    workers  int
}

func NewAsyncLoader(workers int) *AsyncLoader {
    loader := &AsyncLoader{
        requests: make(chan LoadRequest, 100),
        workers:  workers,
    }

    // 워커 고루틴 시작
    for i := 0; i < workers; i++ {
        go loader.worker()
    }

    return loader
}

func (al *AsyncLoader) worker() {
    for req := range al.requests {
        err := al.loadNode(req.Node)
        req.Done <- err
    }
}

func (al *AsyncLoader) loadNode(node *TreeNode) error {
    entries, err := os.ReadDir(node.Path)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        child := &TreeNode{
            Name:   entry.Name(),
            Path:   filepath.Join(node.Path, entry.Name()),
            IsDir:  entry.IsDir(),
            Parent: node,
        }

        if !entry.IsDir() {
            info, _ := entry.Info()
            child.Size = info.Size()
            child.ModTime = info.ModTime()
        }

        node.Children = append(node.Children, child)
    }

    node.Loaded = true
    return nil
}

func (al *AsyncLoader) LoadAsync(node *TreeNode) <-chan error {
    done := make(chan error, 1)
    al.requests <- LoadRequest{
        Node: node,
        Done: done,
    }
    return done
}
```

**사용**:
```go
// cmd/twf/main.go
loader := filetree.NewAsyncLoader(4) // 4 워커

// 비동기 로딩
errChan := loader.LoadAsync(node)

// 로딩 상태 표시
go func() {
    err := <-errChan
    if err != nil {
        // 에러 처리
    } else {
        // UI 업데이트
        app.Render()
    }
}()
```

**학습 포인트**:
- 고루틴 (Goroutines)
- 채널 (Channels)
- 워커 풀 패턴
- 동시성 제어

---

### 3.3 플러그인 시스템

**난이도**: ⭐⭐⭐⭐⭐ (매우 높음)
**소요 시간**: 12-16시간

#### 목표
외부 명령어를 통합 (예: fzf, ripgrep, git)

#### 설계

```go
// internal/plugins/plugin.go
package plugins

type Plugin interface {
    Name() string
    Execute(context Context) (Result, error)
}

type Context struct {
    CurrentNode   *TreeNode
    SelectedNodes []*TreeNode
    WorkingDir    string
}

type Result struct {
    Output   string
    NewNodes []*TreeNode
    Error    error
}
```

**예시 플러그인: fzf 통합**

```go
// internal/plugins/fzf.go
type FzfPlugin struct{}

func (fp *FzfPlugin) Name() string {
    return "fzf"
}

func (fp *FzfPlugin) Execute(ctx Context) (Result, error) {
    cmd := exec.Command("fzf")
    cmd.Dir = ctx.WorkingDir

    output, err := cmd.Output()
    if err != nil {
        return Result{Error: err}, err
    }

    // 선택된 파일 파싱
    selectedPath := strings.TrimSpace(string(output))

    return Result{
        Output: selectedPath,
    }, nil
}
```

**학습 포인트**:
- 인터페이스 고급 활용
- 외부 프로세스 실행
- IPC (Inter-Process Communication)

---

## 🌟 Level 4: 전문가 레벨 (16+시간)

### 4.1 Git 통합

**기능**:
- 파일 변경 상태 표시 (modified, added, deleted)
- 브랜치 정보
- Git blame
- Diff 뷰

**참고 프로젝트**: lazygit

---

### 4.2 원격 파일 시스템 지원

**기능**:
- SFTP
- FTP
- S3

**학습 내용**:
- 네트워크 프로그래밍
- 프로토콜 구현
- 비동기 I/O

---

### 4.3 이미지 미리보기 (Sixel 지원)

**기능**:
- 터미널에서 이미지 표시
- Sixel 프로토콜 사용

---

## 📚 학습 리소스

### 터미널 프로그래밍 심화
- [Termbox](https://github.com/nsf/termbox-go)
- [tcell](https://github.com/gdamore/tcell)
- [bubbletea](https://github.com/charmbracelet/bubbletea) - TUI 프레임워크

### Go 동시성
- "Concurrency in Go" by Katherine Cox-Buday
- [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)

### 성능 최적화
- [pprof](https://pkg.go.dev/net/http/pprof) - 프로파일링
- [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)

---

## ✅ 학습 체크리스트

### Level 1 (빠른 개선)
- [ ] golangci-lint 적용
- [ ] 주석 및 문서화
- [ ] 도움말 화면
- [ ] 상태바 정보 확장

### Level 2 (중급 확장)
- [ ] 파일 미리보기
- [ ] 검색 기능
- [ ] 설정 파일 시스템

### Level 3 (고급 기능)
- [ ] 증분 렌더링
- [ ] 비동기 디렉토리 로딩
- [ ] 플러그인 시스템

### Level 4 (전문가 레벨)
- [ ] Git 통합
- [ ] 원격 파일 시스템
- [ ] 이미지 미리보기

---

## 🎯 추천 학습 경로

### 초보자
```
Level 1 → Level 2.1 (파일 미리보기) → 코드 리뷰
```

### 중급자
```
Level 1 → Level 2 전체 → Level 3.1 (증분 렌더링)
```

### 고급자
```
Level 2 → Level 3 전체 → Level 4 선택
```

---

**이 가이드를 따라가면 TWF Clone을 본격적인 프로덕션 레벨 도구로 발전시킬 수 있습니다!**

*작성일: 2025-10-22*
*목적: 프로젝트 확장 및 심화 학습*
