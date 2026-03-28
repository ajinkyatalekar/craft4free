import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Logo } from "@/components/SiteLogo";
import { useAuth } from "@/context/AuthContext";
import { useNavigate } from "react-router-dom";
import { HelpCircle, Server } from "lucide-react";

export function SiteHeader({ activeTab }: { activeTab?: string }) {
  const { user, signOut } = useAuth();
  const navigate = useNavigate();

  const mainNav = [
    {
      title: "Servers",
      href: "/servers",
      icon: Server,
    },
    {
      title: "Help",
      href: "/help",
      icon: HelpCircle,
    },
  ];

  return (
    <header className="bg-background sticky top-0 z-40 w-full border-b">
      <div className=" flex h-16 items-center space-x-4 sm:justify-between sm:space-x-0 ml-4 sm:ml-5">
        {/* Left Menu */}
        <div className="flex gap-4 sm:gap-6">
          <span className="inline-block font-bold">
            <div className="sm:hidden">
              <Logo otherStyles="h-[2.5rem]" sm stat dashLink />
            </div>
            <div className="hidden sm:block">
              <Logo otherStyles="h-[2.5rem]" stat dashLink />
            </div>
          </span>

          {mainNav?.length ? (
            <nav className="flex items-center space-x-2">
              {mainNav.map(
                (item, index) =>
                  item.href && (
                    <Button key={index} variant="ghost" asChild>
                      <a
                        onClick={() => navigate(item.href)}
                        className={`text-sm font-medium transition-colors cursor-pointer h-full -mb-1 ${activeTab === item.title ? "text-primary" : "text-muted-foreground"}`}
                      >
                        <item.icon />
                        {item.title}
                      </a>
                    </Button>
                  ),
              )}
            </nav>
          ) : null}
        </div>

        {/* Account Settings */}
        <div className="flex flex-1 items-center justify-end space-x-4 mr-4 sm:mx-5">
          <nav className="flex items-center space-x-3">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  className="relative h-8 w-8 rounded-full cursor-pointer"
                >
                  <Avatar className="h-8 w-8">
                    <AvatarImage
                      src={user?.user_metadata?.avatar_url || null}
                      alt={`${user?.email || "User"}'s profile`}
                    />
                    <AvatarFallback>
                      {user?.email?.[0].toUpperCase() ?? "A"}
                    </AvatarFallback>
                  </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                className="w-68 -mr-1"
                align="end"
                forceMount
              >
                <DropdownMenuLabel className="font-normal">
                  <div className="flex flex-col space-y-1">
                    <p className="text-sm leading-none font-bold py-1">
                      {user?.email}
                    </p>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  className="cursor-pointer"
                  disabled
                  onClick={() => console.log("Settings")}
                >
                  Account Settings
                </DropdownMenuItem>
                <DropdownMenuItem
                  className="cursor-pointer"
                  onClick={() => signOut()}
                >
                  Log out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </nav>
        </div>
      </div>
    </header>
  );
}
