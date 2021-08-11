package route53

const (
	keySigningKeyStatusActionNeeded    = "ACTION_NEEDED"
	keySigningKeyStatusActive          = "ACTIVE"
	keySigningKeyStatusDeleting        = "DELETING"
	keySigningKeyStatusInactive        = "INACTIVE"
	keySigningKeyStatusInternalFailure = "INTERNAL_FAILURE"

	serveSignatureActionNeeded    = "ACTION_NEEDED"
	serveSignatureDeleting        = "DELETING"
	serveSignatureInternalFailure = "INTERNAL_FAILURE"
	serveSignatureNotSigning      = "NOT_SIGNING"
	serveSignatureSigning         = "SIGNING"
)
