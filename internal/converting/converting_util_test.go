package converting

import (
	"reflect"
	"testing"
)

func TestStrToBool(t *testing.T) {
	type args struct {
		toConvert string
	}
	tests := []struct {
		name    string
		args    args
		wantRes bool
		wantErr bool
	}{
		{name: "true expected 1", args: args{toConvert: " true "}, wantRes: true, wantErr: false},
		{name: "true expected 2", args: args{toConvert: " 1 "}, wantRes: true, wantErr: false},
		{name: "false expected", args: args{toConvert: " F "}, wantRes: false, wantErr: false},
		{name: "false expected 2", args: args{toConvert: " 0 "}, wantRes: false, wantErr: false},
		{name: "empty parse val", args: args{toConvert: " "}, wantRes: false, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := StrToBool(tt.args.toConvert)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("StrToBool() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestStrToUint64(t *testing.T) {
	type args struct {
		toConvert string
	}
	tests := []struct {
		name    string
		args    args
		wantRes uint64
		wantErr bool
	}{
		{name: "int 77 expected", args: args{toConvert: " 77 "}, wantRes: 77, wantErr: false},
		{name: "error expected bed input", args: args{toConvert: " opa "}, wantRes: 0, wantErr: true},
		{name: "error expected empty str", args: args{toConvert: "  "}, wantRes: 0, wantErr: true},
		{name: "error negative", args: args{toConvert: "  -3"}, wantRes: 0, wantErr: true},
		{name: "zero", args: args{toConvert: "  0"}, wantRes: 0, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := StrToUint64(tt.args.toConvert)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("StrToInt() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestStrToInt(t *testing.T) {
	type args struct {
		toConvert string
	}
	tests := []struct {
		name    string
		args    args
		wantRes int
		wantErr bool
	}{
		{name: "int 77 expected", args: args{toConvert: " 77 "}, wantRes: 77, wantErr: false},
		{name: "error expected bed input", args: args{toConvert: " opa "}, wantRes: 0, wantErr: true},
		{name: "error expected empty str", args: args{toConvert: "  "}, wantRes: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := StrToInt(tt.args.toConvert)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("StrToInt() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestStrToStr(t *testing.T) {
	type args struct {
		toConvert string
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
		wantErr bool
	}{
		{name: "str ' some str ' expected", args: args{toConvert: " some str "}, wantRes: " some str ", wantErr: false},
		{name: "empty string", args: args{toConvert: "  "}, wantRes: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := StrToStr(tt.args.toConvert)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("StrToStr() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestStrToStrWithTrim(t *testing.T) {
	type args struct {
		toConvert string
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
		wantErr bool
	}{
		{name: "str ' some str ' expected", args: args{toConvert: " some str "}, wantRes: "some str", wantErr: false},
		{name: "empty string", args: args{toConvert: "  "}, wantRes: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := StrToStrWithTrim(tt.args.toConvert)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("StrToStr() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
func TestStrToStrArray(t *testing.T) {
	type args struct {
		toConvert string
	}
	tests := []struct {
		name    string
		args    args
		wantRes []string
		wantErr bool
	}{
		{name: "str ' some str ' expected", args: args{toConvert: " some str1   , some   str 2  ,  "}, wantRes: []string{"some str1", "some   str 2"}, wantErr: false},
		{name: "empty string", args: args{toConvert: "  "}, wantRes: nil, wantErr: true},
		{name: "empty coma string", args: args{toConvert: "  "}, wantRes: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := StrToStrArray(tt.args.toConvert)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToStrArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("StrToStrArray() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestStrToStrArrayWithSeparator(t *testing.T) {
	type args struct {
		toConvert string
		separator string
	}
	tests := []struct {
		name    string
		args    args
		wantRes []string
		wantErr bool
	}{
		{name: "str ' some str ' expected", args: args{toConvert: " some str1   , some   str 2  ,  ", separator: ","}, wantRes: []string{"some str1", "some   str 2"}, wantErr: false},
		{name: "empty string", args: args{toConvert: "  ", separator: ","}, wantRes: nil, wantErr: true},
		{name: "empty coma string", args: args{toConvert: "  ", separator: ","}, wantRes: nil, wantErr: true},
		{name: "str semicolon separated", args: args{toConvert: " some str1   ; some   str 2  ;  ", separator: ";"}, wantRes: []string{"some str1", "some   str 2"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := StrToStrArrayWithSeparator(tt.args.toConvert, tt.args.separator)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToStrArrayWithSeparator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("StrToStrArrayWithSeparator() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestStrToMapBool(t *testing.T) {
	type args struct {
		toConvert string
	}
	tests := []struct {
		name    string
		args    args
		wantRes map[string]bool
		wantErr bool
	}{
		{
			name:    "OK",
			args:    args{toConvert: "a,b,c"},
			wantRes: map[string]bool{"a": true, "b": true, "c": true},
		},
		{
			name:    "OK one key",
			args:    args{toConvert: "    a  b c "},
			wantRes: map[string]bool{"a  b c": true},
		},
		{
			name:    "ERROR empty string",
			args:    args{toConvert: ""},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := StrToMapBool(tt.args.toConvert)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToMapBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("StrToMapBool() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestStrToMap(t *testing.T) {
	type args struct {
		toConvert string
	}
	tests := []struct {
		name    string
		args    args
		wantRes map[string]string
		wantErr bool
	}{
		{
			name:    "OK",
			args:    args{toConvert: "a:true,b:valB,c:valC"},
			wantRes: map[string]string{"a": "true", "b": "valB", "c": "valC"},
		},
		{
			name:    "OK one K:V",
			args:    args{toConvert: "b:valB"},
			wantRes: map[string]string{"b": "valB"},
		},
		{
			name:    "Empty string",
			args:    args{toConvert: ""},
			wantRes: nil,
			wantErr: true,
		},
		{
			name:    "String only",
			args:    args{toConvert: "opa string"},
			wantRes: nil,
			wantErr: true,
		},
		{
			name:    "Complex example of string",
			args:    args{toConvert: "opa string:http://google.com"},
			wantRes: map[string]string{"opa string": "http://google.com"},
			wantErr: false,
		},
		{
			name:    "double column",
			args:    args{toConvert: "opa string::http://google.com"},
			wantRes: map[string]string{"opa string": ":http://google.com"},
			wantErr: false,
		},
		{
			name:    "Complex example of string #2",
			args:    args{toConvert: "opa string:http://google.com , next_key::some:::strange_::value"},
			wantRes: map[string]string{"opa string": "http://google.com", "next_key": ":some:::strange_::value"},
			wantErr: false,
		},
		{
			name:    "key no value",
			args:    args{toConvert: "opa string::http://google.com,key_noval:"},
			wantRes: nil,
			wantErr: true,
		},
		{
			name:    "value no key",
			args:    args{toConvert: "opa string::http://google.com,:val_no_key"},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := StrToMap(tt.args.toConvert)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("StrToMap() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
