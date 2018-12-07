package timebot

import "testing"
import "time"

func TestParseTime(t *testing.T) {
	input := "2018-12-06 19:00 PST"
	laTime, err := time.LoadLocation("America/Los_Angeles")

	if err != nil {
		t.Fatal(err)
	}

	expected := time.Date(2018, 12, 6, 19, 0, 0, 0, laTime)

	if tried, err := ParseTime(input); err != nil || !tried.Equal(expected) {
		t.Fatal("Expected", expected, "but got", tried)
	}

	input = "2018-12-07 12:00 KST"
	expected = time.Date(2018, 12, 07, (12 - 9), 0, 0, 0, time.UTC)

	if tried, err := ParseTime(input); err != nil || !tried.Equal(expected) {
		t.Fatal("Expected", expected, "but got", tried.In(time.UTC))
	}
}

func TestToKoreaTime(t *testing.T) {
	input := "2018-12-06 19:00 PST"
	timed, err := ParseTime(input)

	if err != nil {
		t.Fatal(err)
	}

	expected := "2018-12-07 12:00 KST"

	if tried := ToKoreaTime(timed); tried != expected {
		t.Fatal("Expected", expected, "but got", tried)
	}
}

func TestToCaliforniaTime(t *testing.T) {
	input := "2018-12-07 12:00 KST"
	timed, err := ParseTime(input)

	if err != nil {
		t.Fatal(err)
	}

	expected := "2018-12-06 19:00 PST"

	if tried := ToCaliforniaTime(timed); tried != expected {
		t.Fatal("Expected", expected, "but got", tried)
	}
}

func TestIsTargetMessage(t *testing.T) {
	input := "aljdslf"

	if result, ok := IsTargetMessage(input); ok || result != "" {
		t.Fatal(input, "should be  not interesting :(")
	}

	input = "2018-01-03 10:37 PST"

	if result, ok := IsTargetMessage(input); result != PST || !ok {
		t.Fatal(input, "should be interesting :(")
	}
}
