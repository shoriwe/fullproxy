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
CC=$WIN64CC CXX=$WIN64CXX  make windows-64-static
echo "Windows - 64 _ Dynamic"
CC=$WIN64CC CXX=$WIN64CXX make windows-64-dynamic
echo "Windows - 32 _ Static"
CC=$WIN32CC CXX=$WIN32CXX make windows-32-static
echo "Windows - 32 _ Dynamic"
CC=$WIN32CC CXX=$WIN32CXX make windows-32-dynamic

echo "Linux - 64 _ Static"
CC=$LIN64CC CXX=$LIN64CXX make linux-64-static
echo "Linux - 64 _ Dynamic"
CC=$LIN64CC CXX=$LIN64CXX make linux-64-dynamic
echo "Linux - 32 _ Static"
CC=$LIN32CC CXX=$LIN32CXX make linux-32-static
echo "Linux - 32 _ Dynamic"
CC=$LIN32CC CXX=$LIN32CXX make linux-32-dynamic
