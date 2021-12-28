package api

import (
	"errors"
	"log"
	"main/db_model"
	"main/pkg/user_session"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type NewCommentRequest struct {
	ShareId   uint   `json:"share_id"`
	SessionId string `json:"session_id"`
	Content   string `json:"content"`
}

type GetCommentRequest struct {
	ShareId   int    `json:"share_id"`
	SessionId string `json:"session_id"`
}

// CommentItem 响应
type CommentItem struct {
	CommentId      uint      `json:"comment_id"`
	Username       string    `json:"username"`
	UserAvatar     string    `json:"useravatar"`
	Content        string    `json:"content"`
	LikeNum        uint      `json:"like_num"`
	CreatedAt      time.Time `json:"created_at"`
	HasCommentLike bool      `json:"has_comment_like"`
}

type DeleteCommentRequest struct {
	CommentId uint   `json:"comment_id"`
	SessionId string `json:"session_id"`
}

func NewComment(c *gin.Context) {
	var tmp NewCommentRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("BindJSON error: %v", err)
	}
	// session id 为空 没有登录
	if tmp.SessionId == "" {
		log.Println("not have sessionId")
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "没有登录",
		})
	} else {
		//通过sessionID得到userID再进行下一步操作
		userID, _ := user_session.GetUserID(tmp.SessionId)
		log.Println("userID is", userID)
		if userID == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "没有登录",
			})
			return
		}
		if tmp.Content == "" {
			log.Println("content is null")
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "评论为空",
			})
			return
		}

		var comment db_model.Comment
		comment.Share_id = tmp.ShareId
		comment.User_id = userID
		comment.Content = tmp.Content
		// 插入 commnent 表
		db_model.Db.Create(&comment)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "创建成功",
		})

	}

}

func GetComment(c *gin.Context) {
	var tmp GetCommentRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("BindJSON error: %v", err)
	}
	// 获得用户的id
	userID, _ := user_session.GetUserID(tmp.SessionId)
	log.Println(tmp.ShareId)
	// 根据share id 查找comment
	var comments []db_model.Comment
	//更新时间排序
	if err := db_model.Db.Where("Share_id = ?", tmp.ShareId).Order("Created_At").Find(&comments).Error; err != nil {
		log.Fatalf("find share: %v", err)
	}
	var data []CommentItem
	for _, comment := range comments {
		// user id 找到user name
		var user db_model.Userinfo
		if err := db_model.Db.Where("ID = ?", comment.User_id).First(&user).Error; err != nil {
			log.Fatalf("find user error: %v", err)
		}
		item := CommentItem{
			CommentId:  comment.ID,
			Username:   user.Username,
			UserAvatar: user.User_pic,
			Content:    comment.Content,
			LikeNum:    comment.Like_num,
			CreatedAt:  comment.CreatedAt,
		}
		// 查找是否存在 item在 comment like 表
		var commentlike db_model.CommentLike
		err := db_model.Db.Where("User_id = ? AND Comment_id = ?", userID, comment.ID).First(&commentlike).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			item.HasCommentLike = false
		} else {
			item.HasCommentLike = true
		}

		data = append(data, item)
		// append
	}
	c.JSON(http.StatusOK, gin.H{
		"code":          0,
		"comment_items": data,
	})

}

func DeleteComment(c *gin.Context) {
	// 从请求中把数据取出来
	var tmp DeleteCommentRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("BindJSON error: %v", err)
	}

	if tmp.SessionId == "" {
		log.Println("not have sessionId")
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "没有登录",
		})
	} else {
		// session id 得到 user id
		userID, _ := user_session.GetUserID(tmp.SessionId)
		log.Println(userID)
		// 验证用户id
		var comment db_model.Comment
		if err := db_model.Db.Where("ID   = ?", tmp.CommentId).First(&comment).Error; err != nil {
			log.Printf("Find comment error: %v\n", err)
		}
		if userID != comment.User_id {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "无权限删除",
			})
		} else {
			// 删除 对应comment like
			var commentlikes []db_model.CommentLike
			if err := db_model.Db.Where("Comment_id = ?", comment.ID).Find(&commentlikes).Error; err != nil {
				log.Fatalf("find commentlikes: %v", err)
			}
			for _, commentlike := range commentlikes {
				db_model.Db.Delete(&commentlike)
			}
			// 删除 comment
			db_model.Db.Delete(&comment)
			// 返回response
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "删除成功",
			})
		}
	}

}

type CommentLikeRequest struct {
	CommentId uint   `json:"comment_id"`
	SessionId string `json:"session_id"`
}

func AlterCommentLike(c *gin.Context) {
	var tmp CommentLikeRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("CommentLike BindJSON error: %v", err)
	}
	//通过sessionID得到userID再进行下一步操作
	if tmp.SessionId == "" {
		log.Println("not have sessionId")
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "没有登录",
		})
		return
	}
	userID, _ := user_session.GetUserID(tmp.SessionId)
	log.Println("userID is", userID)
	if userID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "请登录后再点赞",
		})
		return
	}

	// 查找是否存在 item在 comment like 表
	var commentlike db_model.CommentLike
	err := db_model.Db.Where("User_id = ? AND Comment_id = ?", userID, tmp.CommentId).First(&commentlike).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 不存在
		// 插入commentlike 表
		commentlike.User_id = userID
		commentlike.Comment_id = tmp.CommentId
		db_model.Db.Create(&commentlike)
		// comment like num +1
		var comment db_model.Comment
		if err := db_model.Db.Where("ID = ?", tmp.CommentId).First(&comment).Error; err != nil {
			log.Printf("Find comment error: %v\n", err)
		}
		comment.Like_num += 1
		db_model.Db.Model(&comment).Update("like_num", comment.Like_num)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})

	} else {
		// 存在
		// 从commentlike 表删除
		db_model.Db.Delete(&commentlike)
		// comment like num -1
		var comment db_model.Comment
		if err := db_model.Db.Where("ID = ?", tmp.CommentId).First(&comment).Error; err != nil {
			log.Printf("Find comment error: %v\n", err)
		}
		comment.Like_num -= 1
		db_model.Db.Model(&comment).Update("like_num", comment.Like_num)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})

	}
}

//func AddCommentLike(c *gin.Context) {
//	var tmp CommentLikeRequest
//	if err := c.BindJSON(&tmp); err != nil {
//		log.Fatalf("AddCommentLike BindJSON error: %v", err)
//	}
//	var comment db_model.Comment
//	if err := db_model.Db.Where("ID   = ?", tmp.CommentId).First(&comment).Error; err != nil {
//		log.Printf("Find comment error: %v\n", err)
//	}
//	comment.Like_num += 1
//	db_model.Db.Model(&comment).Update("like_num", comment.Like_num)
//	c.JSON(http.StatusOK, gin.H{
//		"code": 0,
//	})
//}
//
//// ReduceCommentLike api/reducecommentlike
//func ReduceCommentLike(c *gin.Context) {
//	var tmp CommentLikeRequest
//	if err := c.BindJSON(&tmp); err != nil {
//		log.Fatalf("AddCommentLike BindJSON error: %v", err)
//	}
//	var comment db_model.Comment
//	if err := db_model.Db.Where("ID   = ?", tmp.CommentId).First(&comment).Error; err != nil {
//		log.Printf("Find comment error: %v\n", err)
//	}
//	comment.Like_num -= 1
//	db_model.Db.Model(&comment).Update("like_num", comment.Like_num)
//	c.JSON(http.StatusOK, gin.H{
//		"code": 0,
//	})
//}
