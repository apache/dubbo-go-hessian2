package logger

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/apache/dubbo-go-hessian2/toolkit-go/common-go/xlog"
)

var accessLogger *xlog.XLogger
var traceLogger *xlog.XLogger
var errLogger *xlog.XLogger

// AccessLogger access log handler
func AccessLogger() *xlog.XLogger {
	return accessLogger
}

// TraceLogger trace log handler
func TraceLogger() *xlog.XLogger {
	return traceLogger
}

// ErrLogger trace error log handler
func ErrLogger() *xlog.XLogger {
	return errLogger
}

// Sync union sync of all loggers
func Sync() {
	accessLogger.Sync()
	traceLogger.Sync()
	errLogger.Sync()
}

/**
* callLevel  第几个调用堆栈,1 表示当前堆栈，2表示第二级调用堆栈, callLevel必须大于0
* 返回调用栈的 协程id，源文件名称，方法名
 */
func getCallInfo(callLevel int) (Id string, fileLocation string, methodName string) {
	b := make([]byte, 1000)
	b = b[:runtime.Stack(b, false)]
	lines := bytes.Split(b, []byte("\n"))

	//for i,v := range lines {
	//	fmt.Println("line:"+strconv.FormatInt(int64(i+1),10)," ",string(v))
	//}

	Id = string(regexp.MustCompile("[^\\d+]").ReplaceAll(lines[0], []byte("")))

	begin := bytes.LastIndex(lines[2*callLevel], []byte("/")) + 1
	end := bytes.LastIndex(lines[2*callLevel], []byte(" "))
	fileLocation = string(lines[2*callLevel][begin:end])

	mbegin := bytes.LastIndex(lines[2*callLevel-1], []byte("/")) + 1
	mend := bytes.LastIndex(lines[2*callLevel-1], []byte("("))
	methodName = string(lines[2*callLevel-1][mbegin:mend])

	//fmt.Printf("goroutine=%s\n",Id)
	//fmt.Printf("location=%s\n",fileLocation)
	//fmt.Printf("method=%s\n",methodName)

	return
}

/**
* 获取当前的协程id
 */
func GetCurrentCoroutineId() string {
	b := make([]byte, 100)
	b = b[:runtime.Stack(b, false)]
	lines := bytes.Split(b, []byte("\n"))
	id := string(regexp.MustCompile("[^\\d+]").ReplaceAll(lines[0], []byte("")))
	return id
}

type syncMap struct {
	Data map[string]string
	Lock *sync.RWMutex
}

func (d *syncMap) Load(k string) (string, bool) {
	d.Lock.RLock()
	defer d.Lock.RUnlock()
	ret, ok := d.Data[k]
	return ret, ok
}

func (d *syncMap) Store(k, v string) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	d.Data[k] = v
}

func (d *syncMap) Delete(k string) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	delete(d.Data, k)
}

var coroutineLocal = &syncMap{
	Data: make(map[string]string),
	Lock: &sync.RWMutex{},
}

//var inheritCoroutineLocal  = make(map[string]map[string]string)

func requestContentId() string {
	rand.Seed(time.Now().UnixNano())
	val := rand.Uint64()
	value := fmt.Sprintf("%X", val)
	return value
}

/**
* 设置日志的上下文id,返回上下文id
 */
func SetLogContext() string {
	id := GetCurrentCoroutineId()
	contextId := requestContentId()
	coroutineLocal.Store(id, contextId)
	return contextId
}

/**
* 兼容老的Context-id
 */
func SetLogContextWithContextId(contextId string) {
	id := GetCurrentCoroutineId()
	coroutineLocal.Store(id, contextId)
}

/**
*  删除当前日志上下文
 * 返回上下文id,没有上下文id就返回空字符串
*/
func RemoveLogContext() string {
	id := GetCurrentCoroutineId()
	if contextId, ok := coroutineLocal.Load(id); ok {
		coroutineLocal.Delete(id)
		//保留子携程的数据
		//sub := inheritCoroutineLocal[id]
		//if sub != nil {
		//	for k,_ := range sub{
		//		//删除所有子协程的context
		//		delete(coroutineLocal,k)
		//	}
		//	//删除所有子协程的context映射关系
		//	delete(inheritCoroutineLocal,id)
		//}
		return contextId
	}
	return ""
}

/**
* 获取日志上下文id
 */
func GetLogContextId() string {
	id := GetCurrentCoroutineId()
	if contextId, ok := coroutineLocal.Load(id); ok {
		return contextId
	}
	return ""
}

func getContextId(coroutineId string) string {
	if contextId, ok := coroutineLocal.Load(coroutineId); ok {
		return contextId
	}
	return ""
}

/**
* 该方法必须在父协程代码里面调用，返回的id为父协程ID
 */
func MarkParentCoroutineId() string {
	return GetCurrentCoroutineId()
}

/**
*   继承父协程的日志上下文,
*   parentCoroutineId是调用markParentCoroutineId方法的返回值
 */
