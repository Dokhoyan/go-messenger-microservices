package validate

import (
	"context"
	"strings"
	"time"

	"github.com/Dokhoyan/common/pkg/sys/validate"
)

func ID(id int64) validate.Condition {
	return func(ctx context.Context) error {
		if id <= 0 {
			return validate.NewValidationErrors("id must be greater than 0")
		}

		return nil
	}
}

func Birthday(birthday time.Time) validate.Condition {
	return func(ctx context.Context) error {
		if birthday.IsZero() {
			return validate.NewValidationErrors("birthday must be provided")
		}

		if birthday.After(time.Now()) {
			return validate.NewValidationErrors("birthday cannot be in the future")
		}

		return nil
	}
}

func Email(email string) validate.Condition {
	return func(ctx context.Context) error {
		if email == "" {
			return validate.NewValidationErrors("email must not be empty")
		}

		// можно заменить на regexp при необходимости
		if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
			return validate.NewValidationErrors("email must be a valid address")
		}

		return nil
	}
}
