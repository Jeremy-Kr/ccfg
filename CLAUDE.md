# ccfg — Claude Code Config Viewer

Go + Bubbletea 기반 TUI 대시보드. Claude Code의 분산된 설정 파일들을 한 화면에서 조회.

## 기술 스택

- **언어:** Go 1.26
- **TUI:** Bubbletea v1.3.x + Lipgloss + Bubbles
- **구문 강조:** Glamour + Chroma
- **CLI:** 표준 `flag` 패키지 (별도 프레임워크 없음)

## 빌드 / 테스트 / 린트

```bash
# 빌드
go build -o ccfg ./cmd/ccfg

# 전체 테스트
go test ./...

# 특정 패키지 테스트
go test ./internal/scanner/...

# 린트 (golangci-lint 설치 필요)
golangci-lint run

# 실행
./ccfg
```

## 디렉토리 구조

```
cmd/ccfg/          엔트리포인트
internal/
  scanner/         설정 파일 탐색 (경로, 스캔 로직)
  parser/          설정 파일 파싱 (JSON, JSONC, Markdown)
  model/           공유 타입 정의
  tui/             Bubbletea 모델, 뷰, 컴포넌트
docs/              PRD, 기술 설계, 로드맵
```

## 코딩 컨벤션

- Go 표준 스타일 (`gofmt`)
- `internal/` 패키지로 캡슐화
- 에러는 `fmt.Errorf("context: %w", err)` 패턴으로 래핑
- 인터페이스는 사용하는 쪽에서 정의
- 테스트 파일은 `*_test.go`로 같은 패키지에 위치
- 변수/함수명은 Go 네이밍 컨벤션 준수 (camelCase, exported는 PascalCase)

## 핵심 제약

- 읽기 전용 도구: 설정 파일을 절대 수정하지 않음
- macOS 우선 지원 (Linux는 Phase 5에서)
- 최소 의존성: 필요한 라이브러리만 추가
