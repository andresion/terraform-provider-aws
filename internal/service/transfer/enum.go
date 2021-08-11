package transfer

const (
	securityPolicyName2018_11      = "TransferSecurityPolicy-2018-11"
	securityPolicyName2020_06      = "TransferSecurityPolicy-2020-06"
	securityPolicyNameFIPS_2020_06 = "TransferSecurityPolicy-FIPS-2020-06"
)

func securityPolicyName_Values() []string {
	return []string{
		securityPolicyName2018_11,
		securityPolicyName2020_06,
		securityPolicyNameFIPS_2020_06,
	}
}
