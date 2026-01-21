"use client"

import { useEffect, useState } from "react"
import { JurisdictionTree } from "@/components/dashboard/jurisdiction-tree"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Plus, UserCog, UserMinus, ShieldCheck, Mail, Phone, Loader2, Users } from "lucide-react"
import { api } from "@/pkg/api"
import { toast } from "sonner"

export default function CommitteeManagementPage() {
    const [selectedJuris, setSelectedJuris] = useState<any>(null)
    const [committee, setCommittee] = useState<any>(null)
    const [members, setMembers] = useState<any[]>([])
    const [loading, setLoading] = useState(false)

    const handleSelectJurisdiction = async (node: any) => {
        setSelectedJuris(node)
        setLoading(true)
        try {
            // In our system, one jurisdiction has one active committee
            const committeesRes = await api.getCommittees()
            const active = committeesRes?.data?.find((c: any) => c.jurisdiction_id === node.id && c.status === 'active')

            if (active) {
                setCommittee(active)
                const membersRes = await api.getCommitteeMembers(active.id)
                setMembers(membersRes?.data || [])
            } else {
                setCommittee(null)
                setMembers([])
            }
        } catch (err) {
            toast.error("Failed to load committee details")
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="grid grid-cols-12 gap-8 h-[calc(100vh-140px)]">
            {/* Sidebar Tree */}
            <Card className="col-span-3 border-border/50 bg-card/30 backdrop-blur-md overflow-y-auto">
                <CardContent className="p-4">
                    <JurisdictionTree onSelect={handleSelectJurisdiction} activeId={selectedJuris?.id} />
                </CardContent>
            </Card>

            {/* Main Content */}
            <div className="col-span-9 space-y-6 overflow-y-auto pr-2">
                {!selectedJuris ? (
                    <div className="flex flex-col items-center justify-center h-full opacity-30">
                        <Users className="w-16 h-16 mb-4" />
                        <p className="text-xl font-bold uppercase tracking-widest">Select a Jurisdiction to Manage</p>
                    </div>
                ) : (
                    <>
                        <div className="flex justify-between items-start">
                            <div className="space-y-1">
                                <h1 className="text-3xl font-bold tracking-tight">{selectedJuris.name}</h1>
                                <CardDescription className="flex items-center gap-2">
                                    <ShieldCheck className="w-4 h-4 text-emerald-500" />
                                    {committee ? `${committee.type.toUpperCase()} Committee • Active` : "No active committee found"}
                                </CardDescription>
                            </div>
                            <div className="flex gap-3">
                                <Button variant="outline" className="gap-2">
                                    <UserCog className="w-4 h-4" />
                                    Manage Positions
                                </Button>
                                <Button className="bg-blue-600 hover:bg-blue-700 gap-2">
                                    <Plus className="w-4 h-4" />
                                    Assign Member
                                </Button>
                            </div>
                        </div>

                        {loading ? (
                            <div className="flex justify-center py-24">
                                <Loader2 className="w-8 h-8 animate-spin text-blue-500" />
                            </div>
                        ) : (
                            <>
                                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                                    {members.filter(m => m.position_rank <= 3).map(m => (
                                        <LeaderCard
                                            key={m.id}
                                            name={m.user_name}
                                            position={m.position_name}
                                            phone={m.phone || "+880..."}
                                            email={m.email || "N/A"}
                                        />
                                    ))}
                                    {members.filter(m => m.position_rank <= 3).length === 0 && (
                                        <p className="col-span-full text-center text-sm text-muted-foreground italic py-8 border border-dashed border-border rounded-xl">
                                            No leadership positions assigned yet.
                                        </p>
                                    )}
                                </div>

                                <Card className="border-border/50 bg-card/50">
                                    <CardHeader className="flex flex-row items-center justify-between">
                                        <div>
                                            <CardTitle>Committee Members</CardTitle>
                                            <CardDescription>All designated members of this jurisdiction.</CardDescription>
                                        </div>
                                        <div className="flex items-center gap-4">
                                            <input
                                                type="text"
                                                placeholder="Filter members..."
                                                className="h-9 px-3 rounded-md border border-input bg-background/50 text-sm outline-none focus:ring-1 focus:ring-blue-500"
                                            />
                                        </div>
                                    </CardHeader>
                                    <CardContent>
                                        <div className="space-y-4">
                                            {members.length > 0 ? members.map((m) => (
                                                <div key={m.id} className="flex items-center justify-between p-4 rounded-xl border border-border/50 bg-accent/20 hover:bg-accent/30 transition-colors">
                                                    <div className="flex items-center gap-4">
                                                        <div className="w-10 h-10 rounded-full bg-blue-600/10 flex items-center justify-center text-blue-500 font-bold">
                                                            {m.user_name?.[0] || 'M'}
                                                        </div>
                                                        <div className="space-y-0.5">
                                                            <p className="text-sm font-bold">{m.user_name}</p>
                                                            <p className="text-xs text-muted-foreground">{m.position_name} • Joined {new Date(m.joined_at).toLocaleDateString()}</p>
                                                        </div>
                                                    </div>
                                                    <div className="flex items-center gap-1">
                                                        <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-foreground">
                                                            <Mail className="w-4 h-4" />
                                                        </Button>
                                                        <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-red-500">
                                                            <UserMinus className="w-4 h-4" />
                                                        </Button>
                                                    </div>
                                                </div>
                                            )) : (
                                                <p className="text-center text-muted-foreground py-12 italic">The committee roster is currently empty.</p>
                                            )}
                                        </div>
                                    </CardContent>
                                </Card>
                            </>
                        )}
                    </>
                )}
            </div>
        </div>
    )
}

function LeaderCard({ name, position, phone, email }: any) {
    return (
        <Card className="border-blue-500/20 bg-blue-500/5 backdrop-blur-sm relative overflow-hidden group">
            <div className="absolute top-0 right-0 p-2 opacity-10 group-hover:opacity-20 transition-opacity">
                <ShieldCheck className="w-16 h-16 text-blue-500" />
            </div>
            <CardContent className="pt-6 space-y-4">
                <div className="space-y-1">
                    <p className="text-sm font-bold text-blue-400">{position}</p>
                    <h3 className="text-xl font-bold">{name}</h3>
                </div>
                <div className="space-y-2 text-sm text-muted-foreground">
                    <div className="flex items-center gap-2">
                        <Phone className="w-3 h-3" /> {phone}
                    </div>
                    <div className="flex items-center gap-2">
                        <Mail className="w-3 h-3" /> {email}
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
