"use client"

import { useEffect, useState } from "react"
import { Search, CreditCard, Filter, Download, ArrowUpRight, ArrowDownLeft, Clock, Loader2, Landmark } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { cn } from "@/pkg/utils"
import { api } from "@/pkg/api"
import { toast } from "sonner"

export default function FinancePage() {
    const [loading, setLoading] = useState(true)
    const [statement, setStatement] = useState<any>(null)
    const [jurisdictionId, setJurisdictionId] = useState<string>("")

    useEffect(() => {
        async function loadInitialData() {
            try {
                // Get user's primary jurisdiction first
                const user = await api.getCurrentUser()
                if (user?.data?.jurisdiction_id) {
                    setJurisdictionId(user.data.jurisdiction_id)
                    const res = await api.getFinanceStatement(user.data.jurisdiction_id)
                    setStatement(res?.data)
                }
            } catch (err) {
                console.error("Failed to load finance data:", err)
                // Fallback or toast could go here
            } finally {
                setLoading(false)
            }
        }
        loadInitialData()
    }, [])

    if (loading) {
        return (
            <div className="flex flex-col items-center justify-center h-[60vh]">
                <Loader2 className="w-12 h-12 animate-spin text-blue-500 mb-4" />
                <p className="text-sm font-bold uppercase tracking-widest opacity-50">Auditing Ledger...</p>
            </div>
        )
    }

    const { balance = 0, transactions = [] } = statement || {}

    return (
        <div className="space-y-8">
            <div className="flex justify-between items-end">
                <div className="space-y-1">
                    <h1 className="text-3xl font-bold tracking-tight">Finance Ledger</h1>
                    <p className="text-muted-foreground">Immutable financial tracking for your jurisdiction.</p>
                </div>
                <div className="flex gap-3">
                    <Button variant="outline" className="gap-2">
                        <Download className="w-4 h-4" />
                        Export Statement
                    </Button>
                    <Button className="bg-blue-600 hover:bg-blue-700 gap-2">
                        <CreditCard className="w-4 h-4" />
                        Add Transaction
                    </Button>
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <BalanceCard title="Current Balance" amount={balance.toLocaleString()} type="primary" />
                <BalanceCard title="Monthly Donations" amount="+0" type="success" />
                <BalanceCard title="Monthly Expenses" amount="-0" type="danger" />
            </div>

            <Card className="border-border/50 bg-card/50">
                <CardHeader className="flex flex-row items-center justify-between">
                    <div>
                        <CardTitle>Recent Transactions</CardTitle>
                        <CardDescription>Verified financial movements in the ledger.</CardDescription>
                    </div>
                    <div className="flex gap-2">
                        <Button variant="ghost" size="sm" className="gap-2">
                            <Filter className="w-4 h-4" />
                            Filter
                        </Button>
                    </div>
                </CardHeader>
                <CardContent>
                    <div className="space-y-1">
                        <div className="grid grid-cols-12 px-4 py-2 text-xs font-bold text-muted-foreground uppercase tracking-wider">
                            <div className="col-span-1">Type</div>
                            <div className="col-span-5">Description</div>
                            <div className="col-span-2">Date</div>
                            <div className="col-span-2">Amount</div>
                            <div className="col-span-2">Status</div>
                        </div>
                        {transactions.length > 0 ? transactions.map((tx: any) => (
                            <div key={tx.id} className="grid grid-cols-12 px-4 py-4 rounded-xl hover:bg-accent/30 transition-colors border border-transparent hover:border-border/50 items-center">
                                <div className="col-span-1">
                                    {tx.type === 'income' ? (
                                        <div className="w-8 h-8 rounded-lg bg-emerald-500/10 flex items-center justify-center text-emerald-500">
                                            <ArrowDownLeft className="w-4 h-4" />
                                        </div>
                                    ) : (
                                        <div className="w-8 h-8 rounded-lg bg-red-500/10 flex items-center justify-center text-red-500">
                                            <ArrowUpRight className="w-4 h-4" />
                                        </div>
                                    )}
                                </div>
                                <div className="col-span-5">
                                    <p className="text-sm font-bold">{tx.description || tx.category_name}</p>
                                    <p className="text-xs text-muted-foreground">ID: {tx.id.substring(0, 8).toUpperCase()} • Ref: {tx.reference || 'Internal'}</p>
                                </div>
                                <div className="col-span-2 text-sm text-muted-foreground font-medium">
                                    {new Date(tx.created_at).toLocaleDateString()}
                                </div>
                                <div className={cn(
                                    "col-span-2 text-sm font-bold",
                                    tx.type === 'income' ? "text-emerald-500" : "text-red-500"
                                )}>
                                    {tx.type === 'income' ? "+" : "-"} Tk {tx.amount.toLocaleString()}
                                </div>
                                <div className="col-span-2">
                                    <div className="inline-flex items-center gap-1.5 px-2 py-1 rounded-full bg-blue-500/10 text-blue-500 text-[10px] font-bold uppercase border border-blue-500/20">
                                        <Clock className="w-3 h-3" />
                                        Ledgered
                                    </div>
                                </div>
                            </div>
                        )) : (
                            <div className="flex flex-col items-center justify-center py-12 opacity-30">
                                <Landmark className="w-12 h-12 mb-2" />
                                <p className="text-sm font-black uppercase tracking-widest">No Transactions Found</p>
                            </div>
                        )}
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}

function BalanceCard({ title, amount, type }: { title: string, amount: string, type: 'primary' | 'success' | 'danger' }) {
    const colors = {
        primary: "bg-blue-600 shadow-blue-900/40",
        success: "bg-emerald-600/10 text-emerald-500 border-emerald-500/20",
        danger: "bg-red-600/10 text-red-500 border-red-500/20"
    }

    return (
        <Card className={cn(
            "border-transparent overflow-hidden relative group transition-transform hover:-translate-y-1",
            type === 'primary' ? colors.primary : "bg-card/50 border-border/50"
        )}>
            {type === 'primary' && (
                <div className="absolute -right-4 -top-4 opacity-10 group-hover:scale-125 transition-transform">
                    <CreditCard className="w-32 h-32" />
                </div>
            )}
            <CardContent className="pt-6 space-y-2">
                <p className={cn(
                    "text-xs font-bold uppercase tracking-widest px-2 py-0.5 rounded-full w-fit",
                    type === 'primary' ? "bg-white/20 text-white" : "bg-accent/50"
                )}>
                    {title}
                </p>
                <div className="flex items-baseline gap-2">
                    <span className={cn("text-3xl font-black tracking-tighter", type === 'primary' && "text-white")}>
                        Tk {amount}
                    </span>
                    {type !== 'primary' && <span className="text-xs font-medium">this month</span>}
                </div>
            </CardContent>
        </Card>
    )
}
