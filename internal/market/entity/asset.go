package entity

type Asset struct {
	ID           string
	Name         string
	MarketVolume int
}

func NewAsset(assetID string, name string, marketVolume int) *Asset {
	return &Asset{
		ID:           assetID,
		Name:         name,
		MarketVolume: marketVolume,
	}
}
