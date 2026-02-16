# ccfg — Technical Design

## 핵심 타입 정의

### Scope

설정 파일의 적용 범위를 나타내는 열거형.

```go
type Scope int

const (
    ScopeManaged Scope = iota  // 시스템 관리자 설정
    ScopeUser                   // 사용자 전역 설정
    ScopeProject                // 프로젝트별 설정
)
```

### FileType

설정 파일의 형식.

```go
type FileType int

const (
    FileTypeJSON     FileType = iota  // .json
    FileTypeJSONC                      // .json with comments
    FileTypeMarkdown                   // .md
)
```

### ConfigFile

스캔된 개별 설정 파일.

```go
type ConfigFile struct {
    Path        string    // 절대 경로
    Scope       Scope     // 소속 Scope
    FileType    FileType  // 파일 형식
    Exists      bool      // 파일 존재 여부
    Size        int64     // 바이트 크기 (존재 시)
    ModTime     time.Time // 최종 수정 시간 (존재 시)
    Description string    // 사용자에게 보여줄 설명
}
```

### ScanResult

전체 스캔 결과.

```go
type ScanResult struct {
    Managed  []ConfigFile  // Managed Scope 파일들
    User     []ConfigFile  // User Scope 파일들
    Project  []ConfigFile  // Project Scope 파일들
    RootDir  string        // 감지된 프로젝트 루트
}
```

## 컴포넌트 인터페이스

### Scanner

설정 파일 탐색 담당.

```go
// Scanner는 Claude Code 설정 파일을 탐색한다.
type Scanner interface {
    Scan() (*ScanResult, error)
}
```

구현체: `scanner.DefaultScanner` — OS별 경로를 조회하고 파일 메타데이터 수집.

### Parser

파일 내용 파싱 담당.

```go
// Parser는 설정 파일 내용을 파싱한다.
type Parser interface {
    Parse(file ConfigFile) (any, error)
}
```

- `JSONParser`: JSON/JSONC 파일을 `map[string]any`로 파싱
- `MarkdownParser`: Markdown 파일을 문자열로 읽기

### Merger

여러 Scope의 JSON 설정을 병합.

```go
// Merger는 여러 Scope의 설정을 우선순위에 따라 병합한다.
type Merger interface {
    Merge(result *ScanResult) (*MergedConfig, error)
}
```

병합 우선순위: Project > User > Managed (높은 우선순위가 낮은 것을 오버라이드).

## TUI 아키텍처

Elm Architecture (The Elm Architecture, TEA) 기반 — Bubbletea의 기본 패턴.

```
Model → Update(Msg) → (Model, Cmd) → View() → string
```

### 레이아웃

```
┌─────────────────────────────────────────────┐
│  ccfg v0.1.0 — Claude Code Config Viewer    │  ← 헤더
├──────────────┬──────────────────────────────┤
│              │                              │
│  [Tree View] │  [Content Preview]           │  ← 메인 (30/70)
│  30% width   │  70% width                   │
│              │                              │
│  ▸ Managed   │  {                           │
│    settings  │    "permissions": { ... }    │
│  ▾ User      │    "hooks": { ... }          │
│    settings  │  }                           │
│    CLAUDE.md │                              │
│  ▸ Project   │                              │
│              │                              │
├──────────────┴──────────────────────────────┤
│  [Status Bar]  Tab: switch pane  q: quit    │  ← 풋터
└─────────────────────────────────────────────┘
```

### 키 바인딩

| 키 | 동작 |
|---|---|
| `j/k` 또는 `↑/↓` | 트리 항목 이동 |
| `Enter` | 노드 펼치기/접기 또는 파일 선택 |
| `Tab` | 좌/우 패널 전환 |
| `h/l` | 좌/우 패널 전환 (vim 스타일) |
| `/` | 검색 모드 진입 |
| `m` | 병합 뷰 토글 |
| `q` / `Ctrl+C` | 종료 |

### TUI 모델 구조

```go
type Model struct {
    scan       *ScanResult
    tree       TreeModel      // 좌측 트리 상태
    preview    PreviewModel   // 우측 미리보기 상태
    focus      Pane           // 현재 포커스 패널
    width      int            // 터미널 너비
    height     int            // 터미널 높이
    searchMode bool           // 검색 모드 활성화 여부
    searchText string         // 검색어
}
```

## 데이터 플로우

```
main()
  ├─ flag.Parse()              // CLI 인자 처리
  ├─ scanner.New()             // Scanner 생성
  ├─ scanner.Scan()            // 설정 파일 탐색
  │    ├─ ManagedPaths()       // Managed 경로 조회
  │    ├─ UserPaths()          // User 경로 조회
  │    ├─ ProjectPaths(root)   // Project 경로 조회
  │    └─ stat each file       // 각 파일 메타데이터 수집
  ├─ tui.NewModel(scanResult)  // TUI 모델 초기화
  └─ tea.NewProgram(model)     // Bubbletea 실행
       └─ Run()
```

## 에러 처리 전략

- 파일이 없는 경우: `ConfigFile.Exists = false`로 표시, 에러 아님
- 파싱 실패: 미리보기에 에러 메시지 표시, 원본 텍스트 폴백
- 권한 없음: 해당 파일을 "접근 불가"로 표시
- 프로젝트 루트 미감지: Project Scope를 빈 상태로 표시
