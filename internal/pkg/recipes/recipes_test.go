package recipes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	recipes := []Recipe{
		{
			Name: "drink1",
			Ingredients: []*Ingredient{
				{
					Description: "booze1",
				},
				{
					Description: "booze2",
				},
			},
		},
		{
			Name: "drink2",
			Ingredients: []*Ingredient{
				{
					Description: "booze1",
				},
				{
					Description: "booze3",
				},
			},
		},
	}

	index, err := readAndIndexFiles(recipes)
	require.NoError(t, err)

	recipesService := service{
		recipesIndex: *index,
		recipesMap: map[string]Recipe{
			"drink1": recipes[0],
			"drink2": recipes[0],
		},
	}

	hits, err := recipesService.Search([]string{"booze1"}, nil)
	require.NoError(t, err)
	require.Len(t, hits, 2)
	require.Equal(t, "drink1", hits[0].Name)
	require.Equal(t, "drink1", hits[1].Name)

	hits, err = recipesService.Search([]string{"booze1", "booze2"}, nil)
	require.NoError(t, err)
	require.Len(t, hits, 1)
	require.Equal(t, "drink1", hits[0].Name)

	hits, err = recipesService.Search([]string{"booze4"}, nil)
	require.ErrorIs(t, err, ErrNotFound)
	require.Nil(t, hits)

	hits, err = recipesService.Search([]string{"booze1"}, []string{"booze3"})
	require.NoError(t, err)
	require.Len(t, hits, 1)
	require.Equal(t, "drink1", hits[0].Name)
}
