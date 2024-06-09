// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package storage

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonF642ad3eDecodeGithubComTrb1makerOtusGolangHomeWorkHw12131415CalendarInternalStorage(in *jlexer.Lexer, out *Notify) {
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
		case "ownerId":
			out.OwnerID = string(in.String())
		case "id":
			out.ID = string(in.String())
		case "title":
			out.Title = string(in.String())
		case "startTime":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.StartTime).UnmarshalJSON(data))
			}
		case "notifyTime":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.NotifyTime).UnmarshalJSON(data))
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
func easyjsonF642ad3eEncodeGithubComTrb1makerOtusGolangHomeWorkHw12131415CalendarInternalStorage(out *jwriter.Writer, in Notify) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"ownerId\":"
		out.RawString(prefix[1:])
		out.String(string(in.OwnerID))
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"startTime\":"
		out.RawString(prefix)
		out.Raw((in.StartTime).MarshalJSON())
	}
	{
		const prefix string = ",\"notifyTime\":"
		out.RawString(prefix)
		out.Raw((in.NotifyTime).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Notify) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF642ad3eEncodeGithubComTrb1makerOtusGolangHomeWorkHw12131415CalendarInternalStorage(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Notify) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF642ad3eEncodeGithubComTrb1makerOtusGolangHomeWorkHw12131415CalendarInternalStorage(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Notify) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF642ad3eDecodeGithubComTrb1makerOtusGolangHomeWorkHw12131415CalendarInternalStorage(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Notify) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF642ad3eDecodeGithubComTrb1makerOtusGolangHomeWorkHw12131415CalendarInternalStorage(l, v)
}
func easyjsonF642ad3eDecodeGithubComTrb1makerOtusGolangHomeWorkHw12131415CalendarInternalStorage1(in *jlexer.Lexer, out *Event) {
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
		case "startTime":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.StartTime).UnmarshalJSON(data))
			}
		case "endTime":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.EndTime).UnmarshalJSON(data))
			}
		case "description":
			out.Description = string(in.String())
		case "ownerId":
			out.OwnerID = string(in.String())
		case "notify":
			out.Notify = time.Duration(in.Int64())
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
func easyjsonF642ad3eEncodeGithubComTrb1makerOtusGolangHomeWorkHw12131415CalendarInternalStorage1(out *jwriter.Writer, in Event) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"startTime\":"
		out.RawString(prefix)
		out.Raw((in.StartTime).MarshalJSON())
	}
	{
		const prefix string = ",\"endTime\":"
		out.RawString(prefix)
		out.Raw((in.EndTime).MarshalJSON())
	}
	if in.Description != "" {
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"ownerId\":"
		out.RawString(prefix)
		out.String(string(in.OwnerID))
	}
	if in.Notify != 0 {
		const prefix string = ",\"notify\":"
		out.RawString(prefix)
		out.Int64(int64(in.Notify))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Event) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF642ad3eEncodeGithubComTrb1makerOtusGolangHomeWorkHw12131415CalendarInternalStorage1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Event) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF642ad3eEncodeGithubComTrb1makerOtusGolangHomeWorkHw12131415CalendarInternalStorage1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Event) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF642ad3eDecodeGithubComTrb1makerOtusGolangHomeWorkHw12131415CalendarInternalStorage1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Event) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF642ad3eDecodeGithubComTrb1makerOtusGolangHomeWorkHw12131415CalendarInternalStorage1(l, v)
}
