package common

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func randomFile(t *testing.T) (string, error) {
	dir := t.TempDir()

	if fileName, err := RandomString(20, "abcdefghijklmnopqrstuvwxyz"); err == nil {
		return filepath.Join(dir, fileName), nil
	} else {
		return "", err
	}
}

func parseRPM(fileName string, validatePayloadMd5 bool) (RpmFile, ReaderWithOffset, error) {
	var err error
	var file *os.File
	if file, err = os.Open(fileName); err != nil {
		return RpmFile{}, nil, fmt.Errorf("failed to open %v - %w", fileName, err)
	}
	// NOTE: File must be closed by caller

	inputReader := NewIOReaderWithOffset(file)
	// NOTE: ParseRPMFromReader WITHOUT validation (validatePayloadMd5 == false)
	//       will read up to the beginning of the RPM payload so the returned
	//       ReaderWithOffset is in fact the payload reader
	rpmFile, err := ParseRPMFromReader(inputReader, validatePayloadMd5)
	if err != nil {
		_ = file.Close()
		return RpmFile{}, nil, err
	}
	return rpmFile, inputReader, nil
}

func TestSignRpm(t *testing.T) {

	const fileName = "testdata/my-payload-package-1.0-1.noarch.rpm"
	const privKey = "testdata/my-signing-key.priv.gpg.asc"
	var err error

	outputName, err := randomFile(t)
	if err != nil {
		t.Fatal(err)
	}

	var rpmFile RpmFile
	// NOTE: It's important to NOT perform payload validation here because otherwise
	//       the returned reader will not be at the start of the payload but at the end of the file
	rpmFile, payloadReader, err := parseRPM(fileName, false)
	if err != nil {
		t.Fatalf("Failed to parse %v - %v", fileName, err)
	}

	privateKey, err := LoadPrivateKeyFromFile(privKey, []byte{})
	if err != nil {
		t.Fatalf("Failed to load private key - %v", err)
	}
	if err = SignRpm(&rpmFile, privateKey, true); err != nil {
		CloseQuietly(payloadReader)
		t.Fatalf("Failed to sign %v - %v", fileName, err)
	}

	err = rpmFile.WriteToFile(outputName, payloadReader)
	CloseQuietly(payloadReader)
	if err != nil {
		t.Fatalf("Failed open output file %v - %v", outputName, err)
	}

	// verify by parsing signed output again
	_, payloadReader, err = parseRPM(outputName, true)
	if err != nil {
		t.Fatalf("round-trip failed - %v", err)
	}
	CloseQuietly(payloadReader)
}
