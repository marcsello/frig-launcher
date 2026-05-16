package loader

type AssetType int

const (
	SoundAsset AssetType = iota
	ImageAsset
	FontAsset
)

func (t AssetType) String() string {
	switch t {
	case SoundAsset:
		return "sound"
	case ImageAsset:
		return "image"
	case FontAsset:
		return "font"
	default:
		return ""
	}
}
