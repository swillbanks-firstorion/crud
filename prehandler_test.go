package crud

import (
	"encoding/json"
	"errors"
	"github.com/swillbanks-firstorion/crud/option"
	"net/url"
	"testing"
)

type TestAdapter struct{}

func (t *TestAdapter) Install(router *Router, spec *Spec) error {
	return nil
}

func (t *TestAdapter) Serve(swagger *Swagger, addr string) error {
	return nil
}

func TestQueryValidation(t *testing.T) {
	r := NewRouter("", "", &TestAdapter{})

	tests := []struct {
		Schema   map[string]Field
		Input    string
		Expected error
	}{
		{
			Schema: map[string]Field{
				"testquery": String(),
			},
			Input:    "",
			Expected: nil,
		}, {
			Schema: map[string]Field{
				"testquery": String().Required(),
			},
			Input:    "",
			Expected: errRequired,
		}, {
			Schema: map[string]Field{
				"testquery": String().Required(),
			},
			Input:    "testquery=",
			Expected: errRequired,
		}, {
			Schema: map[string]Field{
				"testquery": String().Required(),
			},
			Input:    "testquery=ok",
			Expected: nil,
		}, {
			Schema: map[string]Field{
				"testquery": Number(),
			},
			Input:    "",
			Expected: nil,
		}, {
			Schema: map[string]Field{
				"testquery": Number().Required(),
			},
			Input:    "",
			Expected: errRequired,
		},
		{
			Schema: map[string]Field{
				"testquery": Number().Required(),
			},
			Input:    "testquery=1",
			Expected: nil,
		},
		{
			Schema: map[string]Field{
				"testquery": Number().Required(),
			},
			Input:    "testquery=1.1",
			Expected: nil,
		},
		{
			Schema: map[string]Field{
				"testquery": Number(),
			},
			Input:    "testquery=a",
			Expected: errWrongType,
		},
		{
			Schema: map[string]Field{
				"testquery": Boolean(),
			},
			Input:    "testquery=true",
			Expected: nil,
		},
		{
			Schema: map[string]Field{
				"testquery": Boolean(),
			},
			Input:    "testquery=false",
			Expected: nil,
		},
		{
			Schema: map[string]Field{
				"testquery": Boolean(),
			},
			Input:    "testquery=1",
			Expected: errWrongType,
		},
		{
			Schema: map[string]Field{
				"testquery": Integer(),
			},
			Input:    "testquery=1",
			Expected: nil,
		},
		{
			Schema: map[string]Field{
				"testquery": Integer().Max(1),
			},
			Input:    "testquery=2",
			Expected: errMaximum,
		},
		{
			Schema: map[string]Field{
				"testquery": Integer().Min(5),
			},
			Input:    "testquery=4",
			Expected: errMinimum,
		},
		{
			Schema: map[string]Field{
				"testquery": Integer(),
			},
			Input:    "testquery=1.1",
			Expected: errWrongType,
		},
		{
			Schema: map[string]Field{
				"testquery": Integer(),
			},
			Input:    "testquery=a",
			Expected: errWrongType,
		},
		{
			Schema: map[string]Field{
				"testquery": Integer().Enum(1, 2),
			},
			Input:    "testquery=2",
			Expected: nil,
		},
		{
			Schema: map[string]Field{
				"testquery": Integer().Enum(1, 2),
			},
			Input:    "testquery=3",
			Expected: errEnumNotFound,
		},
		{
			Schema: map[string]Field{
				"testquery": String().Enum("a"),
			},
			Input:    "testquery=a",
			Expected: nil,
		},
		{
			Schema: map[string]Field{
				"testquery": String().Enum("a"),
			},
			Input:    "testquery=b",
			Expected: errEnumNotFound,
		},
		{
			Schema: map[string]Field{
				"testarray": Array(),
			},
			Input:    "testarray=1&testarray=a",
			Expected: nil,
		},
		{
			Schema: map[string]Field{
				"testquery": Array().Items(Number()),
			},
			Input:    "testquery=1&testquery=2",
			Expected: nil,
		},
		{
			Schema: map[string]Field{
				"testquery": Array().Items(Number()),
			},
			Input:    "testquery=1&testquery=a",
			Expected: errWrongType,
		},
		{
			Schema: map[string]Field{
				"testquery": Array().Min(2),
			},
			Input:    "testquery=z",
			Expected: errMinimum,
		},
		{
			Schema: map[string]Field{
				"testquery": Array().Min(2),
			},
			Input:    "testquery=z&testquery=x",
			Expected: nil,
		},
		{
			Schema: map[string]Field{
				"testquery": Array(),
			},
			Input:    "testquery=d",
			Expected: nil,
		},
	}

	for i, test := range tests {
		query, err := url.ParseQuery(test.Input)
		if err != nil {
			t.Fatal(err)
		}

		err = r.Validate(Validate{Query: Object(test.Schema)}, query, nil, nil)

		if !errors.Is(err, test.Expected) {
			t.Errorf("%v: expected '%v' got '%v'. input: '%v'. schema: '%v'", i, test.Expected, err, test.Input, test.Schema)
		}
	}
}

