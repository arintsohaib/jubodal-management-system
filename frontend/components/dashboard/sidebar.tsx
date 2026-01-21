"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { cn } from "@/pkg/utils"
import {
    LayoutDashboard,
    Users,
    Activity,
    BarChart3,
    Flag,
    Search,
    Landmark,
    UserPlus,
    ShieldCheck,
    Settings,
    ChevronRight,
    MapPin,
    User
} from "lucide-react"

const menuItems = [
    { icon: LayoutDashboard, label: "Dashboard", href: "/dashboard" },
    { icon: Users, label: "Committees", href: "/dashboard/committees" },
    { icon: Activity, label: "Activities", href: "/dashboard/activities" },
    { icon: BarChart3, label: "Analytics", href: "/dashboard/analytics" },
    { icon: UserPlus, label: "Join Requests", href: "/dashboard/join-requests" },
    { icon: Search, label: "Global Search", href: "/dashboard/search" },
    { icon: Landmark, label: "Finance Ledger", href: "/dashboard/finance" },
    { icon: ShieldCheck, label: "Security Audit", href: "/dashboard/audit" },
]

export function Sidebar() {
    const pathname = usePathname()

    return (
        <div className="w-64 border-r border-border/50 bg-card/30 backdrop-blur-xl h-screen sticky top-0 flex flex-col pt-6">
            <div className="px-6 mb-8">
                <div className="flex items-center gap-3 text-blue-500">
                    <div className="w-8 h-8 rounded-lg bg-blue-600 flex items-center justify-center text-white font-bold">B</div>
                    <span className="font-bold text-xl tracking-tight text-foreground">BJDMS</span>
                </div>
            </div>

            <div className="px-4 mb-6">
                <div className="p-3 rounded-xl bg-accent/50 border border-border/50">
                    <div className="flex items-center gap-2 text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-2">
                        <MapPin className="w-3 h-3" />
                        Jurisdiction
                    </div>
                    <div className="flex items-center justify-between text-sm font-medium">
                        <span>Dhaka Division</span>
                        <ChevronRight className="w-4 h-4 text-muted-foreground" />
                    </div>
                </div>
            </div>

            <nav className="flex-1 px-3 space-y-1">
                {menuItems.map((item) => {
                    const isActive = pathname === item.href
                    return (
                        <Link
                            key={item.href}
                            href={item.href}
                            className={cn(
                                "flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-all",
                                isActive
                                    ? "bg-blue-600 text-white shadow-lg shadow-blue-900/40"
                                    : "text-muted-foreground hover:text-foreground hover:bg-accent/50"
                            )}
                        >
                            <item.icon className="w-5 h-5" />
                            {item.label}
                        </Link>
                    )
                })}
            </nav>

            <div className="p-4 border-t border-border/50 mt-auto space-y-1">
                <Link
                    href="/dashboard/profile"
                    className={cn(
                        "flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-all",
                        pathname === "/dashboard/profile" ? "bg-accent text-foreground" : "text-muted-foreground hover:text-foreground hover:bg-accent/50"
                    )}
                >
                    <User className="w-5 h-5" />
                    My Profile
                </Link>
                <Link
                    href="/dashboard/settings"
                    className={cn(
                        "flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-all",
                        pathname === "/dashboard/settings" ? "bg-accent text-foreground" : "text-muted-foreground hover:text-foreground hover:bg-accent/50"
                    )}
                >
                    <Settings className="w-5 h-5" />
                    Settings
                </Link>
            </div>
        </div>
    )
}
