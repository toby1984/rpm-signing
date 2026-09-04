package common

import (
	"io"
	"os"
	"testing"
)

func assertEqualsToFile(actualData []byte, fileName string, t *testing.T) {
	var file *os.File
	var err error
	if file, err = os.Open(fileName); err != nil {
		t.Fatalf("Error opening file %s", fileName)
	}
	var expectedData []byte
	if expectedData, err = io.ReadAll(file); err != nil {
		t.Fatalf("Error reading file %s", fileName)
	}
	if len(expectedData) != len(actualData) {
		t.Fatalf("Data length mismatch, expected %d, got %d", len(expectedData), len(actualData))
	}
	for i := range actualData {
		if actualData[i] != expectedData[i] {
			t.Fatalf("Data mismatch at offset 0x%x", i)
		}
	}
}

func TestRoundTripNoRewriting(t *testing.T) {

	RootLogger().SetCurrentLogLevel(LOG_LEVEL_TRACE)

	testData := "testdata/my-payload-package-1.0-1.noarch.rpm"
	file, err := os.Open(testData)
	if err != nil {
		t.Fatalf("Failed to open file %s", testData)
		return
	}
	defer func() { _ = file.Close() }()

	reader := NewIOReaderWithOffset(file)
	var rpm RpmFile
	if rpm, err = ParseRPMFromReader(reader, false); err != nil {
		t.Fatalf("Failed to parse RPM from %s", testData)
	}

	writer := NewArrayWriter()
	if err := rpm.Write(&writer, reader); err != nil {
		t.Fatalf("Failed to write RPM to %s", testData)
	}
	//if err := os.WriteFile("/tmp/test.out", writer.dst, 0644); err != nil {
	//	t.Fatalf("Failed to write RPM to /tmp/test.out")
	//}
	assertEqualsToFile(writer.Data, testData, t)
}

// assertHeaderConsistent checks the invariants that have to hold after index entries
// of a header got added, removed or replaced: the entry counters agree with the actual
// entries, the signed region trailer covers all of them, no two entries claim the same
// data store bytes and the declared data store size matches what Write() would emit.
//
// It also checks the invariants rpm itself enforces in hdrblobVerifyInfo(): the index
// entries are sorted by tag, their data is laid out in that same order without going
// backwards, every entry is aligned as its format type requires and the region trailer
// occupies the last bytes of the data store.
func assertHeaderConsistent(header *RpmHeader, label string, t *testing.T) {

	if int(header.indexEntryCount) != len(header.IndexEntries) {
		t.Errorf("%s: header claims %d index entries but has %d", label, header.indexEntryCount, len(header.IndexEntries))
	}
	if header.SignedRegion != nil && int(header.SignedRegion.IndexEntryCount) != len(header.IndexEntries) {
		t.Errorf("%s: signed region covers %d index entries but header has %d", label,
			header.SignedRegion.IndexEntryCount, len(header.IndexEntries))
	}

	// NewDynamicArray() performs the same overlap detection Write() relies on
	dataStore := NewDynamicArray()
	for idx := range header.IndexEntries {
		entry := &header.IndexEntries[idx]
		if err := dataStore.Write(entry.payload, entry.offset); err != nil {
			t.Errorf("%s: entry with tag %d at offset %d (%d bytes): %v", label, entry.tag, entry.offset, len(entry.payload), err)
		}
	}
	if actualSize := uint32(len(dataStore.GetData())); actualSize != header.storeSizeInBytes {
		t.Errorf("%s: data store holds %d bytes but header declares %d", label, actualSize, header.storeSizeInBytes)
	}

	trailer := header.regionTrailerEntry()
	previousTag := uint32(0)
	previousEnd := uint32(0)
	for idx := range header.IndexEntries {
		entry := &header.IndexEntries[idx]

		if idx > 0 && entry.tag <= previousTag {
			t.Errorf("%s: index entry %d has tag %d but the entry before it has tag %d, entries must be sorted by tag",
				label, idx, entry.tag, previousTag)
		}
		previousTag = entry.tag

		alignment := uint32(FormatType(entry.formatType).GetRequiredAlignment())
		if entry.offset%alignment != 0 {
			t.Errorf("%s: entry with tag %d is at offset %d, which is not %d byte aligned",
				label, entry.tag, entry.offset, alignment)
		}

		if entry == trailer {
			// the trailer is the one entry whose data is not in index order,
			// it has to cover the very end of the data store instead
			if entry.endOffset() != header.storeSizeInBytes {
				t.Errorf("%s: region trailer ends at %d but the data store is %d bytes",
					label, entry.endOffset(), header.storeSizeInBytes)
			}
			continue
		}
		if entry.offset < previousEnd {
			t.Errorf("%s: entry with tag %d starts at offset %d, before the previous entry's data ends at %d",
				label, entry.tag, entry.offset, previousEnd)
		}
		previousEnd = entry.endOffset()
	}
}

// assertRoundTrips writes the (modified) RPM and parses it back so the written bytes
// are checked by the parser's own validation, not just by the in-memory invariants.
func assertRoundTrips(rpmFile RpmFile, payloadReader ReaderWithOffset, label string, t *testing.T) RpmFile {

	writer := NewArrayWriter()
	if err := rpmFile.Write(&writer, payloadReader); err != nil {
		t.Fatalf("%s: failed to write RPM - %v", label, err)
	}

	reader := NewArrayReader(writer.Data)
	defer CloseQuietly(reader)

	reparsed, err := ParseRPMFromReader(reader, false)
	if err != nil {
		t.Fatalf("%s: failed to parse written RPM back - %v", label, err)
	}
	return reparsed
}

