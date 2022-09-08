package id3v2

type FrameFlags struct {
	PreserveWhenTagAltered  bool // 2.3.0 + 2.4.0
	PreserveWhenFileAltered bool // 2.3.0 + 2.4.0
	IsReadOnly              bool // 2.3.0 + 2.4.0
	IsCompressed            bool // 2.3.0 + 2.4.0
	IsEncrypted             bool // 2.3.0 + 2.4.0
	IsGrouped               bool // 2.3.0 + 2.4.0
	IsUnsynchronised        bool // 2.4.0
	HasDataLength           bool // 2.4.0
}

type Frame struct {
	ID    string
	Size  int         // size of frame, excluding header (always 6 bytes (2.2.0) or 10 bytes (2.3.0 and 2.4.0))
	Flags *FrameFlags // flags from the frame header
	Data  interface{} // the frame data (either: string, []string, Comment, Picture, PositionInSet or UserDefinedText)

	Text        *string // used by text frames, otherwise nil
	Description *string // used by user-defined text frames, otherwise nil

	UnknownData []byte // used to preserve data for otherwise unknown frame-types, otherwise nil (empty = unknown data of zero length)
}

type Comment struct {
	LanguageCode string
	Description  string
	Comment      string
}

type Picture struct {
	MimeType    string
	PictureType PictureType
	Description string
	Data        []byte
}

type PositionInSet struct {
	Part  int
	Total int
}

type UserDefinedText struct {
	Description string
	Text        string
}
