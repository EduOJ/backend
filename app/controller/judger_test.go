package controller_test

import (
	"bytes"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestGetTask(t *testing.T) {
	// Not parallel
	assert.Nil(t, base.DB.Delete(models.Run{}, "id > 0").Error)
	problem, user := createProblemForTest(t, "get_task", 1, nil)
	submission := createSubmissionForTest(t, "test_task", 1, &problem, &user, newFileContent(
		"", "code.test_language", b64Encode("balh"),
	), 1)
	var language models.Language
	assert.Nil(t, base.DB.Model(&submission.Runs[0]).Association("TestCase").Error)
	assert.Nil(t, base.DB.Model(&submission).Association("Language").Find(&language))
	req := makeReq(t, "GET", base.Echo.Reverse("judger.getTask"), "", judgerAuthorize)
	httpResp := makeResp(req)
	assert.Equal(t, http.StatusOK, httpResp.StatusCode)
	buf := bytes.Buffer{}
	_, _ = buf.ReadFrom(httpResp.Body)
	var resp response.GetTaskResponse
	mustJsonDecode(buf.String(), &resp)
	t.Log(buf.String())
	assert.Equal(t, response.GetTaskResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			RunID              uint            `json:"run_id"`
			Language           models.Language `json:"language"`
			TestCaseID         uint            `json:"test_case_id"`
			InputFile          string          `json:"input_file"`
			OutputFile         string          `json:"output_file"`
			CodeFile           string          `json:"code_file"`
			TestCaseUpdatedAt  time.Time       `json:"test_case_updated_at"`
			MemoryLimit        uint64          `json:"memory_limit"`
			TimeLimit          uint            `json:"time_limit"`
			CompileEnvironment string          `json:"compile_environment"`
			CompareScriptName  string          `json:"compare_script_name"`
		}{
			submission.Runs[0].ID,
			language,
			submission.Runs[0].TestCaseID,
			resp.Data.InputFile,
			resp.Data.OutputFile,
			resp.Data.CodeFile,
			problem.TestCases[0].UpdatedAt,
			problem.MemoryLimit,
			problem.TimeLimit,
			problem.CompileEnvironment,
			problem.CompareScriptName,
		},
	}, resp)
	req = makeReq(t, "GET", base.Echo.Reverse("judger.getTask"), "", judgerAuthorize)
	httpResp = makeResp(req)
	assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
	jsonEQ(t, response.Response{
		Message: "NOT_FOUND",
		Error:   nil,
		Data:    nil,
	}, httpResp)
}
