package common

import (
	"fmt"
)

type TagValue interface {
	ID() uint32
	Label() string
}

// MainTagValue defines an RPM header tag that occurs only in the main Header.
type MainTagValue struct {
	id    uint32
	label string
}

// ID returns the uint32 tag ID.
func (t MainTagValue) ID() uint32 { return t.id }

// Label returns the string representation of the tag.
func (t MainTagValue) Label() string { return t.label }

func (t MainTagValue) IsUnknown() bool {
	_, found := tagsById[t.id]
	return found
}

// tagsById allows O(1) lookups of a MainTagValue given an integer ID read from an RPM.
var tagsById map[uint32]MainTagValue

// A comprehensive registry of known standard RPM tags.
var (
	// --- Base package tags ---
	TagName        = MainTagValue{id: 1000, label: "NAME"}
	TagVersion     = MainTagValue{id: 1001, label: "VERSION"}
	TagRelease     = MainTagValue{id: 1002, label: "RELEASE"}
	TagEpoch       = MainTagValue{id: 1003, label: "EPOCH"}
	TagLicense     = MainTagValue{id: 1014, label: "LICENSE"}
	TagSummary     = MainTagValue{id: 1004, label: "SUMMARY"}
	TagDescription = MainTagValue{id: 1005, label: "DESCRIPTION"}
	TagOs          = MainTagValue{id: 1021, label: "OS"}
	TagArch        = MainTagValue{id: 1022, label: "ARCH"}

	// --- Informative package tags ---
	TagBuildhost          = MainTagValue{id: 1007, label: "BUILDHOST"}
	TagBuildtime          = MainTagValue{id: 1006, label: "BUILDTIME"}
	TagBugurl             = MainTagValue{id: 5012, label: "BUGURL"}
	TagChangelogname      = MainTagValue{id: 1081, label: "CHANGELOGNAME"}
	TagChangelogtext      = MainTagValue{id: 1082, label: "CHANGELOGTEXT"}
	TagChangelogtime      = MainTagValue{id: 1080, label: "CHANGELOGTIME"}
	TagCookie             = MainTagValue{id: 1094, label: "COOKIE"}
	TagDistribution       = MainTagValue{id: 1010, label: "DISTRIBUTION"}
	TagDisttag            = MainTagValue{id: 1155, label: "DISTTAG"}
	TagDisturl            = MainTagValue{id: 1123, label: "DISTURL"}
	TagEncoding           = MainTagValue{id: 5062, label: "ENCODING"}
	TagGroup              = MainTagValue{id: 1016, label: "GROUP"}
	TagModularitylabel    = MainTagValue{id: 5096, label: "MODULARITYLABEL"}
	TagOptflags           = MainTagValue{id: 1122, label: "OPTFLAGS"}
	TagPackager           = MainTagValue{id: 1015, label: "PACKAGER"}
	TagPlatform           = MainTagValue{id: 1132, label: "PLATFORM"}
	TagPolicies           = MainTagValue{id: 1150, label: "POLICIES"}
	TagPolicyflags        = MainTagValue{id: 5033, label: "POLICYFLAGS"}
	TagPolicynames        = MainTagValue{id: 5030, label: "POLICYNAMES"}
	TagPolicytypes        = MainTagValue{id: 5031, label: "POLICYTYPES"}
	TagPolicytypesindexes = MainTagValue{id: 5032, label: "POLICYTYPESINDEXES"}
	TagRpmformat          = MainTagValue{id: 5114, label: "RPMFORMAT"}
	TagRpmversion         = MainTagValue{id: 1064, label: "RPMVERSION"}
	TagSourcepkgid        = MainTagValue{id: 1146, label: "SOURCEPKGID"}
	TagSourcerpm          = MainTagValue{id: 1044, label: "SOURCERPM"}
	TagTranslationurl     = MainTagValue{id: 5100, label: "TRANSLATIONURL"}
	TagUpstreamReleases   = MainTagValue{id: 5101, label: "UPSTREAMRELEASES"}
	TagUrl                = MainTagValue{id: 1020, label: "URL"}
	TagVcs                = MainTagValue{id: 5034, label: "VCS"}
	TagVendor             = MainTagValue{id: 1011, label: "VENDOR"}
	TagSourcenevr         = MainTagValue{id: 5120, label: "SOURCENEVR"}

	// --- Packages with files ---
	TagArchivesize       = MainTagValue{id: 1046, label: "ARCHIVESIZE"}
	TagDirnames          = MainTagValue{id: 1118, label: "DIRNAMES"}
	TagFiledigestalgo    = MainTagValue{id: 5011, label: "FILEDIGESTALGO"}
	TagLongarchivesize   = MainTagValue{id: 271, label: "LONGARCHIVESIZE"}
	TagLongsize          = MainTagValue{id: 5009, label: "LONGSIZE"}
	TagMimedict          = MainTagValue{id: 5116, label: "MIMEDICT"}
	TagPayloadcompressor = MainTagValue{id: 1125, label: "PAYLOADCOMPRESSOR"}
	TagPayloadflags      = MainTagValue{id: 1126, label: "PAYLOADFLAGS"}
	TagPayloadformat     = MainTagValue{id: 1124, label: "PAYLOADFORMAT"}
	TagPrefixes          = MainTagValue{id: 1098, label: "PREFIXES"}
	TagSize              = MainTagValue{id: 1009, label: "SIZE"}

	// --- Per-file information ---
	TagBasenames       = MainTagValue{id: 1117, label: "BASENAMES"}
	TagDirindexes      = MainTagValue{id: 1116, label: "DIRINDEXES"}
	TagFiledevices     = MainTagValue{id: 1095, label: "FILEDEVICES"}
	TagFiledigests     = MainTagValue{id: 1035, label: "FILEDIGESTS"}
	TagFileflags       = MainTagValue{id: 1037, label: "FILEFLAGS"}
	TagFilegroupname   = MainTagValue{id: 1040, label: "FILEGROUPNAME"}
	TagFileinodes      = MainTagValue{id: 1096, label: "FILEINODES"}
	TagFilelangs       = MainTagValue{id: 1097, label: "FILELANGS"}
	TagFilelinktos     = MainTagValue{id: 1036, label: "FILELINKTOS"}
	TagFilemimeindex   = MainTagValue{id: 5115, label: "FILEMIMEINDEX"}
	TagFilemodes       = MainTagValue{id: 1030, label: "FILEMODES"}
	TagFilemtimes      = MainTagValue{id: 1034, label: "FILEMTIMES"}
	TagFilerdevs       = MainTagValue{id: 1033, label: "FILERDEVS"}
	TagFilesizes       = MainTagValue{id: 1028, label: "FILESIZES"}
	TagFileusername    = MainTagValue{id: 1039, label: "FILEUSERNAME"}
	TagFileverifyflags = MainTagValue{id: 1045, label: "FILEVERIFYFLAGS"}
	TagLongfilesizes   = MainTagValue{id: 5008, label: "LONGFILESIZES"}

	// --- Optional file information ---
	TagClassdict           = MainTagValue{id: 1142, label: "CLASSDICT"}
	TagDependsdict         = MainTagValue{id: 1145, label: "DEPENDSDICT"}
	TagFilecaps            = MainTagValue{id: 5010, label: "FILECAPS"}
	TagFileclass           = MainTagValue{id: 1141, label: "FILECLASS"}
	TagFilecolors          = MainTagValue{id: 1140, label: "FILECOLORS"}
	TagFiledependsn        = MainTagValue{id: 1144, label: "FILEDEPENDSN"}
	TagFiledependsx        = MainTagValue{id: 1143, label: "FILEDEPENDSX"}
	TagFilesignaturelength = MainTagValue{id: 5091, label: "FILESIGNATURELENGTH"}
	TagFilesignatures      = MainTagValue{id: 5090, label: "FILESIGNATURES"}
	TagVeritysignaturealgo = MainTagValue{id: 277, label: "VERITYSIGNATUREALGO"}
	TagVeritysignatures    = MainTagValue{id: 276, label: "VERITYSIGNATURES"}

	// --- Hard dependencies ---
	TagProvidename     = MainTagValue{id: 1047, label: "PROVIDENAME"}
	TagProvideversion  = MainTagValue{id: 1113, label: "PROVIDEVERSION"}
	TagProvideflags    = MainTagValue{id: 1112, label: "PROVIDEFLAGS"}
	TagRequirename     = MainTagValue{id: 1049, label: "REQUIRENAME"}
	TagRequireversion  = MainTagValue{id: 1050, label: "REQUIREVERSION"}
	TagRequireflags    = MainTagValue{id: 1048, label: "REQUIREFLAGS"}
	TagConflictname    = MainTagValue{id: 1054, label: "CONFLICTNAME"}
	TagConflictversion = MainTagValue{id: 1055, label: "CONFLICTVERSION"}
	TagConflictflags   = MainTagValue{id: 1053, label: "CONFLICTFLAGS"}
	TagObsoletename    = MainTagValue{id: 1090, label: "OBSOLETENAME"}
	TagObsoleteversion = MainTagValue{id: 1115, label: "OBSOLETEVERSION"}
	TagObsoleteflags   = MainTagValue{id: 1114, label: "OBSOLETEFLAGS"}

	// --- Soft dependencies ---
	TagEnhancename       = MainTagValue{id: 5055, label: "ENHANCENAME"}
	TagEnhanceversion    = MainTagValue{id: 5056, label: "ENHANCEVERSION"}
	TagEnhanceflags      = MainTagValue{id: 5057, label: "ENHANCEFLAGS"}
	TagRecommendname     = MainTagValue{id: 5046, label: "RECOMMENDNAME"}
	TagRecommendversion  = MainTagValue{id: 5047, label: "RECOMMENDVERSION"}
	TagRecommendflags    = MainTagValue{id: 5048, label: "RECOMMENDFLAGS"}
	TagSuggestname       = MainTagValue{id: 5049, label: "SUGGESTNAME"}
	TagSuggestversion    = MainTagValue{id: 5050, label: "SUGGESTVERSION"}
	TagSuggestflags      = MainTagValue{id: 5051, label: "SUGGESTFLAGS"}
	TagSupplementname    = MainTagValue{id: 5052, label: "SUPPLEMENTNAME"}
	TagSupplementversion = MainTagValue{id: 5053, label: "SUPPLEMENTVERSION"}
	TagSupplementflags   = MainTagValue{id: 5054, label: "SUPPLEMENTFLAGS"}
	TagOrdername         = MainTagValue{id: 5035, label: "ORDERNAME"}
	TagOrderversion      = MainTagValue{id: 5036, label: "ORDERVERSION"}
	TagOrderflags        = MainTagValue{id: 5037, label: "ORDERFLAGS"}

	// --- Scriptlets ---
	TagPostin            = MainTagValue{id: 1024, label: "POSTIN"}
	TagPostinflags       = MainTagValue{id: 5021, label: "POSTINFLAGS"}
	TagPostinprog        = MainTagValue{id: 1086, label: "POSTINPROG"}
	TagPosttrans         = MainTagValue{id: 1152, label: "POSTTRANS"}
	TagPosttransflags    = MainTagValue{id: 5025, label: "POSTTRANSFLAGS"}
	TagPosttransprog     = MainTagValue{id: 1154, label: "POSTTRANSPROG"}
	TagPostuntrans       = MainTagValue{id: 5104, label: "POSTUNTRANS"}
	TagPostuntransflags  = MainTagValue{id: 5108, label: "POSTUNTRANSFLAGS"}
	TagPostuntransprog   = MainTagValue{id: 5106, label: "POSTUNTRANSPROG"}
	TagPostun            = MainTagValue{id: 1026, label: "POSTUN"}
	TagPostunflags       = MainTagValue{id: 5023, label: "POSTUNFLAGS"}
	TagPostunprog        = MainTagValue{id: 1088, label: "POSTUNPROG"}
	TagPrein             = MainTagValue{id: 1023, label: "PREIN"}
	TagPreinflags        = MainTagValue{id: 5020, label: "PREINFLAGS"}
	TagPreinprog         = MainTagValue{id: 1085, label: "PREINPROG"}
	TagPretrans          = MainTagValue{id: 1151, label: "PRETRANS"}
	TagPretransflags     = MainTagValue{id: 5024, label: "PRETRANSFLAGS"}
	TagPretransprog      = MainTagValue{id: 1153, label: "PRETRANSPROG"}
	TagPreuntrans        = MainTagValue{id: 5103, label: "PREUNTRANS"}
	TagPreuntransflags   = MainTagValue{id: 5107, label: "PREUNTRANSFLAGS"}
	TagPreuntransprog    = MainTagValue{id: 5105, label: "PREUNTRANSPROG"}
	TagPreun             = MainTagValue{id: 1025, label: "PREUN"}
	TagPreunflags        = MainTagValue{id: 5022, label: "PREUNFLAGS"}
	TagPreunprog         = MainTagValue{id: 1087, label: "PREUNPROG"}
	TagVerifyscript      = MainTagValue{id: 1079, label: "VERIFYSCRIPT"}
	TagVerifyscriptflags = MainTagValue{id: 5026, label: "VERIFYSCRIPTFLAGS"}
	TagVerifyscriptprog  = MainTagValue{id: 1091, label: "VERIFYSCRIPTPROG"}

	// --- Triggers ---
	TagTriggerflags       = MainTagValue{id: 1068, label: "TRIGGERFLAGS"}
	TagTriggerindex       = MainTagValue{id: 1069, label: "TRIGGERINDEX"}
	TagTriggername        = MainTagValue{id: 1066, label: "TRIGGERNAME"}
	TagTriggerscriptflags = MainTagValue{id: 5027, label: "TRIGGERSCRIPTFLAGS"}
	TagTriggerscriptprog  = MainTagValue{id: 1092, label: "TRIGGERSCRIPTPROG"}
	TagTriggerscripts     = MainTagValue{id: 1065, label: "TRIGGERSCRIPTS"}
	TagTriggerversion     = MainTagValue{id: 1067, label: "TRIGGERVERSION"}

	// --- File triggers ---
	TagFiletriggerflags            = MainTagValue{id: 5072, label: "FILETRIGGERFLAGS"}
	TagFiletriggerindex            = MainTagValue{id: 5070, label: "FILETRIGGERINDEX"}
	TagFiletriggername             = MainTagValue{id: 5069, label: "FILETRIGGERNAME"}
	TagFiletriggerpriorities       = MainTagValue{id: 5084, label: "FILETRIGGERPRIORITIES"}
	TagFiletriggerscriptflags      = MainTagValue{id: 5068, label: "FILETRIGGERSCRIPTFLAGS"}
	TagFiletriggerscriptprog       = MainTagValue{id: 5067, label: "FILETRIGGERSCRIPTPROG"}
	TagFiletriggerscripts          = MainTagValue{id: 5066, label: "FILETRIGGERSCRIPTS"}
	TagFiletriggerversion          = MainTagValue{id: 5071, label: "FILETRIGGERVERSION"}
	TagTransfiletriggerflags       = MainTagValue{id: 5082, label: "TRANSFILETRIGGERFLAGS"}
	TagTransfiletriggerindex       = MainTagValue{id: 5080, label: "TRANSFILETRIGGERINDEX"}
	TagTransfiletriggername        = MainTagValue{id: 5079, label: "TRANSFILETRIGGERNAME"}
	TagTransfiletriggerpriorities  = MainTagValue{id: 5085, label: "TRANSFILETRIGGERPRIORITIES"}
	TagTransfiletriggerscriptflags = MainTagValue{id: 5078, label: "TRANSFILETRIGGERSCRIPTFLAGS"}
	TagTransfiletriggerscriptprog  = MainTagValue{id: 5077, label: "TRANSFILETRIGGERSCRIPTPROG"}
	TagTransfiletriggerscripts     = MainTagValue{id: 5076, label: "TRANSFILETRIGGERSCRIPTS"}
	TagTransfiletriggerversion     = MainTagValue{id: 5081, label: "TRANSFILETRIGGERVERSION"}

	// --- Signatures and digests ---
	TagDsaheader        = MainTagValue{id: 267, label: "DSAHEADER"}
	TagLongsigsize      = MainTagValue{id: 270, label: "LONGSIGSIZE"}
	TagOpenpgp          = MainTagValue{id: 278, label: "OPENPGP"}
	TagPayloadsha256    = MainTagValue{id: 5092, label: "PAYLOADSHA256"}
	TagPayloadsha256alt = MainTagValue{id: 5097, label: "PAYLOADSHA256ALT"}
	TagPayloadsha512    = MainTagValue{id: 5121, label: "PAYLOADSHA512"}
	TagPayloadsha512alt = MainTagValue{id: 5122, label: "PAYLOADSHA512ALT"}
	TagPayloadsha3_256  = MainTagValue{id: 5123, label: "PAYLOADSHA3_256"}
	TagPayloadsha3_256a = MainTagValue{id: 5124, label: "PAYLOADSHA3_256ALT"}
	TagRsaheader        = MainTagValue{id: 268, label: "RSAHEADER"}
	TagSha1header       = MainTagValue{id: 269, label: "SHA1HEADER"}
	TagSha256header     = MainTagValue{id: 273, label: "SHA256HEADER"}
	TagSha3_256_header  = MainTagValue{id: 279, label: "SHA3_256_HEADER"}
	TagSiggpg           = MainTagValue{id: 262, label: "SIGGPG"}
	TagSigmd5           = MainTagValue{id: 261, label: "SIGMD5"}
	TagSigpgp           = MainTagValue{id: 259, label: "SIGPGP"}
	TagSigsize          = MainTagValue{id: 257, label: "SIGSIZE"}

	// --- Installed package headers only ---
	TagFilestates         = MainTagValue{id: 1029, label: "FILESTATES"}
	TagInstallcolor       = MainTagValue{id: 1127, label: "INSTALLCOLOR"}
	TagInstalltid         = MainTagValue{id: 1128, label: "INSTALLTID"}
	TagInstalltime        = MainTagValue{id: 1008, label: "INSTALLTIME"}
	TagInstprefixes       = MainTagValue{id: 1099, label: "INSTPREFIXES"}
	TagOrigbasenames      = MainTagValue{id: 1120, label: "ORIGBASENAMES"}
	TagOrigdirindexes     = MainTagValue{id: 1119, label: "ORIGDIRINDEXES"}
	TagOrigdirnames       = MainTagValue{id: 1121, label: "ORIGDIRNAMES"}
	TagPackagedigests     = MainTagValue{id: 5118, label: "PACKAGEDIGESTS"}
	TagPackagedigestalgos = MainTagValue{id: 5119, label: "PACKAGEDIGESTALGOS"}

	// --- Source packages ---
	TagBuildarchs = MainTagValue{id: 1089, label: "BUILDARCHS"}

	// Region & Header Integrity Tags
	TagHeaderImmutable = MainTagValue{id: 63, label: "HEADERIMMUTABLE"}
	TagHeaderI18NTable = MainTagValue{id: 100, label: "HEADERI18NTABLE"}

	// Payload Digest Tags
	TagPayloadSHA256Algo = MainTagValue{id: 5093, label: "PAYLOADSHA256ALGO"} // Obsolete
)

