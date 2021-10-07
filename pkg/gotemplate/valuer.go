package gotemplate

var (
	_ Valuer       = &Value{}
	_ Valuer       = DynamicValue(nil)
	_ BoolValuer   = BoolValue(false)
	_ BoolValuer   = DynamicBoolValue(nil)
	_ StringValuer = StringValue("")
	_ StringValuer = DynamicStringValue(nil)
)

type Valuer interface {
	Value(vals *OptionValues) interface{}
}
type Value struct {
	v interface{}
}

func StaticValue(v interface{}) *Value {
	return &Value{v: v}
}

func (v *Value) Value(_ *OptionValues) interface{} {
	return v.v
}

// DefferedValue is a func that calculates the Value based on earlier inputs.
type DynamicValue func(vals *OptionValues) interface{}

func (f DynamicValue) Value(vals *OptionValues) interface{} {
	return f(vals)
}

type BoolValuer interface {
	Value(vals *OptionValues) bool
}
type BoolValue bool

func (v BoolValue) Value(_ *OptionValues) bool {
	return bool(v)
}

// DefferedValue is a func that calculates the Value based on earlier inputs.
type DynamicBoolValue func(vals *OptionValues) bool

func (f DynamicBoolValue) Value(vals *OptionValues) bool {
	return f(vals)
}

type StringValuer interface {
	Value(vals *OptionValues) string
}
type StringValue string

func (v StringValue) Value(_ *OptionValues) string {
	return string(v)
}

// DefferedValue is a func that calculates the Value based on earlier inputs.
type DynamicStringValue func(vals *OptionValues) string

func (f DynamicStringValue) Value(vals *OptionValues) string {
	return f(vals)
}
