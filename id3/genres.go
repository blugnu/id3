package id3

type Genre byte

const (
	Blues Genre = iota
	ClassicRock
	Country
	Dance
	Disco
	Funk
	Grunge
	HipHop
	Jazz
	Metal
	NewAge
	Oldies
	Other
	Pop
	RAndB
	Rap
	Reggae
	Rock
	Techno
	Industrial
	Alternative
	Ska
	DeathMetal
	Pranks
	Soundtrack
	EuroTechno
	Ambient
	TripHop
	Vocal
	JazzFunk
	Fusion
	Trance
	Classical
	Instrumental
	Acid
	House
	Game
	SoundClip
	Gospel
	Noise
	AlternRock
	Bass
	Soul
	Punk
	Space
	Meditative
	InstrumentalPop
	InstrumentalRock
	Ethnic
	Gothic
	Darkwave
	TechnoIndustrial
	Electronic
	PopFolk
	Eurodance
	Dream
	SouthernRock
	Comedy
	Cult
	Gangsta
	Top40
	ChristianRap
	PopFunk
	Jungle
	NativeAmerican
	Cabaret
	NewWave
	Psychadelic
	Rave
	Showtunes
	Trailer
	LoFi
	Tribal
	AcidPunk
	AcidJazz
	Polka
	Retro
	Musical
	RockAndRoll
	HardRock
	// Original WinAmp Extensions
	Folk
	FolkRock
	NationalFolk
	Swing
	FastFusion
	Bebob
	Latin
	Revival
	Celtic
	Bluegrass
	Avantgarde
	GothicRock
	ProgressiveRock
	PsychedelicRock
	SymphonicRock
	SlowRock
	BigBand
	Chorus
	EasyListening
	Acoustic
	Humour
	Speech
	Chanson
	Opera
	ChamberMusic
	Sonata
	Symphony
	BootyBass
	Primus
	PornGroove
	Satire
	SlowJam
	Club
	Tango
	Samba
	Folklore
	Ballad
	PowerBallad
	RhythmicSoul
	Freestyle
	Duet
	PunkRock
	DrumSolo
	ACapella
	EuroHouse
	DanceHall
	GoaTrance
	DrumAndBass
	ClubHouse
	HardcoreTechno
	Terror
	Indie
	BritPop
	AfroPunk
	PolskPunk
	Beat
	ChristianGangstaRap
	HeavyMetal
	BlackMetal
	Crossover
	ContemporaryChristian
	ChristianRock
	// WinAmp 1.91
	Merengue
	Salsa
	ThrashMetal
	Anime
	JPop
	Synthpop
	// WinAmp 5.6
	Abstract
	ArtRock
	Baroque
	Bhangra
	BigBeat
	Breakbeat
	Chillout
	Downtempo
	Dub
	EBM
	Eclectic
	Electro
	Electroclash
	Emo
	Experimental
	Garage
	Global
	IDM
	Illbient
	IndustroGoth
	JamBand
	Krautrock
	Leftfield
	Lounge
	MathRock
	NewRomantic
	NuBreakz
	PostPunk
	PostRock
	Psytrance
	Shoegaze
	SpaceRock
	TropRock
	WorldMusic
	Neoclassical
	Audiobook
	AudioTheatre
	NeueDeutscheWelle
	Podcast
	IndieRock
	GFunk
	Dubstep
	GarageRock
	Psybient
	// No Genre
	NoGenre = 255
)

