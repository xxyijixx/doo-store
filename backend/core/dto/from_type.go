package dto

import (
	"fmt"
	"regexp"
	"strings"
)

type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeSelect   FieldType = "select"
	FieldTypeNumber   FieldType = "number"
	FieldTypePassword FieldType = "password"
	FieldTypeRadio    FieldType = "radio"
	FieldTypeCheckbox FieldType = "checkbox"
)

// Option 定义选项结构
type Option struct {
	Label     string      `json:"label"`
	Value     string      `json:"value"`
	SubFields []FormField `json:"sub_fields,omitempty"` // 该选项特有的子字段
}

// Dependency 定义字段间的依赖关系
type Dependency struct {
	Field    string      `json:"field"`    // 依赖的字段
	Value    interface{} `json:"value"`    // 依赖字段的值
	Operator string      `json:"operator"` // 比较操作符：eq, neq, in, etc.
}

// Validation 定义字段验证规则
type Validation struct {
	Required bool   `json:"required"`
	Pattern  string `json:"pattern,omitempty"`
	MinLen   *int   `json:"min_len,omitempty"`
	MaxLen   *int   `json:"max_len,omitempty"`
}

// FormField 定义通用的表单字段结构
type FormField struct {
	Label       string      `json:"label"`
	EnvKey      string      `json:"env_key"`
	Type        FieldType   `json:"type"`
	Default     interface{} `json:"default,omitempty"`
	Options     []Option    `json:"options,omitempty"`
	Validation  *Validation `json:"validation,omitempty"`
	Dependency  *Dependency `json:"dependency,omitempty"`
	Placeholder string      `json:"placeholder,omitempty"`
	Order       int         `json:"order"`              // 显示顺序
	Hidden      bool        `json:"hidden,omitempty"`   // 是否隐藏
	ReadOnly    bool        `json:"readonly,omitempty"` // 是否只读
}

// FormConfig 定义表单配置结构
type FormConfig struct {
	Fields []FormField `json:"fields"`
}

// ValidationError 定义验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type formValidator struct {
	fields      []*FormField
	params      map[string]interface{}
	fieldValues map[string]string
	errors      []error
}

func FillAndValidateForm(config []*FormField, params map[string]interface{}) []*FormField {
	validator := newFormValidator(config, params)
	validator.preprocessParams()

	// 填充配置值
	filledConfig := validator.fillConfigValues()

	// 执行验证
	validator.validateFields()

	return filledConfig
}

func ValidateFormData(config []*FormField, params map[string]interface{}) []error {
	return newFormValidator(config, params).validate()
}

// newFormValidator 创建新的验证器实例
func newFormValidator(fields []*FormField, params map[string]interface{}) *formValidator {
	return &formValidator{
		fields:      fields,
		params:      params,
		fieldValues: make(map[string]string),
		errors:      make([]error, 0),
	}
}

// fillConfigValues 填充配置中字段的值
func (v *formValidator) fillConfigValues() []*FormField {

	for i, field := range v.fields {
		// 获取参数值
		if value, exists := v.params[field.EnvKey]; exists {
			v.fields[i].Default = value
		}

		// 处理子字段的值（如果存在）
		if len(field.Options) > 0 {
			for j, option := range field.Options {
				if len(option.SubFields) > 0 {
					for k, subField := range option.SubFields {
						if value, exists := v.params[subField.EnvKey]; exists {
							v.fields[i].Options[j].SubFields[k].Default = value
						}
					}
				}
			}
		}
	}

	return v.fields
}

// validate 执行验证流程
func (v *formValidator) validate() []error {
	v.preprocessParams()
	v.validateFields()
	return v.errors
}

// preprocessParams 预处理参数
func (v *formValidator) preprocessParams() {
	for key, value := range v.params {
		if value == nil {
			v.fieldValues[key] = ""
			continue
		}

		switch val := value.(type) {
		case string:
			v.fieldValues[key] = val
		case float64:
			v.fieldValues[key] = fmt.Sprintf("%v", val)
		case int:
			v.fieldValues[key] = fmt.Sprintf("%d", val)
		case bool:
			v.fieldValues[key] = fmt.Sprintf("%v", val)
		case []interface{}:
			values := make([]string, len(val))
			for i, item := range val {
				values[i] = fmt.Sprintf("%v", item)
			}
			v.fieldValues[key] = strings.Join(values, ",")
		default:
			v.fieldValues[key] = fmt.Sprintf("%v", val)
		}
	}
}

