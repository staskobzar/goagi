package goagi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRespOkZero(t *testing.T) {
	str := "200 result=0\n"
	r, err := cmdParse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.Equal(t, "0", r.result)
	assert.Equal(t, r.data, "")
}

func TestRespOkOne(t *testing.T) {
	str := "200 result=1\n"
	r, err := cmdParse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.Equal(t, "1", r.result)
	assert.Equal(t, r.data, "")
}

func TestRespOkMinusOne(t *testing.T) {
	str := "200 result=-1\n"
	r, err := cmdParse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.Equal(t, "-1", r.result)
	assert.Equal(t, r.data, "")
}
