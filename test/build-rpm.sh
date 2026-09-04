#!/bin/bash

set -e -o pipefail

if [ "$1" == "-h" -o "$1" == "--help" ] ; then
  echo "Usage: $0 [sign]"
  exit 1
fi
echo "Building RPM..."
rpmbuild -bb \
        --quiet \
        --define "_sourcedir $(pwd)" \
        --define "_rpmdir $(pwd)" \
        --define "_build_name_fmt %%{NAME}-%%{VERSION}-%%{RELEASE}.%%{ARCH}.rpm" \
        payload.spec

if [ "$1" == "sign" ] ; then
  KEYNAME="Test Signing Key"
  export PRIVATE_KEYFILE="my-signing-key.priv.gpg.asc"
  export PUBLIC_KEYFILE="my-signing-key.pub.gpg.asc"

  if [ ! -e ${PRIVATE_KEYFILE} ] ; then
    echo "Generating new GPG key"
    gpg --batch --generate-key <<EOF
Key-Type: RSA
Key-Length: 4096
Name-Real: ${KEYNAME}
Name-Email: key@example.com
Expire-Date: 0
%no-protection
%commit
EOF
    
    echo "Exporting public key using ASCII armor"
    gpg --export --armor "${KEYNAME}" > ${PUBLIC_KEYFILE}

    echo "Exporting private key using ASCII armor"
    gpg --armor --export-secret-keys "${KEYNAME}" > ${PRIVATE_KEYFILE}
    
    KEY_FPR=$(gpg --with-colons --show-keys ${PUBLIC_KEYFILE} | awk -F: '$1=="fpr"{print $10; exit}' | tr '[:upper:]' '[:lower:]')

    if ! rpm -q "gpg-pubkey-${KEY_FPR: -8}" &>/dev/null; then
      echo "Importing GPG key into RPM key database"
      sudo rpm --import ${PUBLIC_KEYFILE}
    fi
  fi
  

  # find all RPM files and assign them to array
  shopt -s nullglob
  rpm_files=( *.rpm )
  shopt -u nullglob

  if [[ ${#rpm_files[@]} -ne 1 ]]; then
    echo "Error: Expected exactly one RPM file to sign but found ${#rpm_files[@]}." >&2
    exit 1
  fi

  RPM="${rpm_files[0]}"
  echo "Signing ${RPM}"
  SIGNED_RPM="${RPM}.signed"
  cp ${RPM} ${SIGNED_RPM}
  rpmsign --define "__gpg /usr/bin/gpg" --define "_gpg_name ${KEYNAME}" --addsign ${SIGNED_RPM} 
fi
