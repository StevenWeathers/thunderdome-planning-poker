export type HealthCheckArea = {
  name: string;
  healthy: string;
  unhealthy: string;
};

export type HealthCheckTemplate = {
  id: string;
  name: string;
  description: string;
  areas: HealthCheckArea[];
  isPublic: boolean;
  organizationId?: string | null;
  teamId?: string | null;
};

export type HealthCheck = {
  id: string;
  teamId: string;
  name: string;
  areas: HealthCheckArea[];
  deadline: string;
  closed: boolean;
  createdBy: string;
  responseCount: number;
  createdAt: string;
  updatedAt: string;
};

export type HealthCheckRating = 'green' | 'yellow' | 'red';

export type HealthCheckTrend = 'improving' | 'stable' | 'declining';

export type HealthCheckAnswer = {
  area: string;
  rating: HealthCheckRating;
  trend: HealthCheckTrend;
  comment: string;
};

export type HealthCheckAreaResult = HealthCheckArea & {
  ratings: Record<HealthCheckRating, number>;
  trends: Record<HealthCheckTrend, number>;
  comments: string[];
};

export type HealthCheckDetail = {
  healthCheck: HealthCheck;
  canRespond: boolean;
  canManage: boolean;
  minResponses: number;
  myAnswers: HealthCheckAnswer[];
  pending: { id: string; name: string }[] | null;
  results: HealthCheckAreaResult[] | null;
};
