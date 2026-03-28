import { Server } from "@/types/server";
import { IconCopy } from "@tabler/icons-react";
import { toast } from "sonner";

export const ServerStatus = ({
  server,
  styles,
}: {
  server: Server | null;
  styles?: string;
}) => {
  return (
    <>
      <div className={`flex items-center gap-2 ${styles}`}>
        {server?.status === "running" && (
          <>
          <div className="flex flex-col">
          <span className="text-green-400">
              Running
            </span>
            <div className="flex gap-2 items-center">
            <span className="">
              Connect to {server?.name}.craft4free.com
            </span>
            <IconCopy
              size={14}
              className="cursor-pointer hover:text-gray-400 transition-colors -ml-1"
              onClick={() => {
                navigator.clipboard.writeText(`${server?.name}.craft4free.com`);
                toast.success(
                  `Server ${server?.name}'s URL copied to clipboard`,
                );
              }}
            />
            </div>
            <div className="flex gap-2 items-center">
              <span className="text-gray-400 text-sm">
                Temporary IP: {server?.ip}
              </span>
              <IconCopy
                size={10}
                className="cursor-pointer hover:text-gray-400 transition-colors -ml-1"
                onClick={() => {
                  navigator.clipboard.writeText(server?.ip || "");
                  toast.success(
                    `Server ${server?.name}'s URL copied to clipboard`,
                  );
                }}
              />
              </div>
            </div>
          </>
        )}
        {server?.status === "scheduled" && (
          <>
          <div className="flex flex-col">
          <span className="text-orange-400">
              Starting
            </span>
            <div className="flex gap-2 items-center">
            <span className="">
              Connect to {server?.name}.craft4free.com
            </span>
            <IconCopy
              size={14}
              className="cursor-pointer hover:text-gray-400 transition-colors -ml-1"
              onClick={() => {
                navigator.clipboard.writeText(`${server?.name}.craft4free.com`);
                toast.success(
                  `Server ${server?.name}'s URL copied to clipboard`,
                );
              }}
            />
            </div>
            <div className="flex gap-2 items-center">
              <span className="text-gray-400 text-sm">
                Temporary IP: {server?.ip}
              </span>
              <IconCopy
                size={10}
                className="cursor-pointer hover:text-gray-400 transition-colors -ml-1"
                onClick={() => {
                  navigator.clipboard.writeText(server?.ip || "");
                  toast.success(
                    `Server ${server?.name}'s URL copied to clipboard`,
                  );
                }}
              />
              </div>
            </div>
          </>
        )}
        {server?.status === "stopped" && (
          <>
            <span className="text-red-400">
              Not running
              <br/>
              {server?.name}.craft4free.com
            </span>
          </>
        )}
      </div>
    </>
  );
};
