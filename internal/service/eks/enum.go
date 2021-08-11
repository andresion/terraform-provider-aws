package eks

const (
	identityProviderConfigTypeOIDC = "oidc"
)

const (
	resourcesSecrets = "secrets"
)

func resources_Values() []string {
	return []string{
		resourcesSecrets,
	}
}
