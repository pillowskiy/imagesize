package imagerrors_test

import (
	"errors"
	"testing"

	"github.com/pillowskiy/imagesize/imagerrors"
)

func TestJoin(t *testing.T) {
	t.Parallel()

	t.Run("AllInputsNilErrors", func(t *testing.T) {
		err := imagerrors.Join(nil, nil, nil)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("SingleNonNilError", func(t *testing.T) {
		expected := errors.New("single error")
		err := imagerrors.Join(nil, expected, nil)

		if err == nil {
			t.Fatal("expected non-nil error, got nil")
		}

		if err.Error() != expected.Error() {
			t.Errorf("expected error %q, got %q", expected.Error(), err.Error())
		}
	})

	t.Run("MultipleErrors", func(t *testing.T) {
		err1 := errors.New("error 1")
		err2 := errors.New("error 2")
		err3 := errors.New("error 3")

		err := imagerrors.Join(err1, nil, err2, err3)
		if err == nil {
			t.Fatal("expected non-nil error, got nil")
		}

		expected := "error 1\nerror 2\nerror 3"
		if err.Error() != expected {
			t.Errorf("expected error %q, got %q", expected, err.Error())
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		err1 := errors.New("error 1")
		err2 := errors.New("error 2")
		err3 := errors.New("error 3")

		joinedErr := imagerrors.Join(err1, err2, err3)
		if joinedErr == nil {
			t.Fatal("expected non-nil error, got nil")
		}

		var unwrap interface{ Unwrap() []error }
		if !errors.As(joinedErr, &unwrap) {
			t.Fatalf("expected error to implement Unwrap, got %T", joinedErr)
		}

		unwrapped := unwrap.Unwrap()
		if len(unwrapped) != 3 {
			t.Fatalf("expected 3 unwrapped errors, got %d", len(unwrapped))
		}

		if unwrapped[0] != err1 || unwrapped[1] != err2 || unwrapped[2] != err3 {
			t.Errorf("unexpected unwrapped errors: %v", unwrapped)
		}
	})

	t.Run("EmptyInput", func(t *testing.T) {
		err := imagerrors.Join()
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("OneErrorWithoutNewline", func(t *testing.T) {
		err1 := errors.New("no newline error")
		joinedErr := imagerrors.Join(err1)
		expected := "no newline error"
		if joinedErr.Error() != expected {
			t.Errorf("expected %q, got %q", expected, joinedErr.Error())
		}
	})
}
