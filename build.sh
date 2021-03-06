#!/bin/sh

VERSION=1.0

LOCAL_BUILD="true"

UPDATE="false"

if [ $# -gt 0 ]; then
	if [ $1 = "release" ]
	then
		LOCAL_BUILD=-1
	fi
	if [ $1 = "update" ]
	then
		UPDATE="true"
	fi
fi

if [ $LOCAL_BUILD = "true" ]
then
	VERSION=$VERSION:$(git rev-parse --short HEAD)
	echo "Building local version"
else
	echo "Building release version"
fi

echo "package util\n\nconst VERSION = \"$VERSION\"\n" > util/VERSION.go

if [ $UPDATE = "true" ]
then
	echo "Updating dependencies"
	go get
else
	go build
fi

rm util/VERSION.go

