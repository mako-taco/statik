package checkstyle

import (
	"strings"
	"testing"

	"github.com/statik/pkg/plugin"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []plugin.AnalysisResult
		wantErr bool
	}{
		{
			name: "valid checkstyle output",
			input: `<?xml version="1.0" encoding="UTF-8"?>
<checkstyle version="8.44">
  <file name="src/main/java/com/example/Test.java">
    <error line="10" column="5" severity="error" message="Missing a Javadoc comment" source="com.puppycrawl.tools.checkstyle.checks.javadoc.JavadocMethodCheck"/>
    <error line="15" column="20" severity="warning" message="Parameter 'param' should be final" source="com.puppycrawl.tools.checkstyle.checks.FinalParametersCheck"/>
  </file>
</checkstyle>`,
			want: []plugin.AnalysisResult{
				{
					Tool:        "checkstyle",
					File:        "src/main/java/com/example/Test.java",
					Line:        10,
					Column:      5,
					Message:     "Missing a Javadoc comment",
					Severity:    plugin.SeverityError,
					RuleID:      "com.puppycrawl.tools.checkstyle.checks.javadoc.JavadocMethodCheck",
					Description: "Missing a Javadoc comment",
				},
				{
					Tool:        "checkstyle",
					File:        "src/main/java/com/example/Test.java",
					Line:        15,
					Column:      20,
					Message:     "Parameter 'param' should be final",
					Severity:    plugin.SeverityWarning,
					RuleID:      "com.puppycrawl.tools.checkstyle.checks.FinalParametersCheck",
					Description: "Parameter 'param' should be final",
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid XML",
			input:   "not xml",
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty checkstyle output",
			input: `<?xml version="1.0" encoding="UTF-8"?>
<checkstyle version="8.44">
</checkstyle>`,
			want:    []plugin.AnalysisResult{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, err := p.Parse(strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != len(tt.want) {
				t.Errorf("Parser.Parse() got %d results, want %d", len(got), len(tt.want))
				return
			}
			for i, result := range got {
				if result != tt.want[i] {
					t.Errorf("Parser.Parse() result[%d] = %v, want %v", i, result, tt.want[i])
				}
			}
		})
	}
}

func TestParser_Name(t *testing.T) {
	p := &Parser{}
	if got := p.Name(); got != "checkstyle" {
		t.Errorf("Parser.Name() = %v, want %v", got, "checkstyle")
	}
}

func TestParser_SupportedFileExtensions(t *testing.T) {
	p := &Parser{}
	want := []string{".java", ".xml", ".properties"}
	got := p.SupportedFileExtensions()
	if len(got) != len(want) {
		t.Errorf("Parser.SupportedFileExtensions() got %d extensions, want %d", len(got), len(want))
		return
	}
	for i, ext := range got {
		if ext != want[i] {
			t.Errorf("Parser.SupportedFileExtensions()[%d] = %v, want %v", i, ext, want[i])
		}
	}
} 