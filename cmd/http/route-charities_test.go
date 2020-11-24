package main

import (
	"net/http"
	"testing"

	grtest "github.com/Vindexus/go-router-test"

	"github.com/monstercat/golib/expectm"
)

const (
	CharityId = "9d0b23cd-657b-4cc4-8258-a8cabb1f6847"
)

func TestGetCharities(t *testing.T) {
	test := &grtest.RouteTest{
		Path:           "/charities",
		ExpectedStatus: http.StatusOK,
		ExpectedM: &expectm.ExpectedM{
			"Charities.Data.#": 2,
			"Charities.Total":  2,
			"Charities.Limit":  50,
			"Charities.Offset": 0,
		},
	}
	if err := runTest(test); err != nil {
		t.Error(err)
	}
}

func TestGetCharity(t *testing.T) {
	test := &grtest.RouteTest{
		Path:           "/charity/" + CharityId,
		ExpectedStatus: http.StatusOK,
		ExpectedM: &expectm.ExpectedM{
			"Charity.Id": CharityId,
		},
	}
	if err := runTest(test); err != nil {
		t.Error(err)
	}
}
