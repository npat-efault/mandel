#!/bin/bash
set -e
d=$(dirname "$0")

DATADIR="$d"/data

if [ $# -lt 1 ]; then
    cmd=build
else
    cmd="$1"; shift
fi

case "$cmd" in
    build | install | test)
        mkbundle -v -g -skip="*~" -o="$d"/data.go "$DATADIR"
	go $cmd "$@" "$d"
	;;
    clean)
	go clean "$@" "$d"
	rm -f "$d"/data.go
	;;
    *)
	echo "$0: Nothing to do for $cmd"
	;;
esac
