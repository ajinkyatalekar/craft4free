import React, { createContext, ReactNode, useContext } from "react";
import supabase from "@/utils/supabase";
import { Tables } from "@/../database.types";

type Server = Tables<"servers">;
interface ServerContextType {
  getServer: (server_id: string) => Promise<Server>;
}

interface ServerProviderProps {
  children: ReactNode;
}

const ServerContext = createContext<ServerContextType | undefined>(undefined);

export const ServerProvider: React.FC<ServerProviderProps> = ({ children }) => {
  const getServer = async (server_id: string) => {
    const { data, error } = await supabase
      .from("servers")
      .select("*")
      .eq("id", server_id)
      .single();

    if (error) throw error;

    return data as Server;
  };

  const value: ServerContextType = {
    getServer,
  };

  return (
    <ServerContext.Provider value={value}>{children}</ServerContext.Provider>
  );
};

export const useServer = (): ServerContextType => {
  const context = useContext(ServerContext);
  if (context === undefined) {
    throw new Error("useServer must be used within a ServerProvider");
  }
  return context;
};
