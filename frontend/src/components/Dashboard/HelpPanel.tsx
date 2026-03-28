import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";
import { Badge } from "../ui/badge";
import { Copy, Gamepad2, Users, Wifi } from "lucide-react";
import { Server as ServerIcon } from "lucide-react";
import { Button } from "../ui/button";
import { toast } from "sonner";
import { Server } from "@/types/server";

interface HelpPanelProps {
  server: Server | null;
}

const HelpPanel = ({ server }: HelpPanelProps) => {
  const handleCopyConnectionInfo = () => {
    if (server?.ip) {
      const connectionString = server.ip;
      navigator.clipboard.writeText(connectionString);
      toast.success("URL copied to clipboard");
    } else {
      toast.error("Server connection info not available");
    }
  };

  const getServerStatus = () => {
    if (!server) return "Unknown";
    return server.status;
  };

  const getConnectionString = () => {
    if (!server?.ip) {
      return "Server connection info will appear when server is running";
    }
    return server.ip;
  };

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Gamepad2 className="h-5 w-5" />
          How to Connect to Your Minecraft Server
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Server Status Check */}
        <div className="bg-muted/50 p-4 rounded-lg">
          <h3 className="font-semibold mb-2 flex items-center gap-2">
            <ServerIcon className="h-4 w-4" />
            Step 1: Check Server Status
          </h3>
          <p className="text-sm text-muted-foreground mb-3">
            Make sure your server is running before attempting to connect. The
            server status should show "running" above.
          </p>
          <div className="flex items-center gap-2">
            <Badge
              variant={
                getServerStatus() === "running" ? "default" : "secondary"
              }
            >
              Status: {getServerStatus()}
            </Badge>
          </div>
        </div>

        {/* Connection Information */}
        <div className="bg-blue-50 dark:bg-blue-950/20 p-4 rounded-lg border border-blue-200 dark:border-blue-800">
          <h3 className="font-semibold mb-2 flex items-center gap-2">
            <Wifi className="h-4 w-4" />
            Step 2: Get Connection Details
          </h3>
          <p className="text-sm text-muted-foreground mb-3">
            Your server connection URL will appear here when the server is
            running:
          </p>
          <div className="bg-background p-3 rounded border font-mono text-sm">
            <div className="flex items-center justify-between">
              <span>
                <span
                  className={
                    server?.ip
                      ? "text-foreground"
                      : "text-muted-foreground"
                  }
                >
                  {getConnectionString()}
                </span>
              </span>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleCopyConnectionInfo}
                disabled={!server?.ip}
              >
                <Copy className="h-3 w-3" />
              </Button>
            </div>
          </div>
        </div>

        {/* Minecraft Connection Steps */}
        <div className="space-y-4">
          <h3 className="font-semibold flex items-center gap-2">
            <Users className="h-4 w-4" />
            Step 3: Connect in Minecraft Java Edition
          </h3>

          <div className="grid gap-4">
            {/* Step 3a */}
            <div className="flex gap-4">
              <div className="flex-shrink-0 w-8 h-8 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-sm font-semibold">
                3a
              </div>
              <div className="flex-1">
                <h4 className="font-medium mb-2">Open Minecraft</h4>
                <p className="text-sm text-muted-foreground mb-2">
                  Launch Minecraft on your computer. Make sure you launch the
                  correct version (Paper 1.21.4 needs Minecraft 1.21.4).
                </p>
                <div className="bg-muted/50 p-3 rounded text-sm">
                  <strong>Note:</strong> You need to be logged in with an
                  official Mojang account to connect.
                </div>
              </div>
            </div>

            {/* Step 3b */}
            <div className="flex gap-4">
              <div className="flex-shrink-0 w-8 h-8 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-sm font-semibold">
                3b
              </div>
              <div className="flex-1">
                <h4 className="font-medium mb-2">Navigate to Multiplayer</h4>
                <p className="text-sm text-muted-foreground mb-2">
                  From the main menu, click "Multiplayer"
                </p>
                <div className="p-3 rounded text-center bg-[url(assets/dashboard/help/step_3b.jpg)] bg-cover bg-center bg-no-repeat h-160">
                </div>
              </div>
            </div>

            {/* Step 3c */}
            <div className="flex gap-4">
              <div className="flex-shrink-0 w-8 h-8 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-sm font-semibold">
                3c
              </div>
              <div className="flex-1">
                <h4 className="font-medium mb-2">Add Server</h4>
                <p className="text-sm text-muted-foreground mb-2">
                  Click "Add Server"
                </p>
                  <div className="p-3 rounded text-center bg-[url(assets/dashboard/help/step_3c.jpg)] bg-cover bg-center bg-no-repeat h-160">
                </div>
              </div>
            </div>

            {/* Step 3d */}
            <div className="flex gap-4">
              <div className="flex-shrink-0 w-8 h-8 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-sm font-semibold">
                3d
              </div>
              <div className="flex-1">
                <h4 className="font-medium mb-2">Enter Server Details</h4>
                <p className="text-sm text-muted-foreground mb-2">
                  Enter the server address from Step 2 and click "Done"
                </p>
                  <div className="p-3 rounded text-center bg-[url(assets/dashboard/help/step_3d.png)] bg-cover bg-center bg-no-repeat h-160">
                  </div>
              </div>
            </div>

            {/* Step 3e */}
            <div className="flex gap-4">
              <div className="flex-shrink-0 w-8 h-8 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-sm font-semibold">
                3e
              </div>
              <div className="flex-1">
                <h4 className="font-medium mb-2">Connect to Server</h4>
                <p className="text-sm text-muted-foreground mb-2">
                  Select your server and click "Join Server"
                </p>
                  <div className="p-3 rounded text-center bg-[url(assets/dashboard/help/step_3e.jpg)] bg-cover bg-center bg-no-repeat h-160">
                </div>
              </div>
            </div>

            {/* Step 3f */}
            <div className="flex gap-4">
              <div className="flex-shrink-0 w-8 h-8 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-sm font-semibold">
                3f
              </div>
              <div className="flex-1">
                <h4 className="font-medium mb-2">Enjoy!</h4>
                <p className="text-sm text-muted-foreground mb-2">
                  You're now connected to the server! Enjoy!
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Troubleshooting */}
        <div className="bg-yellow-50 dark:bg-yellow-950/20 p-4 rounded-lg border border-yellow-200 dark:border-yellow-800">
          <h3 className="font-semibold mb-2">Troubleshooting</h3>
          <div className="space-y-2 text-sm">
            <div>
              <strong>Can't connect?</strong> Make sure your server is running
              and you're using the correct address and port.
            </div>
            <div>
              <strong>Connection timeout?</strong> Check your internet
              connection and firewall settings.
            </div>
            <div>
              <strong>Wrong version?</strong> Ensure your Minecraft version
              matches the server version shown above.
            </div>
          </div>
        </div>

        {/* Additional Info */}
        <div className="text-xs text-muted-foreground space-y-1">
          <p>
            • Server connection details are only available when the server is
            running
          </p>
          <p>• Make sure you're using Minecraft Java Edition (not Bedrock)</p>
          <p>• The server may take a few minutes to fully start up</p>
        </div>
      </CardContent>
    </Card>
  );
};

export default HelpPanel;
