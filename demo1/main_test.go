package main

import (
	"github.com/kataras/iris/v12/httptest"
	"testing"
)

func TestMvc(t *testing.T) {
	e := httptest.New(t, newApp())

	e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal("当前")
}
