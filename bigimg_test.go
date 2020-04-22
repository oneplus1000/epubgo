package epubgo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	_ "image/jpeg"
	_ "image/png"
)

func TestBigImg1(t *testing.T) {

	eb, err := Open("./testdata/too_bigimage.epub")
	if err != nil {
		t.Fatalf("%+v", err)
		return
	}
	defer eb.Close()
	errs := Validate(eb, &Condition{
		MaxImageSizeByte: 100000000, //~100MB
	})
	if errs != nil {
		var buff bytes.Buffer
		for _, e := range errs {
			fmt.Fprintf(&buff, "%+v\n", e)
		}
		t.Fatalf("%s", buff.String())
		return
	}
}

func TestBigImg2(t *testing.T) {

	data, err := ioutil.ReadFile("./testdata/too_bigimage.epub")
	if err != nil {
		t.Fatalf("%+v", err)
		return
	}

	r := bytes.NewReader(data)
	eb, err := Load(r, int64(r.Len())) //Open("./testdata/too_bigimage.epub")
	if err != nil {
		t.Fatalf("%+v", err)
		return
	}
	defer eb.Close()
	errs := Validate(eb, &Condition{
		MaxImageSizeByte: 100000000, //~100MB
	})
	if errs != nil {
		var buff bytes.Buffer
		for _, e := range errs {
			fmt.Fprintf(&buff, "%+v\n", e)
		}
		t.Fatalf("%s", buff.String())
		return
	}
}
