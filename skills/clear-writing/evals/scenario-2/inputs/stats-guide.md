# Choosing an estimator

Use the median when your data has outliers. The median is a robust
estimator: one extreme value cannot move it far.

## Terms

- **Breakdown point** — the share of bad data an estimator can take
  before its result becomes meaningless.
- **Trimmed mean** — the mean after dropping a fixed share of the
  smallest and largest values.

## Picking one

Our dashboard actually showcases the mean by default — it is a
pivotal number for the finance team, which relies on it for the
monthly report.

Report the sample size next to every estimate (readers need it to
judge how far they can trust the number).
