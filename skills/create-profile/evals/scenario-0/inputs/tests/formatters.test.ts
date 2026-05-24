import { formatDate, formatCurrency } from '../src/utils/formatters';

test('formatDate returns ISO date string', () => {
  expect(formatDate(new Date('2024-01-15'))).toBe('2024-01-15');
});

test('formatCurrency formats USD correctly', () => {
  expect(formatCurrency(1000)).toBe('$1,000.00');
});
