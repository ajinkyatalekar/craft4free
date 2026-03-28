import { useEffect, useState } from "react";
import { useAuth } from "@/context/AuthContext";
import useServerStore from "@/stores/ServerStore";
import { SiteHeader } from "@/components/SiteHeader";
import { Button } from "@/components/ui/button";

import { toast } from "sonner";
import { RefreshCcw, Plus, Server } from "lucide-react";

import { Tables } from "@/../database.types";
import { API_URL } from "@/utils/server";
import {
  CreateServerCard,
  CreateServerDialog,
} from "@/components/Dashboard/CreateServerDialog";
import { ServerCard } from "@/components/Dashboard/ServerCard";
type Server = Tables<"servers">;

function Servers() {
  const { session } = useAuth();
  const { servers, refresh_servers, delete_server } = useServerStore();
  const [createDialogOpen, setCreateDialogOpen] = useState(false);

  useEffect(() => {
    const loop = async () => {
      if (!session) return;
      await refresh_servers(session.access_token);
      // setTimeout(() => {
      //   loop();
      // }, 10000);
    };

    loop();
  }, [session]);

  const handleDeleteServer = async (server_id: string) => {
    const result = await delete_server(
      server_id,
      session?.access_token || "null",
    );

    if (result.success) {
      toast.success("Server Deleted", {
        description: `Server was successfully deleted.`,
      });
      refresh_servers(session?.access_token || "");
    } else {
      toast.error("Error Deleting Server", {
        description: result.error || "An unknown error occurred",
      });
    }
  };

  const handleCopyServer = (server: Server) => {
    console.log("Copying", server.name);
    setCreateDialogOpen(true);
  };

  const handleFetch = async () => {
    try {
      const response = await fetch(`${API_URL}/api/servers`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${session?.access_token}`,
        },
      });

      const result = await response.json();
      // console.log(result.success);

      return result;
    } catch (error) {
      console.log(error);
    }
  };

  useEffect(() => {
    if (!session?.access_token) return;
    handleFetch();
  }, [session?.access_token]);

  return (
    <>
      <SiteHeader activeTab="Servers" />
      <div className="mx-auto lg:mx-10 xl:mx-20 py-10 px-4">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-3xl">Your Servers</h1>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => {
              refresh_servers(session?.access_token || "");
              toast("Servers refreshed");
            }}
            title="Refresh servers"
            className="cursor-pointer"
          >
            <RefreshCcw className="h-5 w-5" />
          </Button>
        </div>

        {/* No Servers Found */}
        {servers.length === 0 && !createDialogOpen && (
          <div className="flex flex-col items-center justify-center pt-40 text-center">
            <div className="rounded-full p-6 bg-muted">
              <Server className="h-12 w-12 text-muted-foreground" />
            </div>
            <h2 className="mt-6 text-xl font-semibold">No servers found</h2>
            <p className="mt-2 text-muted-foreground">
              You don't have any servers yet. Create your first server now!
            </p>
            <Button
              className="mt-6 cursor-pointer"
              onClick={() => setCreateDialogOpen(true)}
            >
              <Plus className="mr-2 h-4 w-4" />
              Create Server
            </Button>
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {/* Server Cards */}
          {servers.map((server) => (
            <div key={server.id}>
              <ServerCard
                server={server}
                handleCopyServer={handleCopyServer}
                handleDeleteServer={handleDeleteServer}
              />
            </div>
          ))}

          {/* Create New Server Card */}
          {servers.length != 0 && (
            <CreateServerCard setCreateDialogOpen={setCreateDialogOpen} />
          )}
        </div>

        {/* Create Server Dialog */}
        <CreateServerDialog
          createDialogOpen={createDialogOpen}
          setCreateDialogOpen={setCreateDialogOpen}
        />
      </div>
    </>
  );
}

export default Servers;
