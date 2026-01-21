"use client"

import { useEffect, useState } from "react"
import { UserPlus, XCircle, Clock, MapPin, UserCheck, ShieldAlert, Loader2 } from "lucide-react"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { cn } from "@/pkg/utils"
import { api } from "@/pkg/api"
import { toast } from "sonner"

export default function JoinRequestsPage() {
    const [filter, setFilter] = useState("pending")
    const [requests, setRequests] = useState<any[]>([])
    const [loading, setLoading] = useState(true)
    const [actionId, setActionId] = useState<string | null>(null)

    const fetchRequests = async () => {
        setLoading(true)
        try {
            // Check current user's jurisdiction in a real app, 
            // but for now we'll fetch general or filtered.
            // backend expects jurisdiction_id, we'll try to find a way to get it or fallback.
            const res = await api.getJoinRequests(`?status=${filter}&jurisdiction_id=00000000-0000-0000-0000-000000000000`)
            setRequests(res?.data || [])
        } catch (err) {
            toast.error("Failed to load join requests")
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        fetchRequests()
    }, [filter])

    const handleApprove = async (id: string) => {
        setActionId(id)
        try {
            await api.approveJoinRequest(id)
            toast.success("Application approved. Member created.")
            fetchRequests()
        } catch (err: any) {
            toast.error(err.message || "Approval failed")
        } finally {
            setActionId(null)
        }
    }

    const handleReject = async (id: string) => {
        const reason = prompt("Enter rejection reason:")
        if (!reason) return

        setActionId(id)
        try {
            await api.rejectJoinRequest(id, reason)
            toast.success("Application rejected")
            fetchRequests()
        } catch (err: any) {
            toast.error(err.message || "Rejection failed")
        } finally {
            setActionId(null)
        }
    }

    return (
        <div className="space-y-8">
            <div className="flex justify-between items-end">
                <div className="space-y-1">
                    <h1 className="text-3xl font-bold tracking-tight">Join Requests</h1>
                    <p className="text-muted-foreground">Review and process new membership applications (Form 'Ka').</p>
                </div>
                <div className="flex gap-2">
                    <div className="flex p-1 bg-accent/50 rounded-lg">
                        <button
                            onClick={() => setFilter("pending")}
                            className={cn(
                                "px-4 py-2 rounded-md text-sm font-medium transition-all",
                                filter === "pending" ? "bg-background shadow-sm text-foreground" : "text-muted-foreground"
                            )}
                        >
                            Pending
                        </button>
                        <button
                            onClick={() => setFilter("approved")}
                            className={cn(
                                "px-4 py-2 rounded-md text-sm font-medium transition-all",
                                filter === "approved" ? "bg-background shadow-sm text-foreground" : "text-muted-foreground"
                            )}
                        >
                            Processed
                        </button>
                    </div>
                </div>
            </div>

            <div className="grid grid-cols-1 gap-4">
                {loading ? (
                    <div className="flex flex-col gap-4">
                        {[1, 2, 3].map(i => <div key={i} className="h-32 bg-accent/20 animate-pulse rounded-xl" />)}
                    </div>
                ) : requests.length > 0 ? requests.map((req) => (
                    <Card key={req.id} className="border-border/50 bg-card/30 hover:bg-accent/10 transition-colors">
                        <CardContent className="p-0">
                            <div className="flex flex-col md:flex-row md:items-center">
                                <div className="p-6 flex-1 space-y-4">
                                    <div className="flex justify-between items-start">
                                        <div className="flex items-center gap-4">
                                            <div className="w-12 h-12 rounded-full bg-blue-600/10 flex items-center justify-center border border-blue-500/20">
                                                <UserPlus className="w-6 h-6 text-blue-500" />
                                            </div>
                                            <div className="space-y-0.5">
                                                <h3 className="text-xl font-bold">{req.full_name_en}</h3>
                                                <p className="text-sm font-medium text-muted-foreground">{req.full_name_bn}</p>
                                            </div>
                                        </div>
                                        <div className={cn(
                                            "px-3 py-1 text-[10px] font-black uppercase rounded-full border flex items-center gap-1.5",
                                            req.status === "pending" ? "bg-amber-500/10 text-amber-500 border-amber-500/20" : "bg-emerald-500/10 text-emerald-500 border-emerald-500/20"
                                        )}>
                                            {req.status === "pending" ? <Clock className="w-3 h-3" /> : <UserCheck className="w-3 h-3" />}
                                            {req.status}
                                        </div>
                                    </div>

                                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 text-sm">
                                        <div className="space-y-1">
                                            <p className="text-muted-foreground uppercase text-[10px] font-bold">Target Jurisdiction</p>
                                            <p className="font-semibold flex items-center gap-2">
                                                <MapPin className="w-4 h-4 text-blue-500" />
                                                {req.ward_id || "Unspecified Unit"}
                                            </p>
                                        </div>
                                        <div className="space-y-1">
                                            <p className="text-muted-foreground uppercase text-[10px] font-bold">National ID (NID)</p>
                                            <p className="font-semibold tracking-wider">{req.nid}</p>
                                        </div>
                                        <div className="space-y-1">
                                            <p className="text-muted-foreground uppercase text-[10px] font-bold">Phone</p>
                                            <p className="font-bold flex items-center gap-2">
                                                <ShieldAlert className="w-4 h-4 text-amber-500" />
                                                {req.phone}
                                            </p>
                                        </div>
                                    </div>
                                </div>

                                {req.status === "pending" && (
                                    <div className="p-6 bg-accent/20 border-t md:border-t-0 md:border-l border-border/50 flex md:flex-col gap-2 justify-center w-full md:w-auto">
                                        <Button
                                            onClick={() => handleApprove(req.id)}
                                            disabled={!!actionId}
                                            className="bg-emerald-600 hover:bg-emerald-700 font-bold gap-2"
                                        >
                                            {actionId === req.id ? <Loader2 className="w-4 h-4 animate-spin" /> : <UserCheck className="w-4 h-4" />}
                                            Approve
                                        </Button>
                                        <Button
                                            onClick={() => handleReject(req.id)}
                                            disabled={!!actionId}
                                            variant="outline"
                                            className="text-red-500 hover:text-red-600 hover:bg-red-500/10 border-red-500/20 gap-2"
                                        >
                                            {actionId === req.id ? <Loader2 className="w-4 h-4 animate-spin" /> : <XCircle className="w-4 h-4" />}
                                            Reject
                                        </Button>
                                    </div>
                                )}
                            </div>
                        </CardContent>
                    </Card>
                )) : (
                    <div className="text-center py-12 border border-dashed border-border rounded-xl">
                        <UserPlus className="w-12 h-12 text-muted-foreground mx-auto mb-4 opacity-20" />
                        <p className="text-muted-foreground font-medium">No join requests found in this category.</p>
                    </div>
                )}
            </div>
        </div>
    )
}
