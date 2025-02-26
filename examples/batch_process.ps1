# Example PowerShell script for batch processing screenshots

# Process all screenshots from Downloads to Pictures
screenshot-sorter -source "$env:USERPROFILE\Downloads" -target "$env:USERPROFILE\Pictures" -recursive

# Process Screenshots folder with verbose output
screenshot-sorter -source "$env:USERPROFILE\Pictures\Screenshots" -verbose

# Dry run on Desktop
screenshot-sorter -source "$env:USERPROFILE\Desktop" -dry-run