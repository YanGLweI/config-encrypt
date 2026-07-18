package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/YanGLweI/config-encrypt/crypto"
)

func newDecryptPage() fyne.CanvasObject {
	// 私钥文件路径
	privKeyEntry := widget.NewEntry()
	privKeyEntry.SetPlaceHolder("选择私钥文件 (.pem)...")

	// 密文输入（自动换行，占据主要空间）
	cipherEntry := widget.NewMultiLineEntry()
	cipherEntry.SetPlaceHolder("请输入密文，如 ENC[base64...]")
	cipherEntry.Wrapping = fyne.TextWrapWord
	cipherEntry.SetMinRowsVisible(6)

	// 结果展示（多行文本框，自动换行，较小）
	resultEntry := widget.NewMultiLineEntry()
	resultEntry.SetPlaceHolder("解密结果将显示在这里...")
	resultEntry.Wrapping = fyne.TextWrapWord
	resultEntry.SetMinRowsVisible(3)

	// 复制按钮
	copyBtn := widget.NewButton("复制结果", func() {
		if resultEntry.Text != "" {
			fyne.CurrentApp().Clipboard().SetContent(resultEntry.Text)
		}
	})
	copyBtn.Importance = widget.LowImportance

	// 解密按钮
	decryptBtn := widget.NewButton("解密", func() {
		privKeyPath := privKeyEntry.Text
		if privKeyPath == "" {
			resultEntry.SetText("❌ 请选择私钥文件")
			return
		}
		ciphertext := cipherEntry.Text
		if ciphertext == "" {
			resultEntry.SetText("❌ 请输入密文")
			return
		}

		inner, ok := crypto.Unwrap(ciphertext)
		if !ok {
			resultEntry.SetText("❌ 密文格式不正确，应为 ENC[base64] 格式")
			return
		}

		privKey, err := crypto.LoadPrivateKey(privKeyPath)
		if err != nil {
			resultEntry.SetText("❌ 加载私钥失败: " + err.Error())
			return
		}

		plaintext, err := crypto.Decrypt(privKey, inner)
		if err != nil {
			resultEntry.SetText("❌ 解密失败: " + err.Error())
			return
		}

		resultEntry.SetText(plaintext)
	})
	decryptBtn.Importance = widget.HighImportance

	// 选择私钥文件按钮（优先使用原生选择器）
	browseBtn := widget.NewButton("浏览...", func() {
		path, ok := nativeOpenFile("选择私钥文件", "", []string{"pem"})
		if !ok {
			// 原生选择器不可用，回退到 Fyne 对话框
			dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
				if err != nil || uri == nil {
					return
				}
				privKeyEntry.SetText(uri.URI().Path())
			}, fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		if path != "" {
			privKeyEntry.SetText(path)
		}
	})

	privKeyRow := container.NewBorder(nil, nil, nil, browseBtn, privKeyEntry)

	// 顶部：私钥文件
	topContent := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("私钥文件", privKeyRow),
		),
	)

	// 中间：密文输入，自动换行，撑满剩余空间
	cipherLabel := widget.NewLabel("密文（ENC[...] 格式）:")
	cipherContainer := container.NewBorder(cipherLabel, nil, nil, nil, cipherEntry)
	cipherScroll := container.NewScroll(cipherContainer)

	// 底部：按钮 + 结果（小框）
	resultLabel := widget.NewLabel("解密结果:")
	resultRow := container.NewBorder(nil, nil, nil, copyBtn, resultEntry)
	bottomContent := container.NewVBox(
		decryptBtn,
		widget.NewSeparator(),
		resultLabel,
		resultRow,
	)

	form := container.NewBorder(topContent, bottomContent, nil, nil, cipherScroll)

	return container.NewPadded(form)
}
