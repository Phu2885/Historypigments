export const ROUTES = {
  LIST: "/",
  DETAIL: "/detail",
} as const;

export type RouteKeyType = keyof typeof ROUTES;

export const ROUTE_LABELS: Record<RouteKeyType, string> = {
  LIST: "Список",
  DETAIL: "Подробнее",
};
