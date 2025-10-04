# I/O 추상화 레이어의 이해

## 개요

MVC 패턴을 학습할 때 종종 간과되는 것이 **I/O 레이어의 역할**입니다. 이 문서는 TUI 애플리케이션 개발을 통해 발견한 I/O 추상화 레이어의 중요성과, 이를 다른 플랫폼(웹)과 비교하여 정리합니다.

## MVC에서 누락되는 것: I/O

### 전통적인 MVC 설명

일반적으로 MVC는 다음과 같이 설명됩니다:

```
User → Controller → Model
              ↓
            View
```

하지만 여기서 **"User가 어떻게 Controller와 소통하는가?"**가 빠져있습니다!

### 실제 MVC: I/O 레이어 포함

```
User (실제 사람)
  ↓ (키보드 입력)
┌─────────────────┐
│  I/O Layer      │ ← 입력을 Event로 변환
│  (Event 생성)    │
└─────────────────┘
  ↓ Event 객체
┌─────────────────┐
│  Controller     │ ← 이벤트 해석 및 처리
│  (이벤트 해석)   │
└─────────────────┘
  ↓ 상태 변경
┌─────────────────┐
│  Model          │ ← 애플리케이션 상태
└─────────────────┘
  ↓ 데이터 읽기
┌─────────────────┐
│  View           │ ← UI 렌더링 로직
│  (렌더링 명령)   │
└─────────────────┘
  ↓ Rendering Primitive
┌─────────────────┐
│  I/O Layer      │ ← 렌더링 명령을 화면 출력으로 변환
│  (ANSI 출력)     │
└─────────────────┘
  ↓ (화면 출력)
User (화면으로 확인)
```

**핵심: I/O 레이어가 양쪽 끝에 존재합니다!**

## 플랫폼별 I/O 추상화

### 1. TUI (Terminal User Interface) - 우리 프로젝트

```
┌─────────────────────────────────┐
│  Terminal (I/O 추상화 레이어)     │
├─────────────────────────────────┤
│  Input:  raw bytes → Event      │
│  Output: Primitive → ANSI       │
└─────────────────────────────────┘
```

**책임:**
- **Input**: 키보드 원시 바이트 → 타입-안전한 Event 객체
- **Output**: 렌더링 프리미티브 (WriteColoredAt) → ANSI 이스케이프 시퀀스

**코드 예시:**
```go
// Input: 원시 입력을 Event로 변환
event := term.ReadEvent()  // KeyPressEvent{Key: KeyArrowDown}

// Output: 렌더링 프리미티브 제공
term.WriteColoredAt(y, x, text, color)  // → \033[33m...
```

**I/O 종류:**
- ✅ User I/O만 존재 (키보드 입력, 화면 출력)
- ❌ Network I/O 없음 (로컬 파일시스템만 사용)

### 2. Web Frontend

```
┌─────────────────────────────────┐
│  Browser (이중 I/O 추상화)        │
├─────────────────────────────────┤
│  User I/O:                      │
│    Input:  mouse/keyboard       │
│            → DOM Event          │
│    Output: DOM API → Pixels     │
│                                 │
│  Network I/O:                   │
│    Input:  HTTP Response        │
│    Output: HTTP Request         │
└─────────────────────────────────┘
```

**Browser = Terminal (User I/O) + Network Client (Network I/O)**

#### User I/O (Terminal과 동일한 역할)

**Input:**
```javascript
// Browser가 원시 입력을 Event로 변환
button.addEventListener('click', (event) => {
  // event = 추상화된 이벤트 객체
});
```

**Output:**
```javascript
// Browser가 렌더링 프리미티브 제공
element.textContent = "Hello";
element.style.color = "red";
```

#### Network I/O (추가적인 책임)

```javascript
// Browser가 Network I/O 추상화
fetch('/api/users')
  .then(response => response.json())
  .then(data => updateState(data));
```

### 3. Backend (API Server)

