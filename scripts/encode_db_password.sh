#!/bin/sh
set -eu

if [ $# -ne 1 ]; then
	echo "Usage: $0 <PASSWORD_FILE>" >&2
	exit 1
fi

PASSWORD_FILE="$1"

password="$(tr -d '\r\n' < "$PASSWORD_FILE")"

LC_ALL=C
export LC_ALL

if [ -z "$password" ]; then
	echo "Error: password in $PASSWORD_FILE is empty" >&2
	exit 1
fi

case "$password" in
	*[!!-~]*)
		echo "Error: password in $PASSWORD_FILE must contain only printable ASCII without spaces" >&2
		exit 1
		;;
esac

rest="$password"
while [ -n "$rest" ]; do
	char="${rest%"${rest#?}"}"
	rest="${rest#?}"
	case "$char" in
		[A-Za-z0-9._~-]) printf '%s' "$char" ;;
		*) printf '%%%02X' "'$char" ;;
	esac
done
