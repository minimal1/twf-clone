# TWF Clone 학습 가이드

이 가이드는 TUI 파일 브라우저를 단계별로 클론 코딩하면서 터미널 프로그래밍과 Go 언어를 학습할 수 있도록 구성되었습니다.

## 🎯 학습 목표

- **터미널 프로그래밍**: raw 모드, 이벤트 처리, 화면 제어
- **Go 시스템 프로그래밍**: 파일 시스템, 동시성, 메모리 관리
- **소프트웨어 아키텍처**: MVC 패턴, 모듈화, 인터페이스 설계
- **사용자 인터페이스**: TUI 디자인, 상호작용, 사용성

## 📚 사전 준비 학습

### 1. Go 기본 개념 복습
- 구조체와 메서드
- 인터페이스와 임베딩
- 에러 처리 패턴
- 패키지와 모듈 시스템

### 2. 터미널 기본 이해
- ANSI 이스케이프 시퀀스
- 터미널 입출력 모델
- Raw 모드와 Cooked 모드
- 신호 처리 (SIGINT, SIGTERM)

## 🏗️ 단계별 구현 가이드

## 1단계: 프로젝트 기초 설정 (1-2시간)

### 목표
- Go 모듈 초기화
- 기본 프로젝트 구조 이해
- 첫 번째 실행 가능한 프로그램 작성

### 구현 내용

#### 1.1 Go 모듈 초기화
```bash
cd twf-clone
go mod init twf-clone
```

#### 1.2 기본 main.go 작성
```go
// cmd/twf/main.go
package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("TWF Clone - TUI File Browser")

    // 명령행 인자 확인
    args := os.Args[1:]
    if len(args) > 0 {
        fmt.Printf("Starting path: %s\n", args[0])
    } else {
        fmt.Println("Starting path: current directory")
    }
}
```

### 학습 포인트
- Go 모듈 시스템의 이해
- 명령행 인자 처리
- 패키지 구조의 중요성

### 실습 과제
1. 프로그램을 실행해보고 다양한 인자로 테스트
2. 현재 작업 디렉토리를 출력하는 기능 추가
3. 파일이 존재하는지 확인하는 로직 추가

---

## 2단계: 터미널 기본 제어 (3-4시간)

### 목표
- 터미널 raw 모드 설정
- 키보드 입력 처리
- 기본적인 화면 제어

### 구현 내용

#### 2.1 터미널 구조체 정의
```go
// internal/terminal/terminal.go
package terminal

import (
    "os"
    "golang.org/x/crypto/ssh/terminal"
)

type Terminal struct {
    input         *os.File
    output        *os.File
    originalState *terminal.State
    rows          int
    cols          int
}

func NewTerminal() (*Terminal, error) {
    // 터미널 초기화 로직
}

func (t *Terminal) EnterRawMode() error {
    // Raw 모드 진입
}

func (t *Terminal) ExitRawMode() error {
    // 원래 상태로 복구
}
```

#### 2.2 키보드 이벤트 처리
```go
// internal/terminal/event.go
package terminal

type EventType int

const (
    KeyPress EventType = iota
    Resize
    Error
)

type Event struct {
    Type EventType
    Key  rune
    Raw  []byte
}

func (t *Terminal) ReadEvent() (Event, error) {
    // 키보드 입력 읽기 및 이벤트 변환
}
```

### 학습 포인트
- 터미널의 cooked 모드 vs raw 모드
- 논블로킹 I/O의 필요성
- 시스템 호출과 Go의 추상화

### 실습 과제
1. 'q' 키를 누르면 종료되는 루프 구현
2. 화살표 키 입력 감지 및 출력
3. Ctrl+C 시그널 처리

---

## 3단계: 파일 시스템 인터페이스 (4-5시간)

### 목표
- 파일과 디렉토리 정보 읽기
- 트리 자료구조 구현
- 재귀적 디렉토리 탐색

### 구현 내용

#### 3.1 FileTree 구조체
```go
// internal/filetree/filetree.go
package filetree

import (
    "os"
    "path/filepath"
)

type FileTree struct {
    Path     string
    Name     string
    IsDir    bool
    Size     int64
    Children []*FileTree
    Parent   *FileTree
    Expanded bool
}

func NewFileTree(path string) (*FileTree, error) {
    // 파일 정보를 읽어서 FileTree 생성
}

func (ft *FileTree) LoadChildren() error {
    // 하위 디렉토리와 파일들을 로드
}
```

#### 3.2 트리 탐색 메서드
```go
func (ft *FileTree) Walk(fn func(*FileTree) error) error {
    // 트리를 순회하며 함수 적용
}

func (ft *FileTree) Find(name string) *FileTree {
    // 특정 이름의 파일/디렉토리 찾기
}
```

### 학습 포인트
- 파일 시스템 API (os.Stat, os.ReadDir)
- 트리 자료구조와 재귀 알고리즘
- 지연 로딩 (lazy loading) 패턴

### 실습 과제
1. 디렉토리 크기 계산 기능
2. 파일 타입별 분류 (확장자 기반)
3. 숨김 파일 필터링 옵션

---

## 4단계: 기본 UI 렌더링 (5-6시간)

### 목표
- 파일 목록을 터미널에 출력
- 커서 이동과 선택 상태 표시
- 스크롤 기능 구현

