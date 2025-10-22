# State Management Layer 학습 노트

> **작성일**: 2025-10-23
> **목적**: State Management Layer의 개념, 역사, 필요성에 대한 이해 정리

---

## 📌 핵심 질문

**"State Management Layer는 일반적으로 사용하는 레이어인가?"**

### 결론

**"State Management Layer"는 표준 용어는 아니지만, 인터랙티브 애플리케이션에서 나타나는 현대적 패턴입니다.**

---

## 🏗️ 전통적인 레이어 아키텍처

### 1. 클래식 3-Tier Architecture

```
Presentation Layer (UI)
    ↓
Business Logic Layer (Domain)
    ↓
Data Access Layer (Persistence)
```

### 2. Clean Architecture / Hexagonal Architecture

```
Presentation (Adapters)
    ↓
Application (Use Cases)
    ↓
Domain (Entities)
    ↓
Infrastructure (Frameworks & Drivers)
```

### 3. MVC Pattern

```
View (Presentation)
    ↓
Controller (Application Logic)
    ↓
Model (Domain + Data)
```

---

## 🤔 전통적 아키텍처에서 State는 어디에?

전통적으로 State는 **별도 레이어가 아니라** 다른 레이어에 분산되어 있었습니다:

- **Domain State**: Model/Entity에 포함 (비즈니스 상태)
- **UI State**: View/Presentation에 포함 (표시 상태)
- **Application State**: Controller/Use Case에 포함 (흐름 상태)

### 전통적 웹 앱 예시

```go
// Domain State는 Model에
type User struct {
    ID   int
    Name string
}

// UI State는 View/Controller에
type UserListController struct {
    currentPage int      // 페이지 상태
    sortBy      string   // 정렬 상태
    users       []User   // 표시할 데이터
}
```

---

## 💡 왜 State Layer를 분리하게 되었나?

### 핵심 통찰

> **"유저 인터랙션을 다룰 일이 많은 애플리케이션에서 나타나는 특징"**

현대 프론트엔드/인터랙티브 앱에서는 **State가 복잡해지면서** 별도 레이어로 분리하는 추세가 생겼습니다.

---

## 📊 애플리케이션 유형별 State 복잡도

### 1. 배치 프로그램 (State 거의 없음)

```
입력 → 처리 → 출력
```

```go
// State가 거의 필요 없음
func main() {
    data := readFile()      // 입력
    result := process(data) // 처리
    writeFile(result)       // 출력
}
```

**특징**: 일방향 흐름, 상태 관리 불필요

---

### 2. 전통적인 웹 서버 (State 적음)

```
Request → Process → Response
(각 요청이 독립적)
```

```go
// 요청마다 독립적, State는 DB에
func HandleGetUser(w http.ResponseWriter, r *http.Request) {
    user := db.GetUser(id)  // DB가 State의 주인
    json.NewEncoder(w).Encode(user)
}
```

**특징**: Stateless, 영속성은 DB가 담당

---

### 3. 인터랙티브 TUI/GUI ⭐ (State 많음)

```
Event Loop:
  User Input → Update State → Render → User Input → ...
  (계속 반복, 상태가 누적되고 변경됨)
```

```go
// State가 핵심!
type AppState struct {
    cursor       Position         // 사용자가 어디를 보고 있나?
    selection    []Node          // 무엇을 선택했나?
    scrollOffset int             // 화면이 어디까지 스크롤됐나?
    mode         Mode            // 어떤 모드인가? (Normal/Search/Help)
    inputBuffer  string          // 입력 중인 텍스트는?
    history      []Command       // Undo를 위한 이력
    bookmarks    map[string]Node // 북마크는?
}

// 이벤트마다 상태가 변경됨
for event := range events {
    state = state.Update(event) // State 변경
    view.Render(state)          // State 기반 렌더링
}
```

**특징**:
- 상태가 **계속 변경**되고 **누적**됨
- 사용자 인터랙션이 상태 변경의 주요 원인
- 렌더링이 상태에 전적으로 의존

---

