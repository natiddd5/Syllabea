package cache

import (
	"Syllybea/UIcomponents"
	"Syllybea/types"
	"encoding/json"
	"fmt"
	"sync"
)

var DraftCache sync.Map

var EditedCache sync.Map

func GetEditedSyllabus(syl *types.Syllabus) (*UIcomponents.Draft, error) {
	if draft, ok := EditedCache.Load(syl.ID); ok {
		if draft, ok := draft.(*UIcomponents.Draft); ok {
			fmt.Println("it's in cache")
			return draft, nil
		}
	}

	var syllabusData UIcomponents.Draft
	err := json.Unmarshal([]byte(syl.Data), &syllabusData)
	syllabusData.ID = syl.ID
	if err != nil {
		return nil, err
	}

	EditedCache.Store(syl.ID, &syllabusData)
	return &syllabusData, nil
}

func GetUserDraft(user *types.User) *UIcomponents.Draft {
	// Check if there's already a draft in the cache.
	if cachedDraft, ok := DraftCache.Load(user.ID); ok {
		if draft, ok := cachedDraft.(*UIcomponents.Draft); ok {
			return draft
		}
	}

	// Create a new draft as a pointer.

	draft := &UIcomponents.Draft{
		ID:                      -1,
		LecturerName:            user.Name,
		LecturerEmail:           user.Email,
		CourseRequirements:      []string{},
		LearningOutcomes:        []string{},
		CourseObjectives:        []string{},
		CourseStructure:         []string{},
		AssignmentsStructure:    []string{},
		SyllabusRows:            []UIcomponents.SyllabusRow{{}},
		GradeComponents:         []UIcomponents.GradeComponent{},
		BibliographyRequired:    []string{},
		BibliographyRecommended: []string{},
	}

	// Store the pointer in the cache.
	DraftCache.Store(user.ID, draft)
	return draft
}
