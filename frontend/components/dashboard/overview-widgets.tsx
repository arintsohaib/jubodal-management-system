"use client"
import { useEffect, useState } from "react"
import { motion } from "framer-motion"
import { Users, Activity, Landmark, Flag, ChevronUp, ChevronDown, ShieldCheck, UserPlus } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { cn } from "@/pkg/utils"
import { api } from "@/pkg/api"

export function DashboardOverview() {
    const [pulse, setPulse] = useState<any>(null)
    const [activities, setActivities] = useState<any[]>([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        async function loadData() {
            try {
                const [pulseData, activitiesData] = await Promise.all([
                    api.getPulse(),
                    api.getActivities("?page_size=5")
                ])
                setPulse(pulseData?.data || pulseData)
                setActivities(activitiesData?.data || [])
            } catch (err) {
                console.error("Failed to load dashboard data:", err)
            } finally {
                setLoading(setLoading as any) // Trick to avoid unused if needed, but actually:
                setLoading(false)
            }
        }
        loadData()
    }, [])

    const container = {
        hidden: { opacity: 0 },
        show: {
            opacity: 1,
            transition: {
                staggerChildren: 0.1
            }
        }
    }

    const item = {
        hidden: { opacity: 0, y: 20 },
        show: { opacity: 1, y: 0 }
    }

    if (loading) {
        return <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 animate-pulse">
            {[1, 2, 3, 4].map(i => <div key={i} className="h-32 bg-accent/20 rounded-xl" />)}
        </div>
    }

    return (
        <div className="space-y-8">
            {/* Upper Grid: Real-time Stats */}
            <motion.div
                variants={container}
                initial="hidden"
                animate="show"
                className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6"
            >
                <motion.div variants={item}>
                    <MetricCard
                        title="Total Members"
                        value={pulse?.total_members?.toLocaleString() || "0"}
                        trend="+12%"
                        trendUp={true}
                        icon={<Users className="w-5 h-5 text-blue-500" />}
                        description="Verified active members"
                    />
                </motion.div>
                <motion.div variants={item}>
                    <MetricCard
                        title="Active Activities"
                        value={pulse?.active_activities || "0"}
                        trend="+5"
                        trendUp={true}
                        icon={<Activity className="w-5 h-5 text-emerald-500" />}
                        description="In-progress field operations"
                    />
                </motion.div>
                <motion.div variants={item}>
                    <MetricCard
                        title="Donations (Jan)"
                        value={`Tk ${pulse?.total_donations?.toLocaleString() || "0"}`}
                        trend="+24%"
                        trendUp={true}
                        icon={<Landmark className="w-5 h-5 text-amber-500" />}
                        description="Total collections this month"
                    />
                </motion.div>
                <motion.div variants={item}>
                    <MetricCard
                        title="Pending Complaints"
                        value={pulse?.pending_complaints || "0"}
                        trend="-2"
                        trendUp={false}
                        icon={<Flag className="w-5 h-5 text-red-500" />}
                        description="Awaiting investigation"
                    />
                </motion.div>
            </motion.div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Main Feed: Recent Jurisdiction Activity */}
                <Card className="lg:col-span-2 border-border/50 bg-card/30 backdrop-blur-md">
                    <CardHeader>
                        <CardTitle>Jurisdictional Pulse</CardTitle>
                        <CardDescription>Real-time feed from units under your command.</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-6">
                        {activities.length > 0 ? activities.map((act: any) => (
                            <ActivityItem
                                key={act.id}
                                user={act.user_name || "Agent"}
                                role={act.category}
                                action="logged an activity"
                                detail={act.description}
                                time={new Date(act.created_at).toLocaleTimeString()}
                                type="activity"
                            />
                        )) : (
                            <p className="text-center text-muted-foreground py-8">No recent activities found.</p>
                        )}
                    </CardContent>
                </Card>

                {/* Action Widgets */}
                <div className="space-y-8">
                    <Card className="border-blue-500/20 bg-blue-500/5 backdrop-blur-sm">
                        <CardHeader className="pb-2">
                            <CardTitle className="text-sm font-bold uppercase tracking-wider text-blue-400">Quick Actions</CardTitle>
                        </CardHeader>
                        <CardContent className="grid grid-cols-2 gap-2">
                            <QuickActionButton icon={<UserPlus className="w-4 h-4" />} label="Review 'Ka'" href="/dashboard/join-requests" />
                            <QuickActionButton icon={<Activity className="w-4 h-4" />} label="Log Field" href="/dashboard/activities" />
                            <QuickActionButton icon={<Flag className="w-4 h-4" />} label="View Grievance" href="/dashboard/complaints" />
                            <QuickActionButton icon={<Landmark className="w-4 h-4" />} label="Ledger Entry" href="/dashboard/finance" />
                        </CardContent>
                    </Card>

                    <Card className="border-border/50 bg-card/50">
                        <CardHeader>
                            <CardTitle className="text-sm font-bold">Committee Strength</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <ProgressWidget label="Executive" current={pulse?.committee_stats?.executive || 0} max={151} color="bg-blue-500" />
                            <ProgressWidget label="Advisors" current={pulse?.committee_stats?.advisors || 0} max={21} color="bg-emerald-500" />
                            <ProgressWidget label="Units" current={pulse?.committee_stats?.units || 0} max={12} color="bg-amber-500" />
                        </CardContent>
                    </Card>
                </div>
            </div>
        </div>
    )
}