## 🚫 State Layer 없이 개발하면?

### Without State Layer ❌

```go
// Domain에 UI 상태가 섞임
type TreeNode struct {
    Path       string
    Children   []*TreeNode
    IsSelected bool        // ❌ 이건 Domain 지식이 아님
    IsCursor   bool        // ❌ UI 관심사
    IsExpanded bool        // ❌ 보기 설정
}

// View가 모든 상태를 관리
type TreeView struct {
    nodes         []*TreeNode
    cursorPos     int
    scrollOffset  int
    searchText    string
    sortBy        SortType
    // ❌ View가 너무 많은 책임
}
```

**문제점**:
- Domain이 오염됨 (파일 시스템 구조에 UI 상태가 섞임)
- View가 비대해짐 (렌더링 + 상태 관리 + 이벤트 처리)
- 테스트 어려움 (UI 없이 상태 로직 테스트 불가)

### With State Layer ✅

```go
// Domain은 순수하게
type TreeNode struct {
    Path     string
    Children []*TreeNode
    // Domain 지식만!
}

// State는 "무엇을 어떻게 보여줄까"
type AppState struct {
    cursor    *CursorState    // 어떤 노드를 보고 있나?
    selection *SelectionState // 무엇이 선택됐나?
    view      *ViewState      // 어떻게 보여줄까?
}

// View는 렌더링만
type TreeView struct {
    // 렌더링 로직만!
}

func (tv *TreeView) Render(state *AppState) {
    // state를 읽어서 그리기만 함
}
```

**장점**:
- ✅ 관심사의 분리
- ✅ 각 레이어가 단일 책임
- ✅ 테스트 용이
- ✅ Domain 순수성 유지
- ✅ View는 렌더링에만 집중

---

## 🎮 다른 인터랙티브 앱들의 예시

### 예시 1: 텍스트 에디터 (Vim, VSCode)

```go
type EditorState struct {
    buffer       [][]rune      // 텍스트 내용
    cursor       CursorPos     // 커서 위치
    selection    Range         // 선택 영역
    mode         Mode          // Normal/Insert/Visual
    viewport     Viewport      // 보이는 영역
    undoStack    []Change      // Undo 이력
    searchQuery  string        // 검색어
    highlights   []Range       // 강조 표시
}
```

### 예시 2: 게임

```go
type GameState struct {
    player       Player        // 플레이어 상태
    enemies      []Enemy       // 적들
    camera       Camera        // 카메라 위치
    input        InputState    // 입력 버퍼
    score        int           // 점수
    paused       bool          // 일시정지
    menu         MenuState     // 메뉴 상태
}
```

### 예시 3: IDE

```go
type IDEState struct {
    openFiles    []File        // 열린 파일들
    activeFile   *File         // 현재 파일
    panels       PanelLayout   // 패널 배치
    terminal     TerminalState // 터미널 상태
    debugger     DebugState    // 디버거 상태
    sidebarMode  SidebarMode   // 사이드바 모드
}
```

---

## 📚 현대 프론트엔드 프레임워크의 영향

### 1. React 생태계

```
View (React Components)
    ↓
State Management (Redux/MobX/Zustand)
    ↓
Domain Logic
```

**Redux의 핵심 개념**:
- Single Source of Truth (단일 진실 공급원)
- State는 읽기 전용
- Pure Function으로만 변경

### 2. Flutter

```
UI Layer (Widgets)
    ↓
State Management (Provider/BLoC/Riverpod)
    ↓
Business Logic
```

**BLoC Pattern**:
- Business Logic Component
- UI와 비즈니스 로직 분리
- Stream 기반 상태 관리

### 3. Elm Architecture

```
Model (State)
    ↓
Update (State Transition)
    ↓
View (Rendering)
```

이 패턴이 Redux, BLoC 등에 영향을 줌

---

## 📅 패턴의 역사적 맥락

### Timeline

