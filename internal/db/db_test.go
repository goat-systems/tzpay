package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Open(t *testing.T) {
	db, err := New(nil, "./tzpay.db")
	assert.Nil(t, err)
	assert.NotNil(t, db)
	err = os.Remove("./tzpay.db")
	assert.Nil(t, err)
}
