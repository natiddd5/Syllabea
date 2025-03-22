package handler

import (
	"Syllybea/UIcomponents"
	"Syllybea/types"
	"sync"
)

var draftCache sync.Map

func GetUserDraft(user *types.User) *UIcomponents.Draft {
	// Check if there's already a draft in the cache.
	if cachedDraft, ok := draftCache.Load(user.ID); ok {
		if draft, ok := cachedDraft.(*UIcomponents.Draft); ok {
			return draft
		}
	}

	// Create a new draft as a pointer.
	draft := &UIcomponents.Draft{
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
	draftCache.Store(user.ID, draft)
	return draft
}
