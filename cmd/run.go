package cmd

import (
	"context"
	"fmt"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/luaprovider"
	"github.com/mangalorg/mangalcli/selector"
	"os"
)

type runCmd struct {
	Path                string            `arg:"" help:"path to the lua script" type:"existingfile"`
	Vars                map[string]string `help:"variables to pass to the selection query"`
	Query               string            `help:"selection query as lua script. See wiki for more" required:""`
	Json                bool              `help:"show json"`
	Download            bool              `help:"download"`
	Format              string            `help:"format to use for downloading chapters" enum:"pdf,cbz,images" default:"pdf"`
	Dir                 string            `help:"directory to download to" default:"." type:"existingdir"`
	CreateMangaDir      bool              `help:"create manga directory"`
	CreateVolumeDir     bool              `help:"create volume directory"`
	Strict              bool              `help:"fail if any metadata write fails"`
	SkipIfExists        bool              `help:"do not download if chapter exists" negatable:"" default:"true"`
	Cover               bool              `help:"download manga cover"`
	Banner              bool              `help:"download manga banner"`
	SeriesJson          bool              `help:"create series.json"`
	ComicInfoXml        bool              `help:"create ComicInfo.xml"`
	ReadAfter           bool              `help:"open the chapter for reading after it was downloaded"`
	ReadIncognito       bool              `help:"do not save chapter to the anilist reading history"`
	ComicInfoXmlAddDate bool              `help:"add date to the ComicInfo.xml" negatable:"" default:"true"`
}

func (r *runCmd) Run() error {
	contents, err := os.ReadFile(r.Path)
	if err != nil {
		return err
	}

	loader, err := luaprovider.NewLoader(contents, luaprovider.DefaultOptions())
	if err != nil {
		return err
	}

	client, err := libmangal.NewClient(context.Background(), loader, libmangal.DefaultClientOptions())
	if err != nil {
		return err
	}

	selected, err := selector.Select(&client, r.Vars, r.Query)
	if err != nil {
		return err
	}

	if r.Json {
		marshalled, err := marshal(selected)
		if err != nil {
			return err
		}

		fmt.Println(marshalled)
	}

	if r.Download {
		var chapters []libmangal.Chapter

		switch selected := selected.(type) {
		case []any:
			for _, s := range selected {
				chapter, ok := s.(libmangal.Chapter)
				if !ok {
					return fmt.Errorf("expected chapters, got: %T", s)
				}

				chapters = append(chapters, chapter)
			}
		case libmangal.Chapter:
			chapters = []libmangal.Chapter{selected}
		default:
			return fmt.Errorf("expected chapters, got: %T", selected)
		}

		for _, chapter := range chapters {
			format, err := libmangal.FormatString(r.Format)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "downloading chapter %q\n", chapter)

			path, err := client.DownloadChapter(context.Background(), chapter, libmangal.DownloadOptions{
				Format:              format,
				Directory:           r.Dir,
				CreateMangaDir:      r.CreateMangaDir,
				CreateVolumeDir:     r.CreateVolumeDir,
				Strict:              r.Strict,
				SkipIfExists:        r.SkipIfExists,
				DownloadMangaCover:  r.Cover,
				DownloadMangaBanner: r.Banner,
				WriteSeriesJson:     r.SeriesJson,
				WriteComicInfoXml:   r.ComicInfoXml,
				ReadAfter:           r.ReadAfter,
				ReadIncognito:       r.ReadIncognito,
				ComicInfoOptions: libmangal.ComicInfoXmlOptions{
					AddDate: r.ComicInfoXmlAddDate,
				},
			})

			if err != nil {
				return err
			}

			fmt.Println(path)
		}
	}

	return nil
}
