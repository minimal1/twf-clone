# TWF Clone 아키텍처 설계 문서

## 📋 개요

이 문서는 TWF Clone 프로젝트의 소프트웨어 아키텍처와 설계 원칙을 설명합니다. 원본 twf 프로젝트의 구조를 참조하면서도 학습 목적에 맞게 단순화하고 명확화한 설계입니다.

## 🏛️ 전체 아키텍처

### 아키텍처 패턴: MVC + Component

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Controller    │───▶│     Model       │───▶│      View       │
│  (Event Loop)   │    │   (State +      │    │  (UI Render)    │
│                 │    │   FileTree)     │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         ▲                       ▲                       │
         │              ┌─────────────────┐              │
         │              │   Terminal      │              │
         └──────────────│   Interface     │◀─────────────┘
                        │                 │
                        └─────────────────┘
```

## 📦 패키지 구조

```
twf-clone/
├── cmd/
│   └── twf/                    # 실행 가능한 바이너리
│       ├── main.go            # 진입점
│       └── app.go             # 애플리케이션 오케스트레이터
├── internal/
│   ├── config/                # 설정 관리
│   │   ├── config.go          # 설정 구조체 및 파싱
│   │   └── defaults.go        # 기본값 정의
│   ├── filetree/              # 파일 시스템 모델
│   │   ├── filetree.go        # 파일 트리 구조체
│   │   ├── node.go            # 트리 노드 조작
│   │   └── walker.go          # 트리 순회 알고리즘
│   ├── state/                 # 애플리케이션 상태
│   │   ├── state.go           # 상태 관리
│   │   ├── cursor.go          # 커서 위치 관리
│   │   └── selection.go       # 선택 상태 관리
│   ├── terminal/              # 터미널 인터페이스
│   │   ├── terminal.go        # 터미널 제어
│   │   ├── event.go           # 이벤트 정의 및 파싱
│   │   ├── renderer.go        # 화면 렌더링
│   │   └── input.go           # 입력 처리
│   └── views/                 # UI 컴포넌트
│       ├── tree_view.go       # 파일 트리 뷰
│       ├── preview_view.go    # 미리보기 뷰
│       ├── status_view.go     # 상태바 뷰
│       └── view.go            # 뷰 인터페이스
└── docs/                      # 문서
```

## 🔧 핵심 컴포넌트

### 1. Terminal Interface Layer

**책임**: 저수준 터미널 제어와 사용자 입력 처리

```go
type Terminal interface {
    // 터미널 초기화 및 정리
    Initialize() error
    Cleanup() error

    // 화면 제어
    Clear() error
    MoveCursor(x, y int) error
    SetSize(rows, cols int)

    // 입력 처리
    ReadEvent() (Event, error)

    // 출력
    Write(data []byte) error
    Flush() error
}
```

**핵심 기능**:
- Raw 모드 전환 및 복구
- 논블로킹 입력 읽기
- ANSI 이스케이프 시퀀스 생성
- 터미널 크기 감지 및 리사이즈 처리

### 2. Event System

**책임**: 사용자 입력을 의미있는 이벤트로 변환

```go
type EventType int

const (
    KeyPress EventType = iota
    Resize
    Mouse
    Signal
)

type Event struct {
    Type      EventType
    Key       Key
    Modifiers Modifier
    Data      interface{}
}

type Key int

const (
    KeyEnter Key = iota
    KeyEscape
    KeyArrowUp
    KeyArrowDown
    // ... 기타 키들
)
```

**핵심 기능**:
- 키 조합 인식 (Ctrl+C, Alt+키 등)
- 방향키 및 기능키 처리
- 터미널 리사이즈 이벤트
- 마우스 이벤트 (선택적)

### 3. FileTree Model

**책임**: 파일 시스템을 메모리 내 트리 구조로 표현

```go
type FileTree struct {
    // 기본 정보
    Path     string
    Name     string
    IsDir    bool
    Size     int64
    ModTime  time.Time

    // 트리 구조
    Parent   *FileTree
    Children []*FileTree

    // 상태
    Expanded bool
    Loaded   bool
}

type FileTreeInterface interface {
    // 트리 조작
    LoadChildren() error
    AddChild(child *FileTree)
    RemoveChild(name string) bool

    // 탐색
    Find(predicate func(*FileTree) bool) *FileTree
    Walk(visitor func(*FileTree) error) error

    // 상태
    ToggleExpansion() error
    IsLeaf() bool
    Depth() int
}
```

**핵심 기능**:
- 지연 로딩 (디렉토리 확장 시에만 자식 로드)
- 트리 순회 알고리즘 (DFS, BFS)
- 경로 기반 검색
- 메모리 효율적인 구조

### 4. Application State

**책임**: 애플리케이션의 현재 상태 관리

```go
type State struct {
    // 파일 시스템
    RootTree    *filetree.FileTree
    CurrentPath string

    // UI 상태
    CursorPos    int
    ScrollOffset int
    ViewMode     ViewMode

    // 선택 및 필터
    Selection    []*filetree.FileTree
    Filter       string
    ShowHidden   bool

    // 뷰 상태
    TreeExpanded  bool
    PreviewActive bool
}