func TestQueryDefaults(t *testing.T) {
	r := NewRouter("", "", &TestAdapter{})

	tests := []struct {
		Schema   map[string]Field
		Input    string
		Expected string
	}{
		{
			Schema: map[string]Field{
				"q": String().Default("hey"),
			},
			Input:    "",
			Expected: "q=hey",
		},
		{
			Schema: map[string]Field{
				"q1": String().Default("1"),
				"q2": String().Default("2"),
			},
			Input:    "",
			Expected: "q1=1&q2=2",
		},
	}

	for i, test := range tests {
		query, err := url.ParseQuery(test.Input)
		if err != nil {
			t.Fatal(err)
		}

		err = r.Validate(Validate{Query: Object(test.Schema)}, query, nil, nil)

		if query.Encode() != test.Expected {
			t.Errorf("%v: expected '%v' got '%v'. input: '%v'. schema: '%v'", i, test.Expected, query.Encode(), test.Input, test.Schema)
		}
		if err != nil {
			t.Error(err)
		}
	}
}

func TestSimpleBodyValidation(t *testing.T) {
	r := NewRouter("", "", &TestAdapter{})

	tests := []struct {
		Schema   Field
		Input    interface{}
		Expected error
	}{
		{
			Schema:   Number(),
			Input:    float64(1),
			Expected: nil,
		},
		{
			Schema:   Number(),
			Input:    1,
			Expected: errWrongType,
		},
		{
			Schema:   Number(),
			Input:    "a",
			Expected: errWrongType,
		},
		{
			Schema:   String(),
			Input:    "2",
			Expected: nil,
		},
		{
			Schema:   Boolean(),
			Input:    true,
			Expected: nil,
		},
		{
			Schema:   Boolean(),
			Input:    false,
			Expected: nil,
		},
		{
			Schema:   Boolean(),
			Input:    `1`,
			Expected: errWrongType,
		},
	}

	for _, test := range tests {
		err := r.Validate(Validate{Body: test.Schema}, nil, test.Input, nil)

		if !errors.Is(err, test.Expected) {
			t.Errorf("expected '%v' got '%v'. input: '%v'. schema: '%v'", test.Expected, err, test.Input, test.Schema)
		}
	}
}