```
┌─────────────────────────────────┐
│  Network I/O만                   │
├─────────────────────────────────┤
│  Input:  HTTP Request           │
│  Output: HTTP Response          │
└─────────────────────────────────┘
```

**책임:**
- ✅ Network I/O만 (HTTP Request/Response)
- ❌ User I/O 없음 (사용자와 직접 상호작용 없음)

## 비교표

| 역할 | TUI (우리) | Web Frontend | Backend |
|------|-----------|--------------|---------|
| **User I/O 추상화** | `Terminal` | `Browser (DOM API)` | ❌ |
| **Input 변환** | raw bytes → Event | mouse/keyboard → DOM Event | - |
| **Output 변환** | WriteColoredAt → ANSI | textContent/style → Pixels | - |
| **Network I/O** | ❌ | `Browser (fetch API)` | ✅ HTTP |
| **App (Controller)** | `App.handleEvent()` | `onClick handler` | `Route Handler` |
| **Model** | `AppState` | `React State / Vuex` | `Database Model` |
| **View** | `TreeView` | `React Component` | `JSON Response` |

## Terminal의 아키텍처 역할

### 레이어 관점

```
┌─────────────────────────────────────┐
│   Application Layer (App)           │  비즈니스 로직
│   - 이벤트 처리                        │
│   - 상태 관리                          │
└─────────────────────────────────────┘
           ↕ (Event, Command)
┌─────────────────────────────────────┐
│   Presentation Layer (Views)        │  UI 로직
│   - TreeView, StatusView            │
│   - Layout                           │
└─────────────────────────────────────┘
           ↕ (Render Command)
┌─────────────────────────────────────┐
│   I/O Abstraction Layer (Terminal)  │  ← Terminal!
│   - 입력: Event 생성                  │
│   - 출력: Rendering Primitive 제공   │
└─────────────────────────────────────┘
           ↕ (Raw I/O)
┌─────────────────────────────────────┐
│   OS Layer (TTY, /dev/tty)          │  OS 리소스
└─────────────────────────────────────┘
```

### Terminal의 책임

**해야 할 일:**
- ✅ OS의 터미널 리소스와 애플리케이션 사이의 추상화
- ✅ Raw 입력을 의미 있는 이벤트로 변환 (Input Adapter)
- ✅ 고수준 렌더링 요청을 ANSI 시퀀스로 변환 (Output Adapter)
- ✅ 터미널 상태 관리 (Raw Mode, 화면 크기 등)

**하지 않아야 할 일:**
- ❌ 어떤 View를 그릴지 결정 (Layout의 책임)
- ❌ 어떤 색상을 쓸지 결정 (View의 책임)
- ❌ 이벤트를 어떻게 처리할지 결정 (App의 책임)

### 효과적인 의사소통 용어

**권장 표현:**
- "Terminal이 렌더링 프리미티브를 제공한다"
- "Terminal이 원시 입력을 이벤트로 변환한다"
- "Terminal이 TTY 리소스를 추상화한다"
- "Terminal이 ANSI 시퀀스를 캡슐화한다"

**피해야 할 표현:**
- ❌ "Terminal이 화면을 렌더링한다" (너무 고수준)
- ❌ "Terminal이 키 입력을 처리한다" (비즈니스 로직 혼동)

## 왜 I/O가 자주 생략되는가?

### 웹 개발의 경우

웹 개발에서는 **Browser가 I/O를 이미 구현**해놓았기 때문에 개발자가 신경 쓸 필요가 없습니다:

```
웹 MVC:
Browser (I/O 담당) → HTTP Request → Controller
                                       ↓
                                     Model
                                       ↓
                                     View
                   ← HTTP Response ←   ↓
Browser (I/O 담당) ← HTML 렌더링  ←────┘
```

- Browser가 User I/O를 다 처리해줌
- 개발자는 Request/Response만 신경 쓰면 됨
- 따라서 "MVC만 알면 된다"는 착각 발생

