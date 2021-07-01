package controller

import (
	"encoding/json"
	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/utils"
	"github.com/EduOJ/backend/database/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)




func CreateComment(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		panic("could not convert my user into type models.User")
	}


	fatherId, err := strconv.Atoi(c.FormValue("father_id"))
	if err != nil {
		panic(err)
	}

	firstId, err := strconv.Atoi(c.FormValue("first_id"))
	if err != nil {
		panic(err)
	}


	content:= c.FormValue("content")

	fatherType:= c.FormValue("father_type")


	if fatherType == "comment" {
		//father is comment
		var father models.Comment
		base.DB.Model(&models.Comment{}).First(&father, uint(fatherId))
		newComment := models.Comment{
			Writer:     user,
			IfDeleted:  false,
			FatherID:   uint(fatherId),
			FatherType: fatherType,
			FirstID:    father.FirstID,
			FirstType:  father.FirstType,
			Content:    content,
		}
		utils.PanicIfDBError(base.DB.Save(&newComment), "could not save comment")

	} else{
		//root node
		newComment := models.Comment{
			Writer:     user,
			IfDeleted:  false,
			FatherType: fatherType,
			FirstID: uint(firstId),
			FirstType:  fatherType,
			Content:    content,
		}
		utils.PanicIfDBError(base.DB.Save(&newComment), "could not save comment")
	}


	return c.JSON(http.StatusCreated, response.CreateCommentResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Id uint
		}{
			uint(1),
		},
	})

}

func GetComment(c echo.Context) error {
	_, ok := c.Get("user").(models.User)
	if !ok {
		panic("could not convert my user into type models.User")
	}

	queryType:= c.FormValue("query_type")
	if queryType == "problem" {
		ids,_ := strconv.Atoi(c.FormValue("query_id"))
		var commentsNoneRoot []models.Comment
		var commentsRoot []models.Comment
		query := base.DB.Model(&models.Comment{}).Preload("Writer")
		query.Order("created_at desc").Find(&commentsRoot," first_type = ? AND first_id = ? AND father_id = ?","problem", uint(ids),0)
		base.DB.Model(&models.Comment{}).Preload("Writer").Order("created_at desc").Find(&commentsNoneRoot," first_type = ? AND first_id = ? AND father_id > ?","problem", uint(ids),0)
		return c.JSON(http.StatusCreated, response.GetCommentResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				ComsRoot []models.Comment
				ComsNoneRoot []models.Comment
			}{
				commentsRoot,
				commentsNoneRoot,
			},
		})
	} else {
		panic(queryType)
		panic("we don't have this function now")
	}

	return c.JSON(http.StatusBadRequest, response.ErrorResp("INVALID_STATUS", nil))

}

func AddReaction(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		panic("could not convert my user into type models.User")
	}
	//userId := user.ID
	commentId,_ := strconv.Atoi(c.FormValue("comment_id"))
	compType := c.FormValue("if_add_action")
	typeId := c.FormValue("type_id")
	query := base.DB.Model(&models.Comment{})
	var comment models.Comment
	query = query.First(&comment, uint(commentId))
	maps := make(map[string] []uint)
	if compType == "1" {
		if comment.Detail == "" {
			maps[typeId] = make([]uint,1)
			maps[typeId][0] = uint(user.ID)
		} else {
			err := json.Unmarshal([]byte(comment.Detail), &maps)
			if err != nil {
				panic(err)
			}
			_, key := maps[typeId]
			if key {
				//operator is not zero
				for _,v:= range maps[typeId] {
					if uint(v) == uint(user.ID){
						return c.JSON(http.StatusBadRequest, response.ErrorResp("can't action twice!", nil))
					}
				}
				maps[typeId] = append(maps[typeId],uint(user.ID))
			} else {
				maps[typeId] = make([]uint,1)
				maps[typeId][0] = uint(user.ID)
			}
		}
		jsonStr, err := json.Marshal(maps)
		if err != nil {
			panic("parse error")
			panic(err)
		}
		comment.Detail = string(jsonStr)
		//query.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&comment)
		query.Save(&comment)
	} else{
		//delete action
		if comment.Detail=="" {
			panic("this should not happen, logic in error!")
		} else {
			json.Unmarshal([]byte(comment.Detail), &maps)
			_, key := maps[typeId]
			if key {
				pos := -1
				for k,v := range maps[typeId] {
					if v == uint(user.ID) {
						pos = k
						break
					}
				}
				if pos == -1{
					panic("should not happen, you have not actived")
				} else {
					maps[typeId] = append( maps[typeId][:pos], maps[typeId][pos+1:]...)
				}
				jsonStr, err := json.Marshal(maps)
				if err != nil {
					panic("parse error")
					panic(err)
				}
				comment.Detail = string(jsonStr)
				query.Save(&comment)
			} else {
				return c.JSON(http.StatusBadRequest, response.ErrorResp("you have not actived!", nil))
			}
		}
	}

	return c.JSON(http.StatusCreated, response.CreateCommentResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Id uint
		}{
			uint(commentId),
		},
	})

}
