import { useState, useCallback } from 'react';

// Lightweight local state: each component manages its own transaction slice.
// As the dashboard grows this pattern is becoming hard to maintain —
// sibling components that need the same transaction data each fetch it
// independently, and keeping them in sync requires prop-drilling through
// several layers.

export function useTransactionStore() {
  const [transactions, setTransactions] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetchTransactions = useCallback(async (filters) => {
    setLoading(true);
    try {
      const response = await fetch(`/api/transactions?${new URLSearchParams(filters)}`);
      const data = await response.json();
      setTransactions(data.items);
    } finally {
      setLoading(false);
    }
  }, []);

  return { transactions, loading, fetchTransactions };
}
