export type Server = {
  id: string;
  user_id: string;
  name: string;
  type: string;
  version: string;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
  ip: string;
  status: ServerStatus;
  worker_id: string;
}

export type StandardResponse = {
  success: boolean;
  data?: unknown;
  error?: string;
};

type ServerStatus = "scheduled" | "running" | "stopped" | "starting";