type StateManager interface {
    // 내비게이션
    MoveCursor(direction Direction) error
    Navigate(path string) error
    GoToParent() error

    // 선택 및 액션
    ToggleSelection() error
    ExpandCurrent() error
    CollapseCurrent() error

    // 필터 및 설정
    SetFilter(pattern string) error
    ToggleHidden() error
}
```

### 5. View Layer

**책임**: 상태를 시각적으로 렌더링

```go
type View interface {
    // 렌더링
    Render(ctx RenderContext) ([]string, error)

    // 레이아웃
    GetRequiredSize() (width, height int)
    SetBounds(x, y, width, height int)

    // 상호작용
    CanFocus() bool
    HandleEvent(event Event) (handled bool, err error)
}

type TreeView struct {
    state       *State
    bounds      Rectangle
    styleConfig *StyleConfig
}
```

**렌더링 파이프라인**:
1. 상태에서 표시할 데이터 추출
2. 가시 영역 계산 (스크롤 오프셋 적용)
3. 각 라인을 문자열로 변환
4. 스타일 적용 (색상, 하이라이팅)
5. 터미널 출력 형식으로 변환

## 🔄 데이터 플로우

### 1. 이벤트 처리 플로우

```
User Input → Terminal → Event → Controller → State Update → View Render → Terminal Output
```

상세 과정:
1. 사용자가 키를 누름
2. Terminal이 raw 바이트를 읽음
3. Event 파서가 의미있는 Event 객체로 변환
4. Controller가 이벤트 타입에 따라 적절한 액션 선택
5. State가 업데이트됨
6. View가 새로운 상태를 기반으로 렌더링
7. Terminal이 화면에 출력

### 2. 파일 시스템 로딩 플로우

```
Path Request → FileTree.LoadChildren() → OS API → FileInfo → TreeNode Creation
```

### 3. 렌더링 플로우

```
State → View.Render() → Line Generation → Style Application → Terminal Output
```

## 🎨 디자인 패턴

### 1. Observer Pattern (상태 변경 알림)

```go
type StateObserver interface {
    OnStateChanged(oldState, newState *State)
}

type StateManager struct {
    observers []StateObserver
}

func (sm *StateManager) NotifyObservers(oldState, newState *State) {
    for _, observer := range sm.observers {
        observer.OnStateChanged(oldState, newState)
    }
}
```

### 2. Command Pattern (액션 처리)

```go
type Command interface {
    Execute(state *State) error
    Undo(state *State) error
}

type MoveCommand struct {
    direction Direction
    oldPos    int
}

func (mc *MoveCommand) Execute(state *State) error {
    mc.oldPos = state.CursorPos
    return state.MoveCursor(mc.direction)
}
```

### 3. Strategy Pattern (렌더링 전략)

```go
type RenderStrategy interface {
    Render(tree *FileTree, bounds Rectangle) ([]string, error)
}

type TreeRenderStrategy struct{}
type ListRenderStrategy struct{}
type IconRenderStrategy struct{}
```

## 🚀 성능 고려사항

### 1. 메모리 관리
- **지연 로딩**: 필요한 디렉토리만 메모리에 로드
- **가비지 컬렉션**: 사용하지 않는 서브트리 해제
- **메모리 풀**: 빈번한 할당/해제되는 객체 재사용

### 2. 렌더링 최적화
- **더티 체킹**: 변경된 부분만 다시 렌더링
- **가상 스크롤링**: 보이는 영역만 렌더링
- **버퍼링**: 출력을 모아서 한 번에 터미널에 전송

### 3. 파일 시스템 접근
- **비동기 로딩**: 대용량 디렉토리를 백그라운드에서 로드
- **캐싱**: 파일 정보 캐시로 중복 시스템 콜 방지
- **워커 풀**: 동시 파일 접근 제한

## 🔧 확장성 설계

### 1. 플러그인 시스템
```go
type Plugin interface {
    Name() string
    Version() string
    Initialize(ctx PluginContext) error
    HandleEvent(event Event) (handled bool, err error)
}
```

### 2. 테마 시스템
```go
type Theme interface {
    GetColor(element string) Color
    GetStyle(element string) Style
}
```

### 3. 키바인딩 시스템
```go
type KeyBinding struct {
    Key     KeyCombination
    Command string
    Args    map[string]interface{}
}
```

이 아키텍처는 모듈화, 테스트 용이성, 확장성을 중심으로 설계되었으며, 각 컴포넌트가 명확한 책임을 가지고 있어 개별적으로 개발하고 테스트할 수 있습니다.