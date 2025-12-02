# Billing Module E2E Tests
# Run with: .\test-billing-endpoints.ps1

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  BILLING MODULE E2E TESTS" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

$headers = @{Authorization="Bearer $global:token"}
$testsPassed = 0
$testsFailed = 0

function Test-Endpoint {
    param($Name, $ScriptBlock)
    try {
        & $ScriptBlock
        Write-Host "✓ $Name" -ForegroundColor Green
        $script:testsPassed++
    } catch {
        Write-Host "✗ $Name - $_" -ForegroundColor Red
        $script:testsFailed++
    }
}

# TEST 1: Create Invoice
Write-Host "`n[1] INVOICE TESTS" -ForegroundColor Yellow
Test-Endpoint "Create invoice with auto VAT calculation" {
    $body = @{
        clientId = $global:clientId
        baseAmount = 100.00
        issueDate = "2025-12-02T10:00:00Z"
        dueDate = "2025-12-16T10:00:00Z"
        description = "Test invoice 1"
    } | ConvertTo-Json
    
    $inv = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/invoices" -Method POST -Headers $headers -Body $body -ContentType "application/json"
    $global:invoice1Id = $inv.id
    
    if ($inv.baseAmount -ne 100 -or $inv.vatAmount -ne 21 -or $inv.totalAmount -ne 121) {
        throw "VAT calculation incorrect"
    }
    if ($inv.status -ne "unpaid") {
        throw "Status should be unpaid"
    }
    Write-Host "    Invoice: $($inv.invoiceNumber) | Total: $($inv.totalAmount) EUR" -ForegroundColor Gray
}

# TEST 2: List Invoices
Test-Endpoint "List all invoices with pagination" {
    $list = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/invoices?page=1&pageSize=10" -Method GET -Headers $headers
    if ($list.total -lt 1) {
        throw "Should have at least 1 invoice"
    }
    Write-Host "    Found $($list.total) invoices" -ForegroundColor Gray
}

# TEST 3: Get Invoice by ID
Test-Endpoint "Get invoice by ID" {
    $inv = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/invoices/$global:invoice1Id" -Method GET -Headers $headers
    if ($inv.id -ne $global:invoice1Id) {
        throw "Wrong invoice returned"
    }
}

# TEST 4: Update Invoice
Test-Endpoint "Update invoice" {
    $body = @{
        issueDate = "2025-12-02T10:00:00Z"
        dueDate = "2025-12-20T10:00:00Z"
        description = "Updated test invoice"
        baseAmount = 150.00
    } | ConvertTo-Json
    
    $inv = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/invoices/$global:invoice1Id" -Method PUT -Headers $headers -Body $body -ContentType "application/json"
    if ($inv.baseAmount -ne 150 -or $inv.totalAmount -ne 181.50) {
        throw "Update failed or VAT recalculation incorrect"
    }
    Write-Host "    New total: $($inv.totalAmount) EUR" -ForegroundColor Gray
}

# TEST 5: Mark as Paid
Test-Endpoint "Mark invoice as paid" {
    $body = @{paymentMethod = "transfer"} | ConvertTo-Json
    $inv = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/invoices/$global:invoice1Id/mark-paid" -Method POST -Headers $headers -Body $body -ContentType "application/json"
    if ($inv.status -ne "paid") {
        throw "Status should be paid"
    }
    Write-Host "    Status: $($inv.status) | Method: $($inv.paymentMethod)" -ForegroundColor Gray
}

# TEST 6: Get Unpaid Invoices
Test-Endpoint "List unpaid invoices" {
    $list = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/invoices/unpaid" -Method GET -Headers $headers
    Write-Host "    Unpaid count: $($list.Count)" -ForegroundColor Gray
}

# TEST 7: Create Expense Category
Write-Host "`n[2] EXPENSE CATEGORY TESTS" -ForegroundColor Yellow
Test-Endpoint "Create parent category" {
    $body = @{name = "Test Category"; parentId = $null} | ConvertTo-Json
    $cat = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/expense-categories" -Method POST -Headers $headers -Body $body -ContentType "application/json"
    $global:categoryId = $cat.id
    Write-Host "    Category: $($cat.name) | ID: $($cat.id)" -ForegroundColor Gray
}

# TEST 8: Create Subcategory
Test-Endpoint "Create subcategory" {
    $body = @{name = "Test Subcategory"; parentId = $global:categoryId} | ConvertTo-Json
    $sub = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/expense-categories" -Method POST -Headers $headers -Body $body -ContentType "application/json"
    $global:subcategoryId = $sub.id
    Write-Host "    Subcategory: $($sub.name) | Parent: $($sub.parentId)" -ForegroundColor Gray
}

# TEST 9: Get Category Tree
Test-Endpoint "Get hierarchical category tree" {
    $tree = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/expense-categories/tree" -Method GET -Headers $headers
    if ($tree.Count -lt 1) {
        throw "Tree should have at least 1 category"
    }
    Write-Host "    Tree has $($tree.Count) root categories" -ForegroundColor Gray
}

