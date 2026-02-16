# ccfg — Roadmap

## Phase 1: 프로젝트 스캐폴딩 ✅

- [x] Go 환경 설정
- [x] 프로젝트 문서 (PRD, 기술 설계, 로드맵)
- [x] go.mod 초기화
- [x] 최소 엔트리포인트 (`cmd/ccfg/main.go`)
- [x] 경로 상수 정의 (`internal/scanner/paths.go`)
- [x] git 초기화 + .gitignore

## Phase 2: 스캐너 구현 ✅

- [x] `internal/model/types.go` — 공유 타입 정의 (Scope, FileType, ConfigFile, ScanResult)
- [x] `internal/scanner/scanner.go` — DefaultScanner 구현
- [x] `internal/scanner/root.go` — 프로젝트 루트 감지 (.git 탐색)
- [x] `internal/scanner/scanner_test.go` — 테스트 (임시 디렉토리 기반)
- [x] Bubbletea, Lipgloss 등 의존성 추가

## Phase 3: 기본 TUI ✅

- [x] `internal/tui/model.go` — 메인 TUI 모델 (Init, Update, View)
- [x] `internal/tui/tree.go` — 트리 뷰 컴포넌트
- [x] `internal/tui/preview.go` — 미리보기 패널 (raw text)
- [x] `internal/tui/keys.go` — 키 바인딩 정의
- [x] `internal/tui/styles.go` — Lipgloss 스타일 정의
- [x] 30/70 split 레이아웃
- [x] 키보드 네비게이션 (j/k/Enter/Tab/q)

## Phase 4: 파서 + 구문 강조 ✅

- [x] `internal/parser/json.go` — JSON/JSONC 파서 (주석/trailing comma 제거)
- [x] `internal/parser/markdown.go` — Markdown 렌더러 (Glamour)
- [x] 구문 강조 (Chroma monokai 테마) 적용
- [x] 파싱 에러 시 원본 텍스트 폴백

## Phase 5: 확장 스캔 대상

- [ ] `~/.claude/commands/`, `.claude/commands/` — 커스텀 슬래시 명령어
- [ ] `~/.claude/skills/`, `.claude/skills/` — 설치된 에이전트 스킬
- [ ] `~/.claude/keybindings.json` — 키바인딩 설정
- [ ] settings.json 내부 `"hooks"` 키 파싱 → 훅 목록 트리 표시
- [ ] settings.json 내부 `"mcpServers"` 키 파싱 → MCP 서버 목록 표시
- [ ] 디렉토리 스캔 지원 (파일뿐 아니라 commands/, skills/ 하위 파일 나열)
- [ ] ConfigCategory 확장 (Commands, Skills, Keybindings 추가)

## Phase 6: 고급 기능

- [ ] 병합 뷰 (Merged View) — 모든 Scope 설정 병합 + 출처 표시
- [ ] 검색 기능 (`/` 키로 진입)
- [ ] Linux 경로 지원
- [ ] 파일 변경 감지 (fsnotify) 및 자동 갱신

## Phase 7: 배포 + 품질

- [ ] golangci-lint 설정 및 CI
- [ ] goreleaser 설정 (크로스 컴파일 바이너리)
- [ ] README.md 작성 (설치 방법, 스크린샷)
- [ ] Homebrew formula (선택)

## Phase 8: 게임 스타일 테마

- [ ] 컬러풀한 테두리 (Scope별 고유 색상)
- [ ] 이모지 아이콘 (Scope, 파일 상태, 카테고리)
- [ ] 섹션 타이틀 스타일 (테두리 위 라벨)
- [ ] 하단 키 힌트 바 (게임 HUD 스타일)
- [ ] 전체 컬러 팔레트 리디자인 (다크 배경 + 네온 액센트)
