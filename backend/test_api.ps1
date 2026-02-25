# API Testing Script for Windows PowerShell
# Запуск: .\test_api.ps1

$BASE_URL = "http://localhost:8080"

Write-Host "=== 1. Health Check ===" -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/health" -Method Get
    Write-Host "✓ Server is running" -ForegroundColor Green
    $response | ConvertTo-Json
} catch {
    Write-Host "✗ Server is not responding" -ForegroundColor Red
    exit 1
}

Write-Host "`n=== 2. Signup (Registration) ===" -ForegroundColor Green
$signupBody = @{
    email = "test$(Get-Random)@example.com"
    password = "password123"
    username = "TestUser$(Get-Random -Maximum 1000)"
    role = "student"
} | ConvertTo-Json

try {
    $signupResponse = Invoke-RestMethod -Uri "$BASE_URL/api/signup" -Method Post -Body $signupBody -ContentType "application/json"
    Write-Host "✓ User created successfully" -ForegroundColor Green
    $signupResponse | ConvertTo-Json
    $userEmail = $signupResponse.email
    $userId = $signupResponse.id
} catch {
    Write-Host "✗ Signup failed: $_" -ForegroundColor Red
}

Write-Host "`n=== 3. Login ===" -ForegroundColor Green
$loginBody = @{
    email = $userEmail
    password = "password123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$BASE_URL/api/login" -Method Post -Body $loginBody -ContentType "application/json"
    Write-Host "✓ Login successful" -ForegroundColor Green
    $token = $loginResponse.token
    Write-Host "Token: $token" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Login failed: $_" -ForegroundColor Red
}

Write-Host "`n=== 4. Create Profile ===" -ForegroundColor Green
$profileBody = @{
    user_id = $userId
    full_name = "Test User Full Name"
    role = "student"
    preferred_lang = "ru"
} | ConvertTo-Json

try {
    $profileResponse = Invoke-RestMethod -Uri "$BASE_URL/api/profiles" -Method Post -Body $profileBody -ContentType "application/json"
    Write-Host "✓ Profile created" -ForegroundColor Green
    $profileResponse | ConvertTo-Json
} catch {
    Write-Host "✗ Profile creation failed: $_" -ForegroundColor Red
}

Write-Host "`n=== 5. Generate Test (Russian) ===" -ForegroundColor Green
$generateBody = @{
    subject = "Математика"
    topic = "Алгебра"
    difficulty = "easy"
    lang = "ru"
} | ConvertTo-Json

try {
    $testResponse = Invoke-RestMethod -Uri "$BASE_URL/api/test/generate" -Method Post -Body $generateBody -ContentType "application/json"
    Write-Host "✓ Test generated successfully" -ForegroundColor Green
    $testResponse | ConvertTo-Json -Depth 5
    $questions = $testResponse.questions
} catch {
    Write-Host "✗ Test generation failed: $_" -ForegroundColor Red
}

Write-Host "`n=== 6. Submit Test Answers ===" -ForegroundColor Green
$submitBody = @{
    user_id = $userId
    subject = "Математика"
    topic = "Алгебра"
    language = "ru"
    questions = @(
        @{
            id = 1
            text = "Sample question 1"
            correct_answer = "A"
            user_answer = "A"
        },
        @{
            id = 2
            text = "Sample question 2"
            correct_answer = "B"
            user_answer = "C"
        }
    )
} | ConvertTo-Json -Depth 5

try {
    $submitResponse = Invoke-RestMethod -Uri "$BASE_URL/api/test/submit" -Method Post -Body $submitBody -ContentType "application/json"
    Write-Host "✓ Test submitted successfully" -ForegroundColor Green
    Write-Host "Score: $($submitResponse.attempt.score)%" -ForegroundColor Cyan
    $submitResponse | ConvertTo-Json -Depth 5
} catch {
    Write-Host "✗ Test submission failed: $_" -ForegroundColor Red
}

Write-Host "`n=== 7. Get User Attempts ===" -ForegroundColor Green
try {
    $attemptsResponse = Invoke-RestMethod -Uri "$BASE_URL/api/attempts/user/$userId" -Method Get
    Write-Host "✓ Retrieved user attempts" -ForegroundColor Green
    $attemptsResponse | ConvertTo-Json -Depth 5
} catch {
    Write-Host "✗ Failed to get attempts: $_" -ForegroundColor Red
}

Write-Host "`n=== 8. Get Group Analytics ===" -ForegroundColor Green
try {
    $analyticsResponse = Invoke-RestMethod -Uri "$BASE_URL/api/analytics/group" -Method Get
    Write-Host "✓ Retrieved analytics" -ForegroundColor Green
    $analyticsResponse | ConvertTo-Json
} catch {
    Write-Host "✗ Failed to get analytics: $_" -ForegroundColor Red
}

Write-Host "`n=== 9. Test Kazakh Language ===" -ForegroundColor Green
$generateKkBody = @{
    subject = "Математика"
    topic = "Геометрия"
    difficulty = "medium"
    lang = "kk"
} | ConvertTo-Json

try {
    $testKkResponse = Invoke-RestMethod -Uri "$BASE_URL/api/test/generate" -Method Post -Body $generateKkBody -ContentType "application/json"
    Write-Host "✓ Kazakh test generated" -ForegroundColor Green
    $testKkResponse | ConvertTo-Json -Depth 5
} catch {
    Write-Host "✗ Kazakh test generation failed: $_" -ForegroundColor Red
}

Write-Host "`n=== All Tests Completed ===" -ForegroundColor Green
Write-Host "Summary:" -ForegroundColor Cyan
Write-Host "- User ID: $userId"
Write-Host "- Email: $userEmail"
Write-Host "- Token: $token"
