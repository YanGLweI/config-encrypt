package gui

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/YanGLweI/config-encrypt/crypto"
)

func newKeygenPage() fyne.CanvasObject {
	// 密钥位数下拉选择
	bitsSelect := widget.NewSelect([]string{"2048", "4096"}, nil)
	bitsSelect.SetSelected("2048")

	// 输出文件名前缀
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("config-key")
	nameEntry.SetText("config-key")

	// 保存目录
	dirEntry := widget.NewEntry()
	dirEntry.SetPlaceHolder("选择保存目录...")
	dirEntry.SetText("")

	// 结果展示
	resultLabel := widget.NewLabel("")
	resultLabel.Wrapping = fyne.TextWrapWord

	// 生成按钮
	generateBtn := widget.NewButton("生成密钥对", func() {
		bits := 2048
		if bitsSelect.Selected == "4096" {
			bits = 4096
		}
		output := nameEntry.Text
		if output == "" {
			output = "config-key"
		}
		saveDir := dirEntry.Text
		if saveDir == "" {
			resultLabel.SetText("❌ 请选择保存目录")
			return
		}

		privateKeyPath := filepath.Join(saveDir, output+".pem")
		publicKeyPath := filepath.Join(saveDir, output+".pub.pem")

		resultLabel.SetText("⏳ 正在生成密钥对...")
		if err := crypto.GenerateKeyPair(bits, privateKeyPath, publicKeyPath); err != nil {
			resultLabel.SetText("❌ 生成失败: " + err.Error())
			return
		}

		resultLabel.SetText("✅ 密钥对生成成功！\n\n" +
			"私钥: " + privateKeyPath + "\n" +
			"公钥: " + publicKeyPath + "\n\n" +
			"请将私钥复制到后端 key/ 目录，使用公钥加密 config.yml 中的敏感字段。")
	})
	generateBtn.Importance = widget.HighImportance

	// 选择目录按钮
	browseBtn := widget.NewButton("浏览...", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			dirEntry.SetText(uri.Path())
		}, fyne.CurrentApp().Driver().AllWindows()[0])
	})

	dirRow := container.NewBorder(nil, nil, nil, browseBtn, dirEntry)

	form := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("密钥位数", bitsSelect),
			widget.NewFormItem("文件名前缀", nameEntry),
			widget.NewFormItem("保存目录", dirRow),
		),
		generateBtn,
		widget.NewSeparator(),
		resultLabel,
		layout.NewSpacer(),
	)

	return container.NewPadded(form)
}
