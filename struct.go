package main

import (
	"fmt"
	"go/ast"
	"strings"
)

const (
	read = 1 << iota
	write
	embedded
	jsonOmitEmpty
	hasPrivateFields
)

type StructDefaults struct {
	BaseRepr
}

type StructRepr struct {
	StructDefaults
	Fields []*StructFieldRepr
}

type StructFieldRepr struct {
	BaseFieldRepr
	Read        string // Method name for reads
	Write       string // Method name for writes
	DefaultExpr string // Default expression
	JsonName    string // Name used for json [un]marshaling
	JsonNameCI  string // Case insensitive version of JsonName
	astField    *ast.Field
}

var structDefaults StructDefaults

func (self *StructDefaults) gatherFlags(tagText string) error {
	return self.genericGatherFlags(tagText, func(flag Flag) error {
		switch flag.Name {
		case "drop_json": // Do not generate JSON marshaling methods
			return self.doBooleanFlag(flag, dropJson)

		default:
			return UnknownFlag
		}
		return nil
	})
}

func (self *StructRepr) GetPrivateTypeName() string {
	return "private_" + self.getUniqueId()
}
func (self *StructRepr) GetJSONTypeName() string {
	return "json_" + self.getUniqueId()
}
func (self *StructRepr) DoJson() bool { return self.flags&dropJson == 0 }
func (self *StructRepr) HasPrivateFields() bool {
	return self.flags&hasPrivateFields == hasPrivateFields
}

func (self *StructFieldRepr) HasJSONOmitEmpty() bool {
	return self.flags&jsonOmitEmpty == jsonOmitEmpty
}
func (self *StructFieldRepr) DoRead() bool  { return len(self.Read) > 0 }
func (self *StructFieldRepr) DoWrite() bool { return len(self.Write) > 0 }
func (self *StructFieldRepr) DoDefaultExpr() bool {
	return len(self.DefaultExpr) > 0
}
func (self *StructFieldRepr) IsEmbedded() bool {
	return self.flags&embedded == embedded
}
func (self *StructFieldRepr) IsPrivate() bool {
	return !self.IsPublic() && !self.IsEmbedded()
}
func (self *StructFieldRepr) IsPublic() bool {
	return self.flags&(read|write) == (read | write)
}
func (self *StructFieldRepr) GetSpaceAndTag() string {
	if len(self.Tag) > 0 {
		return fmt.Sprintf(" `%s`", self.Tag)
	}
	return ""
}

func (self *StructFieldRepr) CouldBeJSON() bool {
	return isExportedIdent(self.Name) && !self.IsEmbedded()
}

func (self *StructFieldRepr) getJSONFieldTagName() string {
	if idx := strings.Index(self.Tag, `json:"`); idx != -1 {
		t := self.Tag[idx+6:] // Found a valid start to the JSON field tag

		if idx = strings.IndexByte(t, '"'); idx != -1 {
			t = t[0:idx] // Found the closing quote, so it's valid

			if idx = strings.IndexByte(t, ','); idx != -1 {
				return t[0:idx] // Found a comma, so the name comes before it
			}
			return t // No comma found, so all we have is a name
		}
	}
	return ""
}

// Gets the Name, which may be the Type for embedded fields. If so, it strips
// away any leading `*`
func (self *StructFieldRepr) GetNameMaybeType() string {
	if self.IsEmbedded() {
		return strings.TrimLeft(self.Type, "*")
	}
	return self.Name
}

func (self *FileData) doStructDefaults(tagText string) error {
	return structDefaults.gatherFlags(tagText)
}

func (self *FileData) newStruct(tagText string, docs []*ast.Comment,
	spec *ast.TypeSpec, strct *ast.StructType) error {

	var err error

	strct_repr := StructRepr{
		StructDefaults: structDefaults, // copy of current defaults
	}

	if err = strct_repr.setDocsAndName(docs, spec); err != nil {
		return err
	}

	if err = strct_repr.gatherFlags(tagText); err != nil {
		return err
	}

	if err = strct_repr.doFields(strct.Fields); err != nil {
		return err
	}

	self.Structs = append(self.Structs, &strct_repr)

	return nil
}

