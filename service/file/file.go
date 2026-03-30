package file

import (
	"GoNexus/common/rag"
	"GoNexus/config"
	"GoNexus/utils"
	"context"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// UploadRagFile 上传Rag相关文件
func UploadRagFile(username string, file *multipart.FileHeader) (string, error) {
	// 1. 校验文件类型和文件名
	if err := utils.ValidateFile(file); err != nil {
		log.Printf("File validation failed. err: %v", err)
		return "", err
	}
	// 2. 创建用户目录
	userDir := filepath.Join("uploads", username)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		log.Printf("Create user directory failed. err: %v", err)
		return "", err
	}
	// 3. 删除用户目录中的所有现有文件及其索引
	// 每个用户只能拥有一个文件,需要将其上一次上传的文件删除才能进行这一次的文件Rag操作
	files, err := os.ReadDir(userDir)
	if err == nil {
		for _, f := range files {
			if !f.IsDir() {
				// 不是目录说明是文件
				fileName := f.Name()
				// 删除该文件对应的 Redis 文件索引
				if err = rag.DeleteIndex(context.Background(), fileName); err != nil {
					log.Printf("Delete index for %s failed. err: %v", fileName, err)
				}
			}
		}
	}
	if err = utils.RemoveAllFilesInDir(userDir); err != nil {
		log.Printf("Remove user directory file failed. err: %v", err)
		return "", err
	}
	// 4. 生成UUID作为唯一文件名
	filename := uuid.NewString() + filepath.Ext(file.Filename)
	filePath := filepath.Join(userDir, filename)
	// 5. 打开上传的文件
	src, err := file.Open()
	if err != nil {
		log.Printf("Open file failed. err: %v", err)
		return "", err
	}
	defer src.Close()
	// 6. 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Create target file failed. err: %v", err)
		return "", err
	}
	defer dst.Close()
	// 7. 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		log.Printf("Copy file content failed. err: %v", err)
		return "", err
	}
	log.Printf("File %s uploaded successfully", filePath)
	// 8. 对上传的文件创建Rag索引
	return CreatRagIndexForFile(filename, filePath)
}

// CreatRagIndexForFile 对文件创建Rag索引
func CreatRagIndexForFile(filename, filepath string) (string, error) {
	// 1. 创建RAG索引器并对文件进行向量化
	indexer, err := rag.NewRAGIndexer(filename, config.GetConfig().RagModelConfig.RagEmbeddingModel)
	if err != nil {
		log.Printf("Create rag indexer failed. err: %v", err)
		// 删除已上传的文件
		os.Remove(filepath)
		return "", err
	}
	// 2. 读取文件内容并创建向量索引
	if err = indexer.IndexFile(context.Background(), filepath); err != nil {
		log.Printf("Index file failed. err: %v", err)
		// 删除已上传的文件和索引
		os.Remove(filepath)
		rag.DeleteIndex(context.Background(), filename)
		return "", err
	}
	// 3 返回文件路径
	log.Printf("File %s indexed successfully", filename)
	return filepath, nil
}
