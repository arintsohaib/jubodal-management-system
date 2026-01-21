"use client"

import { useAuth } from "@/pkg/auth-context"
import { Bell, Search, UserCircle } from "lucide-react"

export function TopNav() {
    const { user } = useAuth()

    return (
        <header className="h-16 border-b border-border/50 bg-background/50 backdrop-blur-md sticky top-0 z-10 px-8 flex items-center justify-between">
            <div className="flex items-center gap-4 bg-accent/50 px-4 py-2 rounded-full w-96 border border-border/50">
                <Search className="w-4 h-4 text-muted-foreground" />
                <input
                    type="text"
                    placeholder="Search members, activities, reports..."
                    className="bg-transparent text-sm outline-none w-full"
                />
            </div>

            <div className="flex items-center gap-6">
                <button className="relative p-2 text-muted-foreground hover:text-foreground transition-colors">
                    <Bell className="w-5 h-5" />
                    <span className="absolute top-1 right-1 w-2 h-2 bg-blue-500 rounded-full border-2 border-background"></span>
                </button>

                <div className="flex items-center gap-3 pl-4 border-l border-border/50">
                    <div className="text-right">
                        <div className="text-sm font-semibold">{user?.name || "Guest User"}</div>
                        <div className="text-xs text-muted-foreground capitalize">{user?.role || "Member"}</div>
                    </div>
                    <UserCircle className="w-8 h-8 text-blue-500" />
                </div>
            </div>
        </header>
    )
}
