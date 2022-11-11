package processors

import (
	"lambda-lark-bot/dao"
	"lambda-lark-bot/model/event"
	"lambda-lark-bot/services/api"
	"time"

	"github.com/sirupsen/logrus"
)

type refreshCommentProcessor struct {
}

func GetRefreshCommentProcessor() api.Processor {
	return &refreshCommentProcessor{}
}

func (r refreshCommentProcessor) Process(e *event.Msg) error {
	logrus.Infof("Ready to refresh comment...")
	err := RefreshComments()
	if err != nil {
		logrus.Errorf("Refresh comment failed %s", err)
	}
	logrus.Infof("Refresh comments complated")
	return err
}

func RefreshComments() error {
	// get all un-closed cases
	cs, err := dao.GetProcessingCases()
	if err != nil {
		logrus.Errorf("refresh failed to get cases %s", err)
		return err
	}
	// for loop get latest comments

	for _, c := range cs {
		comments, err := dao.GetCaseComments(c, c.LastCommentTime)
		if err != nil {
			logrus.Errorf("failed to get all comments %s", err)
			return err
		}
		awscase, err := dao.GetAWSCase(c)
		if err != nil {
			logrus.Errorf("failed to get aws case %s", err)
			return err
		}
		if *awscase.Cases[0].Status == "resolved" {
			c.Status = dao.STATUS_CLOSE
		}
		c.Comments = comments
	}
	for _, c := range cs {
		c.LastCommentTime = time.Now()
		_, err := dao.UpsertCase(c)
		if err != nil {
			logrus.Errorf("update case last comment time failed %s", err)
			return err
		}
		// send all comments to channel
		err = dao.SendMsg(c.ChannelID, c.UserID, dao.FormatComments(c.Comments))
		if err != nil {
			logrus.Errorf("failed to send comments %s", err)
			return err
		}
	}
	return nil
}
