# TWF Clone 프로젝트 학습 회고 가이드

> **목적**: 프로젝트를 통해 학습한 내용을 체계적으로 정리하고 복습하기 위한 가이드
> **대상**: 프로젝트 완료 후 스스로 복습하는 개발자
> **방법**: 각 질문에 대해 코드를 직접 보며 답변을 작성

---

## 📋 회고 진행 방법

### 1단계: 아키텍처 이해

각 레이어의 역할과 의존성 관계를 파악합니다.

### 2단계: 코드 리뷰

주요 구현 부분을 다시 읽으며 이해합니다.

### 3단계: 개념 정리

학습한 핵심 개념을 자신의 언어로 설명합니다.

### 4단계: 개선점 도출

더 나은 구현 방법을 고민합니다.

---

## 🏗️ Part 1: 아키텍처 회고

### 1.1 아키텍처 분석

**핵심 질문**:

- 왜 `terminal`, `filetree`, `state`, `views`로 패키지를 분리했을까?
- 각 레이어는 어떤 책임을 가지고 있는가?
- 패키지 간 의존성 방향은 올바른가?

**코드 위치**:

```
internal/
├── terminal/    # Infrastructure Layer
├── filetree/    # Domain Layer
├── state/       # State Management Layer
└── views/       # Presentation Layer
```

**분석할 내용**:

#### 1) 레이어 역할과 책임

각 레이어의 역할을 자신의 언어로 설명해보세요:

- terminal (Infrastructure): [역할]
- filetree (Domain): [역할]
- state (State Management): [역할]
- views (Presentation): [역할]

#### 2) 의존성 방향

패키지 간 의존성 다이어그램을 그려보고 분석하세요:

- 의존성 순환이 있는가?
- 레이어를 건너뛰는 의존성이 있는가?
- 의존성 방향이 올바른가? (고수준 → 저수준)

#### 3) 설계 결정 이유

다음 설계 결정들의 이유를 생각해보세요:

- 파일 시스템 API를 직접 호출하지 않고 추상화한 이유는?
- 지연 로딩(Lazy Loading)을 구현한 이유는?
- 상태를 별도 파일로 분리한 이유는? (cursor.go, selection.go, view.go)
- View 인터페이스를 만든 이유는?

#### 4) 확장성과 유지보수성

- 다른 터미널 라이브러리(예: tcell, termbox-go)로 교체하기 쉬운가?
- 새로운 뷰를 추가하기 쉬운가?
- 개선 가능한 부분은?

#### 5) 분리의 장점

이러한 레이어 분리가 가져온 실질적 장점은?

- 이해 용이성
- 변경 용이성
- 테스트 용이성
- 재사용성
- 병렬 개발 가능성

---

## 💻 Part 2: 코드 구현 회고

### 2.1 터미널 제어

#### Raw 모드 vs Cooked 모드

**복습할 코드**: `internal/terminal/terminal.go:55-68`

```go
func (t *Terminal) EnableRawMode() error {
    state, err := term.MakeRaw(int(t.file.Fd()))
    // ...
}
```

**회고 질문**:

1. Raw 모드는 정확히 무엇인가?
2. 왜 TUI 애플리케이션에 Raw 모드가 필요한가?
3. Raw 모드를 사용할 때 주의할 점은?

**실험해보기**:

- Raw 모드를 비활성화하면 어떻게 될까?
- 프로그램 종료 시 Raw 모드를 해제하지 않으면?

---

#### `/dev/tty` 사용

**복습할 코드**: `internal/terminal/terminal.go:35-53`

```go
file, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
```

**회고 질문**:

1. 왜 stdin/stdout 대신 `/dev/tty`를 사용했는가?
2. `/dev/tty`의 장점은?
3. Windows에서는 어떻게 해야 할까?

---

#### ANSI 이스케이프 시퀀스

**복습할 코드**: `internal/terminal/renderer.go:8-30`

**회고 질문**:

1. 자주 사용한 ANSI 시퀀스는?
   - 화면 지우기: `\x1b[2J`
   - 커서 이동: `\x1b[H`
   - 색상 변경: `\x1b[38;5;Nm`
2. 대체 화면 버퍼의 목적은?
3. 색상 시스템을 타입 안전하게 만든 이유는?

**실험해보기**:

- 다른 색상 조합 시도
- 굵게, 밑줄 등 텍스트 스타일 추가

---

### 2.2 이벤트 처리

#### 이벤트 인터페이스 설계

**복습할 코드**: `internal/terminal/event.go:7-20`