func TestBodyValidation(t *testing.T) {
	r := NewRouter("", "", &TestAdapter{})

	tests := []struct {
		Schema   map[string]Field
		Input    string
		Expected error
	}{
		{
			Schema: map[string]Field{
				"int": Integer(),
			},
			Input:    `{}`,
			Expected: nil,
		},
		{
			Schema: map[string]Field{
				"int": Integer().Required(),
			},
			Input:    `{}`,
			Expected: errRequired,
		},
		{
			Schema: map[string]Field{
				"int": Integer().Required(),
			},
			Input:    `{"int":1}`,
			Expected: nil,
		},
		{
			Schema: map[string]Field{
				"int": Integer().Required(),
			},
			Input:    `{"int":1.9}`,
			Expected: errWrongType,
		},
		{
			Schema: map[string]Field{
				"float": Number().Required(),
			},
			Input:    `{"float":-1}`,
			Expected: nil,
		},
		{
			Schema: map[string]Field{
				"float": Number().Required(),
			},
			Input:    `{"float":1.1}`,
			Expected: nil,
		}, {
			Schema: map[string]Field{
				"obj1": Object(map[string]Field{
					"inner": Number().Required(),
				}),
			},
			Input:    `{"obj1":{"inner":1}}`,
			Expected: nil,
		}, {
			Schema: map[string]Field{
				"obj2": Object(map[string]Field{
					"inner": Number().Required(),
				}),
			},
			Input:    `{"obj2":{"inner":"not a number"}}`,
			Expected: errWrongType,
		}, {
			Schema: map[string]Field{
				"arr1": Array(),
			},
			Input:    `{"arr1":[1,"a"]}`,
			Expected: nil,
		}, {
			Schema: map[string]Field{
				"arr1": Array().Items(Number()),
			},
			Input:    `{"arr1":[1]}`,
			Expected: nil,
		}, {
			Schema: map[string]Field{
				"arr2": Array().Items(Number()),
			},
			Input:    `{"arr2":["a"]}`,
			Expected: errWrongType,
		}, {
			Schema: map[string]Field{
				"arr3": Array().Min(2),
			},
			Input:    `{"arr3":["a"]}`,
			Expected: errMinimum,
		}, {
			Schema: map[string]Field{
				"complex1": Object(map[string]Field{
					"array": Array().Required().Items(Object(map[string]Field{
						"id": Number().Required(),
					})),
				}).Required(),
			},
			Input:    `{"complex1":{"array":[{"id":1}]}}`,
			Expected: nil,
		}, {
			Schema: map[string]Field{
				"complex2": Object(map[string]Field{
					"array": Array().Required().Items(Object(map[string]Field{
						"id": Number().Required(),
					})),
				}).Required(),
			},
			Input:    `{"complex2":{"array":[{"id":"a"}]}}`,
			Expected: errWrongType,
		},
	}

	for _, test := range tests {
		var input interface{}
		if err := json.Unmarshal([]byte(test.Input), &input); err != nil {
			t.Fatal(err)
		}

		err := r.Validate(Validate{Body: Object(test.Schema)}, nil, input, nil)

		if !errors.Is(err, test.Expected) {
			t.Errorf("expected '%v' got '%v'. input: '%v'. schema: '%v'", test.Expected, err, test.Input, test.Schema)
		}
	}
}

func TestBodyStripUnknown(t *testing.T) {
	r := NewRouter("", "", &TestAdapter{}, option.StripUnknown(true))

	tests := []struct {
		Schema   map[string]Field
		Input    string
		Expected string
	}{
		{
			Schema: map[string]Field{
				"str": String(),
			},
			Input:    `{"str":"ok","unknown1":1}`,
			Expected: `{"str":"ok"}`,
		},
		{
			Schema: map[string]Field{
				"str2": String().Default("Hello"),
			},
			Input:    `{}`,
			Expected: `{"str2":"Hello"}`,
		},
		{
			Schema: map[string]Field{
				"int1": Integer().Default(1),
			},
			Input:    `{}`,
			Expected: `{"int1":1}`,
		},
	}

	for _, test := range tests {
		var input interface{}
		if err := json.Unmarshal([]byte(test.Input), &input); err != nil {
			t.Fatal(err)
		}

		err := r.Validate(Validate{Body: Object(test.Schema)}, nil, input, nil)

		if err != nil {
			t.Error(err)
			continue
		}

		data, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(data) != test.Expected {
			t.Errorf("expected '%v' got '%v'. input: '%v'. schema: '%v'", test.Expected, string(data), test.Input, test.Schema)
		}
	}
}

func TestBodyErrorUnknown(t *testing.T) {
	r := NewRouter("", "", &TestAdapter{}, option.AllowUnknown(false))

	tests := []struct {
		Schema   map[string]Field
		Input    string
		Expected string
	}{
		{
			Schema: map[string]Field{
				"str": String(),
			},
			Input:    `{"str":"ok","unknown1":1}`,
			Expected: `{"str":"ok"}`,
		},
	}

	for i, test := range tests {
		var input interface{}
		if err := json.Unmarshal([]byte(test.Input), &input); err != nil {
			t.Fatal(err)
		}

		err := r.Validate(Validate{Body: Object(test.Schema)}, nil, input, nil)

		if !errors.Is(err, errUnknown) {
			t.Errorf("%v: expected '%v' got '%v'", i, errUnknown, err)
			continue
		}
	}
}

