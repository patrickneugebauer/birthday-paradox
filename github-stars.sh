#!/bin/bash

# Convert web URL to API URL
api_url="https://api.github.com/repos/${1#https://github.com/}"

# Fetch from GitHub API
curl -sL -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer $ghtoken" \
    "$api_url"
