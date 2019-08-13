Get-LocalUser -Name Guest | Select -Property Name, Enabled, Description | ConvertTo-Json -Depth 2 

