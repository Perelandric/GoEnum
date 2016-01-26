package main

func (self *FileData) GatherStructImports() {
	if len(self.Structs) == 0 {
		return
	}
	self.Imports["encoding/json"] = true
}

// If any EnumRepr is `bitflag`, `strings` is needed
func (self *FileData) DoStructSummary() bool {
	return false
}

var struct_tmpl = `
{{- define "generate_struct"}}
{{- range $struct := .}}
{{- $privateType := $struct.GetPrivateTypeName}}
{{- $jsonType := $struct.GetJSONTypeName}}

/*****************************

{{$struct.Name}} struct

******************************/

{{- if $struct.DoCtor}}
func {{$struct.GetCtorName}}() *{{$struct.Name}} {
  return &{{$struct.Name}} {
    private: {{$privateType}} {
      {{range $f := $struct.Fields}}
      {{- if and $f.IsPrivate $f.DoDefaultExpr -}}
      {{printf "%s: %s," $f.Name $f.DefaultExpr}}
      {{end -}}
      {{end -}}
    },
    {{range $f := $struct.Fields}}
    {{- if and (or $f.IsEmbedded $f.IsPublic) $f.DoDefaultExpr -}}
    {{printf "%s: %s," $f.GetNameMaybeType $f.DefaultExpr}}
    {{end -}}
    {{end}}
  }
}
{{end}}

type {{$struct.Name}} struct {
  private {{$privateType}}
  {{- range $f := $struct.Fields}}
	{{- if $f.IsEmbedded}}
	{{printf "%s%s" $f.Name $f.GetSpaceAndTag}}
  {{- else if $f.IsPublic}}
  {{printf "%s %s%s" $f.Name $f.Type $f.GetSpaceAndTag}}
  {{- end -}}
  {{end -}}
}

type {{$privateType}} struct {
  {{- range $f := $struct.Fields}}
  {{- if $f.IsPrivate -}}
  {{printf "%s %s%s" $f.Name $f.Type $f.GetSpaceAndTag}}
  {{end -}}
  {{end -}}
}

type {{$jsonType}} struct {
  *{{- $privateType}}
  {{- range $f := $struct.Fields}}
	{{- if $f.IsEmbedded}}
	{{printf "%s%s" $f.Name $f.GetSpaceAndTag}}
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
  var temp {{$jsonType}}
  if err := json.Unmarshal(j, &temp); err != nil {
    return err
  }
  self.private = *temp.{{$privateType}}
  {{range $f := $struct.Fields -}}
  {{if or $f.IsEmbedded $f.IsPublic -}}
  self.{{$f.GetNameMaybeType}} = temp.{{$f.GetNameMaybeType}}
  {{end -}}
  {{end -}}
  return nil
}
{{end -}}

{{end -}}
{{end -}}
`
