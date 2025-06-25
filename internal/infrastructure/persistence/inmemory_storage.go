package persistence

import (
	"cv-analyzer/internal/application/ports"
	"cv-analyzer/internal/domain/analysis"
	"cv-analyzer/internal/domain/cv"
	"cv-analyzer/internal/domain/jd"
	"fmt"
	"sync"
)

// InMemoryStorage provides an in-memory implementation of the StorageService.
// Renamed from InMemoryStorageService
type InMemoryStorage struct {
	cvs      map[string]*cv.CV
	jds      map[string]*jd.JobDescription
	analyses map[string]*analysis.AnalysisResult
	mu       sync.RWMutex
}

// NewInMemoryStorage creates a new InMemoryStorage.
// Renamed from NewInMemoryStorageService and added error return for wire provider compatibility.
func NewInMemoryStorage() (*InMemoryStorage, error) {
	return &InMemoryStorage{
		cvs:      make(map[string]*cv.CV),
		jds:      make(map[string]*jd.JobDescription),
		analyses: make(map[string]*analysis.AnalysisResult),
	}, nil
}

// SaveCV saves a CV to the in-memory store.
func (s *InMemoryStorage) SaveCV(cvData *cv.CV) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cvData == nil {
		return fmt.Errorf("CV data cannot be nil")
	}
	if cvData.ID == "" {
		return fmt.Errorf("CV ID cannot be empty")
	}
	s.cvs[cvData.ID] = cvData
	// fmt.Printf("In-memory: Saved CV ID: %s\n", cvData.ID) // Optional logging
	return nil
}

// GetCV retrieves a CV from the in-memory store by ID.
func (s *InMemoryStorage) GetCV(id string) (*cv.CV, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cvData, ok := s.cvs[id]
	if !ok {
		return nil, fmt.Errorf("CV with ID %s not found in-memory", id)
	}
	return cvData, nil
}

// SaveJD saves a JobDescription to the in-memory store.
func (s *InMemoryStorage) SaveJD(jdData *jd.JobDescription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if jdData == nil {
		return fmt.Errorf("JD data cannot be nil")
	}
	if jdData.ID == "" {
		return fmt.Errorf("JD ID cannot be empty")
	}
	s.jds[jdData.ID] = jdData
	// fmt.Printf("In-memory: Saved JD ID: %s\n", jdData.ID) // Optional logging
	return nil
}

// GetJD retrieves a JobDescription from the in-memory store by ID.
func (s *InMemoryStorage) GetJD(id string) (*jd.JobDescription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	jdData, ok := s.jds[id]
	if !ok {
		return nil, fmt.Errorf("JobDescription with ID %s not found in-memory", id)
	}
	return jdData, nil
}

// SaveAnalysisResult saves an AnalysisResult to the in-memory store.
func (s *InMemoryStorage) SaveAnalysisResult(result *analysis.AnalysisResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if result == nil {
		return fmt.Errorf("analysis result cannot be nil")
	}
	if result.ID == "" {
		return fmt.Errorf("analysis result ID cannot be empty")
	}
	if _, exists := s.analyses[result.ID]; exists {
		return fmt.Errorf("analysis result with ID %s already exists", result.ID)
	}
	s.analyses[result.ID] = result
	// fmt.Printf("In-memory: Saved AnalysisResult ID: %s\n", result.ID) // Optional logging
	return nil
}

// GetAnalysisResult retrieves an AnalysisResult from the in-memory store by ID.
func (s *InMemoryStorage) GetAnalysisResult(id string) (*analysis.AnalysisResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, ok := s.analyses[id]
	if !ok {
		return nil, fmt.Errorf("AnalysisResult with ID %s not found in-memory", id)
	}
	return result, nil
}

// GetAnalysisResultsByCV retrieves all AnalysisResults associated with a given CV ID.
func (s *InMemoryStorage) GetAnalysisResultsByCV(cvID string) ([]*analysis.AnalysisResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []*analysis.AnalysisResult
	for _, analysisResult := range s.analyses {
		if analysisResult.CV != nil && analysisResult.CV.ID == cvID {
			results = append(results, analysisResult)
		}
	}
	if len(results) == 0 {
		// Consider if returning an empty slice and nil error is better than error for "not found"
		return nil, fmt.Errorf("no analysis results found for CV ID %s in-memory", cvID)
	}
	return results, nil
}

// GetAnalysisResultsByJD retrieves all AnalysisResults associated with a given JD ID.
func (s *InMemoryStorage) GetAnalysisResultsByJD(jdID string) ([]*analysis.AnalysisResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []*analysis.AnalysisResult
	for _, analysisResult := range s.analyses {
		if analysisResult.JobDescription != nil && analysisResult.JobDescription.ID == jdID {
			results = append(results, analysisResult)
		}
	}
	if len(results) == 0 {
		// Consider if returning an empty slice and nil error is better than error for "not found"
		return nil, fmt.Errorf("no analysis results found for JD ID %s in-memory", jdID)
	}
	return results, nil
}

// Compile-time check for interface implementation
var _ ports.StorageService = (*InMemoryStorage)(nil)