var GenreName = map[Genre]string{
	Blues:            "Blues",
	ClassicRock:      "Classic Rock",
	Country:          "Country",
	Dance:            "Dance",
	Disco:            "Disco",
	Funk:             "Funk",
	Grunge:           "Grunge",
	HipHop:           "Hip-Hop",
	Jazz:             "Jazz",
	Metal:            "Metal",
	NewAge:           "New Age",
	Oldies:           "Oldies",
	Other:            "Other",
	Pop:              "Pop",
	RAndB:            "R 'n B",
	Rap:              "Rap",
	Reggae:           "Reggae",
	Rock:             "Rock",
	Techno:           "Techno",
	Industrial:       "Industrial",
	Alternative:      "Alternative",
	Ska:              "Ska",
	DeathMetal:       "Death Metal",
	Pranks:           "Pranks",
	Soundtrack:       "Soundtrack",
	EuroTechno:       "Euro Techno",
	Ambient:          "Ambient",
	TripHop:          "Trip Hop",
	Vocal:            "Vocal",
	JazzFunk:         "Jazz Funk",
	Fusion:           "Fusion",
	Trance:           "Trance",
	Classical:        "Classical",
	Instrumental:     "Instrumental",
	Acid:             "Acid",
	House:            "House",
	Game:             "Game",
	SoundClip:        "Sound Clip",
	Gospel:           "Gospel",
	Noise:            "Noise",
	AlternRock:       "Alternative Rock",
	Bass:             "Bass",
	Soul:             "Soul",
	Punk:             "Punk",
	Space:            "Space",
	Meditative:       "Meditative",
	InstrumentalPop:  "Instrumental Pop",
	InstrumentalRock: "Instrumental Rock",
	Ethnic:           "Ethnic",
	Gothic:           "Gothic",
	Darkwave:         "Darkwave",
	TechnoIndustrial: "Techno Industrial",
	Electronic:       "Electronic",
	PopFolk:          "Pop Folk",
	Eurodance:        "Eurodance",
	Dream:            "Dream",
	SouthernRock:     "Southern Rock",
	Comedy:           "Comedy",
	Cult:             "Cult",
	Gangsta:          "Gangsta",
	Top40:            "Top 40",
	ChristianRap:     "Christian Rap",
	PopFunk:          "Pop Funk",
	Jungle:           "Jungle",
	NativeAmerican:   "Native American",
	Cabaret:          "Cabaret",
	NewWave:          "New Wave",
	Psychadelic:      "Psychadelic",
	Rave:             "Rave",
	Showtunes:        "Showtunes",
	Trailer:          "Trailer",
	LoFi:             "Lo-Fi",
	Tribal:           "Tribal",
	AcidPunk:         "Acid Punk",
	AcidJazz:         "Acid Jazz",
	Polka:            "Polka",
	Retro:            "Retro",
	Musical:          "Musical",
	RockAndRoll:      "Rock & Roll",
	HardRock:         "Hard Rock",
	// Original WinAmp extensions...
	Folk:                  "Folk",
	FolkRock:              "Folk Rock",
	NationalFolk:          "National Folk",
	Swing:                 "Swing",
	FastFusion:            "Fast Fusion",
	Bebob:                 "Bebob",
	Latin:                 "Latin",
	Revival:               "Revival",
	Celtic:                "Celtic",
	Bluegrass:             "Bluegrass",
	Avantgarde:            "Avantgarde",
	GothicRock:            "Gothic Rock",
	ProgressiveRock:       "Progressive Rock",
	PsychedelicRock:       "Psychedelic Rock",
	SymphonicRock:         "Symphonic Rock",
	SlowRock:              "Slow Rock",
	BigBand:               "Big Band",
	Chorus:                "Chorus",
	EasyListening:         "Easy Listening",
	Acoustic:              "Acoustic",
	Humour:                "Humour",
	Speech:                "Speech",
	Chanson:               "Chanson",
	Opera:                 "Opera",
	ChamberMusic:          "Chamber Music",
	Sonata:                "Sonata",
	Symphony:              "Symphony",
	BootyBass:             "Booty Bass",
	Primus:                "Primus",
	PornGroove:            "Porn Groove",
	Satire:                "Satire",
	SlowJam:               "Slow Jam",
	Club:                  "Club",
	Tango:                 "Tango",
	Samba:                 "Samba",
	Folklore:              "Folklore",
	Ballad:                "Ballad",
	PowerBallad:           "Power Ballad",
	RhythmicSoul:          "Rhythmic Soul",
	Freestyle:             "Freestyle",
	Duet:                  "Duet",
	PunkRock:              "Punk Rock",
	DrumSolo:              "Drum Solo",
	ACapella:              "a Capella",
	EuroHouse:             "Euro House",
	DanceHall:             "Dance Hall",
	GoaTrance:             "Goa Trance",
	DrumAndBass:           "Drum & Bass",
	ClubHouse:             "Club House",
	HardcoreTechno:        "Hardcore Techno",
	Terror:                "Terror",
	Indie:                 "Indie",
	BritPop:               "BritPop",
	AfroPunk:              "Afro-Punk",
	PolskPunk:             "Polsk Punk",
	Beat:                  "Beat",
	ChristianGangstaRap:   "Christian Gangsta Rap",
	HeavyMetal:            "Heavy Metal",
	BlackMetal:            "Black Metal",
	Crossover:             "Crossover",
	ContemporaryChristian: "Contemporary Christian",
	ChristianRock:         "Christian Rock",
	// WinAmp 1.91 - June 1998+
	Merengue:    "Merengue",
	Salsa:       "Salsa",
	ThrashMetal: "Thrash Metal",
	Anime:       "Anime",
	JPop:        "JPop",
	Synthpop:    "Synthpop",
	// Winamp 5.6,
	Abstract:          "Abstract",
	ArtRock:           "Art Rock",
	Baroque:           "Baroque",
	Bhangra:           "Bhangra",
	BigBeat:           "Big Beat",
	Breakbeat:         "Breakbeat",
	Chillout:          "Chillout",
	Downtempo:         "Downtempo",
	Dub:               "Dub",
	EBM:               "EBM",
	Eclectic:          "Eclectic",
	Electro:           "Electro",
	Electroclash:      "Electroclash",
	Emo:               "Emo",
	Experimental:      "Experimental",
	Garage:            "Garage",
	Global:            "Global",
	IDM:               "IDM",
	Illbient:          "Illbient",
	IndustroGoth:      "Industro-Goth",
	JamBand:           "Jam Band",
	Krautrock:         "Krautrock",
	Leftfield:         "Leftfield",
	Lounge:            "Lounge",
	MathRock:          "Math Rock",
	NewRomantic:       "New Romantic",
	NuBreakz:          "Nu-Breakz",
	PostPunk:          "Post-Punk",
	PostRock:          "Post-Rock",
	Psytrance:         "Psytrance",
	Shoegaze:          "Shoegaze",
	SpaceRock:         "Space Rock",
	TropRock:          "Trop Rock",
	WorldMusic:        "World Music",
	Neoclassical:      "Neoclassical",
	Audiobook:         "Audiobook",
	AudioTheatre:      "Audio Theatre",
	NeueDeutscheWelle: "Neue Deutsche Welle",
	Podcast:           "Podcast",
	IndieRock:         "Indie Rock",
	GFunk:             "G-Funk",
	Dubstep:           "Dubstep",
	GarageRock:        "Garage Rock",
	Psybient:          "Psybient",
}
