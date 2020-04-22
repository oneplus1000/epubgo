package epubgo

import (
	"archive/zip"
	"fmt"
	"image"
	"image/color"
	"strings"
)

const mineJPEG = "jpeg"
const mimePNG = "png"
const mimeNotCare = "" // mime type นี้เราไม่สนใจ

var imgMimes = map[string]string{
	"image/jpeg": mineJPEG,
	"image/jpg":  mineJPEG,
	"image/png":  mimePNG,
	//gif เราไม่สนคิดว่าขนาดคงไม่ใหญ่นะ
}

//validateImage ตรวจสอบขนาดของรูปว่าไม่ให้ใหญ่เกินไปจนทำให้ client ที่มี ram ต่ำทำงานไม่ได้
func validateImage(epub *Epub, condition *Condition) []error {
	var errs []error
	for _, m := range epub.opf.Manifest {
		mime := findImageFile(m.Href, m.MediaType)
		if mime != mimeNotCare {
			err := analyzeImage(epub, condition, m.Href, mime)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func findImageFile(href string, mediaType string) string {
	mediaType = strings.ToLower(mediaType)
	mediaType = strings.TrimSpace(mediaType)
	if val, ok := imgMimes[mediaType]; ok {
		return val
	}
	return mimeNotCare
}

func analyzeImage(epub *Epub, condition *Condition, href string, mime string) error {
	found := false
	var zipf *zip.File
	for _, z := range epub.zip.File {
		if href == removePathOebps(z.Name) {
			zipf = z
			found = true
			break
		}
	}

	if !found {
		return validateError{
			errType: ValidateErrorTypeFileNotFound,
			errMsg:  fmt.Sprintf("%s not found in epub", href),
		}
	}

	rd, err := zipf.Open()
	if err != nil {
		return validateError{
			errType: ValidateErrorTypeFileDamaged,
			errMsg:  fmt.Sprintf("path %s can not open %v", href, err),
		}
	}
	defer rd.Close()

	cfg, fileType, err := image.DecodeConfig(rd)
	if err != nil {
		return validateError{
			errType: ValidateErrorTypeFileDamaged,
			errMsg:  fmt.Sprintf("path %s can decode : %v", href, err),
		}
	} else if fileType != mime {
		return validateError{
			errType: ValidateErrorTypeFileDamaged,
			errMsg:  fmt.Sprintf("path %s media-type in manifest is %s but real file is %s : %+v", href, mime, fileType, err),
		}
	}

	if condition != nil {
		if condition.MaxImageSizeByte != 0 {
			cm := colorModelColor(cfg.ColorModel)
			bytePerPixel := colorToSizeBytePerPixel(cm)
			sizeinMemory := int64(bytePerPixel * cfg.Height * cfg.Width)
			if sizeinMemory > condition.MaxImageSizeByte {
				//size เกินขนาด
				return validateError{
					errType: ValidateErrorTypeOverMaxImageSize,
					errMsg:  fmt.Sprintf("file %s size over limit file size = %d but after decode size = %d in memory try to reduce width and height of image or set color profile to grayscale if image is grayscale image", href, zipf.CompressedSize, sizeinMemory),
				}
			}
		}
	}

	return nil
}

func colorModelColor(cmodel color.Model) color.Color {
	col := color.RGBA{} // This is the "any" color we convert
	//fmt.Printf("%T\n", cmodel.Convert(col))
	return cmodel.Convert(col)
}

//ขนาดนี้ base บน 8 bits per channel image นะถ้าเป็น 16 bits per channel image จะใหญ่เพิ่มเป็น 2 เท่า
func colorToSizeBytePerPixel(c color.Color) int {
	if _, ok := c.(color.RGBA); ok {
		return 4
	} else if _, ok := c.(color.RGBA64); ok {
		return 8
	} else if _, ok := c.(color.NRGBA64); ok {
		return 8
	} else if _, ok := c.(color.CMYK); ok {
		return 4
	} else if _, ok := c.(color.Gray); ok {
		return 1
	} else if _, ok := c.(color.Gray16); ok {
		return 2
	} else if _, ok := c.(color.YCbCr); ok {
		return 3
	} else if _, ok := c.(color.Alpha16); ok {
		return 2
	} else if _, ok := c.(color.Alpha); ok {
		return 1
	}
	return 2 //อันอื่นๆที่ไม่รู้จักขอมั่วๆตัวเลขเลยนะ
}