### TUI/GUI의 경우

**하지만 TUI/GUI는 직접 I/O를 구현해야 합니다!**

Browser 같은 표준화된 플랫폼이 없기 때문에:
- 개발자가 직접 I/O 레이어 구현
- Terminal 초기화, Raw Mode, ANSI 시퀀스 처리
- 이벤트 파싱, 화면 제어 모두 직접

## 우리 프로젝트의 실제 흐름

```go
// 1. I/O Layer: 원시 입력을 Event로 변환
event := term.ReadEvent()  // Terminal (I/O)

// 2. Controller: Event 해석 및 Model 업데이트
app.handleEvent(event)  // App (Controller)
  → appState.Cursor().MoveTo(...)  // AppState (Model)

// 3. View: Model 읽어서 렌더링 명령 생성
layout.Render(term, appState)  // Layout/Views
  → term.WriteColoredAt(...)  // Terminal (I/O)
```

**Terminal이 양쪽 끝에 있습니다:**
- 입력 끝: Event Producer
- 출력 끝: Rendering Backend

## Frontend가 더 복잡한 이유

Frontend 개발자는 **두 가지 I/O를 조율**해야 합니다:

```javascript
// Frontend: 두 가지 I/O 조율
function handleClick() {
  // 1. User I/O 처리
  setLoading(true);  // UI 업데이트 (동기적)

  // 2. Network I/O 시작
  fetch('/api/data')  // 비동기적
    .then(data => {
      // 3. 다시 User I/O (UI 업데이트)
      setState(data);
      setLoading(false);
    });
}
```

**복잡성의 원인:**
- User I/O는 동기적 (즉각 반응 필요)
- Network I/O는 비동기적 (응답 대기)
- 두 I/O의 타이밍을 조율해야 함
- 비동기 상태 관리 필요
- UI와 데이터의 동기화 필요

**따라서 React/Vue 같은 프레임워크가 필요합니다:**
- 비동기 상태 관리
- UI와 데이터 동기화
- 선언적 렌더링

## 핵심 통찰

> **MVC의 핵심은 "관심사의 분리"이지만,
> 실제로 동작하려면 "외부 세계와의 소통" (I/O)이 필수다.**

**I/O 없이는:**
- ❌ Controller가 아무리 똑똑해도 입력을 받을 수 없음
- ❌ View가 아무리 예뻐도 화면에 표시할 수 없음

**I/O를 먼저 구현한 덕분에:**
- ✅ View는 "어떻게 그릴지" 고민 안 해도 됨 (Terminal이 해줌)
- ✅ Controller는 "어떻게 입력받을지" 고민 안 해도 됨 (Terminal이 해줌)
- ✅ 각자 자기 책임에만 집중 가능

## 요약

### 플랫폼별 I/O

| 플랫폼 | User I/O | Network I/O |
|--------|----------|-------------|
| **TUI (우리)** | ✅ Terminal | ❌ |
| **Web Frontend** | ✅ Browser (DOM) | ✅ Browser (fetch) |
| **Backend** | ❌ | ✅ HTTP Server |

### 핵심 개념

1. **I/O 추상화 레이어는 필수**
   - 외부 세계(User, Network, File)와 애플리케이션 사이의 다리
   - MVC가 동작하려면 반드시 필요

2. **Terminal = User I/O 추상화**
   - Input: raw bytes → Event
   - Output: Rendering Primitive → ANSI

3. **Browser = User I/O + Network I/O**
   - User I/O: DOM Events, Rendering
   - Network I/O: fetch, WebSocket

4. **관심사의 분리**
   - I/O Layer: "어떻게 통신할지"
   - View: "무엇을 보여줄지"
   - Controller: "어떻게 처리할지"
   - Model: "어떤 데이터를 저장할지"

---

*이 문서는 TWF Clone 프로젝트 개발 중 발견한 통찰을 정리한 것입니다.*
*날짜: 2025-10-04*
