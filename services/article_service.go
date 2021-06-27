/*
   Copyright (c) [2021] IT.SOS
   kn is licensed under Mulan PSL v2.
   You can use this software according to the terms and conditions of the Mulan PSL v2.
   You may obtain a copy of Mulan PSL v2 at:
            http://license.coscl.org.cn/MulanPSL2
   THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
   See the Mulan PSL v2 for more details.
*/

package services

import (
	"gitee.com/itsos/golibs/utils"
	"gitee.com/itsos/golibs/utils/validate"
	"gitee.com/itsos/studynotes/caches"
	"gitee.com/itsos/studynotes/datamodels"
	"gitee.com/itsos/studynotes/models/vo"
	"gitee.com/itsos/studynotes/repositories"
)

type ArticleService interface {
	// GetRank 获取前访问前50的文章列表
	GetRank() []vo.ArticleAccessTimesVO
	// GetListPage 获取最新文章列表
	GetListPage(isLogin bool, page int, size int) []vo.ArticleVO
	// GetContent 获取文章详情
	GetContent(title string) vo.ArticleContentVO
}

var SArticle ArticleService = &articleService{repositories.RArticle, caches.CAccessTimes, SCategory}

type articleService struct {
	article  repositories.ArticleRepository
	times    caches.AccessTimes
	category CategoryService
}

func (a articleService) GetRank() []vo.ArticleAccessTimesVO {
	a.times.Rank(10)
	articleVO := make([]vo.ArticleAccessTimesVO, 0)
	// 获取访问量的前50条
	rank := a.times.Rank(50)
	if len(rank) > 0 {
		article := a.article.SelectManyByIds(rank)
		if len(article) > 0 {
			for _, v := range article {
				articleVO = append(articleVO, vo.ArticleAccessTimesVO{Title: v.Title, AccessTimes: a.times.Id(v.Id).Get()})
			}
		}
	}
	return articleVO
}

func (a articleService) GetListPage(isLogin bool, page int, size int) []vo.ArticleVO {
	page = validate.IntRange(page, 0, 100000)
	size = validate.IntRange(size, 0, 100000)
	var state = []uint8{repositories.IsStatePublic}
	if isLogin {
		state = append(state, repositories.IsStatePrivate)
	}
	offset := (page - 1) * size
	article := a.article.SelectMany(state, offset, size)
	return getArticles(article)
}

func (a articleService) GetContent(title string) vo.ArticleContentVO {
	// 获取内容详情
	//a.repo.Content(id)
	// 获取专题列表
	panic("implement me")
}

func getArticles(article []datamodels.Article) []vo.ArticleVO {
	articleVO := make([]vo.ArticleVO, 0)
	if len(article) > 0 {
		for _, v := range article {
			articleVO = append(articleVO, vo.ArticleVO{Article: v, Duration: utils.TimeDuration(v.Utime)})
		}
	}
	return articleVO
}
