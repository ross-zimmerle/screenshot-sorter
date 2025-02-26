#!/bin/bash

# Example script demonstrating batch processing of multiple directories

# Process all screenshots from Downloads to Pictures
screenshot-sorter -source ~/Downloads -target ~/Pictures -recursive

# Process Camera Uploads with verbose output
screenshot-sorter -source ~/Dropbox/Camera\ Uploads -verbose

# Dry run on a specific directory
screenshot-sorter -source ~/Desktop/Screenshots -dry-run