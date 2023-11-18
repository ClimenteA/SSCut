package types

type FrameInfo struct {
	Frame int
	Keep  bool
	Url   string
}

type VideoMetadata struct {
	Status       States
	Message      string
	TotalSeconds int
	LastSkip     int
	Frames       []FrameInfo
	VideoSrc     string
}

type VideoOutputOptions struct {
	OutputVolume int
}
