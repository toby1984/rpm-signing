#!/bin/bash
set -e -o pipefail
make 
bin/client adcli-selinux-0.9.3.1-3.el10.noarch.rpm
