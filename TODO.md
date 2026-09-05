# Stuff that might (or might not) be useful

## Higher priority
- client & server: Add support for GZIP compression to both server and client (to cut down on transmission times and data volume)
- client & server: Add support for signing arbitrary data and returning an ASCII-armored, detached GPG signature to both client and server
- server: Change configuration file syntax to something hierarchical (YAML) and add support for having multiple signing keys
  by defining a mapping from an auth token to a specific GPG private key file

