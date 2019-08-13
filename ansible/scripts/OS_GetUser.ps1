Get-LocalUser | Select -Property Name, Enabled, Description | ConvertTo-Json -Depth 2 
