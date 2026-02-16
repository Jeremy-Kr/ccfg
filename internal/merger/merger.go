package merger

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jeremy-kr/ccfg/internal/model"
	"github.com/jeremy-kr/ccfg/internal/parser"
)

// SourcedValue는 값과 그 출처를 나타낸다.
type SourcedValue struct {
	Key   string      // 점 표기법 경로 (예: "permissions.allow")
	Value any         // 실제 값
	Scope model.Scope // 값의 출처 Scope
}

// MergedConfig는 모든 Scope의 설정을 병합한 결과다.
type MergedConfig struct {
	Values []SourcedValue // 플랫한 키-값 목록
}

// Merge는 ScanResult의 JSON 설정 파일들을 우선순위에 따라 병합한다.
// 우선순위: Project > User > Managed
func Merge(result *model.ScanResult) *MergedConfig {
	merged := make(map[string]SourcedValue)

	// 낮은 우선순위부터 적용 (나중 것이 덮어씀)
	applyScope(merged, result.Managed, model.ScopeManaged)
	applyScope(merged, result.User, model.ScopeUser)
	applyScope(merged, result.Project, model.ScopeProject)

	// 키 기준 정렬
	values := make([]SourcedValue, 0, len(merged))
	for _, v := range merged {
		values = append(values, v)
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i].Key < values[j].Key
	})

	return &MergedConfig{Values: values}
}

func applyScope(merged map[string]SourcedValue, files []model.ConfigFile, scope model.Scope) {
	for _, f := range files {
		if !f.Exists || f.FileType != model.FileTypeJSON || f.IsDir {
			continue
		}
		if f.Category != model.CategorySettings {
			continue
		}

		data, err := os.ReadFile(f.Path)
		if err != nil {
			continue
		}

		cleaned := parser.StripJSONC(string(data))
		var obj map[string]any
		if err := json.Unmarshal([]byte(cleaned), &obj); err != nil {
			continue
		}

		flatten("", obj, scope, merged)
	}
}

func flatten(prefix string, obj map[string]any, scope model.Scope, out map[string]SourcedValue) {
	for k, v := range obj {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]any:
			flatten(key, val, scope, out)
		default:
			out[key] = SourcedValue{Key: key, Value: v, Scope: scope}
			_ = val
		}
	}
}

// Render는 병합된 설정을 읽기 좋은 문자열로 포맷한다.
func (mc *MergedConfig) Render() string {
	if len(mc.Values) == 0 {
		return "(병합할 설정이 없습니다)"
	}

	var b strings.Builder
	b.WriteString("Merged Settings (Project > User > Managed)\n")
	b.WriteString(strings.Repeat("─", 50) + "\n\n")

	for _, v := range mc.Values {
		valStr := fmt.Sprintf("%v", v.Value)
		if len(valStr) > 60 {
			valStr = valStr[:57] + "..."
		}
		b.WriteString(fmt.Sprintf("  %-35s = %-20s [%s]\n", v.Key, valStr, v.Scope))
	}

	return b.String()
}
