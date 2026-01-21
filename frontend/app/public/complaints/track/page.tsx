"use client"

import { useState } from "react"
import { Search, History, Clock, CheckCircle, AlertCircle, ShieldEllipsis } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { cn } from "@/pkg/utils"

export default function ComplaintTrackPage() {
    const [trackId, setTrackId] = useState("")
    const [result, setResult] = useState<any>(null)
    const [loading, setLoading] = useState(false)

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault()
        setLoading(true)
        // Mock search
        setTimeout(() => {
            setResult({
                id: trackId,
                status: "under_review",
                submittedAt: "2026-01-21T10:30:00Z",
                jurisdiction: "Ward 5, Joypurhat",
                updates: [
                    { status: "received", time: "2026-01-21T10:30:00Z", note: "Complaint logged into system." },
                    { status: "under_review", time: "2026-01-21T14:45:00Z", note: "Assigned to District Investigation Committee." }
                ]
            })
            setLoading(false)
        }, 1000)
    }

    return (
        <div className="min-h-screen flex flex-col bg-background">
            <header className="p-6 flex justify-between items-center border-b border-border/50 backdrop-blur-md">
                <div className="flex items-center gap-3">
                    <ShieldEllipsis className="w-6 h-6 text-blue-500" />
                    <span className="font-bold text-xl tracking-tight">Track Your Report</span>
                </div>
                <Button variant="ghost" onClick={() => window.location.href = "/"}>Back</Button>
            </header>

            <main className="flex-1 max-w-4xl mx-auto w-full p-4 py-12 space-y-8">
                <div className="text-center space-y-4 max-w-2xl mx-auto">
                    <h1 className="text-4xl font-bold tracking-tight">Check Status</h1>
                    <p className="text-muted-foreground">
                        Enter the tracking ID provided during submission to see the progress of your complaint.
                    </p>
                    <form onSubmit={handleSearch} className="flex gap-2">
                        <div className="relative flex-1">
                            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                            <input
                                className="w-full h-12 pl-10 pr-4 rounded-xl border border-input bg-accent/50 focus:ring-2 focus:ring-blue-500 outline-none transition-all"
                                placeholder="Ex: C-X7Y2Z9"
                                value={trackId}
                                onChange={(e) => setTrackId(e.target.value)}
                                required
                            />
                        </div>
                        <Button type="submit" className="h-12 px-8 bg-blue-600 hover:bg-blue-700" disabled={loading}>
                            {loading ? "Searching..." : "Track Status"}
                        </Button>
                    </form>
                </div>

                {result && (
                    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                            <StatusCard
                                title="Current Status"
                                value={result.status.replace("_", " ")}
                                icon={<History className="w-5 h-5 text-blue-500" />}
                                className="capitalize"
                            />
                            <StatCard
                                title="Submitted On"
                                value={new Date(result.submittedAt).toLocaleDateString()}
                                icon={<Clock className="w-5 h-5 text-emerald-500" />}
                            />
                            <StatCard
                                title="Jurisdiction"
                                value={result.jurisdiction}
                                icon={<AlertCircle className="w-5 h-5 text-amber-500" />}
                            />
                        </div>

                        <Card className="border-border/50 bg-card/50 backdrop-blur-md">
                            <CardHeader>
                                <CardTitle>Investigation Timeline</CardTitle>
                                <CardDescription>Track the journey of your report as it moves through the hierarchy.</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div className="relative pl-8 space-y-8 before:absolute before:left-[15px] before:top-2 before:bottom-2 before:w-[2px] before:bg-border/50">
                                    {result.updates.map((update: any, i: number) => (
                                        <div key={i} className="relative">
                                            <div className={cn(
                                                "absolute -left-[25px] top-1 w-5 h-5 rounded-full border-4 border-background flex items-center justify-center",
                                                i === 0 ? "bg-emerald-500" : "bg-blue-500"
                                            )} />
                                            <div className="space-y-1">
                                                <p className="text-sm font-bold capitalize">{update.status.replace("_", " ")}</p>
                                                <p className="text-xs text-muted-foreground">{new Date(update.time).toLocaleString()}</p>
                                                <p className="text-sm text-foreground/80 mt-2 p-3 rounded-lg bg-accent/30 border border-border/50 italic">
                                                    &ldquo;{update.note}&rdquo;
                                                </p>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                )}
            </main>
        </div>
    )
}

function StatusCard({ title, value, icon, className }: any) {
    return (
        <Card className="border-blue-500/20 bg-blue-500/5 backdrop-blur-md">
            <CardContent className="pt-6">
                <div className="flex items-center gap-2 mb-2">
                    {icon}
                    <span className="text-sm font-medium text-blue-400">{title}</span>
                </div>
                <p className={cn("text-2xl font-bold text-blue-100", className)}>{value}</p>
            </CardContent>
        </Card>
    )
}

function StatCard({ title, value, icon }: any) {
    return (
        <Card className="border-border/50 bg-card/50">
            <CardContent className="pt-6">
                <div className="flex items-center gap-2 mb-2">
                    {icon}
                    <span className="text-sm font-medium text-muted-foreground">{title}</span>
                </div>
                <p className="text-xl font-semibold">{value}</p>
            </CardContent>
        </Card>
    )
}
