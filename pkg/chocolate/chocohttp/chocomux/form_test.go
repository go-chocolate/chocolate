package chocomux

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"net/url"
	"testing"
)

type A struct {
	Name  string `form:"name"`
	Value int32
}

func TestUnmarshalForm(t *testing.T) {
	form := url.Values{}
	form.Set("name", "zhangsan")
	form.Set("age", "19")
	form.Set("Value", "128")
	form.Add("Tag", "a")
	form.Add("Tag", "b")
	form.Add("Group", "65536")

	{
		type UnmarshalFormTest struct {
			A
			Tag    []string
			Group  [2]int64
			Gender int8
		}
		request := UnmarshalFormTest{}

		if err := UnmarshalForm(form, &request); err != nil {
			t.Error(err)
		} else {
			t.Log(request)
		}
		if request.Name != "zhangsan" {
			t.Error("name error")
		}
		if request.Value != 128 {
			t.Error("value error")
		}
		if len(request.Tag) != 2 || request.Tag[0] != "a" || request.Tag[1] != "b" {
			t.Error("tag error")
		}
		if request.Group[0] != 65536 || request.Group[1] != 0 {
			t.Error("group error")
		}
	}
	{
		var request = make(map[string][]string)
		if err := UnmarshalForm(form, &request); err != nil {
			t.Error(err)
		}
		t.Log(request)

	}
}

func TestUnmarshalMultipartForm(t *testing.T) {
	body := bytes.NewBuffer(nil)
	w := multipart.NewWriter(body)
	w.WriteField("name", "zhangsan")
	w.WriteField("age", "22")
	w.WriteField("Value", "3654")
	w.WriteField("Tag", "a")
	w.WriteField("Tag", "b")
	w.WriteField("Tag", "c")
	w.WriteField("Group", "693.16514")
	fw, _ := w.CreateFormFile("file", "example.txt")
	fw.Write([]byte("hello world"))
	w.Close()
	request := httptest.NewRequest("POST", "http://example", body)
	request.Header.Add("Content-Type", w.FormDataContentType())
	if err := request.ParseMultipartForm(1024); err != nil {
		t.Error(err)
		return
	}

	type UnmarshalMultipartFormTest struct {
		A
		Tag    []string              ``
		Group  [2]float64            ``
		Gender int8                  ``
		File   *multipart.FileHeader `form:"file"`
	}

	form := request.MultipartForm
	binding := UnmarshalMultipartFormTest{}

	if err := UnmarshalMultipartForm(form, &binding); err != nil {
		t.Error(err)
	} else {
		t.Log(binding)
	}

	if binding.Name != "zhangsan" {
		t.Error("name not match")
	}
	if binding.Value != 3654 {
		t.Error("value not match")
	}
	if len(binding.Tag) != 3 || binding.Tag[0] != "a" || binding.Tag[1] != "b" || binding.Tag[2] != "c" {
		t.Error("tag length not match")
	}
	if binding.Group[0] != 693.16514 || binding.Group[1] != 0 {
		t.Error("group not match")
	}
	if binding.File == nil {
		t.Error("file not match")
		if file, err := binding.File.Open(); err != nil {
			t.Error(err)
		} else {
			b, err := io.ReadAll(file)
			if err != nil || string(b) != "hello world" {
				t.Fail()
			}
		}
	}
}
