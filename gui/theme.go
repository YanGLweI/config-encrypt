package gui

import (
	"embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed assets/NotoSansCJKsc-Regular.otf
var chineseFontData embed.FS

// customTheme 自定义主题，使用 Noto Sans CJK SC 中文字体
type customTheme struct {
	chineseFont fyne.Resource
}

func newCustomTheme() *customTheme {
	t := &customTheme{}
	data, err := chineseFontData.ReadFile("assets/NotoSansCJKsc-Regular.otf")
	if err == nil {
		t.chineseFont = fyne.NewStaticResource("NotoSansCJKsc-Regular.otf", data)
	}
	return t
}

func (t *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (t *customTheme) Font(style fyne.TextStyle) fyne.Resource {
	if t.chineseFont != nil {
		return t.chineseFont
	}
	return theme.DefaultTheme().Font(style)
}

func (t *customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *customTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
