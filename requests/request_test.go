package requests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRequest(t *testing.T) {
	expected := Query{
		"name":     []string{"fancxxy"},
		"language": []string{"golang"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		for key, val := range expected {
			if v := values.Get(key); v != val[0] {
				t.Errorf("expected is %s, actual is %s", val, v)
			}
		}

		bs, _ := json.Marshal(values)
		w.Write(bs)
	}))

	client := New()
	resp, err := client.Get(server.URL, expected)
	if err != nil {
		t.Errorf("expected is %v, actual is %v", nil, err)
	}

	var actual Query
	resp.ToJSON(&actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected is %v, actual is %v", expected, actual)
	}
}
