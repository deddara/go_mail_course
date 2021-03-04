package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson3486653aDecodeCourseraWeek3Hw3BenchZhopa(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "name":
			out.Name = string(in.String())
		case "email":
			out.Email = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3486653aEncodeCourseraWeek3Hw3BenchZhopa(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"browsers\":"
		out.RawString(prefix[1:])
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3486653aEncodeCourseraWeek3Hw3BenchZhopa(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3486653aEncodeCourseraWeek3Hw3BenchZhopa(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3486653aDecodeCourseraWeek3Hw3BenchZhopa(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3486653aDecodeCourseraWeek3Hw3BenchZhopa(l, v)
}

type User struct {
	Browsers []string `json:"browsers"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
}

var dataPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 64))
	},
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	seenBrowsers := make(map[string]bool, 1000)
	scanner := bufio.NewScanner(file)
	io.WriteString(out, "found users:\n")
	for i := 0; scanner.Scan(); i++ {
		user := User{}
		r := jlexer.Lexer{Data: scanner.Bytes()}
		user.UnmarshalEasyJSON(&r)
		isAndroid := false
		isMSIE := false
		browsers := user.Browsers

		for _, browser := range browsers {

			andr := strings.Contains(browser, "Android")
			msie := strings.Contains(browser, "MSIE")
			if andr || msie {
				if andr {
					isAndroid = true
				}
				if msie {
					isMSIE = true
				}
				_, seen := seenBrowsers[browser]
				if !seen {
					seenBrowsers[browser] = true
				}
			}

		}

		if !(isAndroid && isMSIE) {
			continue
		}

		email := dataPool.Get().(*bytes.Buffer)
		email.WriteString(strings.Replace(user.Email, "@", " [at] ", 1))
		io.WriteString(out, "["+strconv.Itoa(i)+"] "+user.Name+" <"+email.String()+">\n")
		email.Reset()
		dataPool.Put(email)
	}
	io.WriteString(out, "\nTotal unique browsers "+strconv.Itoa(len(seenBrowsers))+"\n")
}
