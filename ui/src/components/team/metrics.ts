import type { Component } from 'svelte';
import {
  ChartNoAxesColumn,
  Briefcase,
  CheckSquare,
  FileText,
  Layout,
  Play,
  RefreshCcw,
  Users,
  type IconProps,
} from '@lucide/svelte';

export type MetricItem = {
  key: string;
  name: string;
  value: number;
  icon: Component<IconProps>;
};

export type MetricsResponse = {
  data: {
    [key: string]: number;
  };
};

export const iconMap: { [key: string]: Component<IconProps> } = {
  user_count: Users,
  department_count: Briefcase,
  team_count: Users,
  retro_count: RefreshCcw,
  poker_count: Play,
  storyboard_count: Layout,
  team_checkin_count: CheckSquare,
  estimation_scale_count: ChartNoAxesColumn,
  retro_template_count: FileText,
};

export async function fetchAndUpdateMetrics(apiPrefix: string, currentMetrics: MetricItem[]): Promise<MetricItem[]> {
  try {
    const response = await fetch(`${apiPrefix}/metrics`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const res = (await response.json()) as MetricsResponse;

    // Update the current metrics with the fetched values
    return currentMetrics.map(metric => ({
      ...metric,
      value: res.data[metric.key] || 0, // Use 0 if the key doesn't exist in the response
    }));
  } catch (error) {
    console.error('Error fetching metrics:', error);
    // In case of error, return the current metrics unchanged
    return currentMetrics;
  }
}
