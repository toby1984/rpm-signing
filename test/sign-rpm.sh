#!/bin/bash
set -e -o pipefail

if [ "$#" != "1" ] ; then
  echo "Usage: $0 <RPM file>"
  exit 1
fi

RPM_FILE="$1"

KEY_NAME="Dummy test key"
KEY_EMAIL="test@test.de"

PRIV_KEY="my-signing-key.priv.gpg.asc"
PUB_KEY="my-signing-key.pub.gpg.asc"

GPG_BINARY=`which gpg`
if [ -z "$GPG_BINARY" ] ; then
  echo "Failed to find GPG binary"
  exit 1
fi
echo "GPG is at $GPG_BINARY"

# 1. Create a temporary GPG directory to keep everything isolated
export GNUPGHOME=$(mktemp -d)

gpg --batch --passphrase '' --quick-generate-key "$KEY_NAME <$KEY_EMAIL>" rsa4096 sign 0

gpg --armor --export "$KEY_EMAIL" > ${PUB_KEY}
gpg --armor --export-secret-keys "$KEY_EMAIL" > ${PRIV_KEY}

# Sign the RPM using the temporary keyring

SIGNED_RPM_FILE="$RPM_FILE.signed"
cp $RPM_FILE $RPM_FILE.signed

rpmsign --define "_gpg_name $KEY_EMAIL" \
        --define "_gpg_path $GNUPGHOME" \
        --define "%__gpg $GPG_BINARY" \
        --define "%__gpg_sign_cmd %{__gpg} gpg --batch --no-armor --no-secmem-warning -u \"%{_gpg_name}\" --s2k-mode 3 --digest-algo sha256 --detach-sign --output %{__signature_filename} %{__plaintext_filename}" \
        --addsign "$SIGNED_RPM_FILE"

# Clean up the temporary GPG directory
rm -rf "$GNUPGHOME"

echo "Successfully signed $SIGNED_RPM_FILE"
