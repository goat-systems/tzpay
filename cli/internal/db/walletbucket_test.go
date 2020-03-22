package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InitWallet(t *testing.T) {
	type input struct {
		secret   string
		password string
	}

	type want struct {
		err bool
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				"edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
				"password12345",
			},
			want{
				false,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			db, err := New(nil, "./tzpay.db")
			assert.Nil(t, err)
			assert.NotNil(t, db)

			err = db.InitWallet(tt.input.password, tt.input.secret)
			checkErr(t, tt.want.err, "", err)

			err = os.Remove("./tzpay.db")
			assert.Nil(t, err)
		})
	}
}

func Test_GetSecret(t *testing.T) {
	type store struct {
		secret   string
		password string
	}

	type input struct {
		store    store
		password string
	}

	type want struct {
		err         bool
		errContains string
		esesk       string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				store{
					"edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
					"password12345",
				},
				"password12345",
			},
			want{
				false,
				"",
				"edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
			},
		},
		{
			"handles unauthorized",
			input{
				store{
					"edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
					"password12345",
				},
				"wrong password",
			},
			want{
				true,
				"failed to authorize",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			db, err := New(nil, "./tzpay.db")
			assert.Nil(t, err)
			assert.NotNil(t, db)

			err = db.InitWallet(tt.input.store.password, tt.input.store.secret)
			checkErr(t, false, "", err)

			edesk, err := db.GetSecret(tt.input.password)
			checkErr(t, tt.want.err, tt.want.errContains, err)

			assert.Equal(t, tt.want.esesk, edesk)

			err = os.Remove("./tzpay.db")
			assert.Nil(t, err)
		})
	}
}

func Test_IsWalletInitialized(t *testing.T) {
	type store struct {
		init     bool
		secret   string
		password string
	}

	type input struct {
		store    store
		password string
	}

	type want struct {
		ok bool
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is initialized",
			input{
				store{
					true,
					"edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
					"password12345",
				},
				"password12345",
			},
			want{
				true,
			},
		},
		{
			"is not initialized",
			input{
				store{
					false,
					"",
					"",
				},
				"wrong password",
			},
			want{
				false,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			db, err := New(nil, "./tzpay.db")
			assert.Nil(t, err)
			assert.NotNil(t, db)

			if tt.input.store.init {
				err = db.InitWallet(tt.input.store.password, tt.input.store.secret)
				checkErr(t, false, "", err)
			}

			ok := db.IsWalletInitialized()
			assert.Equal(t, tt.want.ok, ok)

			err = os.Remove("./tzpay.db")
			assert.Nil(t, err)
		})
	}
}

func checkErr(t *testing.T, wantErr bool, errContains string, err error) {
	if wantErr {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errContains)
	} else {
		assert.Nil(t, err)
	}
}
