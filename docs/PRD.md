# ccfg — Product Requirements Document

## 제품 개요

**ccfg**는 Claude Code의 분산된 설정 파일들을 단일 TUI 대시보드에서 조회하는 읽기 전용 도구다.

Claude Code는 설정이 8개 이상의 파일/디렉토리에 분산되어 있다:
- 시스템 관리자가 관리하는 managed 설정
- 사용자 홈 디렉토리의 전역 설정
- 프로젝트별 로컬 설정

어떤 설정이 어디서 오는지, 현재 활성화된 값이 무엇인지 파악하기 어렵다. ccfg는 이 문제를 해결한다.

## 타겟 사용자

- Claude Code를 일상적으로 사용하는 개발자
- 여러 프로젝트에서 서로 다른 Claude Code 설정을 사용하는 사용자
- MCP 서버, 훅, 퍼미션 설정을 자주 변경하는 파워 유저

## 스캔 대상 파일 (3개 Scope)

### Scope 1: Managed (시스템 관리자 설정)
| 파일 | 설명 |
|---|---|
| `/Library/Application Support/ClaudeCode/managed_settings.json` | 관리형 설정 (macOS) |
| `/Library/Application Support/ClaudeCode/policies.json` | 정책 파일 (macOS) |

### Scope 2: User (사용자 전역 설정)
| 파일 | 설명 |
|---|---|
| `~/.claude/settings.json` | 사용자 전역 설정 |
| `~/.claude/settings.local.json` | 사용자 로컬 설정 (gitignore 대상) |
| `~/.claude.json` | 레거시 전역 설정 |
| `~/.claude/CLAUDE.md` | 사용자 전역 지시사항 |
| `~/.mcp.json` | MCP 서버 전역 설정 |

### Scope 3: Project (프로젝트별 설정)
| 파일 | 설명 |
|---|---|
| `<root>/.claude/settings.json` | 프로젝트 설정 |
| `<root>/.claude/settings.local.json` | 프로젝트 로컬 설정 |
| `<root>/CLAUDE.md` | 프로젝트 지시사항 |
| `<root>/.claude/CLAUDE.md` | 프로젝트 지시사항 (대체 위치) |
| `<root>/.mcp.json` | MCP 서버 프로젝트 설정 |

## 기능 요구사항

### FR-1: 설정 파일 탐색
- 3개 Scope의 모든 설정 파일을 자동 탐색
- 파일 존재 여부, 크기, 수정 시간 표시
- 프로젝트 루트는 CWD 기준 자동 감지 (`.git` 디렉토리 탐색)

### FR-2: 트리 뷰 탐색
- 좌측 패널에 Scope > 파일 계층 트리 표시
- 키보드로 트리 탐색 (j/k/Enter/Esc)
- 파일 존재 여부에 따른 시각적 구분

### FR-3: 내용 미리보기
- 우측 패널에 선택한 파일의 내용 표시
- JSON/JSONC 파일: 구문 강조
- Markdown 파일: 렌더링된 형태로 표시

### FR-4: 병합 뷰 (Merged View)
- 모든 Scope의 JSON 설정을 병합한 최종 결과 표시
- 각 값의 출처 Scope 표시
- 우선순위: Project > User > Managed

### FR-5: 검색
- 설정 키/값 기반 검색
- 검색 결과 하이라이트

## 비기능 요구사항

- **성능:** 시작부터 TUI 표시까지 500ms 이내
- **읽기 전용:** 설정 파일을 절대 수정하지 않음
- **단일 바이너리:** 외부 런타임 의존성 없음
- **OS 지원:** macOS 우선, Linux 지원 예정
- **접근성:** 키보드 전용 탐색 가능