# TEST 10: Get Parent Categories
Test-Endpoint "List parent categories only" {
    $parents = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/expense-categories/parents" -Method GET -Headers $headers
    Write-Host "    Found $($parents.Count) parent categories" -ForegroundColor Gray
}

# TEST 11: Get Subcategories
Test-Endpoint "Get subcategories of parent" {
    $subs = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/expense-categories/$global:categoryId/subcategories" -Method GET -Headers $headers
    Write-Host "    Found $($subs.Count) subcategories" -ForegroundColor Gray
}

# TEST 12: Create Expense
Write-Host "`n[3] EXPENSE TESTS" -ForegroundColor Yellow
Test-Endpoint "Create expense" {
    $body = @{
        categoryId = $global:categoryId
        subcategoryId = $global:subcategoryId
        amount = 50.00
        expenseDate = "2025-12-01T10:00:00Z"
        supplier = "Test Supplier"
        description = "Test expense"
        hasInvoice = $true
        supplierInvoice = "SUP-001"
    } | ConvertTo-Json
    
    $exp = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/expenses" -Method POST -Headers $headers -Body $body -ContentType "application/json"
    $global:expenseId = $exp.id
    Write-Host "    Expense: $($exp.supplier) | Amount: $($exp.amount) EUR" -ForegroundColor Gray
}

# TEST 13: List Expenses
Test-Endpoint "List all expenses with pagination" {
    $list = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/expenses?page=1&pageSize=10" -Method GET -Headers $headers
    if ($list.total -lt 1) {
        throw "Should have at least 1 expense"
    }
    Write-Host "    Found $($list.total) expenses" -ForegroundColor Gray
}

# TEST 14: Filter Expenses by Category
Test-Endpoint "Filter expenses by category" {
    $list = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/expenses?categoryId=$global:categoryId" -Method GET -Headers $headers
    Write-Host "    Filtered: $($list.total) expenses in category" -ForegroundColor Gray
}

# TEST 15: Get Expense by ID
Test-Endpoint "Get expense by ID" {
    $exp = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/expenses/$global:expenseId" -Method GET -Headers $headers
    if ($exp.id -ne $global:expenseId) {
        throw "Wrong expense returned"
    }
}

# TEST 16: Update Expense
Test-Endpoint "Update expense" {
    $body = @{
        categoryId = $global:categoryId
        subcategoryId = $global:subcategoryId
        amount = 75.00
        expenseDate = "2025-12-01T10:00:00Z"
        supplier = "Updated Supplier"
        description = "Updated expense"
    } | ConvertTo-Json
    
    $exp = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/expenses/$global:expenseId" -Method PUT -Headers $headers -Body $body -ContentType "application/json"
    if ($exp.amount -ne 75) {
        throw "Update failed"
    }
    Write-Host "    New amount: $($exp.amount) EUR" -ForegroundColor Gray
}

# TEST 17: Dashboard Stats
Write-Host "`n[4] BILLING STATS TESTS" -ForegroundColor Yellow
Test-Endpoint "Get dashboard statistics" {
    $stats = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/dashboard?fromDate=2025-12-01&toDate=2025-12-31" -Method GET -Headers $headers
    Write-Host "    Revenue: $($stats.totalRevenue) EUR | Expenses: $($stats.totalExpenses) EUR | Balance: $($stats.balance) EUR" -ForegroundColor Gray
}

# TEST 18: Revenue by Month
Test-Endpoint "Get revenue grouped by month" {
    $revenue = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/revenue-by-month?fromDate=2025-01-01&toDate=2025-12-31" -Method GET -Headers $headers
    Write-Host "    $($revenue.Count) months with revenue data" -ForegroundColor Gray
}

# TEST 19: Expenses by Category
Test-Endpoint "Get expenses grouped by category" {
    $breakdown = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/expenses-by-category?fromDate=2025-12-01&toDate=2025-12-31" -Method GET -Headers $headers
    Write-Host "    $($breakdown.Count) categories with expenses" -ForegroundColor Gray
}

# TEST 20: Balance
Test-Endpoint "Get balance calculation" {
    $result = Invoke-RestMethod -Uri `"http://localhost:8080/api/v1/billing/balance?fromDate=2025-12-01&toDate=2025-12-31" -Method GET -Headers $headers
    Write-Host "    Balance: $($result.balance) EUR" -ForegroundColor Gray
}

# Summary
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  TEST SUMMARY" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Passed: $testsPassed" -ForegroundColor Green
Write-Host "  Failed: $testsFailed" -ForegroundColor Red
Write-Host "  Total:  $($testsPassed + $testsFailed)" -ForegroundColor White
Write-Host "========================================`n" -ForegroundColor Cyan

if ($testsFailed -eq 0) {
    Write-Host "✓ ALL TESTS PASSED!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "✗ SOME TESTS FAILED" -ForegroundColor Red
    exit 1
}