func (self *StructRepr) doFields(fields *ast.FieldList) (err error) {
	const name_conflit = "%q method name conflicts with property name %q"

	if len(fields.List) == 0 {
		return fmt.Errorf("@structs must have at least one field defined")
	}

	for _, field := range fields.List {
		var f = StructFieldRepr{astField: field}

		if err := f.gatherCodeCommentsAndName(field, true); err != nil {
			return err
		}

		f.JsonName = f.Name

		if f.flags&embedded == 0 {
			if err = f.gatherFlags(getFlags(field.Tag)); err != nil {
				return err
			}

			if f.flags&(read|write) == (read | write) { // if `read` AND `write`
				if f.Name == f.Read {
					return fmt.Errorf(name_conflit, "read", f.Name)
				}
				if f.Name == f.Write {
					return fmt.Errorf(name_conflit, "write", f.Name)
				}

				// if `read` OR `write` are set (but not both), set default name if needed
			} else if f.flags&read == read || f.flags&write == write {
				self.flags |= hasPrivateFields

				if f.flags&read == read && len(f.Read) == 0 {
					f.Read = f.Name
				}

				if f.flags&write == write && len(f.Write) == 0 {
					f.Write = "Set" + strings.Title(f.Name)
				}
			}
		}

		f.JsonNameCI = strings.ToLower(f.JsonName)

		self.Fields = append(self.Fields, &f)
	}

	return nil
}

func (self *StructFieldRepr) gatherFlags(tagText string) error {
	return self.genericGatherFlags(tagText, func(flag Flag) error {
		switch flag.Name {

		case "gRead": // Set read access
			self.flags |= read
			if flag.FoundColon {
				self.Read = flag.Value
			}

		case "gWrite": // Set write access
			self.flags |= write
			if flag.FoundColon {
				self.Write = flag.Value
			}

		case "json": // Just to find out if it has `omitempty`
			if len(flag.Value) > 0 {
				if idx := strings.IndexByte(flag.Value, ','); idx == -1 {
					self.JsonName = flag.Value

				} else {
					jsonName := strings.TrimSpace(flag.Value[0:idx])
					if len(jsonName) > 0 {
						self.JsonName = jsonName
					}

					if strings.Contains(flag.Value[idx:], "omitempty") {
						self.flags |= jsonOmitEmpty
					}
				}
			}

			return UnknownFlag

		default:
			return UnknownFlag
		}
		return nil
	})
}

func (self *FileData) GatherStructImports() {
	if len(self.Structs) == 0 {
		return
	}
	self.Imports["encoding/json"] = true
	self.Imports["Golific/gJson"] = true
}

func (self *StructFieldRepr) GetJSONOmitCondition() string {
	if self.HasJSONOmitEmpty() {
		switch n := self.astField.Type.(type) {
		case *ast.ArrayType, *ast.MapType:
			return "len(self." + self.GetNameMaybeType() + ") != 0"

		case *ast.Ident:
			switch n.Name {
			case "bool":
				return "!self." + self.GetNameMaybeType()

			case "string":
				return "len(self." + self.GetNameMaybeType() + ") != 0"

			case "int", "int64", "int32", "int16", "int8", "uint", "uint64", "uint32",
				"uint16", "uint8", "float64", "float32":
				return "self." + self.GetNameMaybeType() + " != 0"
			}
		default:
			return "z, ok := self." + self.GetNameMaybeType() + ".(gJson.Zeroable); !ok || !z.IsZero()"
		}
	}

	return "true"
}

