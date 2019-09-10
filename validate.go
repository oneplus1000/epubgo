package epubgo

import (
	"fmt"
	"strings"
)

//ValidateErrorTypeFileNotFound file not found
var ValidateErrorTypeFileNotFound = 1

//ValidateError validate error
type ValidateError interface {
	error
	Type() int
}

//ValidateError error of validate
type validateError struct {
	errType int
	errMsg  string
}

func (v validateError) Type() int {
	return v.errType
}

func (v validateError) Error() string {
	return v.errMsg
}

//Validate validate epub file
func Validate(epub *Epub) error {

	err := validateNcx(epub)
	if err != nil {
		return err
	}

	return nil
}

func validateNcx(epub *Epub) error {

	navs := epub.ncx.navMap()
	for _, n := range navs {
		url := n.URL()
		err := checkZipContent(epub, url)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkZipContent(epub *Epub, url string) error {
	found := false
	for _, zf := range epub.zip.File {
		zfname := removePathOebps(zf.Name)
		//fmt.Printf("zf.Name:%s %s\n", zfname, url)
		if comparePath(zfname, url) {
			found = true
			break
		}
	}
	if !found {
		return validateError{
			errType: ValidateErrorTypeFileNotFound,
			errMsg:  fmt.Sprintf("file %s not found", url),
		}
	}
	return nil
}

func removePathOebps(src string) string {
	return strings.Replace(src, "OEBPS/", "", -1)
}

func comparePath(a, b string) bool {
	if a == b {
		return true
	}
	return false
}
