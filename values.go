package sdk

// NumberUsageValue builds an EntitlementUsageValue for a numeric entitlement value.
func NumberUsageValue(value float64) (EntitlementUsageValue, error) {
	var usageValue EntitlementUsageValue
	err := usageValue.FromNumberEntitlementValue(NumberEntitlementValue{
		Type:  Number,
		Value: value,
	})
	return usageValue, err
}

// BooleanUsageValue builds an EntitlementUsageValue for a boolean entitlement value.
func BooleanUsageValue(value bool) (EntitlementUsageValue, error) {
	var usageValue EntitlementUsageValue
	err := usageValue.FromBooleanEntitlementValue(BooleanEntitlementValue{
		Type:  Boolean,
		Value: value,
	})
	return usageValue, err
}

// ConfigUsageValue builds an EntitlementUsageValue for a config (object) entitlement value.
func ConfigUsageValue(value map[string]any) (EntitlementUsageValue, error) {
	var usageValue EntitlementUsageValue
	err := usageValue.FromConfigEntitlementValue(ConfigEntitlementValue{
		Type:  Object,
		Value: cloneMap(value),
	})
	return usageValue, err
}

// NumberLicenseValue builds a LicenseEntitlementValue for a numeric entitlement value.
func NumberLicenseValue(value float64) (LicenseEntitlementValue, error) {
	var licenseValue LicenseEntitlementValue
	err := licenseValue.FromNumberEntitlementValue(NumberEntitlementValue{
		Type:  Number,
		Value: value,
	})
	return licenseValue, err
}

// BooleanLicenseValue builds a LicenseEntitlementValue for a boolean entitlement value.
func BooleanLicenseValue(value bool) (LicenseEntitlementValue, error) {
	var licenseValue LicenseEntitlementValue
	err := licenseValue.FromBooleanEntitlementValue(BooleanEntitlementValue{
		Type:  Boolean,
		Value: value,
	})
	return licenseValue, err
}

// ConfigLicenseValue builds a LicenseEntitlementValue for a config (object) entitlement value.
func ConfigLicenseValue(value map[string]any) (LicenseEntitlementValue, error) {
	var licenseValue LicenseEntitlementValue
	err := licenseValue.FromConfigEntitlementValue(ConfigEntitlementValue{
		Type:  Object,
		Value: cloneMap(value),
	})
	return licenseValue, err
}