// TestSetOrAddIndexEntryReplaceSameSize covers replacing an existing index entry with a
// payload of the same size, which must reuse the space the previous payload occupied.
func TestSetOrAddIndexEntryReplaceSameSize(t *testing.T) {

	const testData = "testdata/my-payload-package-1.0-1.noarch.rpm"

	rpmFile, payloadReader, err := parseRPM(testData, false)
	if err != nil {
		t.Fatalf("Failed to parse %s - %v", testData, err)
	}
	defer CloseQuietly(payloadReader)

	header := &rpmFile.SignatureHeader

	header.SetOrAddIndexEntry(SigTagRSAHeader, FORMAT_TYPE_BINARY, 620, make([]byte, 620))
	added := header.FindIndexEntry(SigTagRSAHeader)
	if added == nil {
		t.Fatal("Entry was not added")
	}
	offsetAfterAdd := added.offset
	entryCountAfterAdd := len(header.IndexEntries)

	// replace with a payload of identical size but different content
	replacement := make([]byte, 620)
	for i := range replacement {
		replacement[i] = 0xab
	}
	header.SetOrAddIndexEntry(SigTagRSAHeader, FORMAT_TYPE_BINARY, 620, replacement)

	replaced := header.FindIndexEntry(SigTagRSAHeader)
	if replaced.offset != offsetAfterAdd {
		t.Errorf("Replacing with a same-sized payload moved the entry from offset %d to %d", offsetAfterAdd, replaced.offset)
	}
	if len(header.IndexEntries) != entryCountAfterAdd {
		t.Errorf("Replacing added an index entry, expected %d entries, got %d", entryCountAfterAdd, len(header.IndexEntries))
	}
	if replaced.payload[0] != 0xab {
		t.Error("Payload was not replaced")
	}
	assertHeaderConsistent(header, "replace same size", t)
	assertRoundTrips(rpmFile, payloadReader, "replace same size", t)
}

// TestSetOrAddIndexEntryReplaceLargerPayload covers replacing an index entry that is NOT
// at the end of the data store with a larger payload: the entry cannot grow in place
// without overwriting the entries behind it, so everything behind it has to move up.
func TestSetOrAddIndexEntryReplaceLargerPayload(t *testing.T) {

	const testData = "testdata/my-payload-package-1.0-1.noarch.rpm"

	rpmFile, payloadReader, err := parseRPM(testData, false)
	if err != nil {
		t.Fatalf("Failed to parse %s - %v", testData, err)
	}
	defer CloseQuietly(payloadReader)

	header := &rpmFile.SignatureHeader

	header.SetOrAddIndexEntry(SigTagRSAHeader, FORMAT_TYPE_BINARY, 620, make([]byte, 620))
	before := *header.FindIndexEntry(SigTagRSAHeader)
	if header.isAtEndOfDataStore(before) {
		t.Fatalf("Test requires an entry that is not at the end of the data store, got one at offset %d", before.offset)
	}
	// the entry the grown payload must not be allowed to overwrite
	followingTag := SigTagSHA1Header
	followingBefore := *header.FindIndexEntry(followingTag)
	if followingBefore.offset < before.offset {
		t.Fatalf("Test requires an entry behind the one being replaced, tag %d is at offset %d", followingTag.ID(), followingBefore.offset)
	}

	header.SetOrAddIndexEntry(SigTagRSAHeader, FORMAT_TYPE_BINARY, 2000, make([]byte, 2000))

	after := *header.FindIndexEntry(SigTagRSAHeader)
	if len(after.payload) != 2000 {
		t.Errorf("Payload was not replaced, expected 2000 bytes, got %d", len(after.payload))
	}
	followingAfter := *header.FindIndexEntry(followingTag)
	if followingAfter.offset < after.endOffset() {
		t.Errorf("Entry behind the replaced one stayed at offset %d and is overwritten by the payload ending at %d",
			followingAfter.offset, after.endOffset())
	}
	assertHeaderConsistent(header, "replace larger payload", t)
	assertRoundTrips(rpmFile, payloadReader, "replace larger payload", t)
}

// TestAddIndexEntryWithoutReservedSpace covers adding an index entry to a header that has
// no reserved space entry to allocate from, so one has to be created first.
func TestAddIndexEntryWithoutReservedSpace(t *testing.T) {

	const testData = "testdata/my-payload-package-1.0-1.noarch.rpm"

	rpmFile, payloadReader, err := parseRPM(testData, false)
	if err != nil {
		t.Fatalf("Failed to parse %s - %v", testData, err)
	}
	defer CloseQuietly(payloadReader)

	header := &rpmFile.SignatureHeader

	header.DeleteIndexEntryIfExists(SigTagReservedSpace)
	if header.FindIndexEntry(SigTagReservedSpace) != nil {
		t.Fatal("Reserved space entry was not deleted")
	}
	assertHeaderConsistent(header, "reserved space deleted", t)

	header.SetOrAddIndexEntry(SigTagRSAHeader, FORMAT_TYPE_BINARY, 620, make([]byte, 620))

	added := header.FindIndexEntry(SigTagRSAHeader)
	if added == nil {
		t.Fatal("Entry was not added")
	}
	reservedSpace := header.FindIndexEntry(SigTagReservedSpace)
	if reservedSpace == nil {
		t.Fatal("No reserved space entry was created")
	}
	if reservedSpace.offset == 0 {
		t.Error("Newly created reserved space entry was left at offset 0")
	}
	assertHeaderConsistent(header, "reserved space recreated", t)

	reparsed := assertRoundTrips(rpmFile, payloadReader, "reserved space recreated", t)
	if reparsed.SignatureHeader.FindIndexEntry(SigTagRSAHeader) == nil {
		t.Error("Added entry did not survive the round trip")
	}
}
