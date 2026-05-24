import { create } from 'zustand';

export const useStore = create((set) => ({
  selectedSegment: 'all',
  dateRange: { start: null, end: null },
  setSegment: (segment) => set({ selectedSegment: segment }),
  setDateRange: (range) => set({ dateRange: range }),
}));
