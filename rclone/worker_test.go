package rclone

import "testing"

func TestRenameFile(t *testing.T) {
	var renamed, want string

	renamed = RenameFile("[Lilith-Raws] Karakai Jouzu no Takagi-san S03 - 01v2 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4")
	want = "[Lilith-Raws] Karakai Jouzu no Takagi-san III - 01v2 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4"
	if renamed != want {
		t.Errorf("New filename was incorrect, got: %s, want: %s.", renamed, want)
	}

	renamed = RenameFile("[Nekomoe kissaten][Sono Bisque Doll wa Koi wo Suru][01V2][1080p][CHS].mp4")
	want = "[Nekomoe kissaten] Sono Bisque Doll wa Koi wo Suru - 01v2 [1080p][CHS].mp4"
	if renamed != want {
		t.Errorf("New filename was incorrect, got: %s, want: %s.", renamed, want)
	}

	renamed = RenameFile("[SBE][KAMEN RIDER REVICE][EP01][x264_AAC][1080P].mp4")
	want = "[SBE] KAMEN RIDER REVICE - 01 [x264_AAC][1080P].mp4"
	if renamed != want {
		t.Errorf("New filename was incorrect, got: %s, want: %s.", renamed, want)
	}

	renamed = RenameFile("[Ohys-Raws] Shingeki no Kyojin The Final Season Part 2 - 01 (NHKG 1280x720 x264 AAC JP).mp4")
	want = "[Ohys-Raws] Shingeki no Kyojin The Final Season Part II - 01 (NHKG 1280x720 x264 AAC JP).mp4"
	if renamed != want {
		t.Errorf("New filename was incorrect, got: %s, want: %s.", renamed, want)
	}

	renamed = RenameFile("[Nekomoe kissaten] Slow Loop 01 [WebRip 1080p HEVC-10bit AAC ASSx2].mkv")
	want = "[Nekomoe kissaten] Slow Loop - 01 [WebRip 1080p HEVC-10bit AAC ASSx2].mkv"
	if renamed != want {
		t.Errorf("New filename was incorrect, got: %s, want: %s.", renamed, want)
	}
}
