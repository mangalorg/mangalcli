# MangalCLI

Frontend for the [libmangal](https://github.com/mangalorg/libmangal) and
[luaprovider](https://github.com/mangalorg/luaprovider) wrapped
as a CLI app.

## Example

```bash
mangalcli run mangapill.lua --vars="title=chainsaw man" --query "$(cat query.lua)" --json
```

**...where**:

`mangapill.lua` looks like [this](https://github.com/mangalorg/saturno/blob/261c5739eacb73525fbe52705b8862a11c14040f/luas/mangapill.lua);

`query.lua` looks like this:

```lua
local mangas = SearchMangas(Vars.title) -- search with the given title
local volumes = MangaVolumes(mangas[1]) -- select the first manga

local chapters = {}

-- get all chapters of the manga
for _, volume in ipairs(volumes) do
  for _, chapter in ipairs(VolumeChapters(volume)) do
    table.insert(chapters, chapter)
  end
end

-- return chapters
return chapters
```

and the output would be like the following (because we passed the `--json` flag):

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

## Queries

Queries use Lua5.1(+ goto statement from Lua5.2)

`require('sdk')` from [luaprovider](https://github.com/mangalorg/luaprovider)
isn't available (subject of change)

It has the following functions available:

- `SearchMangas(query: string): []Manga`
- `MangaVolumes(manga: Manga): []Volume`
- `VolumeChapters(volume: Volume): []Chapter`
- `ChapterPages(chapter: Chapter): []Page`

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

Queries must end with `return` statement.

Queries can return any type of data,
but if the `--download` flag is provided only `Chapter` or `[]Chapter` are allowed
to be returned.