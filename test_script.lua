function execute()
  local http = require("http")
  local client = http.client()
  local request = http.request("GET", "https://google.com")
  local result, err = client:do_request(request)
  if err then error(err) end

  return {
    success = true,
    counters = {"fart"}
  } 
end