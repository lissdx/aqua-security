package utils

import (
	//"github.cloudoptix.io/Bayonet/to-data-pipeline/internal/pkg/consumer"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sort"
	"testing"
)

func TestMapAny_String(t *testing.T) {
	type args struct {
		v any
	}
	type testCase[E any] struct {
		name    string
		args    args
		want    *E
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name: "String OK",
			args: args{v: "test string"},
			want: func() *string {
				res := "test string"
				return &res
			}(),
			wantErr: false,
		},
		{
			name: "String nil",
			args: args{v: nil},
			want: func() *string {
				return nil
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapAny[string](tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapAny() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, *tt.want, *got)
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("MapAny() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestMapAny_Struct(t *testing.T) {
	type internalStruct1 struct {
		is  string
		isp *string
	}
	type internalStruct2 struct {
		is   string
		isp  *string
		i2ip *int
		il2  []int
	}

	type testStruct struct {
		s    *string
		i    int
		is1  internalStruct1
		is2p *internalStruct2
	}

	type args struct {
		v any
	}
	type testCase[E any] struct {
		name    string
		args    args
		want    *E
		wantErr bool
	}
	tests := []testCase[testStruct]{
		{
			name: "struct OK",
			args: args{v: testStruct{}},
			want: func() *testStruct {
				return &testStruct{}
			}(),
			wantErr: false,
		},
		{
			name: "struct nil",
			args: args{v: nil},
			want: func() *testStruct {
				return nil
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapAny[testStruct](tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapAny() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MapAny() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

//func TestMapAny_ConsumerConsumerMessage(t *testing.T) {
//
//	type args struct {
//		v any
//	}
//	type testCase[E any] struct {
//		name    string
//		args    args
//		want    *E
//		wantErr bool
//	}
//	tests := []testCase[*consumer.ConsumerMessage]{
//		{
//			name: "struct OK",
//			args: args{v: &consumer.ConsumerMessage{}},
//			want: func() **consumer.ConsumerMessage {
//				res := &consumer.ConsumerMessage{}
//				return &res
//			}(),
//			wantErr: false,
//		},
//		{
//			name:    "struct nil",
//			args:    args{v: nil},
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := MapAny[*consumer.ConsumerMessage](tt.args.v)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("MapAny() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//
//			if !tt.wantErr {
//				if !reflect.DeepEqual(got, tt.want) {
//					t.Errorf("MapAny() got = %v, want %v", *got, *tt.want)
//				}
//			}
//		})
//	}
//}

func TestCopyMapNil(t *testing.T) {
	var nilMap map[interface{}]interface{} = nil
	var want map[interface{}]interface{} = nil
	assert.Equal(t, want, CopyMap(nilMap))
}

func TestCopyMap(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want map[K]V
	}
	tests := []testCase[string, interface{}]{
		{
			name: "OK empty",
			args: args[string, interface{}]{m: make(map[string]interface{})},
			want: make(map[string]interface{}),
		},
		{
			name: "OK with nil",
			args: args[string, interface{}]{m: func() map[string]interface{} {
				m := make(map[string]interface{})
				m["nil"] = nil
				m["string"] = "string"
				return m
			}()},
			want: func() map[string]interface{} {
				m := make(map[string]interface{})
				m["nil"] = nil
				m["string"] = "string"
				return m
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CopyMap(tt.args.m), "CopyMap(%v)", tt.args.m)
		})
	}
}

func TestRemoveDuplicatesString(t *testing.T) {
	type args[E comparable] struct {
		l []E
	}
	type testCase[E comparable] struct {
		name string
		args args[E]
		want []E
	}
	tests := []testCase[string]{
		{name: "OK",
			//args: args[]{},
			args: args[string]{
				l: func() []string {
					return []string{"A", "B", "C"}
				}(),
			},
			want: []string{"A", "B", "C"},
		},
		{name: "filter",
			//args: args[]{},
			args: args[string]{
				l: func() []string {
					return []string{"A", "A", "A", "B", "C", "C"}
				}(),
			},
			want: []string{"A", "B", "C"},
		},
		{name: "empty",
			//args: args[]{},
			args: args[string]{
				l: func() []string {
					return []string{}
				}(),
			},
			want: []string{},
		},
		{name: "null",
			//args: args[]{},
			args: args[string]{
				l: func() []string {
					return nil
				}(),
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Strings(tt.want)
			sRes := RemoveDuplicates(tt.args.l)
			sort.Strings(sRes)
			assert.Equalf(t, tt.want, sRes, "RemoveDuplicates(%v)", tt.args.l)
		})
	}
}

func TestRemoveDuplicatesInt(t *testing.T) {
	type args[E comparable] struct {
		l []E
	}
	type testCase[E comparable] struct {
		name string
		args args[E]
		want []E
	}
	tests := []testCase[int]{
		{name: "OK",
			//args: args[]{},
			args: args[int]{
				l: func() []int {
					return []int{-1, 5, -100}
				}(),
			},
			want: []int{-1, 5, -100},
		},
		{name: "filter",
			//args: args[]{},
			args: args[int]{
				l: func() []int {
					return []int{-1, 5, -100, -100, 5, 5, 1, 0, -100, 5, 0}
				}(),
			},
			want: []int{-1, 5, -100, 1, 0},
		},
		{name: "empty",
			//args: args[]{},
			args: args[int]{
				l: func() []int {
					return []int{}
				}(),
			},
			want: []int{},
		},
		{name: "null",
			//args: args[]{},
			args: args[int]{
				l: func() []int {
					return nil
				}(),
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Ints(tt.want)
			sRes := RemoveDuplicates(tt.args.l)
			sort.Ints(sRes)
			assert.Equalf(t, tt.want, sRes, "RemoveDuplicates(%v)", tt.args.l)
		})
	}
}
