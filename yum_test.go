package yum

import (
	"testing"
)

func TestProvides(t *testing.T) {
	const providedFile string = "/usr/sbin/postfix"
	result := Package{
		Name:         "postfix",
		Architecture: "x86_64",
		Epoch:        2,
		Version:      "2.11.3",
		Release:      "1.fc21",
		Repository:   "fedora",
		Summary:      "Postfix Mail Transport Agent",
	}

	yumPackages, err := Provides(providedFile)
	if err != nil {
		t.Fatal(err)
	}

	for key, _ := range yumPackages {
		if key != providedFile {
			t.Errorf("Retrieved an extra provided file: %v\n", key)
		}
	}
	if t.Failed() {
		t.FailNow()
	}

	if len(yumPackages[providedFile]) != 1 {
		t.Fatalf("Expected %v package(s) but got %v\n", 1, len(yumPackages[providedFile]))
	}
	if yumPackages[providedFile][0] != result {
		t.Fatalf("%v does not match expected %v\n", yumPackages[providedFile], result)
	}
}
