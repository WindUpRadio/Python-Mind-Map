package controller

import (
	"GinTest/common"
	"GinTest/model"
	"GinTest/response"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

func ImportStudents(c *gin.Context) {
	db := common.GetDataBase()
	file, err := c.FormFile("students")
	if err != nil {
		response.Response(c, 400, gin.H{"error": err}, "上传失败")
		return
	}
	filename := file.Filename
	if err != nil {
		response.Response(c, 400, gin.H{"error": err}, "上传失败")
		return
	}
	err = c.SaveUploadedFile(file, "./upload/"+filename)
	f, err := excelize.OpenFile("./upload/" + filename)
	if err != nil {
		panic(err.Error())
	}
	rows, err := f.GetRows("Sheet1")
	if err == nil {
		for _, row := range rows {
			studentID := row[0]
			studentName := row[1]
			classID := row[2]
			password := studentID[6:]
			uintID, _ := strconv.Atoi(studentID)
			classid, _ := strconv.Atoi(classID)
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				response.Response(c, 500, gin.H{"error": err}, "加密错误")
				return
			}
			newUser := model.User{
				StudentID:   uint(uintID),
				StudentName: studentName,
				Password:    string(hashedPassword),
				ClassID:     uint(classid),
				Admin:       false,
			}
			err = db.Create(&newUser).Error
			if err != nil {
				response.Response(c, 500, gin.H{"error": err}, studentID+"注册失败")
				return
			}
			response.Response(c, 200, nil, studentID+"注册成功")
		}
	}
	f.Close()
}

func UploadResource(c *gin.Context) {
	f, err := c.FormFile("resource")
	if err != nil {
		response.Response(c, 400, gin.H{"error": err}, "上传失败")
		return
	}
	err = c.SaveUploadedFile(f, "./vue/public/"+f.Filename)
	if err != nil {
		response.Response(c, 500, gin.H{"error": err}, "文件保存出错")
		return
	}
	response.Response(c, 200, gin.H{"f": f.Filename}, "上传成功")

}

func UploadCodes(c *gin.Context) {
	f, err := c.FormFile("codes")
	if err != nil {
		response.Response(c, 400, gin.H{"error": err}, "上传失败")
		return
	}
	chap := c.Query("ChapterID")
	test := c.Query("TestID")
	path := "./upload/c" + chap + "t" + test + "/"
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			response.Response(c, 500, gin.H{"error": err}, "创建文件夹出错")
			return
		}
	}
	err = c.SaveUploadedFile(f, path+f.Filename)
	if err != nil {
		response.Response(c, 500, gin.H{"error": err}, "文件保存出错")
		return
	}
	response.Response(c, 200, gin.H{"f": f.Filename}, "上传成功")
}
