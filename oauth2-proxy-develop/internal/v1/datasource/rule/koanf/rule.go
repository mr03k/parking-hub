package koanf

import (
	"application/internal/v1/datasource/rule"
	rule2 "application/internal/v1/entity/rule"
	"github.com/knadh/koanf/v2"
)

type RuleDS struct {
	k *koanf.Koanf
}

func NewRuleDS(k *koanf.Koanf) rule.Rule {
	return &RuleDS{k: k}
}

func (r *RuleDS) GetAll() ([]*rule2.Rule, error) {
	var rules []*rule2.Rule
	if err := r.k.Unmarshal("rules", &rules); err != nil {
		return nil, err
	}
	return rules, nil
}
