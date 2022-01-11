package rclone

import "testing"

func TestRenameFile(t *testing.T) {
	var renamed, want string

	renamed = RenameFile("[Lilith-Raws] Karakai Jouzu no Takagi-san S03 - 01v2 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4")
	want = "[Lilith-Raws] Karakai Jouzu no Takagi-san III - 01v2 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4"
	if renamed != want {
		t.Errorf("New filename was incorrect, got: %s, want: %s.", renamed, want)
	}

	renamed = RenameFile("[Nekomoe kissaten][Sono Bisque Doll wa Koi wo Suru][10.5(OVA)][1080p][CHS].mp4")
	want = "[Nekomoe kissaten] Sono Bisque Doll wa Koi wo Suru - S10.5 [1080p][CHS].mp4"
	if renamed != want {
		t.Errorf("New filename was incorrect, got: %s, want: %s.", renamed, want)
	}

	renamed = RenameFile("[SBE][KAMEN RIDER REVICE][EP01V2][x264_AAC][1080P].mp4")
	want = "[SBE] KAMEN RIDER REVICE - 01v2 [x264_AAC][1080P].mp4"
	if renamed != want {
		t.Errorf("New filename was incorrect, got: %s, want: %s.", renamed, want)
	}

	renamed = RenameFile("[Ohys-Raws] Shingeki no Kyojin The Final Season Part 2 - 12 END (NHKG 1280x720 x264 AAC JP).mp4")
	want = "[Ohys-Raws] Shingeki no Kyojin The Final Season Part II - 12 END (NHKG 1280x720 x264 AAC JP).mp4"
	if renamed != want {
		t.Errorf("New filename was incorrect, got: %s, want: %s.", renamed, want)
	}

	renamed = RenameFile("[Nekomoe kissaten] Slow Loop 01 [WebRip 1080p HEVC-10bit AAC ASSx2].mkv")
	want = "[Nekomoe kissaten] Slow Loop - 01 [WebRip 1080p HEVC-10bit AAC ASSx2].mkv"
	if renamed != want {
		t.Errorf("New filename was incorrect, got: %s, want: %s.", renamed, want)
	}

	renamed = RenameFile("[Lilith-Raws] 86 - Eighty Six - 01 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4")
	want = "[Lilith-Raws] Eighty Six - 01 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4"
	if renamed != want {
		t.Errorf("New filename was incorrect, got: %s, want: %s.", renamed, want)
	}

	renamed = RenameFile("[Lilith-Raws] Tokyo 24-ku - 01 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4")
	want = "[Lilith-Raws] Tokyo ku - 01 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4"
	if renamed != want {
		t.Errorf("New filename was incorrect, got: %s, want: %s.", renamed, want)
	}

	renamed = RenameFile("[Kawaiika-Raws] Kobayashi-san (2021) 01 [BDRip 1920x1080 HEVC FLAC].mkv")
	want = "[Kawaiika-Raws] Kobayashi-san - 01 [BDRip 1920x1080 HEVC FLAC].mkv"
	if renamed != want {
		t.Errorf("New filename was incorrect, got: %s, want: %s.", renamed, want)
	}

	renamed = RenameFile("[GM-Team][国漫][Dou Luo Da Lu][Douro Mainland][2019][190][AVC][GB][1080P].mp4")
	want = "[GM-Team] 国漫 Dou Luo Da Lu Douro Mainland - 190 [AVC][GB][1080P].mp4"
	if renamed != want {
		t.Errorf("New filename was incorrect, got: %s, want: %s.", renamed, want)
	}
}
