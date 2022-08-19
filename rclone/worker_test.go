package rclone

import "testing"

func TestRenameFileInAnimeStyle(t *testing.T) {
	type TestCase struct {
		Name       string
		Filename   string
		AnimeStyle string
	}

	cases := []TestCase{
		{
			Name:       "case 01",
			Filename:   "[Lilith-Raws] Karakai Jouzu no Takagi-san S03 - 01v2 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4",
			AnimeStyle: "[Lilith-Raws] Karakai Jouzu no Takagi-san III - 01v2 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4",
		},
		{
			Name:       "case 02",
			Filename:   "[Nekomoe kissaten][Sono Bisque Doll wa Koi wo Suru][10.5(OVA)][1080p][CHS].mp4",
			AnimeStyle: "[Nekomoe kissaten] Sono Bisque Doll wa Koi wo Suru - S10.5 [1080p][CHS].mp4",
		},
		{
			Name:       "case 03",
			Filename:   "[SBE][KAMEN RIDER REVICE][EP01V2][x264_AAC][1080P].mp4",
			AnimeStyle: "[SBE] KAMEN RIDER REVICE - 01v2 [x264_AAC][1080P].mp4",
		},
		{
			Name:       "case 04",
			Filename:   "[Ohys-Raws] Shingeki no Kyojin The Final Season Part 2 - 12 END (NHKG 1280x720 x264 AAC JP).mp4",
			AnimeStyle: "[Ohys-Raws] Shingeki no Kyojin The Final Season II - 12 END (NHKG 1280x720 x264 AAC JP).mp4",
		},
		{
			Name:       "case 05",
			Filename:   "[Ohys-Raws] Princess Connect! ReDive Season 2 - 02 (BS11 1280x720 x264 AAC).mp4",
			AnimeStyle: "[Ohys-Raws] Princess Connect! ReDive II - 02 (BS11 1280x720 x264 AAC).mp4",
		},
		{
			Name:       "case 06",
			Filename:   "[Ohys-Raws] Arifureta Shokugyou de Sekai Saikyou 2nd Season - 01 (AT-X 1280x720 x264 AAC).mp4",
			AnimeStyle: "[Ohys-Raws] Arifureta Shokugyou de Sekai Saikyou II - 01 (AT-X 1280x720 x264 AAC).mp4",
		},
		{
			Name:       "case 07",
			Filename:   "[Nekomoe kissaten] Slow Loop 01 [WebRip 1080p HEVC-10bit AAC ASSx2].mkv",
			AnimeStyle: "[Nekomoe kissaten] Slow Loop - 01 [WebRip 1080p HEVC-10bit AAC ASSx2].mkv",
		},
		{
			Name:       "case 08",
			Filename:   "[Lilith-Raws] 86 - Eighty Six - 01 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4",
			AnimeStyle: "[Lilith-Raws] Eighty Six - 01 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4",
		},
		{
			Name:       "case 09",
			Filename:   "[Lilith-Raws] Tokyo 24-ku - 01 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4",
			AnimeStyle: "[Lilith-Raws] Tokyo ku - 01 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4",
		},
		{
			Name:       "case 10",
			Filename:   "[Kawaiika-Raws] Kobayashi-san (2021) 01 [BDRip 1920x1080 HEVC FLAC].mkv",
			AnimeStyle: "[Kawaiika-Raws] Kobayashi-san - 01 [BDRip 1920x1080 HEVC FLAC].mkv",
		},
		{
			Name:       "case 11",
			Filename:   "[GM-Team][国漫][Dou Luo Da Lu][Douro Mainland][2019][190][AVC][GB][1080P].mp4",
			AnimeStyle: "[GM-Team] 国漫 Dou Luo Da Lu Douro Mainland - 190 [AVC][GB][1080P].mp4",
		},
		{
			Name:       "case 12",
			Filename:   "[Pussub&VCB-Studio] White Album 2 [01][Hi10p_1080p][x264_2flac].mkv",
			AnimeStyle: "[Pussub&VCB-Studio] White Album II - 01 [Hi10p_1080p][x264_2flac].mkv",
		},
		{
			Name:       "case 13",
			Filename:   "[VCB-Studio] PERSONA5 the Animation [26.5(Summary)][Ma10p_1080p][x265_flac].mkv",
			AnimeStyle: "[VCB-Studio] PERSONA5 the Animation - S26.5 [Ma10p_1080p][x265_flac].mkv",
		},
		{
			Name:       "case 14",
			Filename:   "[HYSUB]Komi-san wa, Komyushou Desu.[13][GB_MP4][1920X1080].mp4",
			AnimeStyle: "[HYSUB] Komi-san wa, Komyushou Desu. - 13 [GB_MP4][1920X1080].mp4",
		},
		{
			Name:       "case 15",
			Filename:   "[UHA-WINGS] [RPG Fudousan] [01] [x264 1080p] [CHS].mp4",
			AnimeStyle: "[UHA-WINGS] RPG Fudousan - 01 [x264 1080p][CHS].mp4",
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			result := RenameFileInAnimeStyle(c.Filename)
			if result != c.AnimeStyle {
				t.Errorf("New filename was incorrect, got: %s, want: %s.", result, c.AnimeStyle)
			}
		})
	}
}

func TestRenameFileInTvStyle(t *testing.T) {
	type TestCase struct {
		Name     string
		Filename string
		TvStyle  string
	}

	cases := []TestCase{
		{
			Name:     "case 01",
			Filename: "[DMG][Yofukashi_no_Uta][01][1080P][GB].mp4",
			TvStyle:  "Yofukashi no Uta - S01E01 - [DMG][1080P][GB].mp4",
		},
		{
			Name:     "case 02",
			Filename: "[Ohys-Raws] Kawaii Dake ja Nai Shikimori-san - 12 END (EX 1280x720 x264 AAC).mp4",
			TvStyle:  "Kawaii Dake ja Nai Shikimori-san - S01E12 - [Ohys-Raws] (EX 1280x720 x264 AAC).mp4",
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			result := RenameFileInTvStyle(c.Filename)
			if result != c.TvStyle {
				t.Errorf("New filename was incorrect, got: %s, want: %s.", result, c.TvStyle)
			}
		})
	}
}
