package collection

import collectionSvc "github.com/board-box/backend/internal/service/collection"

type CreateCollectionRequest struct {
	Name string `json:"name" example:"my collection"`
}

type UpdateCollectionRequest struct {
	Name   string `json:"name" example:"my collection"`
	Pinned bool   `json:"pinned" example:"false"`
}

func convertCreateReqToDTO(req CreateCollectionRequest) collectionSvc.Collection {
	return collectionSvc.Collection{
		Name: req.Name,
	}
}

func convertUpdateReqToDTO(req UpdateCollectionRequest) collectionSvc.Collection {
	return collectionSvc.Collection{
		Name:   req.Name,
		Pinned: req.Pinned,
	}
}
