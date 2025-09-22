# TWF Clone - Claude Task Management Board

Claude Code가 프로젝트 진행을 체계적으로 관리하기 위한 태스크 보드입니다.

## 🤝 역할 분담

### 👨‍💻 사용자 (개발자)
- **실제 코드 구현**: 모든 파일 작성 및 코딩 작업
- **과제 수행**: 제시된 학습 과제 완료
- **문제 해결**: 구현 중 발생하는 기술적 문제 해결

### 🤖 Claude Code
- **과제 제시**: 다음 단계 학습 과제 및 구현 가이드 제공
- **진행 관리**: 전체 프로젝트 진도 추적 및 우선순위 관리
- **검증 지원**: 구현된 코드 리뷰 및 다음 단계 제안
- **학습 가이드**: 기술적 개념 설명 및 참고 자료 제공

## 🎯 프로젝트 개요
- **목표**: Go 언어 기반 TUI 파일 브라우저 클론 구현
- **학습 중심**: 단계별 점진적 구현을 통한 터미널 프로그래밍 학습
- **상태**: 초기 단계

---

## 📊 전체 진행 상황

```
전체 진도: ██████░░░░ 60%

1단계 (기본 구조)     : ██████████ 100%
2단계 (터미널 제어)   : ██████████ 100%
3단계 (파일시스템)    : █████████░ 90%
4단계 (상태 관리)     : ░░░░░░░░░░ 0%
5단계 (UI 렌더링)     : ░░░░░░░░░░ 0%
6단계 (고급 기능)     : ░░░░░░░░░░ 0%
```

---

## 🔄 현재 활성 태스크

### 진행 중 (IN_PROGRESS)
- 없음

### 다음 우선순위 (NEXT)
1. **Walker Phase 2: 검색 기능 구현**
   - FindByName(): 파일명 패턴 검색
   - FindByExtension(): 파일 확장자별 필터링
   - FilterHidden(): 숨김 파일 처리

2. **Walker Phase 3: 고급 순회 기능**
   - Walk(): 범용 트리 순회 시스템
   - CollectWhere(): 조건부 노드 수집

---

## 📋 단계별 태스크 상세

### 1️⃣ 1단계: 프로젝트 기본 구조
```
상태: 🟡 진행중 (80%)
예상 소요: 0.5일
```

#### ✅ 완료
- [x] 프로젝트 디렉토리 구조 생성
- [x] CLAUDE.md 프로젝트 가이드 작성
- [x] 학습 문서 구조 설계
- [x] **Go 모듈 초기화 완료**
  - 파일: `go.mod`, `go.sum` 생성
  - 의존성: `golang.org/x/term`, `golang.org/x/sys`, `go.uber.org/zap` 추가
  - 터미널 라이브러리: `golang.org/x/term` 선택 완료

#### ✅ 추가 완료
- [x] `cmd/twf/main.go` 기본 스켈레톤 구현
- [x] 기본 터미널 제어 테스트 코드 작성 및 검증

---

### 2️⃣ 2단계: 터미널 기본 제어
```
상태: ✅ 완료 (100%)
예상 소요: 1.5일
선행 조건: 1단계 완료
```

#### ✅ 완료된 모듈들
- [x] **`internal/terminal/terminal.go`** (핵심)
  - 터미널 초기화/정리 (`NewTerminal`, `Cleanup`)
  - Raw 모드 설정/해제 (`EnableRawMode`, `DisableRawMode`)
  - 화면 크기 감지 (`GetSize`)
  - `/dev/tty` 기반 효율적인 파일 디스크립터 관리

- [x] **`internal/terminal/event.go`** (핵심)
  - 키보드 입력 파싱 (`parseInputData` 함수)
  - 이벤트 타입 정의 (인터페이스 기반 시스템)
  - 특수 키 (화살표, ESC, Enter, Tab, Ctrl+C/D 등) 처리
  - UTF-8 문자 지원 및 실시간 이벤트 읽기 (`ReadEvent`)

- [x] **`internal/terminal/renderer.go`** (렌더링)
  - ANSI 이스케이프 시퀀스 상수 정의
  - 화면 제어 (`ClearScreen`, `ClearLine`, 대체 화면)
  - 커서 제어 (`MoveCursorTo`, `HideCursor`, `ShowCursor`)
  - 타입 안전한 색상 시스템 (`Color`, `Style` 타입)
  - 색상 출력 함수들 (`WriteColored`, `WriteColoredAt`)

