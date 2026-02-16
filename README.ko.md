# ccfg

Claude Code 설정 파일을 한눈에 조회할 수 있는 TUI 대시보드.

Claude Code는 관리(managed), 사용자(user), 프로젝트(project) 세 가지 스코프에 걸쳐 8개 이상의 설정 파일을 분산 저장합니다. **ccfg**는 이 모든 파일을 하나의 터미널 인터페이스에 모아, 어떤 설정이 어디에서 오는지 즉시 확인할 수 있게 해줍니다.

## 기능

- **통합 뷰** — 관리·사용자·프로젝트 설정 파일을 나란히 조회
- **트리 탐색** — 스코프별로 정리된 설정 파일을 접고 펼치며 탐색
- **구문 강조** — JSON/JSONC는 Chroma로, Markdown은 Glamour로 렌더링
- **병합 뷰** — 최종 병합된 설정을 출처 표시와 함께 조회
- **검색** — 모든 파일에서 키 또는 값으로 설정 검색
- **자동 갱신** — fsnotify로 파일 변경을 감지하여 실시간 업데이트
- **확장 스캔** — 커스텀 명령어, 에이전트 스킬, 훅, MCP 서버, 키바인딩
- **읽기 전용** — 설정 파일을 절대 수정하지 않음

## 설치

### Homebrew

```bash
brew install jeremy-kr/tap/ccfg
```

### Go Install

```bash
go install github.com/jeremy-kr/ccfg/cmd/ccfg@latest
```

### 소스에서 빌드

```bash
git clone https://github.com/jeremy-kr/ccfg.git
cd ccfg
go build -o ccfg ./cmd/ccfg
```

## 사용법

프로젝트 디렉토리에서 `ccfg`를 실행합니다:

```bash
ccfg
```

### 키 바인딩

| 키 | 동작 |
|---|---|
| `j/k` 또는 `Up/Down` | 트리 항목 이동 |
| `Enter` | 노드 펼치기/접기 또는 파일 선택 |
| `Tab` | 좌/우 패널 전환 |
| `h/l` | 패널 전환 (vim 스타일) |
| `/` | 검색 모드 진입 |
| `m` | 병합 뷰 토글 |
| `q` / `Ctrl+C` | 종료 |

### 플래그

```bash
ccfg --version    # 버전 출력
```

## 스크린샷

<!-- TODO: 스크린샷 추가 -->

## 기술 스택

- **언어:** [Go](https://go.dev/) 1.26
- **TUI 프레임워크:** [Bubbletea](https://github.com/charmbracelet/bubbletea) v1.3.x
- **스타일링:** [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **구문 강조:** [Glamour](https://github.com/charmbracelet/glamour) + [Chroma](https://github.com/alecthomas/chroma)
- **파일 감시:** [fsnotify](https://github.com/fsnotify/fsnotify)

## 스캔 대상 파일

ccfg는 세 가지 스코프에서 설정 파일을 탐색합니다:

| 스코프 | 예시 |
|---|---|
| **Managed** | `/Library/Application Support/ClaudeCode/managed_settings.json`, `policies.json` |
| **User** | `~/.claude/settings.json`, `~/.claude/CLAUDE.md`, `~/.mcp.json` |
| **Project** | `.claude/settings.json`, `CLAUDE.md`, `.mcp.json` |

전체 목록은 [docs/PRD.md](docs/PRD.md)를 참고하세요.

## 기여하기

개발 환경 설정, 코딩 컨벤션, PR 프로세스는 [CONTRIBUTING.ko.md](CONTRIBUTING.ko.md)를 참고하세요.

## 라이선스

[MIT](LICENSE)