```go
type Event interface {
    EventType() EventType
}

type KeyPressEvent struct {
    Key  Key
    Rune rune
}
```

**회고 질문**:

1. 왜 구조체가 아닌 인터페이스를 사용했는가?
2. 새로운 이벤트 타입 추가가 쉬운가?
3. 타입 안전성은 어떻게 보장되는가?

---

#### 키보드 입력 파싱

**복습할 코드**: `internal/terminal/event.go:45-124`

**회고 질문**:

1. 화살표 키는 어떻게 파싱되는가? (예: `\x1b[A`)
2. UTF-8 문자는 어떻게 처리되는가?
3. Ctrl 조합 키는?

**디버깅 실험**:

```go
// 입력 데이터를 로그로 출력해보기
fmt.Fprintf(logFile, "Input: %#v\n", data)
```

---

#### 비동기 이벤트 루프

**복습할 코드**: `cmd/twf/main.go:111-130`

```go
for app.running {
    select {
    case <-sigCh:
        app.handleResize()
    default:
        event, _ := app.term.ReadEvent()
        app.handleEvent(event)
    }
}
```

**회고 질문**:

1. `select` 문은 어떻게 동작하는가?
2. `default` 케이스의 역할은?
3. 리사이즈 시그널과 키보드 입력을 동시에 처리하려면?

---

### 2.3 파일 시스템

#### 트리 구조 설계

**복습할 코드**: `internal/filetree/node.go`

```go
type TreeNode struct {
    Path     string
    Parent   *TreeNode
    Children []*TreeNode
    Expanded bool
    Loaded   bool
}
```

**회고 질문**:

1. 왜 포인터를 사용했는가? (`*TreeNode`)
2. `Expanded`와 `Loaded`의 차이는?
3. 순환 참조 가능성은?

**실험해보기**:

- 깊은 복사 vs 얕은 복사
- 메모리 사용량 측정

---

#### 지연 로딩 (Lazy Loading)

**복습할 코드**: `internal/filetree/filetree.go:76-106`

```go
func (ft *FileTreeImpl) loadChildren(node *TreeNode) error {
    if node.Loaded {
        return nil // 이미 로드됨
    }
    // 실제 디렉토리 읽기
    entries, _ := os.ReadDir(node.Path)
    // ...
    node.Loaded = true
}
```

**회고 질문**:

1. 지연 로딩의 장점은?
2. 단점이나 trade-off는?
3. 언제 로딩이 실제로 발생하는가?

**성능 측정**:

- 큰 디렉토리에서 초기 로딩 시간
- 메모리 사용량

---

#### 트리 순회 알고리즘

**복습할 코드**: `internal/filetree/walker.go`

```go
func (w *Walker) GetVisibleNodes(root *TreeNode) []*TreeNode {
    var result []*TreeNode
    w.collectVisible(root, &result, 0)
    return result
}
```

**회고 질문**:

1. DFS(깊이 우선 탐색)인가 BFS(너비 우선 탐색)인가?
2. 왜 이 방식을 선택했는가?
3. 확장되지 않은 노드는 어떻게 처리되는가?

---

### 2.4 상태 관리

#### 중앙 집중식 상태

**복습할 코드**: `internal/state/state.go`

```go
type AppState struct {
    cursor    *CursorState
    selection *SelectionState
    view      *ViewState
    config    *ConfigState
}
```

**회고 질문**:

1. 왜 모든 상태를 하나의 구조체에 모았는가?
2. 접근자 패턴 (`Cursor()`, `Selection()`)의 장점은?
3. 상태를 직접 노출하지 않는 이유는?

---

#### 북마크 시스템

**복습할 코드**: `internal/state/selection.go:35-50`

```go
type SelectionState struct {
    marks map[string]*TreeNode
}

func (s *SelectionState) SetMark(mark string, node *TreeNode) {
    s.marks[mark] = node
}
```

**회고 질문**:

1. Vim의 북마크와 동일한 방식인가?
2. `map[string]*TreeNode`를 사용한 이유는?
3. 북마크가 가리키는 노드가 삭제되면?

---

### 2.5 UI 렌더링

#### View 인터페이스

**복습할 코드**: `internal/views/view.go`

```go
type View interface {
    Render(term *terminal.Terminal, state *state.AppState, area Rect) error
    GetMinSize() (int, int)
}
```

**회고 질문**:

1. 왜 인터페이스로 정의했는가?
2. `Render` 메서드에 필요한 인자는?
3. 새로운 뷰를 추가하려면?

---

#### 레이아웃 시스템

**복습할 코드**: `internal/views/layout.go`

