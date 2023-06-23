// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package api

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api(in *jlexer.Lexer, out *articlesListReq) {
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
		case "id":
			out.ID = string(in.String())
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
func easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api(out *jwriter.Writer, in articlesListReq) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v articlesListReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v articlesListReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *articlesListReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *articlesListReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api(l, v)
}
func easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api1(in *jlexer.Lexer, out *Void) {
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
func easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api1(out *jwriter.Writer, in Void) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Void) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Void) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Void) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Void) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api1(l, v)
}
func easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api2(in *jlexer.Lexer, out *UsersListReq) {
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
		case "userIds":
			if in.IsNull() {
				in.Skip()
				out.UserIds = nil
			} else {
				in.Delim('[')
				if out.UserIds == nil {
					if !in.IsDelim(']') {
						out.UserIds = make([]string, 0, 4)
					} else {
						out.UserIds = []string{}
					}
				} else {
					out.UserIds = (out.UserIds)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.UserIds = append(out.UserIds, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api2(out *jwriter.Writer, in UsersListReq) {
	out.RawByte('{')
	first := true
	_ = first
	if len(in.UserIds) != 0 {
		const prefix string = ",\"userIds\":"
		first = false
		out.RawString(prefix[1:])
		{
			out.RawByte('[')
			for v2, v3 := range in.UserIds {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UsersListReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UsersListReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UsersListReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UsersListReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api2(l, v)
}
func easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api3(in *jlexer.Lexer, out *ParagraphText) {
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
		case "id":
			out.ID = string(in.String())
		case "text":
			out.Text = string(in.String())
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
func easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api3(out *jwriter.Writer, in ParagraphText) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	if in.Text != "" {
		const prefix string = ",\"text\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Text))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ParagraphText) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ParagraphText) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ParagraphText) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ParagraphText) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api3(l, v)
}
func easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api4(in *jlexer.Lexer, out *ParagraphList) {
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
		case "id":
			out.ID = string(in.String())
		case "items":
			if in.IsNull() {
				in.Skip()
				out.Items = nil
			} else {
				in.Delim('[')
				if out.Items == nil {
					if !in.IsDelim(']') {
						out.Items = make([]string, 0, 4)
					} else {
						out.Items = []string{}
					}
				} else {
					out.Items = (out.Items)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.Items = append(out.Items, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "type":
			out.Type = ListType(in.String())
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
func easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api4(out *jwriter.Writer, in ParagraphList) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	if len(in.Items) != 0 {
		const prefix string = ",\"items\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v5, v6 := range in.Items {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
	}
	if in.Type != "" {
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ParagraphList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ParagraphList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ParagraphList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ParagraphList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api4(l, v)
}
func easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api5(in *jlexer.Lexer, out *ParagraphImage) {
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
		case "id":
			out.ID = string(in.String())
		case "url":
			out.URL = string(in.String())
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
func easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api5(out *jwriter.Writer, in ParagraphImage) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	if in.URL != "" {
		const prefix string = ",\"url\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.URL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ParagraphImage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ParagraphImage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ParagraphImage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ParagraphImage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api5(l, v)
}
func easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api6(in *jlexer.Lexer, out *Paragraph) {
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
		case "text":
			if in.IsNull() {
				in.Skip()
				out.Text = nil
			} else {
				if out.Text == nil {
					out.Text = new(ParagraphText)
				}
				(*out.Text).UnmarshalEasyJSON(in)
			}
		case "image":
			if in.IsNull() {
				in.Skip()
				out.Image = nil
			} else {
				if out.Image == nil {
					out.Image = new(ParagraphImage)
				}
				(*out.Image).UnmarshalEasyJSON(in)
			}
		case "list":
			if in.IsNull() {
				in.Skip()
				out.List = nil
			} else {
				if out.List == nil {
					out.List = new(ParagraphList)
				}
				(*out.List).UnmarshalEasyJSON(in)
			}
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
func easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api6(out *jwriter.Writer, in Paragraph) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Text != nil {
		const prefix string = ",\"text\":"
		first = false
		out.RawString(prefix[1:])
		(*in.Text).MarshalEasyJSON(out)
	}
	if in.Image != nil {
		const prefix string = ",\"image\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Image).MarshalEasyJSON(out)
	}
	if in.List != nil {
		const prefix string = ",\"list\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.List).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Paragraph) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Paragraph) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Paragraph) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Paragraph) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api6(l, v)
}
func easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api7(in *jlexer.Lexer, out *ListPostedByUserReq) {
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
		case "userId":
			out.UserId = string(in.String())
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
func easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api7(out *jwriter.Writer, in ListPostedByUserReq) {
	out.RawByte('{')
	first := true
	_ = first
	if in.UserId != "" {
		const prefix string = ",\"userId\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.UserId))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ListPostedByUserReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ListPostedByUserReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ListPostedByUserReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ListPostedByUserReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api7(l, v)
}
func easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api8(in *jlexer.Lexer, out *InputParagraphText) {
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
		case "text":
			out.Text = string(in.String())
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
func easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api8(out *jwriter.Writer, in InputParagraphText) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Text != "" {
		const prefix string = ",\"text\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Text))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InputParagraphText) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InputParagraphText) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InputParagraphText) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InputParagraphText) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api8(l, v)
}
func easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api9(in *jlexer.Lexer, out *InputParagraphImage) {
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
		case "url":
			out.URL = string(in.String())
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
func easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api9(out *jwriter.Writer, in InputParagraphImage) {
	out.RawByte('{')
	first := true
	_ = first
	if in.URL != "" {
		const prefix string = ",\"url\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.URL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InputParagraphImage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InputParagraphImage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InputParagraphImage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InputParagraphImage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api9(l, v)
}
func easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api10(in *jlexer.Lexer, out *InputParagraph) {
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
		case "inputParagraphText":
			if in.IsNull() {
				in.Skip()
				out.InputParagraphText = nil
			} else {
				if out.InputParagraphText == nil {
					out.InputParagraphText = new(InputParagraphText)
				}
				(*out.InputParagraphText).UnmarshalEasyJSON(in)
			}
		case "inputParagraphImage":
			if in.IsNull() {
				in.Skip()
				out.InputParagraphImage = nil
			} else {
				if out.InputParagraphImage == nil {
					out.InputParagraphImage = new(InputParagraphImage)
				}
				(*out.InputParagraphImage).UnmarshalEasyJSON(in)
			}
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
func easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api10(out *jwriter.Writer, in InputParagraph) {
	out.RawByte('{')
	first := true
	_ = first
	if in.InputParagraphText != nil {
		const prefix string = ",\"inputParagraphText\":"
		first = false
		out.RawString(prefix[1:])
		(*in.InputParagraphText).MarshalEasyJSON(out)
	}
	if in.InputParagraphImage != nil {
		const prefix string = ",\"inputParagraphImage\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.InputParagraphImage).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InputParagraph) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InputParagraph) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InputParagraph) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InputParagraph) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api10(l, v)
}
func easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api11(in *jlexer.Lexer, out *InputArticle) {
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
		case "id":
			out.ID = string(in.String())
		case "title":
			out.Title = string(in.String())
		case "paragraphs":
			if in.IsNull() {
				in.Skip()
				out.Paragraphs = nil
			} else {
				in.Delim('[')
				if out.Paragraphs == nil {
					if !in.IsDelim(']') {
						out.Paragraphs = make([]*InputParagraph, 0, 8)
					} else {
						out.Paragraphs = []*InputParagraph{}
					}
				} else {
					out.Paragraphs = (out.Paragraphs)[:0]
				}
				for !in.IsDelim(']') {
					var v7 *InputParagraph
					if in.IsNull() {
						in.Skip()
						v7 = nil
					} else {
						if v7 == nil {
							v7 = new(InputParagraph)
						}
						(*v7).UnmarshalEasyJSON(in)
					}
					out.Paragraphs = append(out.Paragraphs, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api11(out *jwriter.Writer, in InputArticle) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	if in.Title != "" {
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	if len(in.Paragraphs) != 0 {
		const prefix string = ",\"paragraphs\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v8, v9 := range in.Paragraphs {
				if v8 > 0 {
					out.RawByte(',')
				}
				if v9 == nil {
					out.RawString("null")
				} else {
					(*v9).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InputArticle) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api11(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InputArticle) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api11(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InputArticle) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api11(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InputArticle) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api11(l, v)
}
func easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api12(in *jlexer.Lexer, out *Article) {
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
		case "id":
			out.ID = string(in.String())
		case "title":
			out.Title = string(in.String())
		case "user":
			if in.IsNull() {
				in.Skip()
				out.User = nil
			} else {
				if out.User == nil {
					out.User = new(User)
				}
				(*out.User).UnmarshalEasyJSON(in)
			}
		case "createdAt":
			out.CreatedAt = string(in.String())
		case "paragraphs":
			if in.IsNull() {
				in.Skip()
				out.Paragraphs = nil
			} else {
				in.Delim('[')
				if out.Paragraphs == nil {
					if !in.IsDelim(']') {
						out.Paragraphs = make([]*Paragraph, 0, 8)
					} else {
						out.Paragraphs = []*Paragraph{}
					}
				} else {
					out.Paragraphs = (out.Paragraphs)[:0]
				}
				for !in.IsDelim(']') {
					var v10 *Paragraph
					if in.IsNull() {
						in.Skip()
						v10 = nil
					} else {
						if v10 == nil {
							v10 = new(Paragraph)
						}
						(*v10).UnmarshalEasyJSON(in)
					}
					out.Paragraphs = append(out.Paragraphs, v10)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api12(out *jwriter.Writer, in Article) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	if in.Title != "" {
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	if in.User != nil {
		const prefix string = ",\"user\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.User).MarshalEasyJSON(out)
	}
	if in.CreatedAt != "" {
		const prefix string = ",\"createdAt\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.CreatedAt))
	}
	if len(in.Paragraphs) != 0 {
		const prefix string = ",\"paragraphs\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v11, v12 := range in.Paragraphs {
				if v11 > 0 {
					out.RawByte(',')
				}
				if v12 == nil {
					out.RawString("null")
				} else {
					(*v12).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Article) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api12(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Article) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8336a586EncodeGithubComMaterkovMeme9Web6Api12(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Article) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api12(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Article) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8336a586DecodeGithubComMaterkovMeme9Web6Api12(l, v)
}
