package http

//import (
//	"go.uber.org/zap"
//)
//
//type Progress struct {
//	TotalSize int64
//	BytesRead int64
//
//	logger Logger
//}
//
//type Logger interface {
//	Info(message string, fields ...zap.Field)
//	Error(message string, fields ...zap.Field)
//}
//
//func (pr *Progress) Write(p []byte) (n int, err error) {
//	n, err = len(p), nil
//	pr.BytesRead += int64(n)
//	pr.Print()
//	return
//}
//
//func (pr *Progress) Print() {
//	if pr.BytesRead == pr.TotalSize {
//		pr.logger.Info("[+] Done")
//		return
//	}
//
//	pr.logger.Info("File upload in progress: %d\n", zap.Int64("bytes", pr.BytesRead))
//}
