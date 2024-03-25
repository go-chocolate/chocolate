package jsonutil

import "testing"

func TestScheme(t *testing.T) {
	scheme := &Scheme{
		Type: OBJECT,
		Schemes: []*Scheme{
			{Name: "a", Type: STRING, Value: "hello"},
			{Name: "b", Type: INTEGER},
			{Name: "c", Type: ARRAY, Schemes: []*Scheme{
				{Type: OBJECT, Schemes: []*Scheme{
					{Name: "c.a", Type: STRING, Value: "hello"},
				}},
			}},
			{Name: "d", Type: BOOLEAN},
		},
	}

	t.Log(scheme.Validate())
	t.Log(string(scheme.Build()))
}
