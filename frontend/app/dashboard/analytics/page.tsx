"use client"

import { useState } from "react"
import { TrendingUp, Users, Activity, BarChart3, Map as MapIcon, Calendar, ArrowUpRight, ArrowDownRight } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { cn } from "@/pkg/utils"

export default function AnalyticsPage() {
    return (
        <div className="space-y-8">
            <div className="flex justify-between items-end">
                <div className="space-y-1">
                    <h1 className="text-3xl font-bold tracking-tight">Organizational Intelligence</h1>
                    <p className="text-muted-foreground">Strategic analytics and performance heatmaps across all jurisdictions.</p>
                </div>
                <div className="flex gap-3">
                    <Button variant="outline" className="gap-2">
                        <Calendar className="w-4 h-4" />
                        Last 30 Days
                    </Button>
                </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <Card className="lg:col-span-2 border-border/50 bg-card/30 backdrop-blur-md overflow-hidden">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <MapIcon className="w-5 h-5 text-blue-500" />
                            Activity Heatmap
                        </CardTitle>
                        <CardDescription>Geographic distribution of organizational activities across Bangladesh.</CardDescription>
                    </CardHeader>
                    <CardContent className="h-[400px] relative">
                        {/* Simplified Heatmap Visualization */}
                        <div className="absolute inset-0 bg-accent/20 rounded-xl overflow-hidden flex items-center justify-center">
                            <div className="relative w-full h-full p-8 opacity-40 grayscale group-hover:grayscale-0 transition-all">
                                <svg viewBox="0 0 400 600" className="w-full h-full fill-blue-500/20 stroke-blue-500/30 stroke-2">
                                    <path d="M100,50 L150,30 L200,60 L250,40 L300,100 L350,150 L300,250 L250,400 L200,550 L100,500 L50,400 L30,250 L50,150 Z" />
                                </svg>
                                {/* Heat spots */}
                                <div className="absolute top-1/4 left-1/3 w-24 h-24 bg-blue-500/40 blur-3xl animate-pulse rounded-full" />
                                <div className="absolute top-1/2 left-1/2 w-16 h-16 bg-emerald-500/40 blur-3xl animate-pulse rounded-full" />
                                <div className="absolute bottom-1/4 right-1/3 w-20 h-20 bg-amber-500/40 blur-3xl animate-pulse rounded-full" />
                            </div>
                            <div className="absolute inset-x-8 bottom-8 p-4 bg-background/80 backdrop-blur-lg border border-border/50 rounded-xl flex justify-between items-center">
                                <div className="space-y-1">
                                    <p className="text-xs font-bold text-muted-foreground uppercase">High Activity Zones</p>
                                    <p className="text-sm font-bold">Dhaka, Chittagong, Sylhet</p>
                                </div>
                                <div className="flex -space-x-2">
                                    {[1, 2, 3, 4].map(i => (
                                        <div key={i} className="w-7 h-7 rounded-full border-2 border-background bg-accent" />
                                    ))}
                                </div>
                            </div>
                        </div>
                    </CardContent>
                </Card>

                <div className="space-y-8">
                    <Card className="border-border/50 bg-card/50">
                        <CardHeader>
                            <CardTitle className="text-sm font-bold uppercase tracking-wider text-muted-foreground">Growth Velocity</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-6">
                            <GrowthIndicator label="Membership" value="+4.2%" trend="up" />
                            <GrowthIndicator label="Activity Log" value="+12.5%" trend="up" />
                            <GrowthIndicator label="Financial Liquidity" value="-2.1%" trend="down" />
                        </CardContent>
                    </Card>

                    <Card className="border-border/50 bg-card/50">
                        <CardHeader>
                            <CardTitle className="text-sm font-bold uppercase tracking-wider text-muted-foreground">Top Performing Units</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <UnitPerformance rank={1} name="Dhaka City North" activities={142} color="bg-blue-500" />
                            <UnitPerformance rank={2} name="Chittagong South" activities={98} color="bg-emerald-500" />
                            <UnitPerformance rank={3} name="Sylhet Central" activities={84} color="bg-amber-500" />
                        </CardContent>
                    </Card>
                </div>
            </div>
        </div>
    )
}

function GrowthIndicator({ label, value, trend }: { label: string, value: string, trend: "up" | "down" }) {
    return (
        <div className="flex items-center justify-between">
            <span className="text-sm font-medium text-muted-foreground">{label}</span>
            <div className={cn(
                "flex items-center gap-1 text-sm font-black",
                trend === "up" ? "text-emerald-500" : "text-red-500"
            )}>
                {trend === "up" ? <ArrowUpRight className="w-4 h-4" /> : <ArrowDownRight className="w-4 h-4" />}
                {value}
            </div>
        </div>
    )
}

function UnitPerformance({ rank, name, activities, color }: { rank: number, name: string, activities: number, color: string }) {
    return (
        <div className="flex items-center gap-3 group translate-all hover:translate-x-1 cursor-default">
            <div className="w-6 h-6 rounded bg-accent flex items-center justify-center text-[10px] font-black text-muted-foreground">
                #{rank}
            </div>
            <div className="flex-1 space-y-1">
                <div className="flex justify-between text-xs font-bold">
                    <span>{name}</span>
                    <span className="text-muted-foreground">{activities} acts</span>
                </div>
                <div className="h-1 w-full bg-accent rounded-full overflow-hidden">
                    <div className={cn("h-full transition-all duration-1000", color)} style={{ width: `${Math.min(100, activities / 1.5)}%` }} />
                </div>
            </div>
        </div>
    )
}