// validateFields 验证所有字段
func (v *formValidator) validateFields() {
	for _, field := range v.fields {
		if field.Hidden {
			continue
		}

		value, exists := v.fieldValues[field.EnvKey]

		if field.Dependency != nil && !v.checkDependency(field.Dependency) {
			continue
		}

		if err := v.validateField(*field, value, exists); err != nil {
			v.errors = append(v.errors, err)
		}
	}
}

// validateField 验证单个字段
func (v *formValidator) validateField(field FormField, value string, exists bool) error {
	// 验证必填
	if field.Validation != nil && field.Validation.Required {
		if !exists || strings.TrimSpace(value) == "" {
			return &ValidationError{
				Field:   field.Label,
				Message: "此字段为必填项",
			}
		}
	}

	// 如果字段没有值且非必填，跳过后续验证
	if !exists || value == "" {
		return nil
	}

	// 验证字段类型
	if err := v.validateFieldType(field, value); err != nil {
		return err
	}

	// 验证字段格式
	return v.validateFieldFormat(field, value)
}

// validateFieldType 验证字段类型
func (v *formValidator) validateFieldType(field FormField, value string) error {
	switch field.Type {
	case FieldTypeNumber:
		if matched, _ := regexp.MatchString(`^-?\d+(\.\d+)?$`, value); !matched {
			return &ValidationError{
				Field:   field.Label,
				Message: "请输入有效的数字",
			}
		}
	case FieldTypeSelect, FieldTypeRadio:
		if len(field.Options) > 0 {
			if !v.isValidOption(field.Options, value) {
				return &ValidationError{
					Field:   field.Label,
					Message: "请选择有效的选项",
				}
			}
		}
	case FieldTypeCheckbox:
		values := strings.Split(value, ",")
		for _, val := range values {
			if !v.isValidOption(field.Options, val) {
				return &ValidationError{
					Field:   field.Label,
					Message: "存在无效的选项值",
				}
			}
		}
	}
	return nil
}

// validateFieldFormat 验证字段格式
func (v *formValidator) validateFieldFormat(field FormField, value string) error {
	if field.Validation == nil {
		return nil
	}

	if field.Validation.Pattern != "" {
		matched, err := regexp.MatchString(field.Validation.Pattern, value)
		if err != nil {
			return &ValidationError{
				Field:   field.Label,
				Message: "正则表达式验证失败",
			}
		}
		if !matched {
			return &ValidationError{
				Field:   field.Label,
				Message: "输入格式不正确",
			}
		}
	}

	if field.Validation.MinLen != nil && len(value) < *field.Validation.MinLen {
		return &ValidationError{
			Field:   field.Label,
			Message: fmt.Sprintf("长度不能小于 %d", *field.Validation.MinLen),
		}
	}

	if field.Validation.MaxLen != nil && len(value) > *field.Validation.MaxLen {
		return &ValidationError{
			Field:   field.Label,
			Message: fmt.Sprintf("长度不能大于 %d", *field.Validation.MaxLen),
		}
	}

	return nil
}

// checkDependency 检查字段依赖关系
func (v *formValidator) checkDependency(dep *Dependency) bool {
	dependentValue, exists := v.fieldValues[dep.Field]
	if !exists {
		return false
	}

	switch dep.Operator {
	case "eq":
		return fmt.Sprintf("%v", dep.Value) == dependentValue
	case "neq":
		return fmt.Sprintf("%v", dep.Value) != dependentValue
	case "in":
		if values, ok := dep.Value.([]interface{}); ok {
			for _, v := range values {
				if fmt.Sprintf("%v", v) == dependentValue {
					return true
				}
			}
		}
	}
	return false
}

// isValidOption 检查选项值是否有效
func (v *formValidator) isValidOption(options []Option, value string) bool {
	for _, opt := range options {
		if opt.Value == value {
			return true
		}
	}
	return false
}
