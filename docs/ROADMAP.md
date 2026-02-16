# ccfg — Roadmap

## Phase 1: 프로젝트 스캐폴딩 ✅

- [x] Go 환경 설정
- [x] 프로젝트 문서 (PRD, 기술 설계, 로드맵)
- [x] go.mod 초기화
- [x] 최소 엔트리포인트 (`cmd/ccfg/main.go`)
- [x] 경로 상수 정의 (`internal/scanner/paths.go`)
- [x] git 초기화 + .gitignore

## Phase 2: 스캐너 구현

- [ ] `internal/model/types.go` — 공유 타입 정의 (Scope, FileType, ConfigFile, ScanResult)
- [ ] `internal/scanner/scanner.go` — Scanner 인터페이스 및 DefaultScanner 구현
- [ ] `internal/scanner/root.go` — 프로젝트 루트 감지 (.git 탐색)
- [ ] `internal/scanner/scanner_test.go` — 테스트 (임시 디렉토리 기반)
- [ ] Bubbletea, Lipgloss 등 의존성 추가

## Phase 3: 기본 TUI

- [ ] `internal/tui/model.go` — 메인 TUI 모델 (Init, Update, View)
- [ ] `internal/tui/tree.go` — 트리 뷰 컴포넌트
- [ ] `internal/tui/preview.go` — 미리보기 패널 (raw text)
- [ ] `internal/tui/keys.go` — 키 바인딩 정의
- [ ] `internal/tui/styles.go` — Lipgloss 스타일 정의
- [ ] 30/70 split 레이아웃
- [ ] 키보드 네비게이션 (j/k/Enter/Tab/q)

## Phase 4: 파서 + 구문 강조

- [ ] `internal/parser/json.go` — JSON/JSONC 파서
- [ ] `internal/parser/markdown.go` — Markdown 렌더러 (Glamour)
- [ ] 구문 강조 (Chroma) 적용
- [ ] 파싱 에러 시 원본 텍스트 폴백

## Phase 5: 고급 기능

- [ ] 병합 뷰 (Merged View) — 모든 Scope 설정 병합 + 출처 표시
- [ ] 검색 기능 (`/` 키로 진입)
- [ ] Linux 경로 지원
- [ ] 파일 변경 감지 (fsnotify) 및 자동 갱신

## Phase 6: 배포 + 품질

- [ ] golangci-lint 설정 및 CI
- [ ] goreleaser 설정 (크로스 컴파일 바이너리)
- [ ] README.md 작성 (설치 방법, 스크린샷)
- [ ] Homebrew formula (선택)