```go
type Rect struct {
    X, Y, Width, Height int
}

func (l *Layout) SetSize(width, height int) {
    l.treeArea = Rect{0, 0, width, height - 1}
    l.statusArea = Rect{0, height - 1, width, 1}
}
```

**회고 질문**:

1. `Rect`의 역할은?
2. 터미널 크기가 변경되면?
3. 3개 이상의 뷰를 표시하려면?

---

## 🎓 Part 3: 핵심 개념 정리

### 3.1 터미널 프로그래밍

**스스로 설명해보기** (답을 보지 않고):

#### Raw 모드란?

```
내 답변:




```

#### ANSI 이스케이프 시퀀스란?

```
내 답변:




```

#### 대체 화면 버퍼란?

```
내 답변:




```

---

### 3.2 이벤트 기반 프로그래밍

#### 이벤트 루프 패턴

```
내 답변:




```

#### 시그널 처리

```
내 답변:




```

---

### 3.3 자료구조

#### 트리 구조의 특징

```
내 답변:




```

#### 지연 로딩 패턴

```
내 답변:




```

---

### 3.4 소프트웨어 아키텍처

#### MVC 패턴이란?

```
내 답변:
- Model:
- View:
- Controller:

이 프로젝트에서는:
- Model:
- View:
- Controller:
```

#### 관심사의 분리 (Separation of Concerns)

```
내 답변:




```

---

## 📝 Part 4: 회고 정리

### 4.1 가장 어려웠던 부분

**1순위**:

```
주제:
어려웠던 이유:
어떻게 해결했는가:
배운 점:
```

**2순위**:

```
주제:
어려웠던 이유:
어떻게 해결했는가:
배운 점:
```

---

### 4.2 가장 흥미로웠던 부분

**1순위**:

```
주제:
흥미로웠던 이유:
새롭게 알게 된 것:
응용 가능성:
```

**2순위**:

```
주제:
흥미로웠던 이유:
새롭게 알게 된 것:
응용 가능성:
```

---

### 4.3 아쉬웠던 부분

**코드 품질**:

```
더 잘할 수 있었던 부분:
다음에는 어떻게 할 것인가:
```

**설계**:

```
더 잘할 수 있었던 부분:
다음에는 어떻게 할 것인가:
```

---

### 4.4 다음 프로젝트에 적용할 점

**기술적 측면**:

1.
2.
3.

**프로세스 측면**:

1.
2.
3.

---

## 🚀 Part 5: 다음 단계

### 5.1 추가 학습 필요 분야

**아직 완전히 이해하지 못한 부분**:

- [ ]
- [ ]
- [ ]

**더 깊이 공부하고 싶은 부분**:

- [ ]
- [ ]
- [ ]

---

### 5.2 유사 프로젝트 아이디어

이 프로젝트에서 배운 것을 응용할 수 있는 프로젝트들:

1. **TUI 텍스트 에디터**
   - 난이도: 높음
   - 적용 기술: 터미널 제어, 이벤트 처리
   - 추가 학습 필요: 텍스트 버퍼 관리, 구문 강조

2. **TUI 시스템 모니터**
   - 난이도: 중간
   - 적용 기술: 레이아웃, 실시간 업데이트
   - 추가 학습 필요: 시스템 API, 차트 렌더링

3. **TUI Git 클라이언트**
   - 난이도: 높음
   - 적용 기술: 파일 트리, 상태 관리
   - 추가 학습 필요: Git API, 색상 diff

---

## ✅ 회고 완료 체크리스트

- [x] Part 1: 아키텍처 이해 완료
- [x] Part 2: 코드 구현 복습 완료
- [x] Part 3: 핵심 개념 정리 완료
- [ ] Part 4: 회고 정리 완료
- [ ] Part 5: 다음 단계 계획 완료

---

## 📚 참고 자료

### 추가 학습 자료

**터미널 프로그래밍**:

- [ANSI Escape Codes Wiki](https://en.wikipedia.org/wiki/ANSI_escape_code)
- [termbox-go](https://github.com/nsf/termbox-go) - 참고 라이브러리
- [The Linux Programming Interface](https://man7.org/tlpi/) - Chapter 62: Terminals

**Go 언어**:

- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/)

**소프트웨어 설계**:

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Design Patterns in Go](https://refactoring.guru/design-patterns/go)

---

**이 회고 가이드를 완료하면 프로젝트에 대한 깊은 이해를 얻을 수 있습니다!**

_작성일: 2025-10-22_
_목적: 학습 내용 체계적 정리 및 복습_
