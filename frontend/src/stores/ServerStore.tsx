import { create } from "zustand";

import { components } from "@/../api.types";
import { API_URL } from "@/utils/server";
import { Server, StandardResponse } from "@/types/server";

type ServerCreateRequest = components["schemas"]["ServerCreateRequest"];

interface ServerState {
  servers: Server[];
  isLoading: boolean;
  create_server: (
    body: ServerCreateRequest,
    token: string,
  ) => Promise<StandardResponse>;
  refresh_servers: (token: string) => Promise<void>;
  fetch_server: (server_id: string, token: string) => Promise<StandardResponse>;
  delete_server: (
    server_id: string,
    token: string,
  ) => Promise<StandardResponse>;
  start_server: (server_id: string, token: string) => Promise<StandardResponse>;
  stop_server: (server_id: string, token: string) => Promise<StandardResponse>;
}

const useServerStore = create<ServerState>((set) => ({
  servers: [],
  isLoading: false,
  error: null,

  create_server: async (body: ServerCreateRequest, token: string) => {
    set({ isLoading: true });
    try {
      const response = await fetch(`${API_URL}/api/servers`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      });

      const result: StandardResponse = await response.json();

      if (result.success && result.data) {
        set({
          isLoading: false,
        });
      } else {
        set({
          isLoading: false,
        });
      }

      return result;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to create server";
      set({ isLoading: false });
      return {
        success: false,
        error: errorMessage,
      };
    }
  },

  fetch_server: async (server_id: string, token: string) => {
    try {
      const response = await fetch(`${API_URL}/api/servers/${server_id}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        // body: JSON.stringify({ server_id: server_id }),
      });

      const result: StandardResponse = await response.json();
      console.log(result);
      return result;
    } catch (error) {
      return {
        success: false,
        error:
          error instanceof Error ? error.message : "Failed to fetch server",
      };
    }
  },

  refresh_servers: async (token: string) => {
    set({ isLoading: true });
    try {
      const response = await fetch(`${API_URL}/api/servers`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });

      const result: StandardResponse = await response.json();
      if (!result.success) {
        throw new Error(result.error);
      }

      set({
        servers: result.data as Server[],
        isLoading: false,
      });
    } catch {
      set({
        isLoading: false,
      });
    }
  },

  stop_server: async (server_id: string, token: string) => {
    set({ isLoading: true });
    try {
      const response = await fetch(`${API_URL}/api/servers/${server_id}/stop`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ server_id: server_id }),
      });

      const result: StandardResponse = await response.json();

      return result;
    } catch (error) {
      console.log(error);
      return { success: false, error: "Unknown error occured" };
    } finally {
      set({ isLoading: false });
    }
  },

  start_server: async (server_id: string, token: string) => {
    set({ isLoading: true });
    try {
      const response = await fetch(`${API_URL}/api/servers/${server_id}/start`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ server_id: server_id }),
      });

      const result: StandardResponse = await response.json();

      return result;
    } catch (error) {
      console.log(error);
      return { success: false, error: "Unknown error occured" };
    } finally {
      set({ isLoading: false });
    }
  },

  delete_server: async (server_id: string, token: string) => {
    set({ isLoading: true });
    try {
      const response = await fetch(`${API_URL}/servers/${server_id}/delete`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });

      const result: StandardResponse = await response.json();

      if (result.success && result.data) {
        set({
          isLoading: false,
        });
      } else {
        set({
          isLoading: false,
        });
      }

      return result;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to delete server";
      set({ isLoading: false });
      return {
        success: false,
        error: errorMessage,
      };
    }
  },
}));

export default useServerStore;
