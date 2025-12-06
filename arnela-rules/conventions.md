# 🔠 Naming Conventions

## 4.1. Go (Backend Services)

| Convention | Context | Correct Example |
| :--- | :--- | :--- |
| **PascalCase** | Exported (public) Structs, Exported Functions | `type UserService struct {}`, `func GetUserByID() {}` |
| **camelCase** | Unexported (private) Structs, Variables, Parameters | `type carService struct {}`, `var userName string` |
| **CONST_CASE** | Public constants (acronyms in uppercase) | `const MaxRetries = 3` |

## 4.2. TypeScript (Frontend with Zustand)

| Convention | Context | Correct Example |
| :--- | :--- | :--- |
| **PascalCase** | Interfaces, Types, React Components | `interface UserProps`, `const UserList` |
| **camelCase** | Interface properties, variables, functions | `firstName: string`, `const fetchUsers = () => {}` |

## 4.3. JSON API Conventions

All keys used in API request and response bodies must use **camelCase**.

* **Go Structs:** The `json:"camelCase"` tag must be used to ensure Go serialization respects the TypeScript/JavaScript convention.

```go
// ✅ Correct - camelCase in JSON tag
type CreateUserRequest struct {
    FirstName string `json:"firstName"` 
    LastName  string `json:"lastName"` 
    Email     string `json:"email"`
}
```
