package goson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Reader(t *testing.T) {
	jsonStr := []byte(`{
		"id": "0101",
		"name": "agungdp",
		"mustFloat": "2.23423",
		"mustInt": 2.23423,
		"uint": 11,
		"nullable": null,
		"isExist": "true",
		"obj": {
		  "n": 2,
		  "testing": {
			"ss": "23840923849284"
		  },
		  "s": "Hitomi"
		},
		"slice": [
		  {
			"fieldA": "100",
			"fieldB": 3000,
			"exist": true,
			"test": "3.14"
		  },
		  {
			"fieldA": 50000,
			"fieldB": "123000",
			"exist": 3000,
			"test": 1323.123
		  }
		],
		"str": ["a", true],
		"ints": ["2", 3],
		"bools": [1, "true", "0", true],
		"NoTag": 19283091832
	  }`)

	reader, err := Read(jsonStr)
	assert.NoError(t, err)

	t.Run("Testcase #1: Get string field", func(t *testing.T) {
		str, err := reader.GetField("slice", "[1]", "fieldA").String()
		assert.NoError(t, err)
		assert.Equal(t, str, "50000")

		str, err = reader.GetField("slice", "[5]", "fieldA").String()
		assert.Error(t, err)
	})

	t.Run("Testcase #2: Get int field", func(t *testing.T) {
		i, err := reader.GetField("bools", "[3]").Int()
		assert.NoError(t, err)
		assert.Equal(t, i, 1)

		i, err = reader.GetField("name").Int()
		assert.Error(t, err)
	})

	t.Run("Testcase #3: Get float field", func(t *testing.T) {
		fl, err := reader.GetField("slice", "[0]", "test").Float()
		assert.NoError(t, err)
		assert.Equal(t, fl, 3.14)

		fl, err = reader.GetField("str", "[1]").Float()
		assert.NoError(t, err)
		assert.Equal(t, float64(1), fl)
	})

	t.Run("Testcase #4: Get uint field", func(t *testing.T) {
		u, err := reader.GetField("uint").Uint()
		assert.NoError(t, err)
		assert.Equal(t, u, uint(11))

		u, err = reader.GetField("str", "[1]").Uint()
		assert.NoError(t, err)
	})

	t.Run("Testcase #5: Get bool field", func(t *testing.T) {
		b, err := reader.GetField("isExist").Bool()
		assert.NoError(t, err)
		assert.Equal(t, b, true)

		b, err = reader.GetField("obj", "a", "b").Bool()
		assert.Error(t, err)
	})

	t.Run("Testcase #6: Set to model", func(t *testing.T) {
		var m struct {
			SS string `json:"ss"`
		}

		err := reader.GetField("obj", "testing").SetTo(&m)
		assert.NoError(t, err)
		assert.Equal(t, m.SS, "23840923849284")

		err = reader.GetField("slice", "a", "b").SetTo(&m)
		assert.Error(t, err)
	})

	t.Run("Testcase #6: Invalid JSON format", func(t *testing.T) {
		_, err := Read([]byte(`{a:3}`))
		assert.Error(t, err)
	})

	t.Run("Testcase #7: GetString Reader", func(t *testing.T) {
		b := reader.GetString("obj.s")
		assert.Equal(t, "Hitomi", b)
	})
}
