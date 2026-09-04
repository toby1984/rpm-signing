package common

import "fmt"

// SignatureTagValue defines an RPM header tag that occurs only in the Signature Header.
type SignatureTagValue struct {
	id    uint32
	label string
}

// ID returns the uint32 tag ID.
func (t SignatureTagValue) ID() uint32 { return t.id }

// Label returns the string representation of the tag.
func (t SignatureTagValue) Label() string { return t.label }

// signatureTagsByID allows O(1) lookups of a SignatureTagValue given an integer ID
// read from an RPM's Signature Header.
var signatureTagsByID map[uint32]SignatureTagValue

var (
	// --- Legacy V3 Signature Tags (1000 - 1008) ---
	// These are the tags that share IDs with Main Header tags (like NAME, VERSION)
	// but represent completely different data types and values in the Signature Header.
	SigTagSize          = SignatureTagValue{id: 1000, label: "SIG_SIZE"}          // Type: INT32 (Collides with NAME)
	SigTagLEMD5_1       = SignatureTagValue{id: 1001, label: "SIG_LEMD5_1"}       // Obsolete  (Collides with VERSION)
	SigTagPGP           = SignatureTagValue{id: 1002, label: "SIG_PGP"}           // Obsolete  (Collides with RELEASE)
	SigTagLEMD5_2       = SignatureTagValue{id: 1003, label: "SIG_LEMD5_2"}       // Obsolete  (Collides with EPOCH)
	SigTagMD5           = SignatureTagValue{id: 1004, label: "SIG_MD5"}           // Type: BIN (Collides with SUMMARY)
	SigTagGPG           = SignatureTagValue{id: 1005, label: "SIG_GPG"}           // Type: BIN (Collides with DESCRIPTION)
	SigTagPGP5          = SignatureTagValue{id: 1006, label: "SIG_PGP5"}          // Obsolete  (Collides with BUILDTIME)
	SigTagPayloadSize   = SignatureTagValue{id: 1007, label: "SIG_PAYLOADSIZE"}   // UNCOMPRESSED payload size, Type: INT32 (Collides with BUILDHOST)
	SigTagReservedSpace = SignatureTagValue{id: 1008, label: "SIG_RESERVEDSPACE"} // Type: BIN (Collides with INSTALLTIME)

	// --- V4+ Signature & Digest Tags (256+) ---
	// These are universally reserved for cryptographic signatures and digests.
	SigTagHeaderSignatures = SignatureTagValue{id: 62, label: "HEADERSIGNATURES"}
	SigTagSigSize          = SignatureTagValue{id: 257, label: "SIGSIZE"}
	SigTagSigPGP           = SignatureTagValue{id: 259, label: "SIGPGP"}
	SigTagSigMD5           = SignatureTagValue{id: 261, label: "SIGMD5"}
	SigTagSigGPG           = SignatureTagValue{id: 262, label: "SIGGPG"}
	SigTagDSAHeader        = SignatureTagValue{id: 267, label: "DSAHEADER"}
	SigTagRSAHeader        = SignatureTagValue{id: 268, label: "RSAHEADER"}
	SigTagSHA1Header       = SignatureTagValue{id: 269, label: "SHA1HEADER"}
	SigTagLongSigSize      = SignatureTagValue{id: 270, label: "LONGSIGSIZE"}
	SigTagLongArchiveSize  = SignatureTagValue{id: 271, label: "LONGARCHIVESIZE"}
	SigTagSHA256Header     = SignatureTagValue{id: 273, label: "SHA256HEADER"}
	SigTagFileSignTag      = SignatureTagValue{id: 274, label: "FILESIGNTAG"}
	SigTagPayloadDigest    = SignatureTagValue{id: 275, label: "PAYLOADDIGEST"}
	SigTagVeritySignatures = SignatureTagValue{id: 276, label: "VERITYSIGNATURES"}
	SigTagVeritySigAlgo    = SignatureTagValue{id: 277, label: "VERITYSIGNATUREALGO"}
	SigTagOpenPGP          = SignatureTagValue{id: 278, label: "OPENPGP"}
	SigTagSHA3_256Header   = SignatureTagValue{id: 279, label: "SHA3_256_HEADER"}
)

func LookupSignatureTag(id uint32) SignatureTagValue {
	result, found := signatureTagsByID[id]
	if !found {
		label := fmt.Sprintf("Unknown (%d / 0x%08x)", id, id)
		result = SignatureTagValue{id: id, label: label}
		signatureTagsByID[id] = result
	}
	return result
}

// init builds the signatureTagsByID lookup map automatically at runtime startup.
func init() {
	allTags := []SignatureTagValue{
		SigTagSize, SigTagLEMD5_1, SigTagPGP, SigTagLEMD5_2,
		SigTagMD5, SigTagGPG, SigTagPGP5, SigTagPayloadSize, SigTagReservedSpace,
		SigTagHeaderSignatures, SigTagSigSize, SigTagSigPGP, SigTagSigMD5,
		SigTagSigGPG, SigTagDSAHeader, SigTagRSAHeader, SigTagSHA1Header,
		SigTagLongSigSize, SigTagLongArchiveSize, SigTagSHA256Header,
		SigTagVeritySignatures, SigTagVeritySigAlgo, SigTagOpenPGP, SigTagSHA3_256Header,
		SigTagFileSignTag, SigTagPayloadDigest,
	}

	signatureTagsByID = make(map[uint32]SignatureTagValue, len(allTags))
	for _, t := range allTags {
		signatureTagsByID[t.id] = t
	}
}
