package server

import (
	"context"

	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
	"github.com/perfect-panel/server/pkg/logger"
	"github.com/perfect-panel/server/pkg/xerr"
	"github.com/pkg/errors"
)

type DeleteNodeLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteNodeLogic Delete Node
func NewDeleteNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNodeLogic {
	return &DeleteNodeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteNodeLogic) DeleteNode(req *types.DeleteNodeRequest) error {
	err := l.svcCtx.NodeModel.DeleteNode(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[DeleteNode] Delete Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "[DeleteNode] Delete Database Error")
	}
	return nil
}
