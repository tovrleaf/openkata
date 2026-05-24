import { useStore } from '../store/dashboardStore';

// Uses Zustand store for selected filters
export function UserFilter() {
  const { selectedSegment, setSegment } = useStore();

  return (
    <select value={selectedSegment} onChange={e => setSegment(e.target.value)}>
      <option value="all">All Users</option>
      <option value="free">Free Tier</option>
      <option value="pro">Pro</option>
    </select>
  );
}