func LookupMainTag(id uint32) MainTagValue {
	result, found := tagsById[id]
	if !found {
		label := fmt.Sprintf("Unknown (%d / 0x%08x)", id, id)
		result = MainTagValue{id: id, label: label}
		tagsById[id] = result
	}
	return result
}

// init builds the tagsById lookup map automatically at runtime startup.
func init() {
	allTags := []MainTagValue{
		TagName, TagVersion, TagRelease, TagEpoch, TagLicense, TagSummary,
		TagDescription, TagOs, TagArch, TagBuildhost, TagBuildtime, TagBugurl,
		TagChangelogname, TagChangelogtext, TagChangelogtime, TagCookie,
		TagDistribution, TagDisttag, TagDisturl, TagEncoding, TagGroup,
		TagModularitylabel, TagOptflags, TagPackager, TagPlatform, TagPolicies,
		TagPolicyflags, TagPolicynames, TagPolicytypes, TagPolicytypesindexes,
		TagRpmformat, TagRpmversion, TagSourcepkgid, TagSourcerpm, TagTranslationurl,
		TagUpstreamReleases, TagUrl, TagVcs, TagVendor, TagSourcenevr, TagArchivesize,
		TagDirnames, TagFiledigestalgo, TagLongarchivesize, TagLongsize, TagMimedict,
		TagPayloadcompressor, TagPayloadflags, TagPayloadformat, TagPrefixes,
		TagSize, TagBasenames, TagDirindexes, TagFiledevices, TagFiledigests,
		TagFileflags, TagFilegroupname, TagFileinodes, TagFilelangs, TagFilelinktos,
		TagFilemimeindex, TagFilemodes, TagFilemtimes, TagFilerdevs, TagFilesizes,
		TagFileusername, TagFileverifyflags, TagLongfilesizes, TagClassdict,
		TagDependsdict, TagFilecaps, TagFileclass, TagFilecolors, TagFiledependsn,
		TagFiledependsx, TagFilesignaturelength, TagFilesignatures,
		TagVeritysignaturealgo, TagVeritysignatures, TagProvidename,
		TagProvideversion, TagProvideflags, TagRequirename, TagRequireversion,
		TagRequireflags, TagConflictname, TagConflictversion, TagConflictflags,
		TagObsoletename, TagObsoleteversion, TagObsoleteflags, TagEnhancename,
		TagEnhanceversion, TagEnhanceflags, TagRecommendname, TagRecommendversion,
		TagRecommendflags, TagSuggestname, TagSuggestversion, TagSuggestflags,
		TagSupplementname, TagSupplementversion, TagSupplementflags, TagOrdername,
		TagOrderversion, TagOrderflags, TagPostin, TagPostinflags, TagPostinprog,
		TagPosttrans, TagPosttransflags, TagPosttransprog, TagPostuntrans,
		TagPostuntransflags, TagPostuntransprog, TagPostun, TagPostunflags,
		TagPostunprog, TagPrein, TagPreinflags, TagPreinprog, TagPretrans,
		TagPretransflags, TagPretransprog, TagPreuntrans, TagPreuntransflags,
		TagPreuntransprog, TagPreun, TagPreunflags, TagPreunprog, TagVerifyscript,
		TagVerifyscriptflags, TagVerifyscriptprog, TagTriggerflags, TagTriggerindex,
		TagTriggername, TagTriggerscriptflags, TagTriggerscriptprog, TagTriggerscripts,
		TagTriggerversion, TagFiletriggerflags, TagFiletriggerindex,
		TagFiletriggername, TagFiletriggerpriorities, TagFiletriggerscriptflags,
		TagFiletriggerscriptprog, TagFiletriggerscripts, TagFiletriggerversion,
		TagTransfiletriggerflags, TagTransfiletriggerindex, TagTransfiletriggername,
		TagTransfiletriggerpriorities, TagTransfiletriggerscriptflags,
		TagTransfiletriggerscriptprog, TagTransfiletriggerscripts,
		TagTransfiletriggerversion, TagDsaheader, TagLongsigsize, TagOpenpgp,
		TagPayloadsha256, TagPayloadsha256alt, TagPayloadsha512, TagPayloadsha512alt,
		TagPayloadsha3_256, TagPayloadsha3_256a, TagRsaheader, TagSha1header,
		TagSha256header, TagSha3_256_header, TagSiggpg, TagSigmd5, TagSigpgp,
		TagSigsize, TagFilestates, TagInstallcolor, TagInstalltid, TagInstalltime,
		TagInstprefixes, TagOrigbasenames, TagOrigdirindexes, TagOrigdirnames,
		TagPackagedigests, TagPackagedigestalgos, TagBuildarchs, TagHeaderImmutable,
		TagHeaderI18NTable, TagPayloadSHA256Algo,
	}

	tagsById = make(map[uint32]MainTagValue, len(allTags))
	for _, t := range allTags {
		tagsById[t.id] = t
	}
}
