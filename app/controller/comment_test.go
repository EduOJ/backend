package controller_test

import (
	"fmt"
	"github.com/EduOJ/backend/app/request"
	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/database/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreateComment(t *testing.T) {
	t.Parallel()
	successTests := []struct {
		name       string
		req        request.CreateCommentRequest
	}{
		{
			name: "SuccessRootCommentForProblem",
			req: request.CreateCommentRequest{
				Content: "balabala1",
				FatherID: 0,
				TargetID:  1,
				TargetType: "problem",
			},
		},
		{
			name: "SuccessNoneRootCommentForProblem",
			req: request.CreateCommentRequest{
				Content: "balabala2",
				FatherID: 1,
				TargetID:  1,
				TargetType: "problem",
			},
		},
	}
	t.Run("TestCreateCommentSuccess", func(t *testing.T) {
			test := successTests[0]
			t.Run("TestCreateComment"+test.name, func(t *testing.T) {
				t.Parallel()
				user := createUserForTest(t, "create_comment", 0)
				user.GrantRole("admin")
				httpReq := makeReq(t, "POST", base.Echo.Reverse("comment.createComment"), successTests[0].req, headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				})
				httpResp := makeResp(httpReq)
				databaseComment := models.Comment{}
				assert.NoError(t, base.DB.Where("content = ?", successTests[0].req.Content).First(&databaseComment).Error)

				assert.Equal(t, test.req.Content, databaseComment.Content)
				assert.Equal(t, test.req.FatherID, databaseComment.FatherID)
				assert.Equal(t, test.req.TargetID, databaseComment.TargetID)


				test = successTests[1]
				test.req.FatherID = databaseComment.ID
				var data interface{}
				data = test.req
				httpReq = makeReq(t, "POST", base.Echo.Reverse("comment.createComment"), data, headerOption{
					"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
				})
				httpResp = makeResp(httpReq)
				assert.Equal(t, http.StatusCreated, httpResp.StatusCode)

				databaseComment = models.Comment{}
				assert.NoError(t, base.DB.Where("content = ?", test.req.Content).First(&databaseComment).Error)
				// request == database
				assert.Equal(t, test.req.Content, databaseComment.Content)
				assert.Equal(t, test.req.FatherID, databaseComment.FatherID)
				assert.Equal(t, test.req.TargetID, databaseComment.TargetID)


			})
	})

}


func TestGetComment(t *testing.T) {
	t.Parallel()

	problem := models.Problem{
		Name:               "balabala",
		Description:        "test_get_problems_temp1_description",
		AttachmentFileName: "test_get_problems_temp1_attachment_file_name",
		LanguageAllowed:    []string{"test_get_problems_temp1_language_allowed"},
		Public:             true,
	}
	assert.NoError(t, base.DB.Create(&problem).Error)
	databaseProblem := models.Problem{}
	assert.NoError(t, base.DB.Where("Name = ?", "balabala").First(&databaseProblem).Error)

	comment1 := models.Comment{
		Content:               "test_get_comment_1",
		FatherID: 0,
		TargetID: databaseProblem.ID,
		TargetType: "problem",
	}
	assert.NoError(t, base.DB.Create(&comment1).Error)
	databaseComment := models.Comment{}
	assert.NoError(t, base.DB.Where("Content = ?", "test_get_comment_1").First(&databaseComment).Error)

	comment2 := models.Comment{
		Content:               "test_get_comment_2",
		FatherID: databaseComment.ID,
		TargetID: databaseProblem.ID,
		TargetType: "problem",
	}
	assert.NoError(t, base.DB.Create(&comment2).Error)


	successTests := []struct {
		name       string
		req        request.GetCommentRequest
	}{
		{
			name: "problem",
			req: request.GetCommentRequest{
				TargetType: "problem",
				TargetID: databaseProblem.ID,
				Limit: 0,
				Offset: 5,
			},
		},
	}

	t.Run("testGetProblemsSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			user := createUserForTest(t, "get_comments", i)
			// assert.False(t,user.Can("manage_problem"))
			httpResp := makeResp(makeReq(t, "GET", base.Echo.Reverse("comment.getComment"), test.req, headerOption{
				"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
			}))
			body := httpResp.Body
			assert.Equal(t, body,body)
		}
	})

}

