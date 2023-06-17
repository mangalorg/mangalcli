local fzf = require('fzf')
local json = require('json')
local anilist = require('anilist')

local mangas = anilist.search_mangas(Vars.search)

local manga = fzf.select_one(mangas, function(manga)
  for _, title in ipairs({
    manga.title.english,
    manga.title.romaji,
    manga.title.native,
  }) do
    if title ~= "" then
      return title
    end
  end
end)

json.print(manga)
