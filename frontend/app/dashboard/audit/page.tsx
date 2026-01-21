"use client"

import { useState } from "react"
import { Shield, Search, Filter, Download, Clock, User, HardDrive, AlertTriangle } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { cn } from "@/pkg/utils"

export default function AuditLogPage() {
    return (
        <div className="space-y-8">
            <div className="flex justify-between items-end">
                <div className="space-y-1">
                    <h1 className="text-3xl font-bold tracking-tight">Security Audit Log</h1>
                    <p className="text-muted-foreground">Immutable record of all administrative actions and system events.</p>
                </div>
                <div className="flex gap-3">
                    <Button variant="outline" className="gap-2">
                        <Download className="w-4 h-4" />
                        Export Audit Trail
                    </Button>
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                <AuditStat label="Total Events" value="1.2M" />
                <AuditStat label="Security Alerts" value="24" color="text-red-500" />
                <AuditStat label="User Actions" value="842" />
                <AuditStat label="Storage Health" value="99.9%" />
            </div>

            <Card className="border-border/50 bg-card/50">
                <CardHeader className="flex flex-row items-center justify-between pb-2">
                    <div>
                        <CardTitle>System Activity</CardTitle>
                        <CardDescription>Tracing actions across all jurisdictions.</CardDescription>
                    </div>
                    <div className="flex gap-2">
                        <div className="flex items-center gap-2 bg-accent/50 px-3 py-1.5 rounded-lg border border-border/50">
                            <Search className="w-4 h-4 text-muted-foreground" />
                            <input type="text" placeholder="Search logs..." className="bg-transparent outline-none text-xs w-48" />
                        </div>
                        <Button variant="ghost" size="sm" className="gap-2">
                            <Filter className="w-4 h-4" />
                            Filters
                        </Button>
                    </div>
                </CardHeader>
                <CardContent>
                    <div className="space-y-1">
                        <div className="grid grid-cols-12 px-4 py-2 text-[10px] font-black text-muted-foreground uppercase tracking-widest border-b border-border/50 mb-2">
                            <div className="col-span-3">Timestamp</div>
                            <div className="col-span-2">User / Actor</div>
                            <div className="col-span-2">Action / Event</div>
                            <div className="col-span-4">Details</div>
                            <div className="col-span-1 text-right">IP</div>
                        </div>
                        {[1, 2, 3, 4, 5, 6, 7, 8].map((i) => (
                            <div key={i} className="grid grid-cols-12 px-4 py-3 rounded-lg hover:bg-accent/30 transition-colors items-center text-sm border-b border-border/10 last:border-0 font-medium">
                                <div className="col-span-3 flex items-center gap-2 text-muted-foreground">
                                    <Clock className="w-3 h-3" />
                                    Jan 21, 2026 • 15:42:{i}0
                                </div>
                                <div className="col-span-2 flex items-center gap-2">
                                    <div className="w-6 h-6 rounded-full bg-blue-500/10 flex items-center justify-center text-[10px] font-bold text-blue-500">
                                        JD
                                    </div>
                                    <span>Sabbir Hasan</span>
                                </div>
                                <div className="col-span-2">
                                    <span className={cn(
                                        "px-2 py-0.5 rounded-full text-[10px] font-black uppercase border",
                                        i % 3 === 0 ? "bg-red-500/10 text-red-500 border-red-500/20" : "bg-blue-500/10 text-blue-500 border-blue-500/20"
                                    )}>
                                        {i % 3 === 0 ? "Permission Denied" : i % 2 === 0 ? "Member Created" : "Committee Approved"}
                                    </span>
                                </div>
                                <div className="col-span-4 text-muted-foreground truncate">
                                    {i % 3 === 0 ? "Attempted unauthorized access to Central Finance Ledger" : `Successfully processed request for ${i} members in Ward 19`}
                                </div>
                                <div className="col-span-1 text-right text-[10px] font-mono text-muted-foreground">
                                    103.42.11.{i}
                                </div>
                            </div>
                        ))}
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}

function AuditStat({ label, value, color }: { label: string, value: string, color?: string }) {
    return (
        <Card className="border-border/50 bg-card/50 p-4 space-y-1">
            <p className="text-[10px] font-black uppercase text-muted-foreground tracking-widest">{label}</p>
            <p className={cn("text-2xl font-black tracking-tighter", color)}>{value}</p>
        </Card>
    )
}
