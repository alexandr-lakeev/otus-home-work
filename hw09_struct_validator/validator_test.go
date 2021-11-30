package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Person struct {
		Name     string   `validate:"require"`
		Document Document `validate:"nested"`
	}

	Document struct {
		Type   string `validate:"in:passport,birth_certificate"`
		Number string `validate:"require"`
	}

	UnknownValidator struct {
		Field string `validate:"unkown:123"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "123e4567-e89b-12d3-a456-426655440000",
				Name:   "Test",
				Age:    36,
				Email:  "email@email.com",
				Role:   "admin",
				Phones: []string{"89100000000"},
				meta:   json.RawMessage{},
			},
			nil,
		},
		{
			User{
				ID:     "123e4567",
				Age:    12,
				Email:  "email@email",
				Role:   "hacker",
				Phones: []string{"9100000000"},
			},
			ValidationErrors([]ValidationError{
				{Field: "ID", Err: errors.New("expected len 36, got 8")},
				{Field: "Age", Err: errors.New("expected greater or equal 18, got 12")},
				{Field: "Email", Err: ErrRegexpNotMatched},
				{Field: "Role", Err: errors.New("expected [admin stuff], got hacker")},
				{Field: "Phones[0]", Err: errors.New("expected len 11, got 10")},
			}),
		},
		{
			App{
				Version: "12345",
			},
			nil,
		},
		{
			App{
				Version: "1234",
			},
			ValidationErrors([]ValidationError{
				{Field: "Version", Err: errors.New("expected len 5, got 4")},
			}),
		},
		{
			Token{},
			nil,
		},
		{
			Response{
				Code: 200,
				Body: "{}",
			},
			nil,
		},
		{
			Response{
				Code: 401,
			},
			ValidationErrors([]ValidationError{
				{Field: "Code", Err: errors.New("expected [200 404 500], got 401")},
			}),
		},
		{
			Person{
				Document: Document{
					Type:   "none",
					Number: "",
				},
			},
			ValidationErrors([]ValidationError{
				{Field: "Name", Err: errors.New("can't be empty")},
				{Field: "Document.Type", Err: errors.New("expected [passport birth_certificate], got none")},
				{Field: "Document.Number", Err: errors.New("can't be empty")},
			}),
		},
		{
			UnknownValidator{},
			ErrUnknownValidator,
		},
		{
			123,
			ErrNotStruct,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			require.Equal(t, tt.expectedErr, Validate(tt.in))
		})
	}
}