function MetricCard({ title, value, trend, trendUp, icon, description }: any) {
    return (
        <motion.div
            whileHover={{ scale: 1.02, translateY: -4 }}
            transition={{ type: "spring", stiffness: 400, damping: 10 }}
        >
            <Card className="border-border/50 bg-card/30 hover:bg-accent/10 transition-colors shadow-xl shadow-blue-900/5">
                <CardContent className="pt-6">
                    <div className="flex items-center justify-between mb-4">
                        <div className="p-2 bg-background rounded-xl border border-border/50">
                            {icon}
                        </div>
                        <div className={cn(
                            "flex items-center gap-1 text-xs font-bold",
                            trendUp ? "text-emerald-500" : "text-red-500"
                        )}>
                            {trendUp ? <ChevronUp className="w-3 h-3" /> : <ChevronDown className="w-3 h-3" />}
                            {trend}
                        </div>
                    </div>
                    <div className="space-y-1">
                        <h3 className="text-2xl font-black tracking-tighter">{value}</h3>
                        <p className="text-xs font-bold text-muted-foreground uppercase tracking-widest">{title}</p>
                    </div>
                    <p className="text-[10px] text-muted-foreground/50 mt-4 italic">{description}</p>
                </CardContent>
            </Card>
        </motion.div>
    )
}

function ActivityItem({ user, role, action, detail, time, type }: any) {
    const icons = {
        activity: <Activity className="w-4 h-4 text-emerald-500" />,
        alert: <ShieldCheck className="w-4 h-4 text-red-500" />,
        member: <UserPlus className="w-4 h-4 text-blue-500" />
    }
    return (
        <div className="flex gap-4 group">
            <div className="flex flex-col items-center gap-2">
                <div className="w-10 h-10 rounded-full bg-accent flex items-center justify-center border border-border/50 group-hover:border-blue-500/50 transition-colors">
                    {icons[type as keyof typeof icons]}
                </div>
                <div className="w-[1px] h-full bg-border/50" />
            </div>
            <div className="flex-1 pb-6">
                <div className="flex justify-between items-start">
                    <div className="space-y-0.5">
                        <p className="text-sm font-bold">
                            <span className="text-blue-400">{user}</span> ({role}) <span className="text-muted-foreground font-medium">{action}</span>
                        </p>
                        <p className="text-sm text-foreground/80 leading-relaxed">{detail}</p>
                    </div>
                    <span className="text-[10px] text-muted-foreground font-bold uppercase whitespace-nowrap">{time}</span>
                </div>
            </div>
        </div>
    )
}

function QuickActionButton({ icon, label, href }: any) {
    return (
        <motion.a
            href={href}
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            className="flex flex-col items-center justify-center p-3 rounded-lg border border-border/50 bg-background/50 hover:bg-accent/50 transition-all gap-2 text-center group"
        >
            <div className="w-8 h-8 rounded-full bg-accent flex items-center justify-center group-hover:bg-blue-600 group-hover:text-white transition-colors">
                {icon}
            </div>
            <span className="text-[10px] font-black uppercase tracking-tighter">{label}</span>
        </motion.a>
    )
}

function ProgressWidget({ label, current, max, color }: any) {
    const percent = Math.min(100, (current / max) * 100)
    return (
        <div className="space-y-2">
            <div className="flex justify-between text-[10px] font-black uppercase tracking-widest text-muted-foreground">
                <span>{label}</span>
                <span>{current} / {max}</span>
            </div>
            <div className="h-2 w-full bg-accent rounded-full overflow-hidden border border-border/50 p-[1px]">
                <div
                    className={cn("h-full rounded-full transition-all duration-1000", color)}
                    style={{ width: `${percent}%` }}
                />
            </div>
        </div>
    )
}
