package main

import (
	"context"
	"encoding/json"
	"msg-event/model/event"
	"msg-event/model/response"
	"msg-event/services"
	"runtime/debug"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

func HandleRequest(ctx context.Context, e *event.Msg) (event *response.MsgResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Infof("panic is %v", string(debug.Stack()))
		}
	}()
	s, _ := json.Marshal(e)

	logrus.Infof("event is %s", string(s))

	r, err := services.Serve(ctx, e)
	if err != nil {
		logrus.Errorf("handle err %v", err)
	}
	return r, nil
}

func main() {

	lambda.Start(HandleRequest)

}
