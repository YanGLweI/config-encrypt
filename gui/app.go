package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const version = "1.0.0"

// navItem 侧边栏导航项
type navItem struct {
	label string
	icon  fyne.Resource
	build func() fyne.CanvasObject
}

// Run 启动 GUI 应用
func Run() {
	a := app.NewWithID("com.config-encrypt.gui")
	a.Settings().SetTheme(newCustomTheme())
	w := a.NewWindow("配置加密工具")
	w.SetMaster()
	w.Resize(fyne.NewSize(780, 560))

	// 设置窗口图标
	if icon := AppIcon(); icon != nil {
		w.SetIcon(icon)
	}

	// 定义导航项
	items := []navItem{
		{label: "密钥生成", icon: theme.FolderNewIcon(), build: newKeygenPage},
		{label: "加密", icon: theme.FileApplicationIcon(), build: newEncryptPage},
		{label: "解密", icon: theme.FileTextIcon(), build: newDecryptPage},
	}

	// 内容容器（右侧）
	content := container.NewStack()

	// 构建侧边栏按钮
	var navButtons []*widget.Button
	for i, item := range items {
		btn := widget.NewButtonWithIcon(item.label, item.icon, func(idx int) func() {
			return func() {
				// 高亮选中按钮
				for j, b := range navButtons {
					if j == idx {
						b.Importance = widget.HighImportance
					} else {
						b.Importance = widget.MediumImportance
					}
					b.Refresh()
				}
				// 切换内容
				content.Objects = []fyne.CanvasObject{items[idx].build()}
				content.Refresh()
			}
		}(i))
		if i == 0 {
			btn.Importance = widget.HighImportance
		}
		navButtons = append(navButtons, btn)
	}

	// 侧边栏容器
	sidebarItems := make([]fyne.CanvasObject, len(navButtons))
	for i, btn := range navButtons {
		sidebarItems[i] = btn
	}
	sidebar := container.NewVBox(sidebarItems...)
	sidebarBox := container.NewVBox(
		widget.NewLabelWithStyle("配置加密工具", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		sidebar,
		layout.NewSpacer(),
		widget.NewLabelWithStyle("v"+version, fyne.TextAlignCenter, fyne.TextStyle{}),
	)
	sidebarContainer := container.NewVScroll(sidebarBox)
	sidebarContainer.Resize(fyne.NewSize(160, 500))

	// 默认显示第一个页面
	content.Objects = []fyne.CanvasObject{items[0].build()}
	content.Refresh()

	// 主布局：左侧边栏 + 右侧内容
	split := container.NewHSplit(sidebarContainer, content)
	split.SetOffset(0.22)

	w.SetContent(container.NewBorder(nil, nil, nil, nil, split))
	w.ShowAndRun()
}