var struct_tmpl = `
{{- define "JSONEncodeField" -}}
	{{- $f := . -}}
	{{- $cond := $f.GetJSONOmitCondition -}}

	{{if ne $cond "true" -}}
	if {{$cond}} {
	{{end -}}

		{{if $f.IsEmbedded -}}
		{
			// Embedded field
			var wasFirst = first
		{{end -}}

		if first {
			first = false
		} else {
			encoder.WriteRawByte(',')
		}

		{{- if $f.IsEmbedded}}
			var didWrite bool

			if je, ok := interface{}(self.{{$f.GetNameMaybeType}}).(gJson.JSONEncodable); ok {
				didWrite = encoder.EmbedEncodedStruct(je)
			} else {
				didWrite = encoder.EmbedMarshaledStruct(self.{{$f.GetNameMaybeType}})
			}

			if !didWrite {
				if wasFirst { // Still haven't had first write
					first = true
				} else { // Drop the comma that was added
					encoder.Truncate(encoder.Len()-1)
				}
			}
		}

		{{- else}}
		encoder.WriteRawString("\"{{$f.JsonName}}\":")
		encoder.Encode(self.{{$f.Name}})
		{{- end}}

	{{- if ne $cond "true" -}}
	}
	{{- end -}}
{{end}}



{{- define "generate_struct"}}
{{- range $struct := .}}
{{- $privateType := $struct.GetPrivateTypeName}}
{{- $jsonType := $struct.GetJSONTypeName}}

/*****************************

{{$struct.Name}} struct

******************************/

{{$struct.DoDocs -}}
type {{$struct.Name}} struct {
  private {{$privateType}}
  {{- range $f := $struct.Fields}}
	{{- if $f.IsEmbedded}}
	{{printf "%s%s%s" $f.DoDocs $f.Type $f.GetSpaceAndTag}}
  {{- else if $f.IsPublic}}
  {{printf "%s%s %s%s" $f.DoDocs $f.Name $f.Type $f.GetSpaceAndTag}}
  {{- end -}}
  {{end -}}
}


type {{$privateType}} struct {
  {{- range $f := $struct.Fields}}
  {{- if $f.IsPrivate}}
  {{printf "%s%s %s%s" $f.DoDocs $f.Name $f.Type $f.GetSpaceAndTag}}
  {{- end -}}
  {{end -}}
}


// JSONEncode implements part of Golific's JSONEncodable interface.
func (self *{{$privateType}}) JSONEncode(encoder *gJson.Encoder) {
	var first = true

	{{- range $f := $struct.Fields -}}
	{{- if $f.IsPrivate}}
	{{template "JSONEncodeField" $f}}
	{{end -}}
	{{end -}}
}


type {{$jsonType}} struct {
  *{{- $privateType}}
  {{- range $f := $struct.Fields}}
	{{- if $f.IsEmbedded}}
	{{printf "%s%s" $f.Type $f.GetSpaceAndTag}}
  {{- else if $f.IsPublic}}
  {{printf "%s %s%s" $f.Name $f.Type $f.GetSpaceAndTag}}
  {{- end -}}
  {{end -}}
}


{{- range $f := $struct.Fields}}
{{- if $f.DoRead}}
func (self *{{$struct.Name}}) {{$f.Read}} () {{$f.Type}} {
  {{- if $f.IsPrivate}}
  return self.private.{{$f.Name}}
  {{- else -}}
  return self.{{$f.Name}}
  {{end -}}
}
{{end -}}
{{if $f.DoWrite}}
func (self *{{$struct.Name}}) {{$f.Write}} ( v {{$f.Type}} ) {
  {{- if $f.IsPrivate}}
  self.private.{{$f.Name}} = v
  {{- else -}}
  self.{{$f.Name}} = v
  {{end -}}
}
{{end -}}
{{end -}}


// JSONEncode implements part of Golific's JSONEncodable interface.
func (self *{{$struct.Name}}) JSONEncode(encoder *gJson.Encoder) {
	encoder.WriteRawByte('{')

	{{if $struct.HasPrivateFields}}
	var startPos = encoder.Len()

	// Encodes only the fields of the struct, without curly braces
	self.private.JSONEncode(encoder)

	var first = encoder.Len() == startPos

	{{else}}
	var first = true
	{{end}}

	{{- range $f := $struct.Fields -}}
	{{- if or $f.IsEmbedded $f.IsPublic}}
	{{template "JSONEncodeField" $f}}
	{{end -}}
	{{end -}}

	encoder.WriteRawByte('}')
}


{{- if $struct.DoJson}}
func (self *{{$struct.Name}}) MarshalJSON() ([]byte, error) {
  return json.Marshal({{$jsonType}} {
    &self.private,
    {{range $f := $struct.Fields -}}
    {{if or $f.IsEmbedded $f.IsPublic -}}
    self.{{$f.GetNameMaybeType}},
    {{end -}}
    {{end -}}
  })
}

func (self *{{$struct.Name}}) UnmarshalJSON(j []byte) error {
	if len(j) == 4 && string(j) == "null" {
		return nil
	}

	m := make(map[string]json.RawMessage)

	err := json.Unmarshal(j, &m)
	if err != nil {
		return err
	}

	// For every property found, perform a separate UnmarshalJSON operation. This
	// prevents overwrite of values in 'self' where properties are absent.
	for key, rawMsg := range m {
		// The anon structs in each case are needed for field tags

		switch strings.ToLower(key) {
		{{- range $f := $struct.Fields -}}
		{{if $f.CouldBeJSON}}
		case {{printf "%q" $f.JsonNameCI}}:

		var x struct {
			F {{$f.Type}}{{$f.GetSpaceAndTag}}
		}

		var msgForStruct = append(append(append(append(
			[]byte("{\""), key...), "\":"...), rawMsg...), '}')

		if err = json.Unmarshal(msgForStruct, &x); err != nil {
			return err
		}

		{{if $f.IsPublic}}
			self.{{$f.Name}} = x.F
		{{else -}}
			self.private.{{$f.Name}} = x.F
		{{- end -}}

		{{end -}}
		{{end}}
		default:
			// Ignoring unknown property
		}
	}
  return nil
}
{{end -}}

{{end -}}
{{end -}}
`
