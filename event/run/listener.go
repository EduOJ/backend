package submission

import (
	"context"
	"fmt"
	"github.com/leoleoasd/EduOJBackend/base"
)

func NotifyGetSubmissionPoll(r EventArgs) EventRst {
	base.Redis.Publish(context.Background(), fmt.Sprintf("submission_update:%d", r.Submission.ID), nil)
	return EventRst{}
}