func TestPathValidation(t *testing.T) {
	r := NewRouter("", "", &TestAdapter{})

	tests := []struct {
		Schema   map[string]Field
		Input    string
		Expected error
	}{
		{
			Schema: map[string]Field{
				"id": Integer(),
			},
			Input:    ``,
			Expected: nil,
		},
		{
			Schema: map[string]Field{
				"int": Integer().Required(),
			},
			Input:    ``,
			Expected: errRequired,
		},
		{
			Schema: map[string]Field{
				"id": Integer().Required(),
			},
			Input:    `a`,
			Expected: errWrongType,
		},
		{
			Schema: map[string]Field{
				"id": Integer().Required(),
			},
			Input:    `1`,
			Expected: nil,
		},
	}

	for _, test := range tests {
		input := map[string]string{
			"id": test.Input,
		}

		err := r.Validate(Validate{Path: Object(test.Schema)}, nil, nil, input)

		if errors.Unwrap(err) != test.Expected {
			t.Errorf("expected '%v' got '%v'. input: '%v'. schema: '%v'", test.Expected, err, test.Input, test.Schema)
		}
	}
}

func TestStrip_Body(t *testing.T) {
	r := NewRouter("", "", &TestAdapter{}, option.StripUnknown(true))
	var input interface{}
	input = map[string]interface{}{
		"id": "blah",
	}

	err := r.Validate(Validate{Body: Object(map[string]Field{})}, nil, input, nil)

	result := input.(map[string]interface{})
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if _, ok := result["id"]; ok {
		t.Errorf("expected the value to be stripped %v", result["id"])
	}

	input = map[string]interface{}{
		"id": "blah",
	}

	err = r.Validate(Validate{Body: Object(map[string]Field{}).Strip(false)}, nil, input, nil)
	result = input.(map[string]interface{})
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if _, ok := result["id"]; !ok {
		t.Errorf("expected the value not to be stripped")
	}
}

func TestUnknown_Body(t *testing.T) {
	r := NewRouter("", "", &TestAdapter{}, option.AllowUnknown(true))
	var input interface{}
	input = map[string]interface{}{
		"id": "blah",
	}

	err := r.Validate(Validate{Body: Object(map[string]Field{})}, nil, input, nil)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	input = map[string]interface{}{
		"id": "blah",
	}

	err = r.Validate(Validate{Body: Object(map[string]Field{}).Unknown(false)}, nil, input, nil)
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestUnknown_Query(t *testing.T) {
	r := NewRouter("", "", &TestAdapter{}, option.AllowUnknown(false))
	query := url.Values{"unknown": []string{"value"}}
	spec := Object(map[string]Field{})

	err := r.Validate(Validate{Query: spec}, query, nil, nil)
	if !errors.Is(err, errUnknown) {
		t.Errorf("Expected '%s' but got '%s'", errUnknown, err)
	}

	spec = spec.Unknown(true)
	err = r.Validate(Validate{Query: spec}, query, nil, nil)
	if err != nil {
		t.Error("Unexpected error", err)
	}
}

func TestStrip_Query(t *testing.T) {
	r := NewRouter("", "", &TestAdapter{})
	query := url.Values{"unknown": []string{"value"}}
	spec := Object(map[string]Field{}).Strip(false)

	err := r.Validate(Validate{Query: spec}, query, nil, nil)
	if err != nil {
		t.Error("Unexpected error", err)
	}

	if _, ok := query["unknown"]; !ok {
		t.Error("Expected the value to not have been stripped")
	}

	spec = spec.Strip(true)
	err = r.Validate(Validate{Query: spec}, query, nil, nil)
	if err != nil {
		t.Error("Unexpected error", err)
	}

	if _, ok := query["unknown"]; ok {
		t.Error("Expected the value to have been stripped")
	}
}

func Test_BodyValidateRequiredAutomatically(t *testing.T) {
	r := NewRouter("", "", &TestAdapter{}, option.AllowUnknown(false))

	err := r.Validate(Validate{Body: Object(map[string]Field{})}, nil, nil, nil)

	if !errors.Is(err, errRequired) {
		t.Error("Expected errRequired got", err)
	}
}
