#!/usr/bin/env bash
set -e

buildgo() {
  echo "$(tput setaf 2)*** $1 ***$(tput sgr0)"
  cd $1
  go build
  go test
  ./$1 || true
  echo
  cd ..
}

buildc() {
  echo "$(tput setaf 2)*** $1 ***$(tput sgr0)"
  cd $1
  gcc -o $1 $1.c
  ./$1 || true
  echo
  cd ..
}

buildcs() {
  echo "$(tput setaf 2)*** $1 ***$(tput sgr0)"
  cd $1
  dotnet build
  bin/*/*/$1 || true
  echo
  cd ..
}

for dir in */ ; do
  dir=${dir%/}
  [ -d "$dir" ] || continue
  case "$dir" in
    .*) continue ;;
  esac

  if [ -f "$dir/go.mod" ]; then
    buildgo "$dir"
    continue
  fi

  if [ -f "$dir/$dir.csproj" ]; then
    buildcs "$dir"
    continue
  fi

  if [ -f "$dir/$dir.c" ]; then
    buildc "$dir"
    continue
  fi

  echo "Skipping $dir (no recognized project files)"
  echo
done