```
1970s: MVC 패턴
  └─ State는 Model에 포함

1990s: 3-Tier Architecture
  └─ State는 각 레이어에 분산

2000s: AJAX / SPA 시대
  └─ 클라이언트 State가 복잡해짐

2010s: React + Redux (2015)
  └─ State Management의 독립
  └─ "Single Source of Truth" 개념

2015+: Flutter, Elm Architecture
  └─ State Layer가 명시적으로 분리
  └─ Functional Reactive Programming

현재: 모든 인터랙티브 앱에서 State 관리가 핵심 과제
```

---

## 🏷️ 업계에서 사용하는 다양한 용어

| 용어 | 사용처 | 의미 |
|------|--------|------|
| **State Management Layer** | 현대 프론트엔드 | 상태 관리 전용 |
| **Application Layer** | DDD, Clean Arch | Use Case + 상태 흐름 |
| **Service Layer** | 전통적 백엔드 | 비즈니스 로직 조율 |
| **Controller Layer** | MVC | 요청 처리 + 상태 변경 |
| **Store** | Redux/Vuex | 중앙 집중식 상태 |
| **ViewModel** | MVVM | View를 위한 상태 |

---

## 🎯 TWF 프로젝트에서의 State Layer

### 당신의 State Layer는?

**"Application Layer (Clean Architecture)" + "Store (Redux)"의 혼합형**

```
전통적 Application Layer 역할:
- Use Case 조율 ✓ (커서 이동, 선택 등)
- Domain과 Presentation 연결 ✓

Redux Store 역할:
- 중앙 집중식 상태 관리 ✓
- 상태 변경의 단일 진실 공급원 ✓
```

### 실제 구조

```go
// TWF의 AppState
type AppState struct {
    cursor    *CursorState    // 커서 위치, 이동
    selection *SelectionState // 선택, 북마크
    view      *ViewState      // 스크롤, 정렬, 필터, 모드
    config    *ConfigState    // 설정
}
```

### 분리한 이유

```go
// filetree (Domain): 파일 시스템 구조
type TreeNode struct {
    Path     string
    Children []*TreeNode
    IsDir    bool
}

// state: UI와 관련된 "어떻게 보여줄까" 상태
type AppState struct {
    cursor    *CursorState    // 어느 노드를 가리키고 있나?
    selection *SelectionState // 어떤 노드들이 선택되었나?
    view      *ViewState      // 스크롤, 정렬, 필터는?
}
```

**만약 State Layer가 없었다면?**

- **Option 1**: State를 filetree에 → Domain 오염
- **Option 2**: State를 views에 → View가 비대해짐
- **Option 3**: State를 별도 레이어로 (현재 선택) ✅

---

## ✅ State Layer가 필요한 애플리케이션

### 특징

1. **Event Loop가 있는 앱**
   - TUI (터미널 앱)
   - GUI (데스크톱/모바일 앱)
   - 게임
   - IDE

2. **상태가 계속 변경되는 앱**
   - 사용자 입력에 반응
   - 상태 간 전이(Transition)
   - 이력 관리 (Undo/Redo)

3. **복잡한 UI 상태를 가진 앱**
   - 여러 모드 (Normal/Search/Edit...)
   - 다중 선택, 북마크
   - 스크롤, 필터, 정렬

---

## ❌ State Layer가 덜 필요한 경우

### 1. Stateless 서버

- RESTful API 서버
- 람다 함수
- 배치 프로세서

### 2. 단순한 CLI 도구

- `cat`, `grep`, `ls`
- 입력 → 처리 → 출력만
- 일회성 실행

---

## 💭 회고에 추가할 내용 제안

```markdown
### State Management Layer에 대한 보충 설명

"State Management Layer"는 전통적인 3-Tier Architecture에는
명시적으로 존재하지 않는 레이어입니다.

이 레이어는 **인터랙티브 애플리케이션**(TUI, GUI, 게임 등)에서
복잡한 사용자 인터랙션 상태를 관리하기 위해 분리되었습니다.

#### 왜 필요했나?

- 전통적으로 State는 Domain(Model)이나 View에 포함
- 하지만 TUI는 커서, 선택, 스크롤, 모드 등 UI 상태가 복잡
- 이를 Domain에 넣으면 Domain이 오염되고
- View에 넣으면 View가 비대해짐
- 따라서 별도 레이어로 분리

#### 유사한 개념

- Redux의 Store
- Flutter의 BLoC
- Clean Architecture의 Application Layer
- MVVM의 ViewModel
```

