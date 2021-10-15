package gotemplate

//go:generate go run cmd/dotembed/main.go -target _template -o embed_gen.go -pkg gotemplate -var FS
const Key = "_template"
