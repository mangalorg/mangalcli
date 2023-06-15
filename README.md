# MangalCLI

Frontend for the [libmangal](https://github.com/mangalorg/libmangal) and
[luaprovider](https://github.com/mangalorg/luaprovider) wrapped
as a CLI app.

## Example

[See more examples here](./examples)

> Documentation is still work in progress... 😪

```bash
mangalcli run mangapill.lua --vars="title=chainsaw man" --exec "$(cat exec.lua)"
```

**...where**:

`mangapill.lua` looks like [this](https://github.com/mangalorg/saturno/blob/261c5739eacb73525fbe52705b8862a11c14040f/luas/mangapill.lua);

`exec.lua` looks like this:

```lua
local json = require('json')

local mangas = SearchMangas(Vars.title) -- search with the given title
local volumes = MangaVolumes(mangas[1]) -- select the first manga

local chapters = {}

-- get all chapters of the manga
for _, volume in ipairs(volumes) do
  for _, chapter in ipairs(VolumeChapters(volume)) do
    table.insert(chapters, chapter)
  end
end

-- chapters encoded in json format for later use
json.print(chapters)
```

and the output would be like the following

```json
[
  {
    "title": "Chapter 1",
    "url": "https://mangapill.com/chapters/723-10001000/chainsaw-man-chapter-1",
    "number": 1
  },
  {
    "title": "Chapter 2",
    "url": "https://mangapill.com/chapters/723-10002000/chainsaw-man-chapter-2",
    "number": 2
  },
  // etc
]
```

## Install

It's not packaged for any package manager *yet*
so the only option is to build from source.

Either like this:

```bash
go install github.com/mangalorg/mangalcli
```

Or like this:

```bash
git clone github.com/mangalorg/mangalcli
cd mangalcli
go install .
```

## Exec Scripts

Scripts use Lua5.1(+ goto statement from Lua5.2)

`sdk` package from [luaprovider](https://github.com/mangalorg/luaprovider)
isn't available (subject of change).

Available functions:

```lua
--- @alias Manga userdata
--- @alias Volume userdata
--- @alias Chapter userdata
--- @alias Page userdata

--- @alias MangaInfo { title: string, id: number, url: string, cover: string, banner: string }
--- @alias VolumeInfo { number: number }
--- @alias ChapterInfo { title: string, url: string, number: number }

--- @alias Format "pdf" | "cbz" | "images"

--- @alias DownloadOptions { format: Format, directory: string, create_manga_dir: boolean, create_volume_dir: boolean, strict: boolean, skip_if_exists: boolean, download_manga_cover: boolean, download_manga_banner: boolean, write_series_json: boolean, write_comic_info_xml: boolean, read_after: boolean, read_incognito: boolean }

--- @param query string
--- @return []MangaInfo
function SearchMangas(query) end

--- @param manga Manga
--- @return []Volume
function MangaVolumes(manga) end

--- @param volume Volume
--- @return []Chapter
function VolumeChapters(volume) end

--- @param chapter Chapter
--- @return []Page
function ChapterPages(chapter) end


--- @param chapter Chapter
--- @param options DownloadOptions?
function DownloadChapter(chapter, options) end

--- Manga, Volume and Chapter has :info() method that would
--- return appropriate info table as defined above
---
--- e.g. manga:info().title
```

Comes with `anilist` package that has the following functions available:

```lua
local anilist

--- @alias AnilistManga TODO, see https://pkg.go.dev/github.com/mangalorg/libmangal#AnilistManga

--- @param title string
--- @return AnilistManga?
function anilist.find_closest_manga(title) end

--- @param title string'
--- @return []AnilistManga
function anilist.search_mangas(title) end

--- @param id number
--- @return AnilistManga?
function anilist.get_by_id(id) end

return anilist
```

And `json` package

```lua
local json

--- Prints data in json format to stdout
--- @param data any
function json.print(data) end

return json
```

Example:

```lua
local anilist = require("anilist")
local json = require("json")

json.print(anilist.search_mangas("one piece"))
```

And `Vars` global table with variables passed from `--vars` flag.

Used like this:

```bash
--vars="key1=value1;key2=value2;key3=18.2"
```

```lua
print(Vars.key1) -- value1
print(Vars.key2) -- value2

-- Note, that all values are passed as strings.
-- If you want to get number value use

print(tonumber(Vars.key3)) -- 18.2
```