### 구현 내용

#### 4.1 뷰 인터페이스
```go
// internal/views/view.go
package views

type View interface {
    Render() ([]string, error)
    HandleEvent(event terminal.Event) error
    GetCursor() (int, int)
}
```

#### 4.2 트리 뷰 구현
```go
// internal/views/tree_view.go
package views

type TreeView struct {
    fileTree *filetree.FileTree
    cursor   int
    scroll   int
    height   int
    width    int
}

func (tv *TreeView) Render() ([]string, error) {
    // 트리를 문자열 라인들로 변환
}

func (tv *TreeView) MoveCursor(delta int) {
    // 커서 이동 및 스크롤 조정
}
```

### 학습 포인트
- 인터페이스 기반 설계
- 터미널 좌표계와 렌더링
- 뷰포트와 스크롤링 개념

### 실습 과제
1. 트리 들여쓰기와 연결선 표시
2. 파일과 디렉토리 아이콘 추가
3. 색상 코딩 (디렉토리, 실행파일 등)

---

## 5단계: 상태 관리 (3-4시간)

### 목표
- 애플리케이션 전체 상태 관리
- 이벤트와 상태 변경의 분리
- 커서 위치와 선택 상태 추적

### 구현 내용

#### 5.1 상태 구조체
```go
// internal/state/state.go
package state

type State struct {
    RootTree     *filetree.FileTree
    CurrentNode  *filetree.FileTree
    CursorPos    int
    ScrollOffset int
    Selection    []*filetree.FileTree
}

func NewState(rootPath string) (*State, error) {
    // 초기 상태 생성
}
```

#### 5.2 상태 변경 메서드
```go
func (s *State) MoveCursor(direction int) {
    // 커서 이동 로직
}

func (s *State) ToggleExpand() {
    // 디렉토리 확장/축소
}

func (s *State) SelectCurrent() {
    // 현재 항목 선택
}
```

### 학습 포인트
- 상태 관리 패턴
- 불변성과 상태 복사
- 이벤트 소싱 개념

---

## 6단계: 통합 및 메인 루프 (4-5시간)

### 목표
- 모든 컴포넌트를 연결
- 이벤트 루프 구현
- 사용자 입력에 따른 동작 연결

### 구현 내용

#### 6.1 애플리케이션 구조체
```go
// cmd/twf/app.go
package main

type App struct {
    terminal *terminal.Terminal
    state    *state.State
    treeView *views.TreeView
}

func (app *App) Run() error {
    // 메인 이벤트 루프
    for {
        // 1. 이벤트 읽기
        // 2. 상태 업데이트
        // 3. 화면 렌더링
        // 4. 종료 조건 확인
    }
}
```

### 학습 포인트
- 이벤트 루프 패턴
- 관심사의 분리 (SoC)
- 에러 전파와 복구

---

## 7단계: 고급 기능 (각 2-3시간)

### 7.1 파일 미리보기
- 텍스트 파일 내용 표시
- 이미지/바이너리 파일 정보
- 분할 화면 레이아웃

### 7.2 검색 기능
- 파일명 필터링
- 정규표현식 지원
- 실시간 검색

### 7.3 설정 시스템
- 키바인딩 커스터마이징
- 색상 테마
- 설정 파일 로드/저장

### 7.4 성능 최적화
- 지연 로딩 최적화
- 메모리 사용량 모니터링
- 대용량 디렉토리 처리

## 🔧 디버깅 및 테스트

### 로깅 시스템
```go
import "go.uber.org/zap"

logger, _ := zap.NewDevelopment()
logger.Info("Debug information",
    zap.String("path", currentPath),
    zap.Int("cursor", cursorPos))
```

### 단위 테스트
```go
// internal/filetree/filetree_test.go
func TestFileTreeCreation(t *testing.T) {
    tree, err := NewFileTree("/tmp")
    assert.NoError(t, err)
    assert.NotNil(t, tree)
}
```

## 🎨 확장 아이디어

1. **플러그인 시스템**: 외부 명령 통합 (fzf, ripgrep)
2. **네트워크 파일 시스템**: SFTP, FTP 지원
3. **압축 파일 지원**: ZIP, TAR 내부 탐색
4. **Git 통합**: 변경 상태, 브랜치 정보 표시
5. **북마크 기능**: 자주 사용하는 경로 저장

## 📖 추가 학습 자료

- **터미널 프로그래밍**: "The Linux Programming Interface" 24장
- **Go 동시성**: "Concurrency in Go" by Katherine Cox-Buday
- **TUI 디자인**: ncurses 라이브러리 문서
- **파일 시스템**: "Understanding the Linux Kernel" 12장

## 🤝 프로젝트 완료 후

1. **코드 리뷰**: 원본 twf와 구조 비교
2. **성능 측정**: 벤치마크 테스트 작성
3. **문서화**: API 문서와 사용 가이드
4. **배포**: 크로스 컴파일과 바이너리 배포

이 가이드를 따라가면서 단계별로 구현하면, TUI 프로그래밍의 핵심 개념들을 체계적으로 학습할 수 있습니다. 각 단계마다 충분한 시간을 투자하여 이해하고 넘어가는 것이 중요합니다.