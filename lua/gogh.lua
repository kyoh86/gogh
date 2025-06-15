---@class gogh.Repo
---@field full_path string
---@field path string
---@field host string
---@field owner string
---@field name string

---@class gogh
---@field repo gogh.Repo

---@type gogh.Repo
local repo = {
  full_path = "",
  path = "",
  host = "",
  owner = "",
  name = "",
}
---@type gogh
_G.gogh = { repo = repo }
