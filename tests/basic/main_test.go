package basic

import (
	"testing"
)

// cd ./tests/basic
// go test --coverprofile=coverage.out
// go tool cover -html=coverage.out -o coverage.html
// xdg-open coverage.html

func TestAddOne(t *testing.T) {
	var (
		input  = 1
		output = 3
	)

	actual := AddOne(1)
	if actual != output {
		t.Errorf("AddOne(%d), output %d, actual %d", input, output, actual)
	}

	// assert.Equal(t, 4, AddOne(2), "AddOne(2) should equal 3")
	// assert.NotEqual(t, 2, 3)
	// assert.Nil(t, nil, nil)
}

func TestAddTwo(t *testing.T) {
	var (
		input  = 1
		output = 3
	)

	actual := AddTwo(1)
	if actual != output {
		t.Errorf("AddTwo(%d), output %d, actual %d", input, output, actual)
	}
}

// func TestRequire(t *testing.T) {
// 	require.Equal(t, 2, 3)
// 	fmt.Println("Not executing")
// }

// func TestArrert(t *testing.T) {
// 	assert.Equal(t, 2, 3)
// 	fmt.Println("executing")
// }
