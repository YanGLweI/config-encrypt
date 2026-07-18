package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/YanGLweI/config-encrypt/crypto"
)

func newEncryptPage() fyne.CanvasObject {
	// 公钥文件路径
	pubKeyEntry := widget.NewEntry()
	pubKeyEntry.SetPlaceHolder("选择公钥文件 (.pem)...")

	// 密码输入
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("请输入要加密的密码")

	// 结果展示
	resultEntry := widget.NewEntry()
	resultEntry.SetPlaceHolder("加密结果将显示在这里...")
	resultEntry.Disable() // 只读

	// 复制按钮
	copyBtn := widget.NewButton("复制结果", func() {
		if resultEntry.Text != "" {
			fyne.CurrentApp().Clipboard().SetContent(resultEntry.Text)
		}
	})
	copyBtn.Importance = widget.LowImportance

	// 加密按钮
	encryptBtn := widget.NewButton("加密", func() {
		pubKeyPath := pubKeyEntry.Text
		if pubKeyPath == "" {
			resultEntry.SetText("❌ 请选择公钥文件")
			return
		}
		password := passwordEntry.Text
		if password == "" {
			resultEntry.SetText("❌ 请输入密码")
			return
		}

		pubKey, err := crypto.LoadPublicKey(pubKeyPath)
		if err != nil {
			resultEntry.SetText("❌ 加载公钥失败: " + err.Error())
			return
		}

		ciphertext, err := crypto.Encrypt(pubKey, password)
		if err != nil {
			resultEntry.SetText("❌ 加密失败: " + err.Error())
			return
		}

		result := crypto.Wrap(ciphertext)
		resultEntry.SetText(result)
	})
	encryptBtn.Importance = widget.HighImportance

	// 选择公钥文件按钮（优先使用 macOS 原生选择器）
	browseBtn := widget.NewButton("浏览...", func() {
		path := nativeOpenFile("选择公钥文件", "", []string{"pem"})
		if path != "" {
			pubKeyEntry.SetText(path)
			return
		}
		// 非 macOS 回退到 Fyne 对话框
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err != nil || uri == nil {
				return
			}
			pubKeyEntry.SetText(uri.URI().Path())
		}, fyne.CurrentApp().Driver().AllWindows()[0])
	})

	pubKeyRow := container.NewBorder(nil, nil, nil, browseBtn, pubKeyEntry)

	resultRow := container.NewBorder(nil, nil, nil, copyBtn, resultEntry)

	form := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("公钥文件", pubKeyRow),
			widget.NewFormItem("密码", passwordEntry),
		),
		encryptBtn,
		widget.NewSeparator(),
		widget.NewLabel("加密结果（复制到 config.yml 中）:"),
		resultRow,
		layout.NewSpacer(),
	)

	return container.NewPadded(form)
}
