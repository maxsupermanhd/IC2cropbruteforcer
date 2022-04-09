# IC2 crop bruteforcer

By FlexCoral

## Usage

Each field "cell" or "block" can be one of 3 things: crop, solid block, air/other.\
Green for crop, block for solid block, white for everything else.\
Input field configuration, parameters and crop data.\
"Water supplied" means you will have auto-watering of plants.

### Biome nutrient bonus

Was hardcoded in IC2/core/crop/IC2Crops.java

| biome | bonus |
|-------|-------|
| JUNGLE | 10 |
| SWAMP | 10 |
| MUSHROOM | 5 |
| FOREST | 5 |
| RIVER | 2 |
| PLAINS | 0 |
| SAVANNA | -2 |
| HILLS | -5 |
| MOUNTAIN | -5 |
| WASTELAND | -8 |
| END | -10 |
| NETHER | -10 |
| DEAD | -10 |

### Biome humidity bonus

Is set in scripts/Biomes.zs and defaults to

| biome | bonus |
|-------|-------|
| RIVER | 10 |
| SWAMP | 10 |
| JUNGLE | 7 |
| MUSHROOM | 5 |
| FOREST | 5 |
| PLAINS | 2 |
