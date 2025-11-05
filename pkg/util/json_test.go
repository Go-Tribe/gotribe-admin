package util

import (
	"testing"
)

type TestStruct struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Email   string `json:"email"`
	Active  bool   `json:"active"`
	Details struct {
		City    string `json:"city"`
		Country string `json:"country"`
	} `json:"details"`
}

func TestJSON_Struct2Json(t *testing.T) {
	tests := []struct {
		name      string
		obj       interface{}
		expectErr bool
	}{
		{
			"正常结构体",
			TestStruct{
				Name:   "John Doe",
				Age:    30,
				Email:  "john@example.com",
				Active: true,
				Details: struct {
					City    string `json:"city"`
					Country string `json:"country"`
				}{
					City:    "New York",
					Country: "USA",
				},
			},
			false,
		},
		{
			"空结构体",
			TestStruct{},
			false,
		},
		{
			"简单类型",
			"hello world",
			false,
		},
		{
			"数字",
			123,
			false,
		},
		{
			"布尔值",
			true,
			false,
		},
		{
			"切片",
			[]string{"a", "b", "c"},
			false,
		},
		{
			"映射",
			map[string]interface{}{
				"key1": "value1",
				"key2": 123,
				"key3": true,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := JSONUtil.Struct2Json(tt.obj)
			if (err != nil) != tt.expectErr {
				t.Errorf("Struct2Json() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && result == "" {
				t.Error("Struct2Json() returned empty string")
			}
		})
	}
}

func TestJSON_Json2Struct(t *testing.T) {
	tests := []struct {
		name      string
		jsonStr   string
		obj       interface{}
		expectErr bool
	}{
		{
			"正常JSON",
			`{"name":"John Doe","age":30,"email":"john@example.com","active":true,"details":{"city":"New York","country":"USA"}}`,
			&TestStruct{},
			false,
		},
		{
			"空JSON对象",
			`{}`,
			&TestStruct{},
			false,
		},
		{
			"简单类型",
			`"hello world"`,
			new(string),
			false,
		},
		{
			"数字",
			`123`,
			new(int),
			false,
		},
		{
			"布尔值",
			`true`,
			new(bool),
			false,
		},
		{
			"切片",
			`["a","b","c"]`,
			new([]string),
			false,
		},
		{
			"无效JSON",
			`{invalid json}`,
			&TestStruct{},
			true,
		},
		{
			"空字符串",
			``,
			&TestStruct{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := JSONUtil.Json2Struct(tt.jsonStr, tt.obj)
			if (err != nil) != tt.expectErr {
				t.Errorf("Json2Struct() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestJSON_JsonI2Struct(t *testing.T) {
	tests := []struct {
		name          string
		jsonInterface interface{}
		obj           interface{}
		expectErr     bool
	}{
		{
			"字符串类型",
			`{"name":"John Doe","age":30}`,
			&TestStruct{},
			false,
		},
		{
			"非字符串类型",
			123,
			&TestStruct{},
			true,
		},
		{
			"空接口",
			nil,
			&TestStruct{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := JSONUtil.JsonI2Struct(tt.jsonInterface, tt.obj)
			if (err != nil) != tt.expectErr {
				t.Errorf("JsonI2Struct() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestJSON_RoundTrip(t *testing.T) {
	original := TestStruct{
		Name:   "John Doe",
		Age:    30,
		Email:  "john@example.com",
		Active: true,
		Details: struct {
			City    string `json:"city"`
			Country string `json:"country"`
		}{
			City:    "New York",
			Country: "USA",
		},
	}

	// 结构体转JSON
	jsonStr, err := JSONUtil.Struct2Json(original)
	if err != nil {
		t.Fatalf("Struct2Json() error = %v", err)
	}

	// JSON转结构体
	var result TestStruct
	err = JSONUtil.Json2Struct(jsonStr, &result)
	if err != nil {
		t.Fatalf("Json2Struct() error = %v", err)
	}

	// 比较结果
	if result.Name != original.Name {
		t.Errorf("Name = %v, want %v", result.Name, original.Name)
	}
	if result.Age != original.Age {
		t.Errorf("Age = %v, want %v", result.Age, original.Age)
	}
	if result.Email != original.Email {
		t.Errorf("Email = %v, want %v", result.Email, original.Email)
	}
	if result.Active != original.Active {
		t.Errorf("Active = %v, want %v", result.Active, original.Active)
	}
	if result.Details.City != original.Details.City {
		t.Errorf("Details.City = %v, want %v", result.Details.City, original.Details.City)
	}
	if result.Details.Country != original.Details.Country {
		t.Errorf("Details.Country = %v, want %v", result.Details.Country, original.Details.Country)
	}
}

func BenchmarkJSON_Struct2Json(b *testing.B) {
	obj := TestStruct{
		Name:   "John Doe",
		Age:    30,
		Email:  "john@example.com",
		Active: true,
		Details: struct {
			City    string `json:"city"`
			Country string `json:"country"`
		}{
			City:    "New York",
			Country: "USA",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := JSONUtil.Struct2Json(obj)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSON_Json2Struct(b *testing.B) {
	jsonStr := `{"name":"John Doe","age":30,"email":"john@example.com","active":true,"details":{"city":"New York","country":"USA"}}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var obj TestStruct
		err := JSONUtil.Json2Struct(jsonStr, &obj)
		if err != nil {
			b.Fatal(err)
		}
	}
}
