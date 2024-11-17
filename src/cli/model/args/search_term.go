package args

import (
	"github.com/rohitramu/kpm/src/cli/model/utils/types"
	"github.com/rohitramu/kpm/src/pkg/utils/validation"
)

func SearchTerm(shortDescription string) *types.Arg {
	return &types.Arg{
		Name:             "search-term",
		ShortDescription: shortDescription,
		Value:            "",
		IsValidFunc:      validation.ValidateSearchTerm,
	}
}