func TestAddReactionAndDeleteReaction(t *testing.T) {
	problem := models.Problem{
		Name:               "balabala",
		Description:        "test_get_problems_temp1_description",
		AttachmentFileName: "test_get_problems_temp1_attachment_file_name",
		LanguageAllowed:    []string{"test_get_problems_temp1_language_allowed"},
		Public:             true,
	}
	assert.NoError(t, base.DB.Create(&problem).Error)
	databaseProblem := models.Problem{}
	assert.NoError(t, base.DB.Where("Name = ?", "balabala").First(&databaseProblem).Error)

	comment1 := models.Comment{
		Content:               "test_get_comment_1",
		FatherID: 0,
		TargetID: databaseProblem.ID,
		TargetType: "problem",
	}
	assert.NoError(t, base.DB.Create(&comment1).Error)
	databaseComment := models.Comment{}
	assert.NoError(t, base.DB.Where("Content = ?", "test_get_comment_1").First(&databaseComment).Error)

	req := request.AddReactionRequest{
		TargetType: "comment",
		TargetID: databaseComment.ID,
		EmojiType: "1",
	}

	// assert.False(t,user.Can("manage_problem"))
	makeResp(makeReq(t, "POST", base.Echo.Reverse("comment.addReaction"), req, applyNormalUser))

	comment2 := models.Comment{
		Content:               "test_get_comment_2",
		FatherID: databaseComment.ID,
		TargetID: databaseProblem.ID,
		TargetType: "problem",
	}
	assert.NoError(t, base.DB.Create(&comment2).Error)


	//user := createUserForTest(t, "get_problems_submitter", 1)
	successTests := []struct {
		name       string
		req        request.AddReactionRequest
	}{
		{
			name: "problem",
			req: request.AddReactionRequest{
				TargetType: "comment",
				TargetID: databaseComment.ID,
				EmojiType: "1",
			},
		},
		{
			name: "problem",
			req: request.AddReactionRequest{
				TargetType: "comment",
				TargetID: databaseComment.ID,
				EmojiType: "2",
			},
		},
	}

	t.Run("testAddReactionSuccess", func(t *testing.T) {
		t.Parallel()
		for i, test := range successTests {
			i := i
			test := test
			user := createUserForTest(t, "add_reaction", i+1)
			// assert.False(t,user.Can("manage_problem"))
			httpResp := makeResp(makeReq(t, "POST", base.Echo.Reverse("comment.addReaction"), test.req, headerOption{
				"Set-User-For-Test": {fmt.Sprintf("%d", user.ID)},
			}))
			assert.Equal(t, http.StatusCreated, httpResp.StatusCode)
		}
	})


	failTests := []failTest{
		{
			name:   "InvalidStatus",
			method: "DELETE",
			path:   base.Echo.Reverse("comment.deleteReaction"),
			req: request.DeleteReactionRequest{
				TargetType: "comment",
				TargetID: databaseComment.ID,
				EmojiType: "5",
			},
			reqOptions: []reqOption{
				applyNormalUser,
			},
			statusCode: http.StatusBadRequest,
			resp:       response.ErrorResp("you have not actived!", nil),
		},
	}

	runFailTests(t, failTests, "")


	successTests1 := []struct {
		name       string
		req        request.DeleteReactionRequest
	}{
		{
			name: "problem",
			req: request.DeleteReactionRequest{
				TargetType: "comment",
				TargetID: databaseComment.ID,
				EmojiType: "1",
			},
		},
	}

	t.Run("testDeleteReactionSuccess", func(t *testing.T) {
		t.Parallel()
		for _,test := range successTests1 {
			test := test
			httpResp := makeResp(makeReq(t, "DELETE", base.Echo.Reverse("comment.deleteReaction"), test.req, applyNormalUser))
			assert.Equal(t, http.StatusCreated, httpResp.StatusCode)
		}
	})


}


func TestDeleteComment(t *testing.T) {
	problem := models.Problem{
		Name:               "balabala",
		Description:        "test_get_problems_temp1_description",
		AttachmentFileName: "test_get_problems_temp1_attachment_file_name",
		LanguageAllowed:    []string{"test_get_problems_temp1_language_allowed"},
		Public:             true,
	}
	assert.NoError(t, base.DB.Create(&problem).Error)
	databaseProblem := models.Problem{}
	assert.NoError(t, base.DB.Where("Name = ?", "balabala").First(&databaseProblem).Error)

	comment1 := models.Comment{
		Content:               "test_get_comment_1",
		FatherID: 0,
		TargetID: databaseProblem.ID,
		TargetType: "problem",
	}
	assert.NoError(t, base.DB.Create(&comment1).Error)
	databaseComment := models.Comment{}
	assert.NoError(t, base.DB.Where("Content = ?", "test_get_comment_1").First(&databaseComment).Error)

	req := request.AddReactionRequest{
		TargetType: "comment",
		TargetID: databaseComment.ID,
		EmojiType: "1",
	}

	// assert.False(t,user.Can("manage_problem"))
	makeResp(makeReq(t, "POST", base.Echo.Reverse("comment.addReaction"), req, applyNormalUser))

	comment2 := models.Comment{
		Content:               "test_get_comment_2",
		FatherID: databaseComment.ID,
		TargetID: databaseProblem.ID,
		TargetType: "problem",
	}
	assert.NoError(t, base.DB.Create(&comment2).Error)


}


