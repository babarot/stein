#!/bin/bash

set -e

ask() {
    message="${1:-Are you sure?}"
    while true
    do
        read -r -p "${message} [y/n] " input
        case $input in
            [yY][eE][sS]|[yY])
                return 0
                ;;
            [nN][oO]|[nN])
                echo "[INFO] canceled" >&2
                return 1
                ;;
            *)
                echo "[ERROR] Invalid input...again"
                ;;
        esac
    done
}

main() {
    if [[ -n "$(git status -s)" ]]; then
        echo "[ERROR] there are untracked or unstaged files" >&2
        return 1
    fi

    current_version="$(gobump show -r)"
    echo "[INFO] current version: ${current_version}"

    while true
    do
        read -p "Specify [major | minor | patch]: " version
        case "${version}" in
            major | minor | patch )
                gobump "${version}" -w
                next_version="$(gobump show -r)"
                break
                ;;
            "")
                echo "[INFO] canceled" >&2
                return 1
                ;;
            *)
                echo "[ERROR] ${version}: invalid semver type" >&2
                continue
                ;;
        esac
        shift
    done

    git-chglog -o CHANGELOG.md --next-tag "v${next_version}"
    git --no-pager diff

    ask "OK to commit/push these changes?" || return 1
    git commit -am "Bump version ${next_version} and update changelog"
    git tag "v${next_version}"
    git push && git push --tags

    ask "OK to release?" || return 1
    goreleaser --release-notes CHANGELOG.md --rm-dist
}

main "${@}"
exit
