#!/usr/bin/env bash
set -euo pipefail
#
# sumeru_custom_addons — run from this directory. Loads defaults, then optional
# INI under SUMERU_CONF_LOCATION (default: ./conf), then env overrides.
# Builds make + go run commands and executes them.
#
# Environment:
#   SUMERU_CONF_LOCATION  Directory holding sumeru.conf (default: conf)
#   SUMERU_CONF / CONF    Explicit INI path (overrides location search)
#   SUMERU_ROOT           Standard sumeru tree for import-gen (default: from INI
#                         key sumeru_home if set, else ../sumeru)
#   OUT                   Optional; passed to make generate when set
#   SUMERU_INSTALL        Comma-separated modules for -i (only if argv has no -i)
#   SUMERU_UPDATE         Comma-separated modules for -u, or "all" (only if argv has no -u)
#
# CLI: any args are forwarded to the app after make generate. If argv omits -c,
#      -c <resolved-ini> is prepended (so -i / -u / -p-only invocations work).
#
# Examples:
#   ./sumeru-workspace.sh
#   ./sumeru-workspace.sh -i my_module --stop-after-init
#   ./sumeru-workspace.sh -u sales,student
#   SUMERU_INSTALL=my_module ./sumeru-workspace.sh --stop-after-init
#   SUMERU_UPDATE=all ./sumeru-workspace.sh --stop-after-init
#

readonly DEFAULT_SUMERU_ROOT="../sumeru"
readonly DEFAULT_INI_BASENAME="sumeru.conf"

workspace_root=$(cd "$(dirname "$0")" && pwd)
cd "$workspace_root"

: "${SUMERU_CONF_LOCATION:=conf}"

# --- INI helpers: read [options] key = value (trim; skip # ; blank) ---
ini_get() {
	local file="$1" key="$2"
	awk -F '=' -v k="$key" '
		/^\[options\]/ { ino=1; next }
		/^\[/ { ino=0 }
		ino && $0 !~ /^[[:space:]]*(#|;|$)/ {
			k0=$1
			gsub(/^[[:space:]]+|[[:space:]]+$/, "", k0)
			if (k0 == k) {
				$1=""
				sub(/^=/, "")
				gsub(/^[[:space:]]+|[[:space:]]+$/, "", $0)
				print
				exit
			}
		}
	' "$file" 2>/dev/null || true
}

# Pick INI path: explicit env > -c in argv > conf/<basename> if exists > ./<basename>.
resolve_ini_path() {
	local explicit="${SUMERU_CONF:-${CONF:-}}"
	if [[ -n "$explicit" ]]; then
		if [[ "$explicit" == /* ]]; then
			echo "$explicit"
		else
			echo "$workspace_root/$explicit"
		fi
		return
	fi

	local i=0 n=${#args[@]}
	while [[ "$i" -lt "$n" ]]; do
		case "${args[i]}" in
		-c)
			if [[ $((i + 1)) -lt "$n" ]]; then
				local p="${args[i + 1]}"
				if [[ "$p" == /* ]]; then
					echo "$p"
				else
					echo "$workspace_root/$p"
				fi
				return
			fi
			;;
		-c=*)
			local p="${args[i]#-c=}"
			if [[ "$p" == /* ]]; then
				echo "$p"
			else
				echo "$workspace_root/$p"
			fi
			return
			;;
		esac
		i=$((i + 1))
	done

	local loc="$workspace_root/$SUMERU_CONF_LOCATION/$DEFAULT_INI_BASENAME"
	local root_ini="$workspace_root/$DEFAULT_INI_BASENAME"
	if [[ -f "$loc" ]]; then
		echo "$loc"
	elif [[ -f "$root_ini" ]]; then
		echo "$root_ini"
	else
		echo "$loc"
	fi
}

# True if argv already passes -c (value is next arg or -c=…).
args_have_c() {
	local i=0 n=${#args[@]}
	while [[ "$i" -lt "$n" ]]; do
		case "${args[i]}" in
		-c | -c=*) return 0 ;;
		esac
		i=$((i + 1))
	done
	return 1
}

# True if argv contains -i (install modules flag).
args_have_i() {
	local i=0 n=${#args[@]}
	while [[ "$i" -lt "$n" ]]; do
		case "${args[i]}" in
		-i | -i=*) return 0 ;;
		esac
		i=$((i + 1))
	done
	return 1
}

# True if argv contains -u (update modules flag).
args_have_u() {
	local i=0 n=${#args[@]}
	while [[ "$i" -lt "$n" ]]; do
		case "${args[i]}" in
		-u | -u=*) return 0 ;;
		esac
		i=$((i + 1))
	done
	return 1
}

args=("$@")
resolved_ini=$(resolve_ini_path)

# Path for make CONF= must be relative to workspace when under workspace (makefile CURDIR rule).
if [[ "$resolved_ini" == "$workspace_root"/* ]]; then
	conf_for_make="${resolved_ini#"$workspace_root"/}"
else
	conf_for_make="$resolved_ini"
fi

# Defaults for SUMERU_ROOT
SUMERU_ROOT="${SUMERU_ROOT:-}"

if [[ -f "$resolved_ini" ]]; then
	ini_dir=$(cd "$(dirname "$resolved_ini")" && pwd)
	shome=$(ini_get "$resolved_ini" sumeru_home)
	if [[ -z "$SUMERU_ROOT" && -n "$shome" ]]; then
		if [[ "$shome" == /* ]]; then
			SUMERU_ROOT="$shome"
		else
			SUMERU_ROOT=$(cd "$ini_dir" && cd "$shome" 2>/dev/null && pwd || echo "$ini_dir/$shome")
		fi
	fi
fi

if [[ -z "$SUMERU_ROOT" ]]; then
	SUMERU_ROOT=$(cd "$workspace_root" && cd "$DEFAULT_SUMERU_ROOT" 2>/dev/null && pwd || echo "$workspace_root/$DEFAULT_SUMERU_ROOT")
fi

# Ensure -c is present whenever we run the app (e.g. -i / -u only).
if ! args_have_c; then
	if [[ "$resolved_ini" == "$workspace_root"/* ]]; then
		args=(-c "${resolved_ini#"$workspace_root"/}" "${args[@]}")
	else
		args=(-c "$resolved_ini" "${args[@]}")
	fi
fi

# Optional env-driven -i / -u when not already on the command line.
if ! args_have_i && [[ -n "${SUMERU_INSTALL:-}" ]]; then
	args+=(-i "$SUMERU_INSTALL")
fi
if ! args_have_u && [[ -n "${SUMERU_UPDATE:-}" ]]; then
	args+=(-u "$SUMERU_UPDATE")
fi

make_cmd=(make "SUMERU_ROOT=$SUMERU_ROOT" "CONF=$conf_for_make")
if [[ -n "${OUT:-}" ]]; then
	make_cmd+=("OUT=$OUT")
fi
make_cmd+=(generate)

run_cmd=(go run . -- "${args[@]}")

"${make_cmd[@]}"
exec "${run_cmd[@]}"
