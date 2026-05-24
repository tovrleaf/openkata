import { useState, useEffect } from 'react';
import axios from 'axios';

// Uses local useState — data not shared with other panels
export function MetricsPanel({ userId }) {
  const [metrics, setMetrics] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    axios.get(`/api/metrics/${userId}`)
      .then(res => setMetrics(res.data))
      .finally(() => setLoading(false));
  }, [userId]);

  if (loading) return <div>Loading...</div>;
  return <div>{JSON.stringify(metrics)}</div>;
}