#### ✅ 완료된 주요 결정사항
- [x] 터미널 라이브러리 선택: `golang.org/x/term`
- [x] 파일 디스크립터 전략: 단일 `/dev/tty` 파일 (`O_RDWR`)
- [x] 에러 처리 전략: 체인형 에러 반환
- [x] 이벤트 시스템: 인터페이스 기반 타입 안전한 설계
- [x] 색상 시스템: Custom 타입으로 타입 안전성 보장

---

### 3️⃣ 3단계: 파일 시스템 인터페이스
```
상태: ✅ 완료 (100%)
예상 소요: 2일
선행 조건: 2단계 완료
```

#### ✅ 완료된 컴포넌트
1. **`internal/filetree/node.go`** (완료)
   - TreeNode 구조체 정의
   - 파일/디렉토리 정보 저장 (Path, Name, IsDir, Size, ModTime)
   - 부모-자식 관계 관리 (Parent, Children)
   - 상태 관리 (Expanded, Loaded, Selected)
   - 트리 조작 메서드 (AddChild, RemoveChild, GetChildByName)
   - 유틸리티 메서드 (IsRoot, Depth, CanExpand, GetDisplayName)

2. **`internal/filetree/filetree.go`** (완료)
   - FileTree 인터페이스 정의
   - FileTreeImpl 구조체 구현
   - 루트 디렉토리 로딩 (LoadRoot)
   - 노드 확장/축소 (ExpandNode, CollapseNode)
   - 지연 로딩 구현 (loadChildren)
   - 노드 새로고침 (RefreshNode)
   - 현재 노드 관리 (GetCurrentNode, SetCurrentNode)

3. **`internal/filetree/walker.go`** (진행중 - Phase 1 완료)
   - **Phase 1**: UI 지원 메서드 (완료 ✅)
     - GetVisibleNodes(): 화면에 표시할 노드들 수집
     - GetNextVisibleNode()/GetPrevVisibleNode(): 키보드 네비게이션
     - collectVisible(): 재귀적 가시 노드 수집
   - **Phase 2**: 검색 기능 (다음 우선순위)
     - FindByName(): 파일명 검색
     - FindByExtension(): 확장자별 검색
     - FilterHidden(): 숨김 파일 필터링
   - **Phase 3**: 고급 순회 (확장 기능)
     - Walk(): 범용 트리 순회
     - CollectWhere(): 조건부 노드 수집

#### ✅ 추가 완성: 실용적 파일 브라우저 데모
- **`cmd/filebrowser/main.go`** (완료)
  - FileBrowserApp 구조체: 완전한 TUI 파일 브라우저
  - 실시간 파일 트리 탐색: 현재 작업 디렉토리 로딩
  - 키보드 네비게이션: 화살표 키 + Vim 스타일 (j/k/l)
  - 디렉토리 확장/축소: Enter 키로 토글
  - 시각적 피드백: 현재 선택 항목 하이라이트
  - 들여쓰기 표시: 트리 구조 깊이별 시각화
  - 안전한 종료: ESC/Ctrl+C 지원

#### ✅ 완료된 성능 고려사항
- [x] 지연 로딩 구현 (Loaded 플래그 활용)
- [x] 메모리 효율적인 트리 구조 설계
- [x] 권한 오류 처리 (에러 체인 패턴)

---

### 4️⃣ 4단계: 애플리케이션 상태 관리
```
상태: ⚪ 대기중 (0%)
예상 소요: 1.5일
선행 조건: 3단계 완료
```

#### 상태 컴포넌트
1. **`internal/state/state.go`** - 전역 상태
2. **`internal/state/cursor.go`** - 커서 위치 관리
3. **`internal/state/selection.go`** - 선택 상태

#### 아키텍처 결정
- [ ] 상태 변경 패턴 (이벤트 기반 vs 직접 변경)
- [ ] 상태 불변성 보장 방법
- [ ] 상태 검증 로직

---

### 5️⃣ 5단계: UI 렌더링 시스템
```
상태: ⚪ 대기중 (0%)
예상 소요: 2일
선행 조건: 4단계 완료
```

#### 뷰 컴포넌트
1. **`internal/views/view.go`** - 뷰 인터페이스
2. **`internal/views/tree_view.go`** - 트리 뷰
3. **`internal/views/status_view.go`** - 상태바

#### UI 설계 결정
- [ ] 레이아웃 시스템 (고정 vs 동적)
- [ ] 색상 스키마
- [ ] 반응형 처리

---

### 6️⃣ 6단계: 고급 기능
```
상태: ⚪ 대기중 (0%)
예상 소요: 1.5일
선행 조건: 5단계 완료
```

#### 확장 기능
1. **`internal/views/preview_view.go`** - 파일 미리보기
2. **`internal/config/config.go`** - 설정 시스템

---

## 🚨 이슈 및 블로커

