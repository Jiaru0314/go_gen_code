package service

import (
	"context"
	"go_gen_code/internal/model"
)

type (
	IDocument interface {
		Add(ctx context.Context, in model.AddDocumentInput) (out *model.AddDocumentOutput, err error)
		Update(ctx context.Context, in model.UpdateDocumentInput) (err error)
		Delete(ctx context.Context, id int) (err error)
		Page(ctx context.Context, in model.PageDocumentInput) (out *model.PageDocumentOutput, err error)
	}
)

var (
	localDocument IDocument
)

func Document() IDocument {
	if localDocument == nil {
		panic("implement not found for interface IDocument, forgot register?")
	}
	return localDocument
}

func RegisterDocument(i IDocument) {
	localDocument = i
}
