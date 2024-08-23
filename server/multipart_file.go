package server

import (
	"bytes"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"main-xyz/config"
	"main-xyz/dto/in"
	"main-xyz/utils/text"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type MultipartReader interface {
	ParseMultipartForm(
		r *http.Request,
	) error

	ReadMultipartFile(
		request *http.Request,
		key string,
	) (
		file []in.MultipartFile,
		errs error,
	)

	SetFileNameAliasingFunction(
		f filenameAliasing,
	) *multipartReader
}

type filenameAliasing func(int, string) string

func SplitFileNameAndExtension(filename string) (string, string) {

	split := strings.Split(filename, ".")
	if len(split) == 1 {
		return filename, ""
	}

	return strings.Join(split[0:len(split)-2], ""), split[len(split)-1]
}

func DefaultAliasingFunction(idx int, fileName string) string {
	var addition string
	if idx > 0 {
		addition = fmt.Sprintf("_%d", idx)
	}

	withoutExt, ext := SplitFileNameAndExtension(fileName)
	return fmt.Sprintf("%s%s.%s", withoutExt, addition, ext)
}

func UUIDAliasingFunction(idx int, fileName string) string {
	var addition string
	if idx > 0 {
		addition = fmt.Sprintf("_%d", idx)
	}

	_, ext := SplitFileNameAndExtension(fileName)

	return fmt.Sprintf("%s%s.%s", text.GetUUID(), addition, ext)
}

func NewMultipartReader(
	maxMemory int64,
	tempFile string,
	logger *config.LoggerCustom,
) MultipartReader {

	_ = os.MkdirAll(tempFile, 0770)

	return &multipartReader{
		maxMemory:        maxMemory,
		tempFile:         tempFile,
		filenameAliasing: DefaultAliasingFunction,
		logger:           logger,
	}
}

type multipartReader struct {
	maxMemory        int64
	tempFile         string
	filenameAliasing filenameAliasing
	state            chan (struct{})
	logger           *config.LoggerCustom
}

func (m multipartReader) ReadMultipartFile(
	request *http.Request,
	key string,
) (
	result []in.MultipartFile,
	err error,
) {

	files, ok := request.MultipartForm.File[key]
	if !ok {
		return
	}

	defer func() {
		if result != nil {
			for i := 0; i < len(result); i++ {
				errs := result[i].File.Close()
				if errs != nil {
					err = errs
					m.logger.Logger.Error("Error Found When Close multipart file", zap.Error(err))
				}
			}
		}
	}()

	for i := 0; i < len(files); i++ {
		file := files[i]
		if file.Header != nil {
			tempResult := in.MultipartFile{
				Header: file,
			}

			tempResult.Alias = tempResult.Header.Header.Get("alias_name")
			tempResult.FullPath = filepath.Join(m.tempFile, tempResult.Alias)
			tempResult.File, err = os.OpenFile(tempResult.FullPath, os.O_RDONLY, os.ModeAppend)
			if err != nil {
				return
			}

			result = append(result, tempResult)
		}
	}

	return result, nil
}

func (m *multipartReader) SetFileNameAliasingFunction(f filenameAliasing) *multipartReader {
	m.filenameAliasing = f
	return m
}

var multipartByReader = &multipart.Form{
	Value: make(map[string][]string),
	File:  make(map[string][]*multipart.FileHeader),
}

func (m multipartReader) ParseMultipartForm(
	r *http.Request,
) error {

	if r.MultipartForm == multipartByReader {
		return errors.New("http: multipart handled by MultipartReader")
	}

	if r.Form == nil {
		err := r.ParseForm()
		if err != nil {
			return err
		}
	}

	if r.MultipartForm != nil {
		return nil
	}

	mr, err := m.readMultipart(r, false)
	if err != nil {
		return err
	}

	f, err := m.readForm(mr)
	if err != nil {
		return err
	}

	if r.PostForm == nil {
		r.PostForm = make(url.Values)
	}

	for k, v := range f.Value {
		r.Form[k] = append(r.Form[k], v...)
		// r.PostForm should also be populated. See Issue 9305.
		r.PostForm[k] = append(r.PostForm[k], v...)
	}

	r.MultipartForm = f

	return nil
}

func (m multipartReader) readMultipart(
	r *http.Request,
	allowMixed bool,
) (
	*multipart.Reader,
	error,
) {

	v := r.Header.Get("Content-Type")
	if v == "" {
		return nil, http.ErrNotMultipart
	}

	d, params, err := mime.ParseMediaType(v)
	if err != nil || !(d == "multipart/form-data" || allowMixed && d == "multipart/mixed") {
		return nil, http.ErrNotMultipart
	}

	boundary, ok := params["boundary"]
	if !ok {
		return nil, http.ErrMissingBoundary
	}

	return multipart.NewReader(r.Body, boundary), nil
}

func (m multipartReader) readForm(
	r *multipart.Reader,
) (
	form *multipart.Form,
	err error,
) {
	maxMemory := m.maxMemory
	tempFile := m.tempFile

	form = &multipart.Form{
		Value: make(map[string][]string),
		File:  make(map[string][]*multipart.FileHeader),
	}

	defer func() {
		if err != nil {
			if form != nil {
				_ = form.RemoveAll()
			}
		}
	}()

	maxValueBytes := maxMemory + int64(10<<20)

	var first = true

	for {
		p, err := r.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		name := p.FormName()
		if name == "" {
			continue
		}

		filename := p.FileName()

		var b bytes.Buffer

		if filename == "" {
			n, err := io.CopyN(&b, p, maxValueBytes+1)
			if err != nil && err != io.EOF {
				return nil, err
			}

			maxValueBytes -= n
			if maxValueBytes < 0 {
				return nil, multipart.ErrMessageTooLarge
			}

			form.Value[name] = append(form.Value[name], b.String())
			continue
		}

		var idx = 0

		//check form.File[name] is nil (prevent panic)
		if form.File[name] != nil {
			idx = len(form.File[name])
		}

		fileAliasing := m.filenameAliasing(idx, p.FileName())
		p.Header.Add("alias_name", fileAliasing)

		fh := &multipart.FileHeader{
			Filename: filename,
			Header:   p.Header,
		}

		err = func() error {
			var fileTemp *os.File

			defer func() {
				if fileTemp != nil {
					_ = fileTemp.Close()
				}
			}()

			fullPath := fmt.Sprintf("%s/%s", tempFile, fileAliasing)

			if first {
				_ = os.Remove(fullPath)
				first = false
			}

			fileTemp, err = os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
			if err != nil {
				fileTemp, err = os.CreateTemp("", "multipart-")
				if err != nil {
					return err
				}
				return err
			}

			size, err := io.Copy(fileTemp, io.MultiReader(&b, p))
			if cerr := fileTemp.Close(); err == nil {
				err = cerr
			}

			if err != nil {
				_ = os.Remove(fileTemp.Name())
				return err
			}

			fh.Size = size

			return nil
		}()

		if err != nil {
			return nil, err
		}

		form.File[name] = append(form.File[name], fh)
	}

	return form, nil
}
