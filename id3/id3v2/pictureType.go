package id3v2

type PictureType byte

const (
	Other   PictureType = iota
	PngIcon             // 32x32
	OtherIcon
	FrontCover
	BackCover
	LeafletPage
	Media
	LeadArtist
	Artist
	Conductor
	BandOrOrchestra
	Composer
	Lyricist
	RecordingLocation
	DuringRecording
	DuringPerformance
	ScreenCapture
	BrightColouredFish
	Illustration
	ArtistLogo
	StudioLogo
)
const maxPictureType = StudioLogo
