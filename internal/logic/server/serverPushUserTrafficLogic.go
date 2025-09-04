package server

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
	"github.com/perfect-panel/server/pkg/logger"
	"github.com/perfect-panel/server/pkg/tool"
	task "github.com/perfect-panel/server/queue/types"
	"github.com/pkg/errors"
)

//goland:noinspection GoNameStartsWithPackageName
type ServerPushUserTrafficLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewServerPushUserTrafficLogic Push user Traffic
func NewServerPushUserTrafficLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServerPushUserTrafficLogic {
	return &ServerPushUserTrafficLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ServerPushUserTrafficLogic) ServerPushUserTraffic(req *types.ServerPushUserTrafficRequest) error {
	// Find server info
	serverInfo, err := l.svcCtx.NodeModel.FindOneServer(l.ctx, req.ServerId)
	if err != nil {
		l.Errorw("[PushOnlineUsers] FindOne error", logger.Field("error", err))
		return errors.New("server not found")
	}

	// Create traffic task
	var request task.TrafficStatistics
	request.ServerId = serverInfo.Id
	tool.DeepCopy(&request.Logs, req.Traffic)

	// Push traffic task
	val, _ := json.Marshal(request)
	t := asynq.NewTask(task.ForthwithTrafficStatistics, val, asynq.MaxRetry(3))
	info, err := l.svcCtx.Queue.EnqueueContext(l.ctx, t)
	if err != nil {
		l.Errorw("[ServerPushUserTraffic] Push traffic task error", logger.Field("error", err.Error()), logger.Field("task", t))
	} else {
		l.Infow("[ServerPushUserTraffic] Push traffic task success", logger.Field("task", t), logger.Field("info", info))
	}

	// Update server last reported time
	now := time.Now()
	serverInfo.LastReportedAt = &now

	err = l.svcCtx.NodeModel.UpdateServer(l.ctx, serverInfo)
	if err != nil {
		l.Errorw("[ServerPushUserTraffic] UpdateServer error", logger.Field("error", err))
		return nil
	}
	return nil
}
