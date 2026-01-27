# ===== Debugo shell hook =====

DEBUGO_LAST_CMD=""
DEBUGO_ERROR_FILE=".debugo/last_error.json"

# ---------- ZSH ----------
if [ -n "$ZSH_VERSION" ]; then
  preexec() {
    DEBUGO_LAST_CMD="$1"
  }

  precmd() {
    local exit_code=$?

    if [ -f ".debugo/metadata.json" ] && [ $exit_code -ne 0 ]; then
      printf '{"cmd":"%s","exit":%d,"cwd":"%s"}\n' \
        "$DEBUGO_LAST_CMD" "$exit_code" "$PWD" > "$DEBUGO_ERROR_FILE"
    fi

    if [ -f "$DEBUGO_ERROR_FILE" ]; then
      echo
      echo "⚠️  debugo: last command failed (exit $exit_code)"
      read -q "ans?Record? [y/N] "
      echo
      if [[ "$ans" == "y" || "$ans" == "Y" ]]; then
        debugo record 
      fi
      rm -f "$DEBUGO_ERROR_FILE"
    fi
  }
fi

# ---------- BASH ----------
if [ -n "$BASH_VERSION" ]; then
  DEBUGO_LAST_CMD=""
  DEBUGO_ERROR_FILE=".debugo/last_error.json"

  # ---- preexec (before command runs) ----
  debugo_preexec() {
    DEBUGO_LAST_CMD="$1"
  }

  # ---- precmd (before prompt) ----
  debugo_precmd() {
    local exit_code=$?

    if [ -f ".debugo/metadata.json" ] && [ "$exit_code" -ne 0 ]; then
      printf '{"cmd":"%s","exit":%d,"cwd":"%s"}\n' \
        "$DEBUGO_LAST_CMD" "$exit_code" "$PWD" > "$DEBUGO_ERROR_FILE"
    fi

    if [ -f "$DEBUGO_ERROR_FILE" ]; then
      echo
      echo "⚠️  debugo: last command failed (exit $exit_code)"
      read -p "Record? [y/N] " ans
      if [[ "$ans" =~ ^[yY]$ ]]; then
        debugo_cli record last
      fi
      rm -f "$DEBUGO_ERROR_FILE"
    fi
  }

  # ---- register hooks safely ----
  if declare -p preexec_functions &>/dev/null; then
    preexec_functions+=(debugo_preexec)
  else
    preexec_functions=(debugo_preexec)
  fi

  if declare -p precmd_functions &>/dev/null; then
    precmd_functions+=(debugo_precmd)
  else
    precmd_functions=(debugo_precmd)
  fi
fi
