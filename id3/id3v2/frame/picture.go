package frame

type Picture struct {
	MimeType    string
	PictureType PictureType
	Description string
	Data        []byte
}
