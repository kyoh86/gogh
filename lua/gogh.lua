---@class gogh.Repo
---@field full_path string
---@field path string
---@field host string
---@field owner string
---@field name string

---@class gogh.Hook
---@field id string Hook UUID
---@field name string Hook name
---@field repoPattern string Pattern that matched
---@field triggerEvent string Event that triggered the hook
---@field operationType string Type of operation that triggered the hook: "script" always.
---@field operationId string Operation UUID (script UUID)

---@class gogh
---@field repo gogh.Repo
---@field hook gogh.Hook
---@field parent? gogh.Repo|nil

---@type gogh.Repo
local repo = {
  full_path = "",
  path = "",
  host = "",
  owner = "",
  name = "",
}

---@type gogh.Hook
local hook = {
  id = "",
  name = "",
  repoPattern = "",
  triggerEvent = "",
  operationType = "script",
  operationId = "",
}

---@type gogh
_G.gogh = { repo = repo, hook = hook }
