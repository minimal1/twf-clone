# 터미널 프로그래밍 가이드

## 📋 개요

이 문서는 TUI (Terminal User Interface) 애플리케이션 개발을 위한 터미널 프로그래밍의 핵심 개념과 기법을 설명합니다.

## 🖥️ 터미널 기본 개념

### 1. 터미널 모드

#### Cooked Mode (기본 모드)
- 줄 단위 입력 처리
- 백스페이스, Ctrl+C 등이 자동으로 처리됨
- 엔터키를 눌러야 입력이 프로그램으로 전달됨

```go
// 기본적인 입력 읽기
scanner := bufio.NewScanner(os.Stdin)
for scanner.Scan() {
    text := scanner.Text()
    fmt.Println("입력:", text)
}
```

#### Raw Mode (원시 모드)
- 키를 누르는 즉시 프로그램으로 전달
- 특수키 처리를 프로그램이 직접 담당
- TUI 애플리케이션에 필수

```go
import "golang.org/x/crypto/ssh/terminal"

// Raw 모드 진입
oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
if err != nil {
    panic(err)
}

// 프로그램 종료 시 복구
defer terminal.Restore(int(os.Stdin.Fd()), oldState)
```

### 2. ANSI 이스케이프 시퀀스

터미널 제어를 위한 특수 문자 조합

#### 기본 제어
```go
const (
    // 커서 제어
    CursorUp    = "\033[A"
    CursorDown  = "\033[B"
    CursorRight = "\033[C"
    CursorLeft  = "\033[D"

    // 화면 제어
    ClearScreen = "\033[2J"
    ClearLine   = "\033[K"
    Home        = "\033[H"

    // 커서 위치
    SaveCursor    = "\033[s"
    RestoreCursor = "\033[u"
)

// 특정 위치로 커서 이동
func MoveCursor(x, y int) string {
    return fmt.Sprintf("\033[%d;%dH", y+1, x+1)
}
```

#### 색상 제어
```go
const (
    // 전경색
    FgBlack   = "\033[30m"
    FgRed     = "\033[31m"
    FgGreen   = "\033[32m"
    FgYellow  = "\033[33m"
    FgBlue    = "\033[34m"
    FgMagenta = "\033[35m"
    FgCyan    = "\033[36m"
    FgWhite   = "\033[37m"

    // 배경색
    BgBlack   = "\033[40m"
    BgRed     = "\033[41m"
    // ...

    // 스타일
    Bold      = "\033[1m"
    Dim       = "\033[2m"
    Underline = "\033[4m"
    Reverse   = "\033[7m"
    Reset     = "\033[0m"
)
```

## ⌨️ 키보드 입력 처리

### 1. 기본 키 읽기

```go
func ReadKey() ([]byte, error) {
    // 1바이트씩 읽기
    buffer := make([]byte, 1)
    _, err := os.Stdin.Read(buffer)
    return buffer, err
}

// 멀티바이트 키 처리 (화살표 키 등)
func ReadKeySequence() ([]byte, error) {
    first, err := ReadKey()
    if err != nil {
        return nil, err
    }

    // ESC 시퀀스인지 확인
    if first[0] == 27 { // ESC
        // 추가 바이트 읽기
        second, _ := ReadKey()
        if second[0] == '[' {
            third, _ := ReadKey()
            return []byte{first[0], second[0], third[0]}, nil
        }
    }

    return first, nil
}
```

### 2. 이벤트 기반 입력 처리

```go
type KeyEvent struct {
    Key       rune
    Special   SpecialKey
    Modifiers KeyModifier
}

type SpecialKey int

const (
    KeyNone SpecialKey = iota
    KeyArrowUp
    KeyArrowDown
    KeyArrowLeft
    KeyArrowRight
    KeyEnter
    KeyEscape
    KeyTab
    KeyBackspace
    KeyDelete
)

type KeyModifier int

const (
    ModNone KeyModifier = 0
    ModCtrl KeyModifier = 1 << iota
    ModAlt
    ModShift
)

func ParseKeyEvent(data []byte) KeyEvent {
    if len(data) == 1 {
        ch := data[0]

        // Ctrl 키 조합 확인
        if ch < 32 {
            return KeyEvent{
                Key:       rune(ch + 'a' - 1),
                Modifiers: ModCtrl,
            }
        }

        // 일반 문자
        return KeyEvent{Key: rune(ch)}
    }

    // ESC 시퀀스 처리
    if len(data) >= 3 && data[0] == 27 && data[1] == '[' {
        switch data[2] {
        case 'A':
            return KeyEvent{Special: KeyArrowUp}
        case 'B':
            return KeyEvent{Special: KeyArrowDown}
        case 'C':
            return KeyEvent{Special: KeyArrowRight}
        case 'D':
            return KeyEvent{Special: KeyArrowLeft}
        }
    }

    return KeyEvent{}
}
```

## 🎨 화면 렌더링

### 1. 버퍼링된 렌더링

```go
type ScreenBuffer struct {
    width   int
    height  int
    buffer  [][]rune
    colors  [][]Color
    dirty   bool
}

func NewScreenBuffer(width, height int) *ScreenBuffer {
    sb := &ScreenBuffer{
        width:  width,
        height: height,
        buffer: make([][]rune, height),
        colors: make([][]Color, height),
    }

    for i := range sb.buffer {
        sb.buffer[i] = make([]rune, width)
        sb.colors[i] = make([]Color, width)
    }

    return sb
}

func (sb *ScreenBuffer) SetCell(x, y int, ch rune, color Color) {
    if x >= 0 && x < sb.width && y >= 0 && y < sb.height {
        sb.buffer[y][x] = ch
        sb.colors[y][x] = color
        sb.dirty = true
    }
}

func (sb *ScreenBuffer) Render() string {
    if !sb.dirty {
        return ""
    }

    var output strings.Builder
    output.WriteString(ClearScreen)
    output.WriteString(Home)

    for y := 0; y < sb.height; y++ {
        for x := 0; x < sb.width; x++ {
            if sb.colors[y][x] != ColorDefault {
                output.WriteString(sb.colors[y][x].String())
            }
            output.WriteRune(sb.buffer[y][x])
        }
        output.WriteString(Reset)
        output.WriteString("\n")
    }

    sb.dirty = false
    return output.String()
}
```

