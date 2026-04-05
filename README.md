# A TREE-WALK INTERPRETER

Explanation on the structure of the Go Lox interpreter and the rationale behind key design decisions: **scanner**, **parser**, **AST**, and **interpreter**.

---

## Pipeline Overview

The interpreter works as a **pipeline**:

```
Source Code (string)
        │
        ▼
      Scanner
        │
        ▼
      Tokens
        │
        ▼
      Parser
        │
        ▼
       AST
        │
        ▼
     Interpreter // I have yet to implement this
        │
        ▼
     Evaluated Result
```

### Summary:

- **Scanner**: Converts raw characters into **tokens**.
- **Parser**: Converts tokens into **AST nodes** (Binary, Unary, Literal, Grouping).
- **AST**: Represents the program as a recursive tree.
- **Interpreter/Visitor**: Traverses the AST and executes code.

---

## Scanner – Tokens

**Goal**: Flatten the source into manageable pieces.

Example:

```
1 + 2 * (3 - 4)
```

Becomes tokens:

```
[Number(1)] [Plus] [Number(2)] [Star] [LeftParen] [Number(3)] [Minus] [Number(4)] [RightParen] [EOF]
```

**Implementation notes:**

```go
tokens []token.Token   // flat slice of values
```

- Tokens are **flat data**, so **use values** instead of pointers.
- Flat slice = better cache locality, fewer heap allocations.
- Millions of tokens? Value slices = single contiguous array. Pointers = millions of heap objects.

**Visual**:

```
tokens slice:  ┌────────┬────────┬────────┐
               │Token 1 │Token 2 │Token 3 │ ... millions
               └────────┴────────┴────────┘
```

---

## Parser – AST Construction

**Goal**: Build a **tree** representing expressions.

```
1 + 2 * (3 - 4)
```

AST:

```
        Binary(+)
       /        \
  Literal(1)   Binary(*)
              /        \
         Literal(2)  Grouping
                       │
                      Binary(-)
                     /       \
                Literal(3)  Literal(4)
```

**AST Nodes (recursive)**:

- `Binary` → left/right expressions + operator
- `Unary` → operator + expression
- `Literal` → value
- `Grouping` → expression inside parentheses

**Implementation notes:**

- Use **pointers for AST nodes**:

```go
type Expr interface {
    Accept(Visitor) any
}

type Binary struct {
    Left  Expr
    Operator token.Token
    Right Expr
}

func (b *Binary) Accept(v Visitor) any { ... }
```

**Why pointers?**

| Reason                        | Explanation                          |
| ----------------------------- | ------------------------------------ |
| Recursive structure           | Trees naturally hold references      |
| Avoid copying                 | Each node may be large or grow later |
| Mutation/future optimizations | Pointers allow in-place changes      |
| Identity                      | Same node can be shared across tree  |

---

## Visitor / Interpreter

**Visitor pattern** decouples behavior from data:

```
AST Node → Accept(visitor) → visitor executes node-specific logic
```

Example:

```go
Binary.Accept(interpreter)
```

Dispatch:

```
Binary node
  │
  ▼
VisitBinaryExpr
```

**Why visitor?**

- Avoids `switch expr.(type)` everywhere
- Encapsulates behavior (evaluation, printing, type checking)
- Easy to extend with new node types

**Example: Evaluating `1 + 2 * 3`**

```
Binary(+)
├─ Literal(1)
└─ Binary(*)
   ├─ Literal(2)
   └─ Literal(3)
```

Evaluation order (recursive):

```
VisitBinaryExpr(+)
  VisitLiteral(1) => 1
  VisitBinaryExpr(*)
    VisitLiteral(2) => 2
    VisitLiteral(3) => 3
    compute 2*3 => 6
compute 1+6 => 7
```

---

## Pointers vs Values – AST vs Tokens

| Component | Use          | Why                                                                      |
| --------- | ------------ | ------------------------------------------------------------------------ |
| Tokens    | **Values**   | Flat slice, millions of small objects, contiguous memory                 |
| AST nodes | **Pointers** | Recursive, shared references, mutation possible, avoids repeated copying |

**Benchmark evidence:**

```
BenchmarkValueTokens-14       1422    844845 ns/op   32MB  1 alloc
BenchmarkPointerTokens-14       73   15.5ms        40MB  1,000,001 allocs
```

- Tokens: value slice = 1 alloc
- Tokens: pointer slice = 1M+ allocs, slower
- AST: pointers = minimal copy, recursion efficient
