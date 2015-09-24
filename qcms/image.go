package qcms

import (
	"fmt"
	"image"
	"image/color"
)

func ImageTransformByProfile(src_image image.Image, src_prof, dst_prof *Profile) (image.Image, error) {
	var dst_image image.Image
	rect := src_image.Bounds()
	width := rect.Dx()
	height := rect.Dy()
	colorModel := src_image.ColorModel()
	// 今のところ RGBA, YCbCr のみ対応
	if (colorModel != color.YCbCrModel) && (colorModel != color.RGBAModel) {
		return nil, fmt.Errorf("ImageTransformByProfile: Unsupported ColorModel(%d)", colorModel)
	}
	var src_rgba *image.RGBA
	var src_ycbcr *image.YCbCr
	if colorModel == color.YCbCrModel {
		// YCbCr の場合は RGB に変換する
		src_ycbcr = src_image.(*image.YCbCr)
		src_rgba = image.NewRGBA(rect)
		DrawYCbCr(src_rgba, rect, src_ycbcr, image.Pt(0, 0))
	} else {
		src_rgba = src_image.(*image.RGBA) // type assertions
	}
	transform := CreateTransform(src_prof, DATA_RGB_8, dst_prof, DATA_RGB_8)
	defer transform.DeleteTransform()
	if transform == nil {
		return nil, fmt.Errorf("ImageTransformByProfile: CreateTransform Failedl(%d)", colorModel)
	}
	dst_rgba := image.NewRGBA(rect)
	src_pix := src_rgba.Pix
	dst_pix := dst_rgba.Pix
	len_pix := len(src_pix)
	transform.DoTransform(src_pix, dst_pix, len_pix)
	// YCbCr の場合は RGB から戻す
	if colorModel == color.YCbCrModel {
		dst_ycbcr := image.NewYCbCr(rect, src_ycbcr.SubsampleRatio)
		var x int
		var y int
		for y = 0; y < height; y++ {
			for x = 0; x < width; x++ {
				r, g, b, _ := dst_rgba.At(x, y).RGBA()
				yy, cb, cr := color.RGBToYCbCr(uint8(r), uint8(g), uint8(b))
				yi := dst_ycbcr.YOffset(x, y)
				ci := dst_ycbcr.COffset(x, y)
				dst_ycbcr.Y[yi] = yy
				dst_ycbcr.Cb[ci] = cb
				dst_ycbcr.Cr[ci] = cr
			}
		}
		dst_image = image.Image(dst_ycbcr)
	} else {
		dst_image = image.Image(dst_rgba)
	}

	return dst_image, nil
}

// XXX: copy & paste from golang src/image/draw/draw.go

func DrawYCbCr(dst *image.RGBA, r image.Rectangle, src *image.YCbCr, sp image.Point) (ok bool) {
	// An image.YCbCr is always fully opaque, and so if the mask is implicitly nil
	// (i.e. fully opaque) then the op is effectively always Src.
	x0 := (r.Min.X - dst.Rect.Min.X) * 4
	x1 := (r.Max.X - dst.Rect.Min.X) * 4
	y0 := r.Min.Y - dst.Rect.Min.Y
	y1 := r.Max.Y - dst.Rect.Min.Y
	switch src.SubsampleRatio {
	case image.YCbCrSubsampleRatio444:
		for y, sy := y0, sp.Y; y != y1; y, sy = y+1, sy+1 {
			dpix := dst.Pix[y*dst.Stride:]
			yi := (sy-src.Rect.Min.Y)*src.YStride + (sp.X - src.Rect.Min.X)
			ci := (sy-src.Rect.Min.Y)*src.CStride + (sp.X - src.Rect.Min.X)
			for x := x0; x != x1; x, yi, ci = x+4, yi+1, ci+1 {
				rr, gg, bb := color.YCbCrToRGB(src.Y[yi], src.Cb[ci], src.Cr[ci])
				dpix[x+0] = rr
				dpix[x+1] = gg
				dpix[x+2] = bb
				dpix[x+3] = 255
			}
		}
	case image.YCbCrSubsampleRatio422:
		for y, sy := y0, sp.Y; y != y1; y, sy = y+1, sy+1 {
			dpix := dst.Pix[y*dst.Stride:]
			yi := (sy-src.Rect.Min.Y)*src.YStride + (sp.X - src.Rect.Min.X)
			ciBase := (sy-src.Rect.Min.Y)*src.CStride - src.Rect.Min.X/2
			for x, sx := x0, sp.X; x != x1; x, sx, yi = x+4, sx+1, yi+1 {
				ci := ciBase + sx/2
				rr, gg, bb := color.YCbCrToRGB(src.Y[yi], src.Cb[ci], src.Cr[ci])
				dpix[x+0] = rr
				dpix[x+1] = gg
				dpix[x+2] = bb
				dpix[x+3] = 255
			}
		}
	case image.YCbCrSubsampleRatio420:
		for y, sy := y0, sp.Y; y != y1; y, sy = y+1, sy+1 {
			dpix := dst.Pix[y*dst.Stride:]
			yi := (sy-src.Rect.Min.Y)*src.YStride + (sp.X - src.Rect.Min.X)
			ciBase := (sy/2-src.Rect.Min.Y/2)*src.CStride - src.Rect.Min.X/2
			for x, sx := x0, sp.X; x != x1; x, sx, yi = x+4, sx+1, yi+1 {
				ci := ciBase + sx/2
				rr, gg, bb := color.YCbCrToRGB(src.Y[yi], src.Cb[ci], src.Cr[ci])
				dpix[x+0] = rr
				dpix[x+1] = gg
				dpix[x+2] = bb
				dpix[x+3] = 255
			}
		}
	case image.YCbCrSubsampleRatio440:
		for y, sy := y0, sp.Y; y != y1; y, sy = y+1, sy+1 {
			dpix := dst.Pix[y*dst.Stride:]
			yi := (sy-src.Rect.Min.Y)*src.YStride + (sp.X - src.Rect.Min.X)
			ci := (sy/2-src.Rect.Min.Y/2)*src.CStride + (sp.X - src.Rect.Min.X)
			for x := x0; x != x1; x, yi, ci = x+4, yi+1, ci+1 {
				rr, gg, bb := color.YCbCrToRGB(src.Y[yi], src.Cb[ci], src.Cr[ci])
				dpix[x+0] = rr
				dpix[x+1] = gg
				dpix[x+2] = bb
				dpix[x+3] = 255
			}
		}
	default:
		return false
	}
	return true
}
