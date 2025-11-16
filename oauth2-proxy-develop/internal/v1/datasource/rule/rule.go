package rule

import (
	"application/internal/v1/entity/rule"
)

type Rule interface {
	GetAll() ([]*rule.Rule, error)
}
