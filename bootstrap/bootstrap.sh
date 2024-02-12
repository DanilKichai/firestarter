#!/bin/bash -e

PAYLOAD="/tmp/payload"
read -a ARGUMENTS </proc/cmdline

( # logo
  cat "$(dirname $0)/logo" 2>/dev/null
)

( # mount
  for ARGUMENT in "${ARGUMENTS[@]}"; do
    if [[ "${ARGUMENT}" =~ ^mount=.+@[^@]+$ ]]; then
      TARGET=$(
        sed \
          --silent \
          --regexp-extended \
          --expression="s/^mount=(.+)@[^@]+$/\1/p" \
          <(echo "${ARGUMENT}")
      )

      SOURCE=$(
        sed \
          --silent \
          --regexp-extended \
          --expression="s/^mount=.+@([^@]+)$/\1/p" \
          <(echo "$ARGUMENT")
      )

      if ! DEVICE=$(
        findfs "${SOURCE}" 2>/dev/null
      ); then
        echo "Unable to resolve source to mount: ${SOURCE}"
        exit 1
      fi

      if ! mkdir -p "${TARGET}"; then
        echo "Unable to create a directory for the mount to: ${PREFIX}/${TARGET}"
        exit 1
      fi

      if ! mount "${DEVICE}" "${TARGET}"; then
        echo "Mount failed: ${DEVICE}, ${TARGET}"
      fi
    fi
  done
)

( # download
  URL="$(
    sed \
      --silent \
      --regexp-extended \
      --expression="s/(^|.*[ ])url=([^ ]*).*$/\2/p" \
      <(echo "${ARGUMENTS[@]}")
  )"

  if [ -n "${URL}" ]; then
    if ! links -source "${URL}" >"${PAYLOAD}" 2>/dev/null; then
      echo "Download failed: ${URL}"
      exit 1
    fi
  fi
)

( # checksum
  for PREFIX in md5 sha1 sha224 sha256 sha384 sha512 b2; do
    case "${PREFIX}" in
      "b2")
      ALGORITHM="BLAKE2b"
      ;;

      *)
      ALGORITHM="${PREFIX^^}"
      ;;
    esac

    HASH="$(
      sed \
        --silent \
        --regexp-extended \
        --expression="s/(^|.*[ ])${PREFIX}sum=([^ ]*).*$/\2/p" \
        <(echo "${ARGUMENTS[@]}")
    )"

    [ -z "${HASH}" ] && \
      continue

    if ! cksum \
      --check \
      <(echo "${ALGORITHM} (${PAYLOAD}) = ${HASH}") \
      &>/dev/null
    then
      echo "Checksum failed: ${PREFIX}sum"
      exit 1
    fi
  done
)

( # execute
  SHELL="/bin/bash"

  [ ! -e "${PAYLOAD}" ] && \
    exec "${SHELL}"

  if ! chmod +x "${PAYLOAD}"; then
    echo "Change the mode of the file failed: +x, ${PAYLOAD}"
    exit 1
  fi

  if ! "${PAYLOAD}"; then
    echo "Execute failed: ${PAYLOAD}"
    exit 1
  fi
)
