package sdk

// BasicDefaultVariant builds the DefaultVariant a feature flag serves when no targeting
// rule matches, naming one of the flag's own variants.
//
// DefaultVariant is a union over the rollout strategies, so it cannot be written as a
// struct literal from outside this package -- this is the "basic" arm of it, the only one
// a fixed fallback needs.
func BasicDefaultVariant(variant string) (DefaultVariant, error) {
	var defaultVariant DefaultVariant
	err := defaultVariant.FromBasic(Basic{Type: basicStrategy, Value: variant})

	return defaultVariant, err
}

// BasicVariantTargeting builds a targeting rule that serves variant to every evaluation
// context the CEL expression rule matches, e.g.
//
//	sdk.BasicVariantTargeting("Demo instances", "__kaiten.instance.metadata.demo == true", "on")
//
// Rules are evaluated in order and the first match wins; contexts matching none of them
// fall through to the flag's default variant.
func BasicVariantTargeting(name, rule, variant string) (TargetingsItem, error) {
	var targeting TargetingsItem
	err := targeting.FromBasicTargeting(BasicTargeting{
		Name:    name,
		Rule:    rule,
		Type:    basicTargetingStrategy,
		Variant: variant,
	})

	return targeting, err
}
