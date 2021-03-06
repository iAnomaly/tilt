package webview

// TODO(dmiller): delete these tests once StateToWebView is deleted
import (
	"time"

	"github.com/windmilleng/tilt/internal/container"
	"github.com/windmilleng/tilt/internal/dockercompose"
	"github.com/windmilleng/tilt/pkg/model"
	"github.com/windmilleng/tilt/pkg/model/logstore"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"

	proto_webview "github.com/windmilleng/tilt/pkg/webview"
)

func NewProtoDCResourceInfo(configPaths []string, status dockercompose.Status, cID container.ID, startTime time.Time) (*proto_webview.DCResourceInfo, error) {
	start, err := timeToProto(startTime)
	if err != nil {
		return nil, err
	}
	return &proto_webview.DCResourceInfo{
		ConfigPaths:     configPaths,
		ContainerStatus: string(status),
		ContainerID:     string(cID),
		StartTime:       start,
	}, nil
}

func timeToProto(t time.Time) (*timestamp.Timestamp, error) {
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func ToProtoBuildRecord(br model.BuildRecord, logStore *logstore.LogStore) (*proto_webview.BuildRecord, error) {
	e := ""
	if br.Error != nil {
		e = br.Error.Error()
	}

	start, err := timeToProto(br.StartTime)
	if err != nil {
		return nil, err
	}
	finish, err := timeToProto(br.FinishTime)
	if err != nil {
		return nil, err
	}

	warnings := []string{}
	if br.SpanID != "" {
		warnings = logStore.Warnings(br.SpanID)
	}

	return &proto_webview.BuildRecord{
		Edits: br.Edits,
		Error: e,
		// TODO(nick): Remove this, and compute it client-side.
		Warnings:       warnings,
		StartTime:      start,
		FinishTime:     finish,
		IsCrashRebuild: br.Reason.IsCrashOnly(),
		SpanId:         string(br.SpanID),
	}, nil
}

func ToProtoBuildRecords(brs []model.BuildRecord, logStore *logstore.LogStore) ([]*proto_webview.BuildRecord, error) {
	ret := make([]*proto_webview.BuildRecord, len(brs))
	for i, br := range brs {
		r, err := ToProtoBuildRecord(br, logStore)
		if err != nil {
			return nil, err
		}
		ret[i] = r
	}
	return ret, nil
}
