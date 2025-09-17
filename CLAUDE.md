# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 프로젝트 개요

`twf-clone`은 원본 twf (Tree View Find) 프로젝트를 기반으로 한 TUI 파일 브라우저 클론 코딩 학습용 프로젝트입니다. Go 언어와 터미널 프로그래밍을 단계별로 학습할 수 있도록 설계되었습니다.

## 클론 코딩 학습 가이드

### 프로젝트 구조

이 프로젝트는 학습 목적으로 설계된 모듈화된 구조를 가집니다:

```
twf-clone/
├── cmd/twf/               # 실행 바이너리 (main 패키지)
├── internal/              # 내부 패키지들
│   ├── config/           # 설정 관리 및 명령행 플래그
│   ├── filetree/         # 파일 시스템 트리 구조
│   ├── state/            # 애플리케이션 상태 관리
│   ├── terminal/         # 터미널 제어 및 이벤트 처리
│   └── views/            # UI 컴포넌트 (트리뷰, 프리뷰 등)
└── docs/                 # 학습 문서
```

### 단계별 구현 순서

학습 효과를 위해 다음 순서로 구현하는 것을 권장합니다:

#### 1단계: 기본 구조 (cmd/twf/main.go)
- 프로젝트 초기화 및 기본 명령행 처리
- Go 모듈과 패키지 구조 이해

#### 2단계: 터미널 기본 제어 (internal/terminal/)
- `terminal.go`: 터미널 초기화, raw 모드 설정
- `event.go`: 키보드 이벤트 파싱 및 처리
- `renderer.go`: 기본 화면 출력 및 제어

#### 3단계: 파일 시스템 인터페이스 (internal/filetree/)
- `filetree.go`: FileTree 구조체 및 기본 메서드
- `node.go`: 트리 노드 조작 (확장, 축소, 탐색)
- `walker.go`: 트리 순회 알고리즘

#### 4단계: 상태 관리 (internal/state/)
- `state.go`: 애플리케이션 상태 구조체
- `cursor.go`: 커서 위치 및 이동 로직
- `selection.go`: 파일/디렉토리 선택 관리

#### 5단계: UI 렌더링 (internal/views/)
- `view.go`: View 인터페이스 정의
- `tree_view.go`: 파일 트리 표시 뷰
- `status_view.go`: 상태바 뷰

#### 6단계: 고급 기능
- `preview_view.go`: 파일 미리보기 뷰
- `config/config.go`: 설정 시스템

### 핵심 학습 포인트

#### 터미널 프로그래밍
- Raw 모드와 Cooked 모드의 차이
- ANSI 이스케이프 시퀀스 활용
- 논블로킹 입력 처리
- 화면 버퍼링과 부분 업데이트

#### Go 언어 특화 패턴
- 인터페이스 기반 설계
- 구조체 임베딩 활용
- 에러 처리 패턴
- 동시성 처리 (goroutine, channel)

#### 소프트웨어 아키텍처
- MVC 패턴 적용
- 관심사의 분리 (Separation of Concerns)
- 의존성 주입
- 테스트 가능한 설계

## 디버깅 및 개발 팁

### 로깅
터미널 애플리케이션에서는 stdout/stderr를 직접 사용할 수 없으므로 파일 로깅을 사용:

```go
import "go.uber.org/zap"

// 개발 중 로깅 설정
logger, _ := zap.Config{
    Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
    Encoding:    "console",
    OutputPaths: []string{"debug.log"},
}.Build()
```

### 테스트 전략
- 각 패키지별로 단위 테스트 작성
- 터미널 I/O를 모킹하여 테스트
- 파일 시스템 조작은 임시 디렉토리 사용

### 일반적인 이슈들
1. **터미널 상태 복구**: 프로그램 종료 시 반드시 원래 터미널 상태로 복구
2. **시그널 처리**: SIGINT, SIGTERM 처리로 안전한 종료
3. **메모리 관리**: 대용량 디렉토리 처리 시 지연 로딩 활용
4. **크로스 플랫폼**: Unix 계열 시스템 중심 (Windows는 별도 처리 필요)

## 확장 아이디어

클론 코딩 완료 후 추가할 수 있는 기능들:

- **검색 기능**: 파일명 필터링, 정규표현식 지원
- **북마크**: 자주 사용하는 경로 저장
- **색상 테마**: 커스터마이징 가능한 색상 설정
- **플러그인 시스템**: 외부 명령 통합 (fzf, ripgrep 등)
- **Git 통합**: 파일 변경 상태 표시

## 참고 문서

- `docs/learning-guide.md`: 상세한 7단계 구현 가이드
- `docs/architecture.md`: 소프트웨어 아키텍처 설계
- `docs/terminal-programming.md`: 터미널 프로그래밍 기법
- 원본 프로젝트: `../1-file-tree-browser-clone-by-twf/`

## 주의사항

이 프로젝트는 학습 목적으로 제작되었습니다:
- 단계별 구현을 통한 점진적 학습을 목표로 함
- 성능보다는 코드의 명확성과 이해도를 우선시
- 각 단계마다 충분한 이해 후 다음 단계로 진행 권장
- 테스트 주도 개발(TDD) 접근 방식 추천
