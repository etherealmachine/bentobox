package bento

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/etherealmachine/bento/text"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (n *Box) Draw(img *ebiten.Image) {
	if n.style.hidden() || !n.style.display() {
		return
	}
	if n.parent == nil {
		n.size()
		n.grow()
		n.justify()
	}
	content := n.contentRect()
	inner := n.innerRect()
	outer := n.outerRect()
	if n.debug {
		// Outer
		drawBox(img, outer, color.White, true)
		// Inner
		drawBox(img, inner, &color.RGBA{R: 200, G: 200, B: 200, A: 255}, true)
	}
	if n.style != nil && n.style.Border != nil {
		n.style.Border.Draw(img, inner.Min.X, inner.Min.Y, inner.Dx(), inner.Dy())
	}
	if n.debug {
		// Content
		drawBox(img, content, &color.RGBA{R: 100, G: 100, B: 100, A: 255}, true)
	}
	switch n.tag {
	case "button":
		n.style.Button[int(n.buttonState)].Draw(img, inner.Min.X, inner.Min.Y, inner.Dx(), inner.Dy())
		text.DrawString(img, n.templateContent(), n.style.Font, n.style.Color, content, text.Center, text.Center)
	case "text":
		text.DrawString(img, n.templateContent(), n.style.Font, n.style.Color, content, text.Center, text.Center)
	case "p":
		text.DrawParagraph(img, n.templateContent(), n.style.Font, n.style.Color, content.Min.X, content.Min.Y, n.style.MaxWidth, -n.TextBounds.Min.Y)
	case "img":
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(content.Min.X), float64(content.Min.Y))
		img.DrawImage(n.style.Image, op)
	case "input":
		n.style.Input[int(n.inputState)].Draw(img, inner.Min.X, inner.Min.Y, inner.Dx(), inner.Dy())
		text.DrawString(img, n.attrs["placeholder"], n.style.Font, n.style.Color, content, text.Center, text.Center)
	case "textarea":
		text.DrawParagraph(img, n.attrs["value"], n.style.Font, n.style.Color, content.Min.X, content.Min.Y, n.style.MaxWidth, -n.TextBounds.Min.Y)
	case "row", "col":
	default:
		log.Fatalf("can't draw %s", n.tag)
	}
	if n.debug {
		// Annotate
		text.DrawString(
			img,
			fmt.Sprintf("%s %dx%d", n.tag, n.OuterWidth, n.OuterHeight),
			text.Font("mono", 10), color.Black, outer.Add(image.Pt(4, 4)), text.Start, text.Start)
	}
	for _, c := range n.children {
		c.Draw(img)
	}
}

func drawBox(img *ebiten.Image, rect image.Rectangle, c color.Color, border bool) {
	x, y := float64(rect.Min.X), float64(rect.Min.Y)
	w, h := float64(rect.Dx()), float64(rect.Dy())
	ebitenutil.DrawRect(img, x, y, w, h, c)
	if border {
		ebitenutil.DrawLine(img, x, y, x+w, y, color.Black)
		ebitenutil.DrawLine(img, x, y, x, y+h, color.Black)
		ebitenutil.DrawLine(img, x+w, y+h, x+w, y, color.Black)
		ebitenutil.DrawLine(img, x+w, y+h, x, y+h, color.Black)
	}
}
