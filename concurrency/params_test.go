package concurrency

import (
	"errors"
	"testing"
	"time"
)

func Test_Params_SetInt(t *testing.T) {
	params := Params{}

	params.SetInt("x", 1)
	got := params["x"]

	assertEqual(t, got, "1", "should return proper parameter value")
}

func Test_Params_SetBool(t *testing.T) {
	params := Params{}

	params.SetBool("x", true)
	got := params["x"]

	assertEqual(t, got, "true", "should return proper parameter value")
}

func Test_Params_SetDuration(t *testing.T) {
	params := Params{}

	params.SetDuration("x", time.Second)
	got := params["x"]

	assertEqual(t, got, "1s", "should return proper parameter value")
}

func Test_Params_SetDate(t *testing.T) {
	params := Params{}

	params.SetDate("x", time.Date(2000, 3, 5, 0, 0, 0, 0, time.UTC), "2006-01-02")
	got := params["x"]

	assertEqual(t, got, "2000-03-05", "should return proper parameter value")
}

func Test_Params_SetText_valid(t *testing.T) {
	params := Params{}

	err := params.SetText("x", time.Date(2000, 3, 5, 13, 20, 0, 0, time.UTC))
	got := params["x"]

	assertNoError(t, err, "should not return any error")
	assertEqual(t, got, "2000-03-05T13:20:00Z", "should return proper parameter value")
}

type badTextMarshaler struct{}

func (badTextMarshaler) MarshalText() ([]byte, error) {
	return nil, errors.New("failing")
}

func Test_Params_SetText_invalid(t *testing.T) {
	params := Params{}

	err := params.SetText("x", badTextMarshaler{})
	got := params["x"]

	assertError(t, err, "should tell that it failed to parse the value")
	assertEqual(t, got, "", "should return proper parameter value")
}

func Test_Params_SetText_nil(t *testing.T) {
	params := Params{}

	err := params.SetText("x", nil)
	got := params["x"]

	assertError(t, err, "should tell that it failed to parse the value")
	assertEqual(t, got, "", "should return proper parameter value")
}

func Test_Params_SetJSON_valid(t *testing.T) {
	params := Params{}

	err := params.SetJSON("x", x{A: "abc"})
	got := params["x"]

	assertNoError(t, err, "should not return any error")
	assertEqual(t, got, "{\"A\":\"abc\"}", "should return proper parameter value")
}

type badJSONMarshaler struct{}

func (badJSONMarshaler) MarshalJSON() ([]byte, error) {
	return nil, errors.New("failing")
}

func Test_Params_SetJSON_invalid(t *testing.T) {
	params := Params{}

	err := params.SetJSON("x", badJSONMarshaler{})
	got := params["x"]

	assertError(t, err, "should tell that it failed to parse the value")
	assertEqual(t, got, "", "should return proper parameter value")
}

func Test_Params_SetJSON_nil(t *testing.T) {
	params := Params{}

	err := params.SetJSON("x", nil)
	got := params["x"]

	assertError(t, err, "should tell that it failed to parse the value")
	assertEqual(t, got, "", "should return proper parameter value")
}