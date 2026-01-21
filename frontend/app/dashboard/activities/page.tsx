"use client"

import { useState } from "react"
import {
    Activity,
    Plus,
    MapPin,
    Calendar,
    Image as ImageIcon,
    CheckCircle2,
    Clock,
    AlertCircle,
    Filter,
    MoreVertical,
    Check
} from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { cn } from "@/pkg/utils"

export default function ActivitiesPage() {
    const [view, setView] = useState<"list" | "tasks">("list")

    return (
        <div className="space-y-8">
            <div className="flex justify-between items-end">
                <div className="space-y-1">
                    <h1 className="text-3xl font-bold tracking-tight">Activity & Tasks</h1>
                    <p className="text-muted-foreground">Monitor and manage field operations and assignments.</p>
                </div>
                <div className="flex gap-3">
                    <div className="flex p-1 bg-accent/50 rounded-lg">
                        <button
                            onClick={() => setView("list")}
                            className={cn(
                                "px-4 py-2 rounded-md text-sm font-medium transition-all",
                                view === "list" ? "bg-background shadow-sm text-foreground" : "text-muted-foreground"
                            )}
                        >
                            Activity Log
                        </button>
                        <button
                            onClick={() => setView("tasks")}
                            className={cn(
                                "px-4 py-2 rounded-md text-sm font-medium transition-all",
                                view === "tasks" ? "bg-background shadow-sm text-foreground" : "text-muted-foreground"
                            )}
                        >
                            Task Board
                        </button>
                    </div>
                    <Button className="bg-blue-600 hover:bg-blue-700 gap-2">
                        <Plus className="w-4 h-4" />
                        {view === "list" ? "Log Activity" : "New Task"}
                    </Button>
                </div>
            </div>

            {view === "list" ? <ActivityLog /> : <TaskBoard />}
        </div>
    )
}

function ActivityLog() {
    return (
        <div className="grid grid-cols-1 gap-6">
            <div className="flex items-center gap-4">
                <Button variant="outline" size="sm" className="gap-2">
                    <Filter className="w-4 h-4" />
                    All Jurisdictions
                </Button>
                <Button variant="outline" size="sm" className="gap-2">
                    <Calendar className="w-4 h-4" />
                    This Week
                </Button>
            </div>

            <div className="space-y-4">
                {[1, 2, 3].map((i) => (
                    <Card key={i} className="border-border/50 bg-card/50 overflow-hidden group">
                        <div className="flex flex-col md:flex-row">
                            <div className="w-full md:w-64 h-48 md:h-auto bg-accent relative overflow-hidden">
                                <div className="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1540575861501-7cf05a4b125a?auto=format&fit=crop&q=80&w=800')] bg-cover bg-center transition-transform group-hover:scale-110" />
                                <div className="absolute inset-0 bg-blue-900/40" />
                                <div className="absolute top-4 left-4 px-2 py-1 bg-emerald-500 text-white text-[10px] font-bold rounded uppercase">
                                    Verified
                                </div>
                            </div>
                            <CardContent className="flex-1 p-6">
                                <div className="flex justify-between items-start mb-4">
                                    <div className="space-y-1">
                                        <div className="flex items-center gap-2 text-xs font-semibold text-blue-400 uppercase tracking-tighter">
                                            <MapPin className="w-3 h-3" />
                                            Dhaka City South • Ward 19
                                        </div>
                                        <h3 className="text-xl font-bold">Relief Distribution Drive</h3>
                                    </div>
                                    <Button variant="ghost" size="icon">
                                        <MoreVertical className="w-4 h-4" />
                                    </Button>
                                </div>
                                <p className="text-sm text-muted-foreground line-clamp-2 mb-6">
                                    Conducted relief distribution for 200 families affected by the recent fire incident.
                                    Coordinated with local unit members and ensures transparency in list preparation.
                                </p>
                                <div className="flex items-center justify-between border-t border-border/50 pt-4">
                                    <div className="flex -space-x-2">
                                        {[1, 2, 3].map((u) => (
                                            <div key={u} className="w-8 h-8 rounded-full border-2 border-background bg-accent flex items-center justify-center text-[10px] font-bold uppercase">
                                                JD
                                            </div>
                                        ))}
                                        <div className="w-8 h-8 rounded-full border-2 border-background bg-accent/50 flex items-center justify-center text-[10px] font-bold">
                                            +5
                                        </div>
                                    </div>
                                    <div className="flex items-center gap-4 text-xs text-muted-foreground">
                                        <span className="flex items-center gap-1"><Clock className="w-3 h-3" /> 2 hours ago</span>
                                        <span className="flex items-center gap-1 font-semibold text-foreground"><CheckCircle2 className="w-3 h-3 text-emerald-500" /> 12 Evidences</span>
                                    </div>
                                </div>
                            </CardContent>
                        </div>
                    </Card>
                ))}
            </div>
        </div>
    )
}

function TaskBoard() {
    const columns = [
        { id: "todo", title: "To Do", color: "bg-slate-500" },
        { id: "in_progress", title: "In Progress", color: "bg-blue-500" },
        { id: "review", title: "Review", color: "bg-amber-500" },
        { id: "done", title: "Done", color: "bg-emerald-500" }
    ]

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 h-[calc(100vh-240px)]">
            {columns.map((col) => (
                <div key={col.id} className="flex flex-col gap-4">
                    <div className="flex items-center justify-between px-2">
                        <div className="flex items-center gap-2">
                            <div className={cn("w-2 h-2 rounded-full", col.color)} />
                            <h3 className="font-bold text-sm uppercase tracking-wider">{col.title}</h3>
                        </div>
                        <span className="text-xs bg-accent px-2 py-0.5 rounded-full font-medium">3</span>
                    </div>

                    <div className="flex-1 bg-accent/20 rounded-xl border border-border/50 p-3 space-y-4 overflow-y-auto border-dashed">
                        {[1, 2].map((t) => (
                            <Card key={t} className="border-border/50 bg-card/50 hover:border-blue-500/50 transition-colors cursor-grab active:cursor-grabbing">
                                <CardContent className="p-4 space-y-3">
                                    <div className="flex justify-between">
                                        <div className="px-2 py-0.5 bg-blue-500/10 text-blue-500 text-[10px] font-bold rounded uppercase">
                                            High
                                        </div>
                                        <Clock className="w-3 h-3 text-muted-foreground" />
                                    </div>
                                    <p className="text-sm font-semibold leading-tight">Update member database for Ward 5</p>
                                    <div className="flex items-center justify-between pt-2">
                                        <div className="w-6 h-6 rounded-full bg-accent flex items-center justify-center text-[10px] font-bold">
                                            SA
                                        </div>
                                        <div className="flex items-center gap-1 text-[10px] text-muted-foreground font-medium">
                                            <AlertCircle className="w-3 h-3" />
                                            Jan 25
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>
                        ))}
                    </div>
                </div>
            ))}
        </div>
    )
}
