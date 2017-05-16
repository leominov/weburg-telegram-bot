package bot

import "testing"

func TestFormatCategoryName(t *testing.T) {
	tests := []struct {
		In  string
		Out string
	}{
		{
			"Movie",
			"#Movie",
		},
		{
			"Pop Rock",
			"#Pop_Rock",
		},
		{
			"18+",
			"#18",
		},
	}
	a := Agent{}
	for i, tt := range tests {
		result := a.FormatCategoryName(tt.In)
		if result != tt.Out {
			t.Errorf("#%d: expected result='%s', got '%s'", i, tt.Out, result)
		}
	}
}

func TestClearCategories(t *testing.T) {
	tests := []struct {
		Categories             []string
		SkipCategories         []string
		ResultCategoryListSize int
	}{
		{
			[]string{"skipped1", "nonskipped1", "skipped2", "nonskipped2"},
			[]string{"skipped1", "skipped2"},
			2,
		},
		{
			[]string{"skipped1", "skipped2"},
			[]string{"skipped1", "skipped2"},
			0,
		},
		{
			[]string{"skipped1", "skipped2"},
			[]string{},
			2,
		},
	}
	for i, tt := range tests {
		a := Agent{
			SkipCategories: tt.SkipCategories,
		}
		result := a.ClearCategories(tt.Categories)
		if len(result) != tt.ResultCategoryListSize {
			t.Errorf("#%d: expected len(result)=%d, got %d", i, tt.ResultCategoryListSize, len(result))
		}
	}
}
