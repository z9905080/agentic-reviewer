package ports

import (
	// 假設 domain models 的路徑
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
)

// FileParser 定義了文件解析操作的通用介面
// 我們可以讓它返回 domain object，或者只返回 text content，由 use case 組裝
// 為了與現有 parser (ParseCV, ParseJD 返回 *CV, *JD) 的風格保持一致：
type FileParser interface {
	ParseCV(filePath string) (*cv.CV, error)
	ParseJD(filePath string) (*jd.JobDescription, error)
}
