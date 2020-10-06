package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func getUrlStringPointer(rawUrl string, paras map[string]string) *string {
	thisURL, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		panic(err)
	}
	q, err := url.ParseQuery(thisURL.RawQuery)
	if err != nil {
		panic(err)
	}
	for key := range paras {
		q.Add(key, paras[key])
	}
	thisURL.RawQuery = q.Encode()
	str := thisURL.String()
	return &str
}

func TestFindUser(t *testing.T) {
	t.Run("findUserNonExist", func(t *testing.T) {
		user, err := FindUser("test_find_user_non_exist")
		assert.Nil(t, user)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
	t.Run("findUserSuccessWithId", func(t *testing.T) {
		user := models.User{
			Username: "test_find_user_id_username",
			Nickname: "test_find_user_id_nickname",
			Email:    "test_find_user_id@mail.com",
			Password: "test_find_user_id_password",
		}
		assert.Nil(t, base.DB.Create(&user).Error)
		foundUser, err := FindUser(fmt.Sprintf("%d", user.ID))
		if foundUser != nil {
			assert.Equal(t, user, *foundUser)
		}
		assert.Nil(t, err)
	})
	t.Run("findUserSuccessWithUsername", func(t *testing.T) {
		user := models.User{
			Username: "test_find_user_name_username",
			Nickname: "test_find_user_name_nickname",
			Email:    "test_find_user_name@mail.com",
			Password: "test_find_user_name_password",
		}
		assert.Nil(t, base.DB.Create(&user).Error)
		foundUser, err := FindUser(user.Username)
		if foundUser != nil {
			assert.Equal(t, user, *foundUser)
		}
		assert.Nil(t, err)
	})
}

type TestObject struct {
	ID               uint   `gorm:"private_key" json:"id"`
	Name             string `json:"name"`
	StringAttribute  string `json:"string_attribute"`
	IntegerAttribute int    `json:"integer_attribute"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func limitQuery(query *gorm.DB, min, max int) *gorm.DB {
	return query.Where("integer_attribute >= ? and integer_attribute <= ?", min, max)
}

func getQuery(t *testing.T, search string) *gorm.DB {
	query := base.DB.Model(&TestObject{}).Where("name like ?", "%"+search+"%")
	assert.Nil(t, query.Error)
	return query
}

func TestPaginator(t *testing.T) {
	testObjects := make([]TestObject, 25)
	for i := range testObjects {
		testObjects[i] = TestObject{
			Name:             fmt.Sprintf("test_paginator_object_%d", i),
			StringAttribute:  fmt.Sprintf("tpo%d", i),
			IntegerAttribute: i,
		}
		assert.Nil(t, base.DB.Create(&testObjects[i]).Error)
	}

	requestURL, err := url.Parse("http://test.paginator.request.url/testing/path?test_para1=tp1&test_para2=tp2")
	assert.Nil(t, err)

	tests := []struct {
		name    string
		limit   int
		offset  int
		query   *gorm.DB
		output  []TestObject
		total   int
		prevUrl *string
		nextUrl *string
		err     error
	}{
		{
			name:    "Default10",
			limit:   0,
			offset:  0,
			query:   limitQuery(getQuery(t, "test_paginator_object_"), 0, 9),
			output:  testObjects[:10],
			total:   10,
			prevUrl: nil,
			nextUrl: nil,
			err:     nil,
		},
		{
			name:    "DefaultLimit",
			limit:   0,
			offset:  0,
			query:   getQuery(t, "test_paginator_object_"),
			output:  testObjects[:20],
			total:   25,
			prevUrl: nil,
			nextUrl: getUrlStringPointer(requestURL.String(), map[string]string{
				"limit":  "20",
				"offset": "20",
			}),
			err: nil,
		},
		{
			name:    "SingleFirst",
			limit:   1,
			offset:  0,
			query:   getQuery(t, "test_paginator_object_"),
			output:  testObjects[:1],
			total:   25,
			prevUrl: nil,
			nextUrl: getUrlStringPointer(requestURL.String(), map[string]string{
				"limit":  "1",
				"offset": "1",
			}),
			err: nil,
		},
		{
			name:   "SingleMiddle",
			limit:  1,
			offset: 12,
			query:  getQuery(t, "test_paginator_object_"),
			output: testObjects[12:13],
			total:  25,
			prevUrl: getUrlStringPointer(requestURL.String(), map[string]string{
				"limit":  "1",
				"offset": "11",
			}),
			nextUrl: getUrlStringPointer(requestURL.String(), map[string]string{
				"limit":  "1",
				"offset": "13",
			}),
			err: nil,
		},
		{
			name:   "SingleLast",
			limit:  1,
			offset: 24,
			query:  getQuery(t, "test_paginator_object_"),
			output: testObjects[24:],
			total:  25,
			prevUrl: getUrlStringPointer(requestURL.String(), map[string]string{
				"limit":  "1",
				"offset": "23",
			}),
			nextUrl: nil,
			err:     nil,
		},
		{
			name:    "MultipleHead",
			limit:   5,
			offset:  0,
			query:   getQuery(t, "test_paginator_object_"),
			output:  testObjects[:5],
			total:   25,
			prevUrl: nil,
			nextUrl: getUrlStringPointer(requestURL.String(), map[string]string{
				"limit":  "5",
				"offset": "5",
			}),
			err: nil,
		},
		{
			name:    "MultipleFront",
			limit:   5,
			offset:  3,
			query:   getQuery(t, "test_paginator_object_"),
			output:  testObjects[3:8],
			total:   25,
			prevUrl: nil,
			nextUrl: getUrlStringPointer(requestURL.String(), map[string]string{
				"limit":  "5",
				"offset": "8",
			}),
			err: nil,
		},
		{
			name:   "MultipleMiddle",
			limit:  5,
			offset: 10,
			query:  getQuery(t, "test_paginator_object_"),
			output: testObjects[10:15],
			total:  25,
			prevUrl: getUrlStringPointer(requestURL.String(), map[string]string{
				"limit":  "5",
				"offset": "5",
			}),
			nextUrl: getUrlStringPointer(requestURL.String(), map[string]string{
				"limit":  "5",
				"offset": "15",
			}),
			err: nil,
		},
		{
			name:   "MultipleBack",
			limit:  5,
			offset: 18,
			query:  getQuery(t, "test_paginator_object_"),
			output: testObjects[18:23],
			total:  25,
			prevUrl: getUrlStringPointer(requestURL.String(), map[string]string{
				"limit":  "5",
				"offset": "13",
			}),
			nextUrl: getUrlStringPointer(requestURL.String(), map[string]string{
				"limit":  "5",
				"offset": "23",
			}),
			err: nil,
		},
		{
			name:   "MultipleTail",
			limit:  5,
			offset: 20,
			query:  getQuery(t, "test_paginator_object_"),
			output: testObjects[20:25],
			total:  25,
			prevUrl: getUrlStringPointer(requestURL.String(), map[string]string{
				"limit":  "5",
				"offset": "15",
			}),
			nextUrl: nil,
			err:     nil,
		},
		{
			name:    "Long",
			limit:   20,
			offset:  3,
			query:   getQuery(t, "test_paginator_object_"),
			output:  testObjects[3:23],
			total:   25,
			prevUrl: nil,
			nextUrl: getUrlStringPointer(requestURL.String(), map[string]string{
				"limit":  "20",
				"offset": "23",
			}),
			err: nil,
		},
		{
			name:    "Full",
			limit:   25,
			offset:  0,
			query:   getQuery(t, "test_paginator_object_"),
			output:  testObjects,
			total:   25,
			prevUrl: nil,
			nextUrl: nil,
			err:     nil,
		},
		{
			name:    "Oversize",
			limit:   30,
			offset:  0,
			query:   getQuery(t, "test_paginator_object_"),
			output:  testObjects,
			total:   25,
			prevUrl: nil,
			nextUrl: nil,
			err:     nil,
		},
		{
			name:    "Empty",
			limit:   0,
			offset:  0,
			query:   limitQuery(getQuery(t, "test_paginator_object_"), 1, 0),
			output:  nil,
			total:   0,
			prevUrl: nil,
			nextUrl: nil,
			err:     nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run("testPaginator"+test.name, func(t *testing.T) {
			t.Parallel()
			var output []TestObject
			reqURL, err := url.Parse("http://test.paginator.request.url/testing/path?test_para1=tp1&test_para2=tp2")
			assert.Nil(t, err)
			total, prevUrl, nextUrl, err := Paginator(test.query, test.limit, test.offset, reqURL, &output)
			assert.Equal(t, test.err, err)
			if err == nil {
				assert.Equal(t, test.total, total)
				if test.prevUrl == nil {
					assert.Nil(t, prevUrl)
				} else {
					assert.Equal(t, *test.prevUrl, *prevUrl)
				}
				if test.nextUrl == nil {
					assert.Nil(t, nextUrl)
				} else {
					assert.Equal(t, *test.nextUrl, *nextUrl)
				}
				assert.Equal(t, len(test.output), len(output))
				for i := range output {
					output[i].UpdatedAt = test.output[i].UpdatedAt
					assert.Equal(t, test.output[i], output[i])
				}
			}
		})
	}
}

func TestSorter(t *testing.T) {
	testObject1 := TestObject{
		Name:             "test_sorter_object_1",
		StringAttribute:  "C_tso_1",
		IntegerAttribute: 2,
	}
	testObject2 := TestObject{
		Name:             "test_sorter_object_2",
		StringAttribute:  "A_tso_2",
		IntegerAttribute: 4,
	}
	testObject3 := TestObject{
		Name:             "test_sorter_object_3",
		StringAttribute:  "D_tso_3",
		IntegerAttribute: 1,
	}
	testObject4 := TestObject{
		Name:             "test_sorter_object_4",
		StringAttribute:  "B_tso_4",
		IntegerAttribute: 3,
	}
	assert.Nil(t, base.DB.Create(&testObject1).Error)
	assert.Nil(t, base.DB.Create(&testObject2).Error)
	assert.Nil(t, base.DB.Create(&testObject3).Error)
	assert.Nil(t, base.DB.Create(&testObject4).Error)

	query := base.DB.Model(&TestObject{}).Where("name like ?", "%test_sorter_object_%")
	assert.Nil(t, query.Error)

	invalidError := HttpError{
		Code:    400,
		Message: "INVALID_ORDER",
		Err:     nil,
	}

	tests := []struct {
		name          string
		query         *gorm.DB
		orderBy       string
		columnAllowed []string
		output        []TestObject
		err           error
	}{
		{
			name:          "Default",
			query:         query,
			orderBy:       "",
			columnAllowed: []string{},
			output: []TestObject{
				testObject1,
				testObject2,
				testObject3,
				testObject4,
			},
			err: nil,
		},
		{
			name:    "SuccessIdAsc",
			query:   query,
			orderBy: "id.ASC",
			columnAllowed: []string{
				"id",
			},
			output: []TestObject{
				testObject1,
				testObject2,
				testObject3,
				testObject4,
			},
			err: nil,
		},
		{
			name:    "SuccessIdDesc",
			query:   query,
			orderBy: "id.DESC",
			columnAllowed: []string{
				"id",
				"name",
			},
			output: []TestObject{
				testObject4,
				testObject3,
				testObject2,
				testObject1,
			},
			err: nil,
		},
		{
			name:    "SuccessNameAsc",
			query:   query,
			orderBy: "name.ASC",
			columnAllowed: []string{
				"name",
			},
			output: []TestObject{
				testObject1,
				testObject2,
				testObject3,
				testObject4,
			},
			err: nil,
		},
		{
			name:    "SuccessStringAttributeDesc",
			query:   query,
			orderBy: "string_attribute.DESC",
			columnAllowed: []string{
				"id",
				"name",
				"string_attribute",
				"integer_attribute",
			},
			output: []TestObject{
				testObject3,
				testObject1,
				testObject4,
				testObject2,
			},
			err: nil,
		},
		{
			name:    "SuccessIntegerAttributeAsc",
			query:   query,
			orderBy: "integer_attribute.ASC",
			columnAllowed: []string{
				"integer_attribute",
			},
			output: []TestObject{
				testObject3,
				testObject1,
				testObject4,
				testObject2,
			},
			err: nil,
		},
		{
			name:    "SuccessIntegerAttributeDesc",
			query:   query,
			orderBy: "integer_attribute.DESC",
			columnAllowed: []string{
				"id",
				"name",
				"string_attribute",
				"integer_attribute",
			},
			output: []TestObject{
				testObject2,
				testObject4,
				testObject1,
				testObject3,
			},
			err: nil,
		},
		{
			name:    "InvalidOrder",
			query:   query,
			orderBy: "invalid_order",
			columnAllowed: []string{
				"id",
				"name",
				"string_attribute",
				"integer_attribute",
			},
			output: nil,
			err:    invalidError,
		},
		{
			name:    "ColumnNotAllowed",
			query:   query,
			orderBy: "id.ASC",
			columnAllowed: []string{
				"name",
				"string_attribute",
				"integer_attribute",
			},
			output: nil,
			err:    invalidError,
		},
		{
			name:    "OrderLowercase",
			query:   query,
			orderBy: "id.desc",
			columnAllowed: []string{
				"id",
				"name",
				"string_attribute",
				"integer_attribute",
			},
			output: nil,
			err:    invalidError,
		},
		{
			name:    "OrderNonExist",
			query:   query,
			orderBy: "id.NON_EXIST_ORDER",
			columnAllowed: []string{
				"id",
				"name",
				"string_attribute",
				"integer_attribute",
			},
			output: nil,
			err:    invalidError,
		},
	}

	for _, test := range tests {
		test := test
		t.Run("testSorter"+test.name, func(t *testing.T) {
			t.Parallel()
			var output []TestObject
			resultQuery, err := Sorter(test.query, test.orderBy, test.columnAllowed...)
			assert.Equal(t, test.err, err)
			if err == nil {
				assert.Nil(t, resultQuery.Find(&output).Error)
				assert.Equal(t, len(test.output), len(output))
				for i := range output {
					output[i].UpdatedAt = test.output[i].UpdatedAt
					assert.Equal(t, test.output[i], output[i])
				}
			}
		})
	}
}
