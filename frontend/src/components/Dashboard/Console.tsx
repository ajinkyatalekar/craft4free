import { API_URL } from "@/utils/server";
import React, { useEffect, useState, useRef } from "react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { AlertCircle, CheckCircle2, Terminal } from "lucide-react";

interface ServerLogsProps {
  serverId: string;
  stopped: boolean;
}

const Console: React.FC<ServerLogsProps> = ({ serverId, stopped }) => {
  const [logs, setLogs] = useState<string[]>([]);
  const [connected, setConnected] = useState<boolean>(false);
  const [serverRunning, setServerRunning] = useState<boolean>(true);
  const socketRef = useRef<WebSocket | null>(null);
  const logsEndRef = useRef<HTMLDivElement>(null);
  const consoleContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Create WebSocket connection if not stopped
    if (stopped) {
      setConnected(false);
      setLogs([]);
      return;
    }

    // Convert http:// to ws:// for WebSocket connection
    const wsUrl = API_URL.replace(/^http/, "ws");
    const ws = new WebSocket(`${wsUrl}/ws/console/${serverId}`);
    socketRef.current = ws;

    ws.onopen = () => {
      setConnected(true);
    };

    ws.onmessage = (event) => {
      const message = event.data;
      if (message === "[SERVER_NOT_RUNNING]") {
        setServerRunning(false);
        setLogs((prev) => [...prev, "Server is not running"]);
      } else {
        setLogs((prev) => [...prev, message]);
      }
    };

    ws.onclose = () => {
      setConnected(false);
    };

    // Clean up on unmount or when stopped changes
    return () => {
      ws.close();
    };
  }, [serverId, stopped]);

  // Auto-scroll to bottom when new logs arrive
  useEffect(() => {
    if (consoleContainerRef.current) {
      consoleContainerRef.current.scrollTop =
        consoleContainerRef.current.scrollHeight;
    }
  }, [logs]);

  return (
    <Card className="w-full">
      <CardHeader className="px-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <Terminal className="h-5 w-5 text-gray-500 dark:text-gray-400" />
            <CardTitle className="text-base font-medium">
              Server Console
            </CardTitle>
          </div>

          <div className="flex items-center gap-2">
            {connected ? (
              <Badge
                variant="outline"
                className="flex items-center gap-1 bg-green-50 text-green-700 dark:bg-green-900/20 dark:text-green-400 border-green-200 dark:border-green-800"
              >
                <CheckCircle2 className="h-3.5 w-3.5" />
                <span>Connected</span>
              </Badge>
            ) : (
              <Badge
                variant="outline"
                className="flex items-center gap-1 bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-400 border-red-200 dark:border-red-800"
              >
                <AlertCircle className="h-3.5 w-3.5" />
                <span>Disconnected</span>
              </Badge>
            )}

            {!serverRunning && !stopped && (
              <Badge
                variant="outline"
                className="flex items-center gap-1 bg-yellow-50 text-yellow-700 dark:bg-yellow-900/20 dark:text-yellow-400 border-yellow-200 dark:border-yellow-800"
              >
                <AlertCircle className="h-3.5 w-3.5" />
                <span>Server not running</span>
              </Badge>
            )}
          </div>
        </div>
      </CardHeader>

      <CardContent className="p-0">
        <div
          ref={consoleContainerRef}
          className="h-96 border-y border-gray-800/50
          overflow-x-hidden overflow-y-auto
          [&::-webkit-scrollbar]:w-2
          [&::-webkit-scrollbar-track]:bg-gray-100
          [&::-webkit-scrollbar-thumb]:bg-gray-300
          dark:[&::-webkit-scrollbar-track]:bg-neutral-700
          dark:[&::-webkit-scrollbar-thumb]:bg-neutral-500"
        >
          <div className="min-h-95 font-mono text-sm bg-zinc-950 text-gray-300 p-4">
            {logs.length > 0 ? (
              logs.map((log, index) => (
                <div
                  key={index}
                  className="log-line py-0.5 border-b border-gray-800/40 last:border-0"
                >
                  <span className="text-gray-500 mr-2">{`[${index + 1}]`}</span>
                  <span>{log}</span>
                </div>
              ))
            ) : (
              <div className="text-gray-500 italic text-center py-4">
                The console output will show up here
              </div>
            )}
            <div ref={logsEndRef} />
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default Console;
