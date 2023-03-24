$env:AWS_ACCESS_KEY_ID = ""
$env:AWS_SECRET_ACCESS_KEY = ""

Set-Content -Path amazonpolly\secrets\AWS_ACCESS_KEY_ID.txt -Value $env:AWS_ACCESS_KEY_ID
Set-Content -Path amazonpolly\secrets\AWS_SECRET_ACCESS_KEY.txt -Value $env:AWS_SECRET_ACCESS_KEY

$env:OPENAI_API_KEY = ""
$env:OPENAI_ORGANIZATION = ""

Set-Content -Path openai\secrets\OPENAI_API_KEY.txt -Value $env:OPENAI_API_KEY
Set-Content -Path openai\secrets\OPENAI_ORGANIZATION.txt -Value $env:OPENAI_ORGANIZATION