---

## 🎓 핵심 교훈

### 1. State Management Layer는 표준이 아니다

하지만 **인터랙티브 앱에서 자연스럽게 나타나는 필요**에 의해 생겨난 패턴입니다.

### 2. 이름보다 중요한 것은 역할

- "State Management Layer"
- "Application Layer"
- "Store"

이름은 다르지만 본질은 같습니다:
> **"Domain과 Presentation 사이에서 UI 상태를 관리하는 레이어"**

### 3. 복잡도에 따라 선택

- **단순한 앱**: State를 View나 Controller에 포함해도 OK
- **복잡한 인터랙티브 앱**: State Layer 분리가 필수

### 4. TWF 프로젝트의 선택은 적절했다

TUI는 전형적인 인터랙티브 앱이므로, State Layer를 분리한 것은 **매우 합리적인 설계 결정**이었습니다.

---

## 📖 참고 자료

### 아키텍처

1. **Clean Architecture** (Robert C. Martin)
   - Application Layer 설명
   - https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

2. **Domain-Driven Design** (Eric Evans)
   - Application Service 개념

3. **Hexagonal Architecture** (Alistair Cockburn)
   - Ports and Adapters 패턴

### State Management

4. **Redux Documentation**
   - State Management 철학
   - https://redux.js.org/

5. **The Elm Architecture**
   - Model-Update-View 패턴
   - https://guide.elm-lang.org/architecture/

6. **Flutter BLoC Pattern**
   - Business Logic Component
   - https://bloclibrary.dev/

### 실무 사례

7. **React State Management**
   - Redux, MobX, Zustand 비교

8. **Game Programming Patterns** (Robert Nystrom)
   - Game State 관리 패턴
   - https://gameprogrammingpatterns.com/

---

## 🤔 추가 학습 질문

### 스스로에게 물어볼 질문들

1. **내 프로젝트의 State는 어떻게 변경되는가?**
   - 이벤트 흐름 추적해보기
   - 상태 전이 다이어그램 그려보기

2. **State를 불변(Immutable)으로 관리해야 하는가?**
   - Go에서는 어떻게 구현할 수 있는가?
   - 장단점은?

3. **Undo/Redo를 구현한다면?**
   - State 이력을 어떻게 저장할 것인가?
   - Command Pattern을 적용할 수 있는가?

4. **상태 동기화 문제**
   - State가 여러 곳에 분산되면 어떻게 동기화하는가?
   - Single Source of Truth를 어떻게 보장하는가?

5. **성능 최적화**
   - State 변경 시 전체를 re-render해야 하는가?
   - 부분 업데이트는 어떻게 구현하는가?

---

## 📝 실습 아이디어

### 1. State 변경 로깅

```go
// State 변경을 추적하는 Middleware 추가
func (as *AppState) logStateChange(action string, before, after interface{}) {
    log.Printf("Action: %s, Before: %+v, After: %+v", action, before, after)
}
```

### 2. State History 구현

```go
type StateHistory struct {
    states    []*AppState
    current   int
    maxSize   int
}

func (sh *StateHistory) Undo() *AppState { /* ... */ }
func (sh *StateHistory) Redo() *AppState { /* ... */ }
```

### 3. State 스냅샷

```go
// 현재 State를 JSON으로 저장/복원
func (as *AppState) Save(filename string) error { /* ... */ }
func LoadState(filename string) (*AppState, error) { /* ... */ }
```

---

**이 문서를 통해 State Management Layer의 개념과 필요성을 이해했습니다.**

**핵심**: 인터랙티브 앱에서는 복잡한 UI 상태를 관리하기 위해 State Layer를 분리하는 것이 자연스러운 선택입니다.
