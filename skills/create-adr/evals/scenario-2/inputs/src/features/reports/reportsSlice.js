import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import axios from 'axios';

// NOTE: This was added in a recent PR — conflicts with Zustand pattern elsewhere
export const fetchReports = createAsyncThunk('reports/fetch', async (params) => {
  const response = await axios.get('/api/reports', { params });
  return response.data;
});

const reportsSlice = createSlice({
  name: 'reports',
  initialState: { items: [], status: 'idle', error: null },
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchReports.pending, (state) => { state.status = 'loading'; })
      .addCase(fetchReports.fulfilled, (state, action) => {
        state.status = 'succeeded';
        state.items = action.payload;
      })
      .addCase(fetchReports.rejected, (state, action) => {
        state.status = 'failed';
        state.error = action.error.message;
      });
  },
});

export default reportsSlice.reducer;
