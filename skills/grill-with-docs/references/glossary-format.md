# GLOSSARY.md Format

Location: `docs/context/GLOSSARY.md`

## Structure

```md
# Glossary

{One sentence: what domain this glossary covers.}

## Language

**Order**:
A customer's request to purchase one or more items.
_Avoid_: Purchase, transaction

**Invoice**:
A request for payment sent to a customer after delivery.
_Avoid_: Bill, payment request

**Customer**:
A person or organization that places orders.
_Avoid_: Client, buyer, account
```

## Rules

- **Be opinionated.** When multiple words exist for the same
  concept, pick one and list others under `_Avoid_`.
- **Keep definitions tight.** One or two sentences. Define
  what it IS, not what it does.
- **Only domain-specific terms.** General programming concepts
  don't belong. Ask: is this unique to this project's domain?
- **Group under subheadings** when natural clusters emerge.
  Flat list is fine when all terms belong to one area.

## When to create

Create `docs/context/GLOSSARY.md` lazily — only when the
first term is resolved during a grilling session.

## When to update

Update inline during the session, not batched at the end.
When a term is sharpened or resolved, write it immediately.
