export interface Pigment {
  id: number;
  name: string;
  short_description?: string;
  status?: string;
  image_url?: string;
  video_url?: string;
  reflection?: number;
  absorption?: number;
  created_at?: string;
  creator_id?: number;
  likes?: number;
  is_like?: boolean;
}

export interface PigmentsApiResponse {
  status: string;
  data: Pigment[] | Pigment;
}

export const getPigmentsList = async (filter?: number): Promise<Pigment[]> => {
  let url = "/api/pigments";
  if (filter !== undefined && filter !== null && !isNaN(filter)) {
    url += `?filter=${filter}`;
  }
  return fetch(url)
    .then((res) => res.json())
    .then((data) => (Array.isArray(data.data) ? data.data : []))
    .catch(() => []);
};

export const getPigmentById = async (id: number | string): Promise<Pigment | null> => {
  return fetch(`/api/pigments/${id}`)
    .then((res) => res.json())
    .then((data) => data.data || null)
    .catch(() => null);
};

export const getFeed = async (id?: number | string, next = false): Promise<Pigment | null> => {
  let url = "/api/pigments/feed";
  const params = new URLSearchParams();
  if (id) params.append("id", String(id));
  if (next) params.append("next", "true");
  if (params.toString()) url += `?${params.toString()}`;

  return fetch(url)
    .then((res) => res.json())
    .then((data) => data.data || null)
    .catch(() => null);
};
