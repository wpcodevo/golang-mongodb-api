package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-mongodb/models"
	"github.com/wpcodevo/golang-mongodb/services"
)

type PostController struct {
	postService services.PostService
}

func NewPostController(postService services.PostService) PostController {
	return PostController{postService}
}

func (pc *PostController) CreatePost(ctx *gin.Context) {
	var post *models.CreatePostRequest

	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newPost, err := pc.postService.CreatePost(post)

	if err != nil {
		if strings.Contains(err.Error(), "title already exists") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newPost})
}

func (pc *PostController) UpdatePost(ctx *gin.Context) {
	postId := ctx.Param("postId")

	var post *models.UpdatePost
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updatedPost, err := pc.postService.UpdatePost(postId, post)
	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedPost})
}

func (pc *PostController) FindPostById(ctx *gin.Context) {
	postId := ctx.Param("postId")

	post, err := pc.postService.FindPostById(postId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": post})
}

func (pc *PostController) FindPosts(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	posts, err := pc.postService.FindPosts(intPage, intLimit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(posts), "data": posts})
}

func (pc *PostController) DeletePost(ctx *gin.Context) {
	postId := ctx.Param("postId")

	err := pc.postService.DeletePost(postId)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
