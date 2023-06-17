#!/bin/sh

set -euo pipefail


SCRIPT=$(cat <<-EOF
local mangas = SearchMangas(Vars.search)

if not Vars.manga then
  for i, manga in ipairs(mangas) do
    local title = manga:info().title
    print(i, title)
  end

  return
end

local manga = mangas[tonumber(Vars.manga)]
local volumes = MangaVolumes(manga)

if not Vars.volume then
  for i, volume in ipairs(volumes) do
    local number = volume:info().number
    print(i, number)
  end

  return
end

local volume = volumes[tonumber(Vars.volume)]
local chapters = VolumeChapters(volume)

if not Vars.chapter then
  for i, chapter in ipairs(chapters) do
    local title = chapter:info().title
    print(i, title)
  end

  return
end

local chapter = chapters[tonumber(Vars.chapter)]

DownloadChapter(chapter, {
  directory = Vars.dir,
  read_after = true
})

EOF

)

PROVIDER="$1"
SEARCH="$2"

select_one() {
  fzf --delimiter="\t" --with-nth=2.. | awk '{print $1}'
}

M() {
  mangalcli run "$SCRIPT" --provider "$PROVIDER" --vars="$1"
}

VARS="search=$SEARCH"

echo "Searching for '$SEARCH'"
MANGA=$(M "$VARS" | select_one)
VARS="$VARS;manga=$MANGA"

echo "Getting volumes"
VOLUME=$(M "$VARS" | select_one)
VARS="$VARS;volume=$VOLUME"

echo "Getting chapters"
CHAPTER=$(M "$VARS" | select_one)
VARS="$VARS;chapter=$CHAPTER"


echo "Downloading chapter for reading..."
M "$VARS;dir=$(mktemp -d)"

