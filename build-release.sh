#!/usr/bin/env bash
if [ -z "$WIN64CC" ]
then
  echo "Windows 64 c code compiler not found; See: WIN64CC environment variable not found"
  exit 1
fi
if [ -z "$WIN64CXX" ]
then
  echo "Windows 64 c++ code compiler not found; See: WIN64CXX environment variable not found"
  exit 1
fi

if [ -z "$WIN32CC" ]
then
  echo "Windows 32 c code compiler not found; See: WIN32CC environment variable not found"
  exit 1
fi
if [ -z "$WIN64CXX" ]
then
  echo "Windows 32 c++ code compiler not found; See: WIN32CXX environment variable not found"
  exit 1
fi

if [ -z "$LIN64CC" ]
then
  echo "Linux 64 c code compiler not found; See: LIN64CC environment variable not found"
  exit 1
fi
if [ -z "$LIN64CXX" ]
then
  echo "Linux 64 c++ code compiler not found; See: LIN64CXX environment variable not found"
  exit 1
fi

if [ -z "$LIN32CC" ]
then
  echo "Linux 32 c code compiler not found; See: LIN32CC environment variable not found"
  exit 1
fi
if [ -z "$LIN32CXX" ]
then
  echo "Linux 32 c++ code compiler not found; See: LIN32CXX environment variable not found"
  exit 1
fi

echo "Windows - 64 _ Static"
make windows-64-static CC="$WIN64CC" CXX="$WIN64CXX"
echo "Windows - 64 _ Dynamic"
make windows-64-dynamic CC="$WIN64CC" CXX="$WIN64CXX"
echo "Windows - 32 _ Static"
make windows-32-static CC="$WIN32CC" CXX="$WIN32CXX"
echo "Windows - 32 _ Dynamic"
make windows-32-dynamic CC="$WIN32CC" CXX="$WIN32CXX"

echo "Linux - 64 _ Static"
# shellcheck disable=SC2086
make linux-64-static CC="$LIN64CC" CXX=$LIN64CXX
echo "Linux - 64 _ Dynamic"
# shellcheck disable=SC2086
make linux-64-dynamic CC="$LIN64CC" CXX=$LIN64CXX
echo "Linux - 32 _ Static"
# shellcheck disable=SC2086
make linux-32-static CC="$LIN32CC" CXX=$LIN32CXX
echo "Linux - 32 _ Dynamic"
# shellcheck disable=SC2086
make linux-32-dynamic CC=$LIN32CC CXX="$LIN32CXX"
