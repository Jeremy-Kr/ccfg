# ccfg 기여 가이드

ccfg에 관심을 가져주셔서 감사합니다. 이 문서는 개발 환경 설정, 코딩 컨벤션, PR 프로세스를 안내합니다.

## 개발 환경 설정

### 사전 요구사항

- **Go 1.26+** — [Go 설치](https://go.dev/doc/install)
- **golangci-lint** (선택, 린트용) — [설치](https://golangci-lint.run/welcome/install/)

### 빌드 및 실행

```bash
# 빌드
go build -o ccfg ./cmd/ccfg

# 실행
./ccfg

# 또는 Makefile 사용
make build
make run
```

### 테스트

```bash
# 전체 테스트
go test ./...

# 특정 패키지 테스트
go test ./internal/scanner/...

# 상세 출력
go test -v ./...
```

### 린트

```bash
# 린트 실행
golangci-lint run

# 또는 Makefile 사용
make lint
```

## 프로젝트 구조

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

- **포맷팅:** Go 표준 스타일 (`gofmt`). 커밋 전 `gofmt`를 실행하세요.
- **캡슐화:** `internal/` 패키지를 사용합니다. 구현 세부사항을 외부에 노출하지 마세요.
- **에러 처리:** `fmt.Errorf("context: %w", err)` 패턴으로 에러를 래핑합니다.
- **인터페이스:** 구현하는 쪽이 아닌, 사용하는 쪽에서 인터페이스를 정의합니다.
- **테스트:** 테스트 파일(`*_test.go`)은 테스트 대상과 같은 패키지에 위치합니다.
- **네이밍:** Go 네이밍 컨벤션을 따릅니다 — 비공개는 `camelCase`, 공개는 `PascalCase`.
- **의존성:** 최소한의 의존성을 유지합니다. 꼭 필요한 라이브러리만 추가하세요.

## 커밋 메시지 컨벤션

[Conventional Commits](https://www.conventionalcommits.org/) 형식을 사용합니다:

```
type(scope): description
```

### 타입

| 타입 | 용도 |
|---|---|
| `feat` | 새로운 기능 |
| `fix` | 버그 수정 |
| `refactor` | 버그 수정이나 기능 추가가 아닌 코드 변경 |
| `docs` | 문서만 변경 |
| `test` | 테스트 추가 또는 수정 |
| `chore` | 유지보수 작업 (의존성, CI 등) |
| `style` | 코드 스타일 변경 (포맷팅, 로직 변경 없음) |
| `perf` | 성능 개선 |
| `ci` | CI/CD 변경 |

### 예시

```
feat(scanner): Linux 경로 지원 추가
fix(parser): JSONC 후행 쉼표 처리
docs(readme): 설치 방법 추가
refactor(tui): 트리 렌더링을 별도 컴포넌트로 분리
```

## Pull Request 프로세스

1. 저장소를 **포크**하고 `main`에서 기능 브랜치를 생성합니다.
2. 위의 코딩 컨벤션에 따라 **코드를 작성**합니다.
3. 새로운 기능에 대한 **테스트를 작성**합니다.
4. 제출 전 **테스트와 린트를 실행**합니다:
   ```bash
   go test ./...
   golangci-lint run
   ```
5. 명확한 제목과 설명으로 `main`을 대상으로 **PR을 생성**합니다.
6. 리뷰 피드백에 **신속하게 대응**합니다.

### PR 가이드라인

- PR을 집중적으로 유지하세요 — PR 하나에 논리적 변경 하나.
- PR이 무엇을 하는지, 왜 하는지 간단히 설명하세요.
- 관련 이슈가 있으면 참조하세요 (예: `Fixes #42`).
- 모든 CI 체크가 통과하는지 확인하세요.

## 핵심 제약

- **읽기 전용:** ccfg는 설정 파일을 절대 수정하지 않습니다. 이것은 핵심 설계 원칙입니다.
- **macOS + Linux:** 두 플랫폼을 모두 지원합니다. OS에 따라 경로 해석이 다릅니다.
- **단일 바이너리:** 외부 런타임 의존성이 없습니다.

## 질문?

궁금한 점이나 제안이 있으면 GitHub 이슈를 열어주세요.
