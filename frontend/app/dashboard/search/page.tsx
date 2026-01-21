"use client"

import { useState } from "react"
import { Search, User, Users, Flag, FileText, Calendar, Filter, Command, ChevronRight } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { cn } from "@/pkg/utils"

export default function SearchPage() {
    const [query, setQuery] = useState("")
    const [isSearching, setIsSearching] = useState(false)
    const [activeTab, setActiveTab] = useState("all")

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault()
        setIsSearching(true)
        // Mock search delay
        setTimeout(() => setIsSearching(false), 800)
    }

    const tabs = [
        { id: "all", label: "All Results", icon: Command },
        { id: "members", label: "Members", icon: User },
        { id: "committees", label: "Committees", icon: Users },
        { id: "activities", label: "Activities", icon: Calendar },
        { id: "complaints", label: "Complaints", icon: Flag },
    ]

    return (
        <div className="max-w-5xl mx-auto space-y-8">
            <div className="space-y-4">
                <h1 className="text-4xl font-black tracking-tight flex items-center gap-3">
                    <Search className="w-8 h-8 text-blue-500" />
                    Global Search
                </h1>
                <p className="text-muted-foreground text-lg">
                    Search across the entire BJDMS ecosystem. Find members, committees, and records instantly.
                </p>
            </div>

            <form onSubmit={handleSearch} className="relative group">
                <div className="absolute inset-0 bg-blue-600/20 rounded-2xl blur-xl opacity-0 group-focus-within:opacity-100 transition-opacity" />
                <div className="relative flex gap-2 p-2 bg-card/50 border border-border/50 rounded-2xl backdrop-blur-xl shadow-2xl">
                    <div className="flex-1 relative">
                        <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" />
                        <input
                            type="text"
                            placeholder="Search for names, NIDs, jurisdictions, or topics..."
                            value={query}
                            onChange={(e) => setQuery(e.target.value)}
                            className="w-full h-14 pl-12 pr-4 bg-transparent outline-none text-lg font-medium"
                        />
                    </div>
                    <Button type="submit" className="h-14 px-8 bg-blue-600 hover:bg-blue-700 text-lg font-bold" disabled={isSearching}>
                        {isSearching ? "Searching..." : "Search"}
                    </Button>
                </div>
            </form>

            <div className="flex items-center gap-2 overflow-x-auto pb-2 scrollbar-none">
                {tabs.map((tab) => (
                    <button
                        key={tab.id}
                        onClick={() => setActiveTab(tab.id)}
                        className={cn(
                            "flex items-center gap-2 px-4 py-2 rounded-full text-sm font-semibold whitespace-nowrap transition-all border",
                            activeTab === tab.id
                                ? "bg-blue-600/10 text-blue-500 border-blue-500/20"
                                : "text-muted-foreground border-transparent hover:bg-accent/50"
                        )}
                    >
                        <tab.icon className="w-4 h-4" />
                        {tab.label}
                    </button>
                ))}
            </div>

            <div className="space-y-4">
                {query ? (
                    <div className="space-y-6">
                        <SearchResultItem
                            type="Member"
                            title="Sheikh Mujibur Rahman (Candidate)"
                            subtitle="Dhaka City South • Joined 12 Jan 2026"
                            icon={<User className="w-5 h-5 text-blue-500" />}
                        />
                        <SearchResultItem
                            type="Committee"
                            title="Dhaka Division Central"
                            subtitle="Central Level • 56 Members Active"
                            icon={<Users className="w-5 h-5 text-emerald-500" />}
                        />
                        <SearchResultItem
                            type="Activity"
                            title="Demo for Reform in Joypurhat"
                            subtitle="Jan 15, 2026 • Verified Activity"
                            icon={<Calendar className="w-5 h-5 text-amber-500" />}
                        />
                    </div>
                ) : (
                    <div className="py-20 text-center space-y-4">
                        <div className="w-16 h-16 bg-accent rounded-full flex items-center justify-center mx-auto opacity-50">
                            <Command className="w-8 h-8" />
                        </div>
                        <div className="space-y-1">
                            <p className="font-bold text-muted-foreground">Type to start searching</p>
                            <p className="text-sm text-muted-foreground/50">Search is case-insensitive and supports Bangla keywords.</p>
                        </div>
                    </div>
                )}
            </div>
        </div>
    )
}

function SearchResultItem({ type, title, subtitle, icon }: any) {
    return (
        <Card className="border-border/50 bg-card/30 hover:bg-accent/30 transition-colors cursor-pointer group">
            <CardContent className="p-4 flex items-center justify-between">
                <div className="flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-background flex items-center justify-center border border-border/50 group-hover:border-blue-500/50 transition-colors">
                        {icon}
                    </div>
                    <div className="space-y-0.5">
                        <p className="text-[10px] font-black uppercase text-muted-foreground tracking-widest">{type}</p>
                        <h3 className="font-bold group-hover:text-blue-400 transition-colors">{title}</h3>
                        <p className="text-xs text-muted-foreground">{subtitle}</p>
                    </div>
                </div>
                <ChevronRight className="w-5 h-5 text-muted-foreground group-hover:text-blue-400" />
            </CardContent>
        </Card>
    )
}