func InheritParentContext(parentCoroutineId string) (parentId, currentId string, result bool) {
	currentId = GetCurrentCoroutineId()
	if v, ok := coroutineLocal.Load(parentCoroutineId); ok {
		coroutineLocal.Store(currentId, v)
		result = true

		//sub := inheritCoroutineLocal[parentCoroutineId]
		//if sub == nil {
		//	sub = make(map[string]string)
		//}
		//sub[currentId]="1"
		//inheritCoroutineLocal[parentCoroutineId] = sub
	}
	return
}

func getCurrentTime() string {
	now := time.Now()
	second := now.Unix()
	nano := now.UnixNano()
	millsDelta := nano/1000000 - second*1000
	str := strconv.FormatInt(millsDelta, 10)
	if len(str) == 2 {
		str = "0" + str
	}
	return now.Format("2006-01-02 15:04:05") + "." + str
}

func DebugLog(content string) {
	id, loc, method := getCallInfo(3)
	contextId := getContextId(id)

	now := getCurrentTime()
	log := fmt.Sprintf("[%s] [ContextId:%s] %s [%s] %s() %s", now, contextId, id, loc, method, content)
	traceLogger.Debug(log)
}

func DebugLogFormat(format string, content ...interface{}) {
	id, loc, method := getCallInfo(3)
	contextId := getContextId(id)
	now := getCurrentTime()

	logContent := []interface{}{now, contextId, id, loc, method}
	if len(content) > 0 {
		for _, v := range content {
			logContent = append(logContent, v)
		}
	}
	log := fmt.Sprintf("[%s] [ContextId:%s] %s [%s] %s() "+format, logContent...)
	traceLogger.Debug(log)
}

func InfoLog(content string) {
	id, loc, method := getCallInfo(3)
	contextId := getContextId(id)

	now := getCurrentTime()
	log := fmt.Sprintf("[%s] [ContextId:%s] %s [%s] %s() %s", now, contextId, id, loc, method, content)
	traceLogger.Info(log)
}

func InfoLogFormat(format string, content ...interface{}) {
	id, loc, method := getCallInfo(3)
	contextId := getContextId(id)
	now := getCurrentTime()

	logContent := []interface{}{now, contextId, id, loc, method}
	if len(content) > 0 {
		for _, v := range content {
			logContent = append(logContent, v)
		}
	}
	log := fmt.Sprintf("[%s] [ContextId:%s] %s [%s] %s() "+format, logContent...)
	traceLogger.Info(log)
}

func ErrorLog(content string) {
	id, loc, method := getCallInfo(3)
	contextId := getContextId(id)

	now := getCurrentTime()
	log := fmt.Sprintf("[%s] [ContextId:%s] %s [%s] %s() %s", now, contextId, id, loc, method, content)
	errLogger.Error(log)
}

func ErrorLogFormat(format string, content ...interface{}) {
	id, loc, method := getCallInfo(3)
	contextId := getContextId(id)
	now := getCurrentTime()

	logContent := []interface{}{now, contextId, id, loc, method}
	if len(content) > 0 {
		for _, v := range content {
			logContent = append(logContent, v)
		}
	}
	log := fmt.Sprintf("[%s] [ContextId:%s] %s [%s] %s() "+format, logContent...)
	errLogger.Error(log)
}

func Debug(ctx context.Context, content string, field ...xlog.Field) {
	logContent := fmt.Sprintf("[requestID:%s]%s", getRequestID(ctx), content)
	traceLogger.Debug(logContent, field...)
}

func Debugf(ctx context.Context, format string, a ...interface{}) {
	logContent := fmt.Sprintf("[requestID:%s]%s", getRequestID(ctx), format)
	fmt.Println(logContent)
	//traceLogger.Debugf(logContent, a...)
}

func Info(ctx context.Context, content string, field ...xlog.Field) {
	logContent := fmt.Sprintf("[requestID:%s]%s", getRequestID(ctx), content)
	fmt.Println(logContent)
	//traceLogger.Info(logContent, field...)
}

func Infof(ctx context.Context, format string, a ...interface{}) {
	logContent := fmt.Sprintf("[requestID:%s]%s", getRequestID(ctx), format)
	fmt.Println(logContent)
	//traceLogger.Infof(logContent, a...)
}

func Error(ctx context.Context, content string, field ...xlog.Field) {
	logContent := fmt.Sprintf("[requestID:%s]%s", getRequestID(ctx), content)
	fmt.Println(logContent)
	//traceLogger.Error(logContent, field...)
}

func Errorf(ctx context.Context, format string, a ...interface{}) {
	logContent := fmt.Sprintf("[requestID:%s]%s", getRequestID(ctx), format)
	traceLogger.Errorf(logContent, a...)
}

func getRequestID(ctx context.Context) string {
	if ctx == nil {
		return "-"
	}
	switch ctx.(type) {
	case *gin.Context:
		requestID := ctx.Value("requestID")
		if requestID != nil {
			if requestIDStr, ok := requestID.(string); ok {
				return requestIDStr
			}
		}
	default:
		return ""
	}
	return "-"
}
