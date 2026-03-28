import { useState } from "react";
import { useAuth } from "@/context/AuthContext";
import useServerStore from "@/stores/ServerStore";

import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/components/ui/hover-card";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { toast } from "sonner";
import { CircleHelp, Plus } from "lucide-react";

import { Card, CardContent } from "@/components/ui/card";

import { SERVER_DATA } from "@/utils/server";

export function CreateServerCard({
  setCreateDialogOpen,
}: {
  setCreateDialogOpen: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  return (
    <Card
      className="flex flex-col items-center justify-center h-full border-dashed cursor-pointer hover:bg-gray-800 transition-colors sm:bg-gray-900 sm:bg-[url(assets/background/caves.png)] bg-blend-overlay bg-cover"
      onClick={() => setCreateDialogOpen(true)}
    >
      <CardContent className="flex flex-col items-center justify-center h-full min-h-36">
        <div className="p-3 rounded-full bg-primary/10 mb-4">
          <Plus className="h-8 w-8 text-primary" />
        </div>
        <p className="text-xl font-medium text-center">Create New Server</p>
      </CardContent>
    </Card>
  );
}

export function CreateServerDialog({
  createDialogOpen,
  setCreateDialogOpen,
}: {
  createDialogOpen: boolean;
  setCreateDialogOpen: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  const { refresh_servers, create_server } = useServerStore();
  const { session } = useAuth();

  const [form, setForm] = useState({
    name: "",
    type: "PAPER",
    version: "",
  });

  const handleAddServer = async () => {
    if (!session?.access_token) {
      toast.error("Authentication Error", {
        description: "You are not authenticated.",
      });
      return;
    }

    if (form.name.trim().length < 3 || form.name.trim().length > 20) {
      toast.error("Validation Error", {
        description: "Server name must be 3-20 characters.",
      });
      return;
    }

    if (!form.type.trim()) {
      toast.error("Validation Error", {
        description: "Server type cannot be empty.",
      });
      return;
    }

    if (!form.version.trim()) {
      toast.error("Validation Error", {
        description: "Version cannot be empty.",
      });
      return;
    }

    const result = await create_server(
      {
        name: form.name,
        type: form.type,
        version: form.version,
      },
      session?.access_token,
    );

    if (result.success) {
      refresh_servers(session?.access_token || "");
      setCreateDialogOpen(false);
      setForm({
        name: "",
        type: "PAPER",
        version: "1.21.4",
      });
      toast.success("Server Created", {
        description: `${form.name} was successfully created.`,
      });
    } else {
      toast.error("Error Creating Server", {
        description: result.error || "An unknown error occurred",
      });
    }
  };

  return (
    <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create New Server</DialogTitle>
          <DialogDescription>
            Fill in the details below to create your new server.
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-4 py-4">
          {/* Server Name */}
          <div className="space-y-2">
            <Label htmlFor="name">Server Name</Label>
            <Input
              id="name"
              placeholder="Awesome Server Name"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
            />
          </div>
          {/* Server Type */}
          <div className="space-y-2">
            <ServerTypeHoverCard />
            <Select
              value={form.type}
              onValueChange={(value) =>
                setForm({ ...form, type: value, version: "" })
              }
            >
              <SelectTrigger id="type" className="w-full">
                <SelectValue placeholder="Select type" />
              </SelectTrigger>
              <SelectContent>
                {Object.entries(SERVER_DATA).map(([key, data]) => (
                  <SelectItem value={key} className="py-2" key={key}>
                    <span>{key}</span>
                    <span className="text-xs text-muted-foreground">
                      {data.description}
                    </span>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          {/* Server Version */}
          <div className="space-y-2">
            <ServerVersionHoverCard />
            <Select
              value={form.version}
              onValueChange={(value) => setForm({ ...form, version: value })}
            >
              <SelectTrigger
                id="version"
                disabled={form.type === ""}
                className="w-[50%]"
              >
                <SelectValue placeholder="Select version" />
              </SelectTrigger>
              <SelectContent>
                {form.type !== "" &&
                  SERVER_DATA[
                    form.type as keyof typeof SERVER_DATA
                  ].versions.map((version: string) => (
                    <SelectItem value={version} key={version}>
                      {version}
                    </SelectItem>
                  ))}
              </SelectContent>
            </Select>
          </div>
        </div>
        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => setCreateDialogOpen(false)}
            className="cursor-pointer"
          >
            Cancel
          </Button>
          <Button onClick={handleAddServer} className="cursor-pointer">
            Create Server
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

const ServerTypeHoverCard = () => {
  return (
    <HoverCard>
      <HoverCardTrigger>
        <Label className="mb-2">
          Server Type
          <CircleHelp size={12} className="hover:underline -ml-1" />
        </Label>
      </HoverCardTrigger>
      <HoverCardContent className="text-sm max-w-full">
        <p className="mb-2">
          <b>Server type</b> determines performance, plugin support, and
          gameplay features.
        </p>
        <div className="ml-2">
          {Object.entries(SERVER_DATA).map(([key, data]) => (
            <p className="mb-2" key={key}>
              <b>{key.charAt(0) + key.slice(1).toLowerCase()}: </b>{" "}
              {data.longDescription}
            </p>
          ))}
        </div>
      </HoverCardContent>
    </HoverCard>
  );
};

const ServerVersionHoverCard = () => {
  return (
    <HoverCard>
      <HoverCardTrigger>
        <Label className="mb-2">
          Minecraft Version
          <CircleHelp size={12} className="hover:underline -ml-1" />
        </Label>
      </HoverCardTrigger>
      <HoverCardContent className="text-sm max-w-full">
        <p className="mb-2">
          <b>Server version</b> determines the Minecraft version and features
          available.
        </p>
        <div className="ml-2">
          <p className="mb-2">
            <b>1.21.4: </b> Latest stable Minecraft version.
          </p>
          <p className="mb-2">
            <b>1.20.1: </b> Previous stable Minecraft version.
          </p>
          <p>
            <b>1.19.4: </b> Older Minecraft version.
          </p>
        </div>
      </HoverCardContent>
    </HoverCard>
  );
};
