package server

import (
	"gbl-api/controllers/booth"
	"gbl-api/controllers/problem"
	"gbl-api/controllers/score"
	"github.com/gin-gonic/gin"
	"log"
)

func problemList(c *gin.Context) {
	bid := c.Param("bid")
	problems, err := problem.GetBoothProblems(bid)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"problems": problems,
	})
}

type problemSubmitRequest struct {
	UID           string   `json:"uid"`
	submit_answer []string `json:"submit_answer"`
}

func problemSubmit(c *gin.Context) {
	bid := c.Param("bid")
	var req problemSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}
	uid := req.UID

	var b booth.Booth
	b, err := booth.GetBooth(bid)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}

	var totalScore int
	var scores []int
	for i, p := range b.ProblemOrder {
		s := problem.CheckAnswer(p, req.submit_answer[i])
		totalScore += s
		scores = append(scores, s)
		err := score.AddScore(bid, uid, p, s)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{
				"message": "Internal Server Error",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"score":   totalScore,
		"answers": scores,
	})
}
