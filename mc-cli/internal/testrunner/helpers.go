package testrunner

import (
	"strings"

	"mc-cli/internal/model"
)

// resolveAttaches はAttachesDataを実際のBlockDataに変換します。
// place_blocks.go と同じロジックです。
func resolveAttaches(attaches []model.AttachesData) []model.BlockData {
	var result []model.BlockData
	for _, a := range attaches {
		state := make(map[string]string)
		component := a.Component

		isFaceType := strings.Contains(component, "lever") ||
			strings.Contains(component, "button") ||
			strings.Contains(component, "grindstone")
		isTorch := strings.Contains(component, "redstone_torch") ||
			strings.Contains(component, "redstone_wall_torch")

		if a.Pos[1] > a.Base[1] {
			// 上面 (floor)
			if isFaceType {
				state["face"] = "floor"
				state["facing"] = "north"
			} else if isTorch {
				component = "minecraft:redstone_torch"
			} else {
				state["facing"] = "up"
			}
		} else if a.Pos[1] < a.Base[1] {
			// 下面 (ceiling)
			if isFaceType {
				state["face"] = "ceiling"
				state["facing"] = "north"
			} else if isTorch {
				component = "minecraft:redstone_torch"
			} else {
				state["facing"] = "down"
			}
		} else {
			// 側面 (wall)
			facing := ""
			if a.Pos[0] > a.Base[0] {
				facing = "east"
			} else if a.Pos[0] < a.Base[0] {
				facing = "west"
			} else if a.Pos[2] > a.Base[2] {
				facing = "south"
			} else if a.Pos[2] < a.Base[2] {
				facing = "north"
			}

			if facing != "" {
				if isFaceType {
					state["face"] = "wall"
				} else if isTorch {
					component = "minecraft:redstone_wall_torch"
				}
				state["facing"] = facing
			}
		}

		result = append(result, model.BlockData{
			X:     a.Pos[0],
			Y:     a.Pos[1],
			Z:     a.Pos[2],
			Block: component,
			State: state,
		})
	}
	return result
}

// resolveConnects はConnectsDataを実際のBlockDataに変換します。
// place_blocks.go と同じロジックです。
func resolveConnects(connects []model.ConnectsData) []model.BlockData {
	var result []model.BlockData
	for _, c := range connects {
		facing := ""
		if c.To[0] > c.From[0] {
			facing = "east"
		} else if c.To[0] < c.From[0] {
			facing = "west"
		} else if c.To[2] > c.From[2] {
			facing = "south"
		} else if c.To[2] < c.From[2] {
			facing = "north"
		} else if c.To[1] > c.From[1] {
			facing = "up"
		} else if c.To[1] < c.From[1] {
			facing = "down"
		}

		state := make(map[string]string)
		if facing != "" {
			state["facing"] = facing
		}

		result = append(result, model.BlockData{
			X:     (c.From[0] + c.To[0]) / 2,
			Y:     (c.From[1] + c.To[1]) / 2,
			Z:     (c.From[2] + c.To[2]) / 2,
			Block: c.Component,
			State: state,
		})
	}
	return result
}
