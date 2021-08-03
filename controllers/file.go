package controllers

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/go-emix/utils"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"wangpan/dao/redis"
)

func UploadHandler(c *gin.Context) {
	total := 0
	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		panic(err)
	}
	saveFile, err := os.OpenFile(header.Filename, os.O_RDWR, os.ModePerm)
	if err != nil {
		saveFile, _ = os.Create(header.Filename)
	}

	buf := make([]byte, 2048)
	if result, _ := redis.Rdb.Get(header.Filename).Int64(); result != 0 {
		total = int(result)
		file.Seek(int64(total), io.SeekStart)

		c.JSON(http.StatusOK, gin.H{
			"该文件之前已上传": total,
		})
	}
	for true {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		saveFile.WriteAt(buf, int64(total))

		total += read
		if err = redis.Rdb.Set(header.Filename, total, 0).Err(); err != nil {
			c.String(http.StatusBadRequest, "redis保存失败")
		}
	}
	saveFile.Close()
	c.JSON(http.StatusOK, gin.H{
		"上传完毕 总共": total,
	})
}

func DownloadHandler(c *gin.Context) {

	myUrl := c.PostForm("myUrl")
	dir := c.PostForm("dir")
	fileName := c.PostForm("fileName")
	//文件夹是否存在
	if !utils.FileIsExist(dir) {
		//不存在则创建文件夹
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		return
	}
	//下载文件路径
	dfn := dir + "/" + fileName

	var file *os.File
	var size int64

	if utils.FileIsExist(dfn) {
		//如果文件存在未下载完  则使file指向该文件
		fi, err := os.OpenFile(dfn, os.O_RDWR, os.ModePerm)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		//指针指向末尾
		stat,_:=fi.Stat()
		size=stat.Size()
		sk,err:=fi.Seek(size, 0)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			_=fi.Close()
			return
		}

		if sk!= size {
			c.JSON(http.StatusBadRequest,gin.H{
				"err":  "seek length not equal file size",
				"seek": sk,
				"size":size,
			})
		}

		file=fi
	} else{
		//没有则以路径dfn创建文件 并使file指向新创建的文件
		create,err:=os.Create(dfn)

		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		file=create
	}
	client:=&http.Client{}
	client.Timeout=time.Hour
	request:=http.Request{}
	request.Method=http.MethodGet

	if size!= 0 {
		//指向上次位置
		header:=http.Header{}
		header.Set("Range","bytes="+strconv.FormatInt(size,10)+"-")
		request.Header=header
	}
	parse,err:= url.Parse(myUrl)
	request.URL=parse
	get,err:=client.Do(&request)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	defer func() {
		//关闭打开的流
		err:=get.Body.Close()
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		err=file.Close()
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}()

	if get.ContentLength== 0 {
		//下载完成
			c.String(http.StatusBadRequest, "already downloaded")
			return
	}
	body:=get.Body
	writer:=bufio.NewWriter(file)
	buf:=make([]byte,10*1024*1024)
	for true {
		var read int
		read,err=body.Read(buf)
		if err != nil {
			if err != io.EOF {
				c.String(http.StatusBadRequest, err.Error())
				return
			} else {
				err=nil
			}
			break
		}
		_,err=writer.Write(buf[:read])
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			break
		}

	}
	if err != nil {
		return
	}
	err=writer.Flush()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK,"success")
}