### 2. 부분 업데이트

```go
type DirtyRegion struct {
    X, Y, Width, Height int
}

func (sb *ScreenBuffer) MarkDirty(x, y, width, height int) {
    // 더티 영역 추가
    sb.dirtyRegions = append(sb.dirtyRegions, DirtyRegion{
        X: x, Y: y, Width: width, Height: height,
    })
}

func (sb *ScreenBuffer) RenderPartial() string {
    var output strings.Builder

    for _, region := range sb.dirtyRegions {
        for y := region.Y; y < region.Y+region.Height; y++ {
            output.WriteString(MoveCursor(region.X, y))
            for x := region.X; x < region.X+region.Width; x++ {
                output.WriteRune(sb.buffer[y][x])
            }
        }
    }

    sb.dirtyRegions = nil
    return output.String()
}
```

## 🔧 터미널 크기 및 시그널 처리

### 1. 터미널 크기 감지

```go
import "golang.org/x/crypto/ssh/terminal"

func GetTerminalSize() (width, height int, err error) {
    width, height, err = terminal.GetSize(int(os.Stdout.Fd()))
    return
}
```

### 2. 리사이즈 시그널 처리

```go
import (
    "os"
    "os/signal"
    "syscall"
)

func HandleResize(callback func(width, height int)) {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGWINCH)

    go func() {
        for range sigCh {
            width, height, err := GetTerminalSize()
            if err == nil {
                callback(width, height)
            }
        }
    }()
}
```

### 3. 종료 시그널 처리

```go
func HandleShutdown(cleanup func()) {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigCh
        cleanup()
        os.Exit(0)
    }()
}
```

## 🎭 고급 기법

### 1. Alternative Screen Buffer

```go
const (
    EnterAltScreen = "\033[?1049h"
    ExitAltScreen  = "\033[?1049l"
)

func EnterAlternateScreen() {
    fmt.Print(EnterAltScreen)
}

func ExitAlternateScreen() {
    fmt.Print(ExitAltScreen)
}
```

### 2. 마우스 지원

```go
const (
    EnableMouse  = "\033[?1000h"
    DisableMouse = "\033[?1000l"
)

type MouseEvent struct {
    X, Y   int
    Button MouseButton
    Action MouseAction
}

func ParseMouseEvent(data []byte) MouseEvent {
    // ESC[M + 3 bytes 형식 파싱
    if len(data) >= 6 && data[0] == 27 && data[1] == '[' && data[2] == 'M' {
        button := data[3] - 32
        x := int(data[4]) - 32
        y := int(data[5]) - 32

        return MouseEvent{
            X:      x,
            Y:      y,
            Button: MouseButton(button & 3),
            Action: MouseAction((button >> 2) & 3),
        }
    }

    return MouseEvent{}
}
```

### 3. 성능 최적화

#### 출력 버퍼링
```go
import "bufio"

type Terminal struct {
    writer *bufio.Writer
}

func NewTerminal() *Terminal {
    return &Terminal{
        writer: bufio.NewWriter(os.Stdout),
    }
}

func (t *Terminal) Write(data string) {
    t.writer.WriteString(data)
}

func (t *Terminal) Flush() {
    t.writer.Flush()
}
```

#### 논블로킹 입력
```go
import "syscall"

func SetNonBlocking(fd int) error {
    return syscall.SetNonblock(fd, true)
}

func ReadNonBlocking() ([]byte, error) {
    SetNonBlocking(int(os.Stdin.Fd()))
    defer SetNonBlocking(int(os.Stdin.Fd())) // blocking으로 복구

    buffer := make([]byte, 256)
    n, err := syscall.Read(int(os.Stdin.Fd()), buffer)
    if err != nil {
        return nil, err
    }

    return buffer[:n], nil
}
```

## 🐛 디버깅 기법

### 1. 로그 파일 사용

터미널을 직접 제어하는 프로그램에서는 stdout/stderr 사용이 어려우므로 파일 로깅을 사용합니다.

```go
import (
    "log"
    "os"
)

func InitLogging() {
    logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        panic(err)
    }
    log.SetOutput(logFile)
}

func DebugLog(format string, args ...interface{}) {
    log.Printf(format, args...)
}
```

### 2. 상태 덤프

```go
func (state *State) DumpState() {
    DebugLog("=== State Dump ===")
    DebugLog("Cursor: %d", state.CursorPos)
    DebugLog("Scroll: %d", state.ScrollOffset)
    DebugLog("Selection: %v", state.Selection)
    DebugLog("=================")
}
```

## 📚 참고 자료

- **ANSI 이스케이프 시퀀스**: [Wikipedia ANSI escape code](https://en.wikipedia.org/wiki/ANSI_escape_code)
- **VT100 터미널**: DEC VT100 사용자 가이드
- **ncurses 라이브러리**: C 언어 TUI 라이브러리 (참고용)
- **termbox-go**: Go 터미널 UI 라이브러리 (대안)

이 가이드를 통해 터미널 프로그래밍의 핵심 개념을 이해하고, TUI 애플리케이션 개발에 필요한 기초 지식을 습득할 수 있습니다.