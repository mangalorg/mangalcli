local fzf = require('fzf')

local manga = fzf.select_one(SearchMangas(Vars.title), function(item)
  return item:info().title
end)

local volume = fzf.select_one(MangaVolumes(manga), function(item)
  return "Volume " .. item:info().number
end)

local chapters = fzf.select_multi(VolumeChapters(volume), function(item)
  return item:info().title
end)

for _, chapter in ipairs(chapters) do
  print("Downloading " .. chapter:info().title)
  DownloadChapter(chapter, {
    format = "pdf"
  })
end
