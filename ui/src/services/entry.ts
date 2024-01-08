import request from '@/utils/request';


export const listEntries = (params: Record<string, unknown>): Promise<API.ListEntriesResponse> => (
  request.post('/api/starlite/v1/entries/list', params)
);


export const listEntryCategory = (params: Record<string, unknown>): Promise<API.ListEntryCategoryResponse> => (
  request.post('/api/starlite/v1/entry/category/list', params)
);

export const addEntryCategory = (params: Record<string, unknown>): Promise<API.EntryCategoryResponse> => (
  request.post('/api/starlite/v1/entry/category', params)
);

export const updateEntryCategory = (params: Record<string, unknown>): Promise<API.EntryCategoryResponse> => (
  request.put('/api/starlite/v1/entry/category', params)
);

export const delEntryCategory = (category :string): Promise<Record<string, unknown>> => (
  request.delete(`/api/starlite/v1/entry/category/${category}`)
);

export const addEntry = (params: API.EntryRequest): Promise<API.EntryResponse> => (
  request.post('/api/starlite/v1/entry', params)
);

export const updateEntry = (params: API.EntryData): Promise<API.EntryResponse> => (
  request.put('/api/starlite/v1/entry', params)
);

export const delEntry = (code:  string): Promise<API.DelEntryResponse> => (
  request.delete(`/api/starlite/v1/entry/${code}`)
);