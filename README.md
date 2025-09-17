# TWF Clone - TUI 파일 브라우저 클론 코딩 프로젝트

이 프로젝트는 원본 twf (Tree View Find) 프로젝트를 기반으로 한 TUI 파일 브라우저 클론 코딩 학습용 프로젝트입니다.

## 프로젝트 목표

- TUI (Terminal User Interface) 프로그래밍 이해
- Go를 활용한 시스템 프로그래밍 학습
- 파일 시스템과 터미널 제어 기법 익히기
- 모듈화된 소프트웨어 아키텍처 설계 경험

## 디렉토리 구조

```
twf-clone/
├── cmd/
│   └── twf/                # 실행 가능한 바이너리 메인 패키지
├── internal/
│   ├── config/            # 설정 관리
│   ├── filetree/          # 파일 트리 구조 및 조작
│   ├── state/             # 애플리케이션 상태 관리
│   ├── terminal/          # 터미널 인터페이스 및 이벤트 처리
│   └── views/             # UI 컴포넌트 (트리뷰, 프리뷰, 상태바)
└── docs/                  # 학습 문서 및 가이드
```

## 단계별 구현 가이드

### 1단계: 프로젝트 초기화 및 기본 구조
- [ ] Go 모듈 초기화
- [ ] 기본 디렉토리 구조 생성
- [ ] main.go 뼈대 작성

### 2단계: 터미널 기본 제어
- [ ] 터미널 초기화 및 종료
- [ ] 키보드 입력 처리
- [ ] 화면 클리어 및 커서 제어

### 3단계: 파일 시스템 인터페이스
- [ ] 파일/디렉토리 읽기
- [ ] 트리 구조 구현
- [ ] 파일 정보 표시

### 4단계: 기본 UI 렌더링
- [ ] 파일 목록 표시
- [ ] 커서 이동 및 선택
- [ ] 스크롤 기능

### 5단계: 고급 기능
- [ ] 디렉토리 확장/축소
- [ ] 파일 미리보기
- [ ] 검색 기능
- [ ] 설정 시스템

## 개발 환경 설정

### 필수 요구사항
- Go 1.13 이상
- Unix/Linux/macOS 터미널

### 의존성 설치
```bash
go mod init twf-clone
go get golang.org/x/crypto/ssh/terminal
go get golang.org/x/sys/unix
go get go.uber.org/zap
```

## 실행 방법

```bash
# 개발 중 실행
go run cmd/twf/main.go

# 빌드 후 실행
go build -o twf cmd/twf/main.go
./twf
```

## 학습 리소스

- `docs/learning-guide.md`: 상세한 단계별 학습 가이드
- `docs/architecture.md`: 아키텍처 설계 문서
- `docs/terminal-programming.md`: 터미널 프로그래밍 기법
- 원본 프로젝트: https://github.com/wvanlint/twf

## 참고사항

이 프로젝트는 학습 목적으로 제작되었으며, 원본 twf 프로젝트의 구조와 설계 패턴을 참조하여 단계별로 구현해 나가는 것을 목표로 합니다.