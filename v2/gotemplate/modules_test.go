package gotemplate

import (
	"reflect"
	"testing"
)

func TestModuleBase_GetNextQuestion(t *testing.T) {
	type fields struct {
		ModuleData ModuleData
		questions  []TemplateQuestion
	}
	tests := []struct {
		name   string
		fields fields
		want   *TemplateQuestion
	}{
		{
			name: "get first question",
			fields: fields{
				questions: []TemplateQuestion{
					{
						Name: "project-1",
					},
					{
						Name: "project-2",
					},
				},
			},
			want: &TemplateQuestion{
				Name: "project-1",
			},
		},
		{
			name: "no question left",
			fields: fields{
				questions: []TemplateQuestion{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ModuleBase{
				ModuleData: tt.fields.ModuleData,
				questions:  tt.fields.questions,
			}
			if got := m.GetNextQuestion(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModuleBase.GetNextQuestion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleBase_GetModule(t *testing.T) {
	type fields struct {
		ModuleData ModuleData
		questions  []TemplateQuestion
	}
	tests := []struct {
		name   string
		fields fields
		want   *ModuleData
	}{
		{
			name: "no question left",
			fields: fields{
				ModuleData: ModuleData{
					Name: "project-1",
				},
				questions: []TemplateQuestion{},
			},
			want: &ModuleData{
				Name: "project-1",
			},
		},
		{
			name: "still questions left",
			fields: fields{
				ModuleData: ModuleData{
					Name: "project-1",
				},
				questions: []TemplateQuestion{
					{
						Name: "project-1",
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ModuleBase{
				ModuleData: tt.fields.ModuleData,
				questions:  tt.fields.questions,
			}
			if got := m.GetModule(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModuleBase.GetModule() = %v, want %v", got, tt.want)
			}
		})
	}
}
