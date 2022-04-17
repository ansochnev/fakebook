package backend

import (
	"errors"
	"fmt"
	"testing"
)

func TestRender(t *testing.T) {
	err := NewError(ErrBadRequest)
	fmt.Println(err)
}

func TestRender2(t *testing.T) {
	err := NewError(ErrBadRequest).Wrap(errors.New("myerror"))
	fmt.Println(err)
}

func TestRender3(t *testing.T) {
	err := invalidParamError("hello", "world")
	fmt.Println(err)

}

func TestRender4(t *testing.T) {
	err := emptyParamError("world")
	fmt.Println(err)
}

func TestRender5(t *testing.T) {
	err := invalidParamError("foo", "bar")
	fmt.Println(err)
}

func TestRender6(t *testing.T) {
	err := NewError(ErrBadRequest).WithMessage("%s: %s", "test", "render6")
	fmt.Println(err)
}
