THIS IS AI GENERATED.

A robust validation library should cover several categories of validation to handle real-world application requirements:

## 1. Primitive Type Validations
- **Strings**: length, pattern matching (regex), allowed characters, email/URL/phone formats
- **Numbers**: ranges, min/max, precision, integer/float constraints
- **Booleans**: required true/false values
- **Dates/Times**: ranges, formats, timezone handling

## 2. Structural Validations
- **Collections**: size constraints, uniqueness, sorted requirements
- **Optional/Required**: conditional presence based on other field values
- **Type checking**: ensuring values match expected types

## 3. Cross-field Validations
- **Dependency rules**: Field A requires Field B
- **Exclusivity rules**: Only one of Fields A, B, C can be present
- **Comparison rules**: Field A must be greater than Field B

## 4. Contextual Validations
- **Business rule validations**: domain-specific constraints
- **State-dependent validations**: rules that change based on entity state
- **Role-based validations**: different rules for different user roles

## 5. Advanced Validations
- **Asynchronous validations**: database lookups, API calls
- **Recursive validations**: for nested structures
- **Conditional validation chains**: IF-THEN-ELSE logic

## 6. Meta-validations
- **Rule composition**: combining rules with AND/OR/NOT operators
- **Validation groups**: applying different rule sets based on context

## Additional Considerations
- **Internationalization**: Error messages in multiple languages
- **Custom error codes**: For API responses
- **Self-documenting API**: For generating documentation/schemas
