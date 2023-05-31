#!/usr/bin/env bash

set -e

adduser -h /opt/low -s /bin/bash -D low

echo "low:password" | chpasswd