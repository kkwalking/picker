package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// tappableRect 在 canvas.Rectangle 基础上实现 fyne.Tappable
type tappableRect struct {
	*canvas.Rectangle
	tapped func()
}

func (t *tappableRect) Tapped(*fyne.PointEvent) {
	if t.tapped != nil {
		t.tapped()
	}
}

func (t *tappableRect) TappedSecondary(*fyne.PointEvent) {}

// ModernButtonConfig 现代按钮配置
type ModernButtonConfig struct {
	Text        string
	NormalColor color.Color
	PressedColor color.Color
	TextColor   color.Color
	TextSize    float32
	MinSize     fyne.Size
	CornerRadius float32
	OnClick     func()
}

// NewModernButton 创建现代风格按钮
func NewModernButton(config ModernButtonConfig) fyne.CanvasObject {
	rect := canvas.NewRectangle(config.NormalColor)
	rect.CornerRadius = config.CornerRadius
	rect.SetMinSize(config.MinSize)

	tRect := &tappableRect{Rectangle: rect}

	label := canvas.NewText(config.Text, config.TextColor)
	label.Alignment = fyne.TextAlignCenter
	label.TextSize = config.TextSize

	// 点击回调（带短动画）
	tRect.tapped = func() {
		go func() {
			// 按下效果
			rect.FillColor = config.PressedColor
			label.TextSize = config.TextSize - 1
			canvas.Refresh(rect)
			canvas.Refresh(label)

			time.Sleep(140 * time.Millisecond)

			// 恢复
			rect.FillColor = config.NormalColor
			label.TextSize = config.TextSize
			canvas.Refresh(rect)
			canvas.Refresh(label)

			if config.OnClick != nil {
				config.OnClick()
			}
		}()
	}

	btnContainer := container.NewMax(rect, container.NewCenter(label), tRect)
	return btnContainer
}

// NewPrimaryButton 创建主要的绿色按钮
func NewPrimaryButton(text string, onClick func()) fyne.CanvasObject {
	return NewModernButton(ModernButtonConfig{
		Text:          text,
		NormalColor:   color.NRGBA{R: 33, G: 148, B: 83, A: 255},  // #219453
		PressedColor:  color.NRGBA{R: 27, G: 122, B: 69, A: 255},  // 深一点
		TextColor:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		TextSize:      16,
		MinSize:       fyne.NewSize(140, 44),
		CornerRadius:  8.0,
		OnClick:       onClick,
	})
}