### 현재 블로커
| 우선순위 | 이슈 | 영향도 | 해결 방안 |
|----------|------|--------|-----------|
| HIGH | 터미널 라이브러리 선택 | 전체 아키텍처 | 각 라이브러리 비교 분석 필요 |

### 알려진 리스크
- **크로스 플랫폼 호환성**: Windows 지원 복잡성
- **성능**: 대용량 디렉토리 처리
- **메모리**: 트리 구조 메모리 사용량

---

## 🔧 개발 환경 설정

### 필수 도구
- [ ] Go 1.19+ 설치 확인
- [ ] 터미널 라이브러리 선택 및 설치
- [ ] 개발용 로깅 설정

### 코드 품질 도구
- [ ] `gofmt` 자동 포맷팅
- [ ] `golangci-lint` 정적 분석
- [ ] 테스트 커버리지 도구

---

## 📝 학습 체크포인트

### 완료된 학습
- [x] Go 모듈 시스템 이해
- [x] 프로젝트 구조 설계 원칙
- [x] 기본 터미널 제어 (`golang.org/x/term`)

### 진행중인 학습
- [ ] 터미널 프로그래밍 기초
- [ ] ANSI 이스케이프 시퀀스
- [ ] Go 포인터와 구조체 메서드

### 예정된 학습
- [ ] 트리 자료구조 구현
- [ ] 이벤트 기반 아키텍처
- [ ] 성능 최적화 기법

---

## 🛡️ 터미널 개발 체크리스트

### 환경 호환성 체크
- [ ] **터미널 환경 감지**: `term.IsTerminal()` 사용
- [ ] **TTY 접근 가능**: `/dev/tty` 파일 열기 가능 여부
- [ ] **권한 확인**: 터미널 제어 권한 보유 여부
- [ ] **플랫폼 호환성**: Windows/macOS/Linux 동작 확인

### 실행 환경 예외 상황
- [ ] **파이프 입력**: `echo "data" | program` 형태 실행
- [ ] **백그라운드 실행**: `program &` 형태 실행
- [ ] **Docker 컨테이너**: TTY 없는 환경에서 실행
- [ ] **SSH 원격 접속**: 네트워크를 통한 터미널 제어
- [ ] **IDE 터미널**: VSCode, IntelliJ 등 내장 터미널

### 리소스 관리 체크
- [ ] **파일 디스크립터 누수 방지**: 모든 열린 파일 정리
- [ ] **터미널 상태 복구**: Raw 모드에서 Cooked 모드로 복구
- [ ] **시그널 처리**: SIGINT, SIGTERM 안전한 종료
- [ ] **메모리 누수 방지**: 동적 할당된 리소스 해제

### 에러 처리 전략
- [ ] **Graceful Degradation**: 터미널 기능 불가 시 대체 동작
- [ ] **명확한 에러 메시지**: 사용자가 이해할 수 있는 에러 설명
- [ ] **로그 파일 활용**: 터미널 외부 디버깅 수단 제공
- [ ] **복구 가능성**: 일시적 오류 시 재시도 로직

### 사용자 경험 고려
- [ ] **화면 크기 변경**: 터미널 리사이즈 이벤트 처리
- [ ] **키보드 레이아웃**: 다양한 키보드 지원
- [ ] **색상 지원**: 터미널 색상 능력 감지
- [ ] **접근성**: 스크린 리더 호환성 고려

---

## 🎯 다음 세션 액션 아이템

### 즉시 처리 필요
1. **`internal/filetree/filetree.go` 구현**
   - FileTree 구조체 정의
   - 루트 디렉토리 로딩 기능
   - 기본 트리 관리 메서드

2. **`internal/filetree/node.go` 구현**
   - TreeNode 구조체 정의
   - 파일/디렉토리 정보 저장
   - 부모-자식 관계 관리

3. **`internal/filetree/walker.go` 구현**
   - 디렉토리 순회 로직
   - 지연 로딩 구현
   - 권한 오류 처리

### 중기 계획
- 애플리케이션 상태 관리 시스템 (4단계)
- 기본 UI 렌더링 시스템 구축 (5단계)
- 고급 기능 추가 (6단계)

---

**🎉 3단계 파일시스템 인터페이스 90% 완성!**
- TreeNode, FileTree, Walker Phase 1 모두 구현 완료
- **실용적인 파일 브라우저 데모 앱 완성** (`cmd/filebrowser/main.go`)
- 현재 작업 디렉토리 탐색, 키보드 네비게이션, 디렉토리 확장/축소 지원
- 다음: Walker Phase 2 (검색 기능) 구현

*📅 Last Updated: 2025-09-22*
*🤖 Managed by Claude Code*