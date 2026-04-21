package image

import (
	"bufio"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sync"

	ort "github.com/yalue/onnxruntime_go"
	"golang.org/x/image/draw"
)

const (
	defaultInputName  = "data"
	defaultOutputName = "mobilenetv20_output_flatten0_reshape0"
)

var (
	initOnce sync.Once
	initErr  error
)

// ImageRecognizer 图片识别器
type ImageRecognizer struct {
	session      *ort.Session[float32]
	inputName    string               // 输入张量名称（默认"data"）
	outputName   string               // 输出张量名称（默认"mobilenetv20_output_flatten0_reshape0"）
	inputH       int                  // 输入图像高度（默认224）
	inputW       int                  // 输入图像宽度（默认224）
	labels       []string             // ImageNet 类别标签
	inputTensor  *ort.Tensor[float32] // 预分配输入张量（1*3*H*W）
	outputTensor *ort.Tensor[float32] // 预分配输出张量（1*1000）
}

func NewImageRecognizer(modelPath, labelPath string, inputH, inputW int) (*ImageRecognizer, error) {
	// 1. 输入图像不足默认值设置为默认值
	if inputH <= 0 || inputW <= 0 {
		inputH, inputW = 224, 224
	}
	// 2. 单例模式初始化ONNX环境
	initOnce.Do(func() {
		initErr = ort.InitializeEnvironment()
	})
	if initErr != nil {
		return nil, fmt.Errorf("onnxruntime initialize error: %v", initErr)
	}
	// 3. 预先创建输入输出张量
	inputShape := ort.NewShape(1, 3, int64(inputH), int64(inputW))
	inData := make([]float32, inputShape.FlattenedSize())
	inTensor, err := ort.NewTensor(inputShape, inData)
	if err != nil {
		return nil, fmt.Errorf("create input tensor error: %v", err)
	}
	outputShape := ort.NewShape(1, 1000)
	outTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		inTensor.Destroy()
		return nil, fmt.Errorf("create output tensor error: %v", err)
	}
	// 4. 创建Session
	session, err := ort.NewSession(
		modelPath,
		[]string{defaultInputName},
		[]string{defaultOutputName},
		[]*ort.Tensor[float32]{inTensor},
		[]*ort.Tensor[float32]{outTensor},
	)
	if err != nil {
		inTensor.Destroy()
		outTensor.Destroy()
		return nil, fmt.Errorf("create onnx session error: %v", err)
	}
	// 5. 读取label文件
	labels, err := loadLabels(labelPath)
	if err != nil {
		session.Destroy()
		inTensor.Destroy()
		outTensor.Destroy()
		return nil, fmt.Errorf("load labels error: %v", err)
	}
	return &ImageRecognizer{
		session:      session,
		inputName:    defaultInputName,
		outputName:   defaultOutputName,
		inputH:       inputH,
		inputW:       inputW,
		labels:       labels,
		inputTensor:  inTensor,
		outputTensor: outTensor,
	}, nil
}

// loadLabels 加载标签文件
func loadLabels(labelPath string) ([]string, error) {
	f, err := os.Open(filepath.Clean(labelPath))
	if err != nil {
		return nil, fmt.Errorf("open label file error: %v", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var labels []string
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("read label file error: %v", err)
	}
	if len(labels) == 0 {
		return nil, fmt.Errorf("no labels found in %s", labelPath)
	}
	return labels, nil
}

// PredictFromFile 核心逻辑，从图片文件中预测对应的所属标签
func (r *ImageRecognizer) PredictFromFile(imagePath string) (string, error) {
	// 1. 打开图片文件
	f, err := os.Open(filepath.Clean(imagePath))
	if err != nil {
		return "", fmt.Errorf("open image file error: %v", err)
	}
	defer f.Close()
	// 2. 图片解码
	img, _, err := image.Decode(f)
	if err != nil {
		return "", fmt.Errorf("decode image file error: %v", err)
	}
	// 3. 图片预处理逻辑
	return r.PredictFromImage(img)
}

// PredictFromImage 核心逻辑，从图片对象中预测对应的所属标签
func (r *ImageRecognizer) PredictFromImage(img image.Image) (string, error) {
	// 1. 将图片数据尺寸进行限制
	resizedImg := image.NewRGBA(image.Rect(0, 0, r.inputW, r.inputH))
	// 2. 缩放到指定尺寸
	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	// 3. 将图片数据转换为 NCHW 格式的float32数组
	h, w := r.inputH, r.inputW
	ch := 3 // RGB三色通道纬度
	data := make([]float32, h*w*ch)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := resizedImg.At(x, y)
			r, g, b, _ := c.RGBA()
			rf := float32(r>>8) / 255.0
			gf := float32(g>>8) / 255.0
			bf := float32(b>>8) / 255.0
			// NCHW format
			data[y*w+x] = rf
			data[h*w+y*w+x] = gf
			data[2*h*w+y*w+x] = bf
		}
	}
	inData := r.inputTensor.GetData()
	copy(inData, data)
	// 4. 执行推理
	if err := r.session.Run(); err != nil {
		return "", fmt.Errorf("run onnx session error: %v", err)
	}
	outData := r.outputTensor.GetData()
	if len(outData) == 0 {
		return "", fmt.Errorf("output tensor data is empty")
	}
	// 5. 找到最大可能的标签索引
	maxIdx := 0
	maxVal := outData[0]
	for i := 1; i < len(outData); i++ {
		if outData[i] > maxVal {
			maxIdx = i
			maxVal = outData[i]
		}
	}
	if maxIdx >= 0 && maxIdx < len(r.labels) {
		return r.labels[maxIdx], nil
	}
	return "Unknown", nil
}

// Close 关闭资源
func (r *ImageRecognizer) Close() {
	if r.session != nil {
		_ = r.session.Destroy()
		r.session = nil
	}
	if r.inputTensor != nil {
		_ = r.inputTensor.Destroy()
		r.inputTensor = nil
	}
	if r.outputTensor != nil {
		_ = r.outputTensor.Destroy()
		r.outputTensor = nil
	}
}
